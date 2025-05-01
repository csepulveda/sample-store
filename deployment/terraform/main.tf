################################################################################
# Locals configurations
################################################################################
data "aws_availability_zones" "available" {}

locals {
  name   = format("%s-%s", var.environment, var.project_name)
  region = var.aws_region

  vpc_cidr = var.vpc_cidr
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)

  tags = {
    CreatedBy   = "csepulveda"
    Owner       = "cesar.sepulveda.b@gmail.com"
    Project     = var.project_name
    Environment = var.environment
    OpenTofu    = "true"
  }
}

################################################################################
# VPC Module
################################################################################
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.19.0"

  name = local.name
  cidr = local.vpc_cidr

  azs             = local.azs
  private_subnets = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 4, k)]
  public_subnets  = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 48)]
  intra_subnets   = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 52)]

  enable_nat_gateway = true
  single_nat_gateway = true

  public_subnet_tags = {
    "kubernetes.io/role/elb" = 1
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = 1
    "karpenter.sh/discovery"          = local.name
  }

  tags = local.tags
}

################################################################################
# EKS Module
################################################################################
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "20.35.0"

  cluster_name    = local.name
  cluster_version = 1.32

  enable_cluster_creator_admin_permissions = true
  cluster_endpoint_public_access           = true

  cluster_addons = {
    coredns                = {}
    eks-pod-identity-agent = {}
    kube-proxy             = {}
    vpc-cni                = {}
    metrics-server         = {}
  }

  vpc_id                   = module.vpc.vpc_id
  subnet_ids               = module.vpc.private_subnets
  control_plane_subnet_ids = module.vpc.intra_subnets

  eks_managed_node_groups = {
    karpenter = {
      ami_type       = "BOTTLEROCKET_x86_64"
      instance_types = ["t2.medium"]

      min_size     = 2
      max_size     = 3
      desired_size = 2

      labels = {
        # Used to ensure Karpenter runs on nodes that it does not manage
        "karpenter.sh/controller" = "true"
      }
    }
  }

  node_security_group_additional_rules = {
    allow-all-80-traffic-from-loadbalancers = {
      cidr_blocks = module.vpc.private_subnets_cidr_blocks
      description = "Allow all traffic from load balancers"
      from_port   = 80
      to_port     = 80
      protocol    = "TCP"
      type        = "ingress"
    }
    allow-all-443-traffic-from-loadbalancers = {
      cidr_blocks = module.vpc.private_subnets_cidr_blocks
      description = "Allow all traffic from load balancers"
      from_port   = 443
      to_port     = 443
      protocol    = "TCP"
      type        = "ingress"
    }
  }

  node_security_group_tags = merge(local.tags, {
    "karpenter.sh/discovery" = local.name
  })

  tags       = local.tags
  depends_on = [module.vpc]
}

################################################################################
# Karpenter
################################################################################
module "karpenter" {
  source = "terraform-aws-modules/eks/aws//modules/karpenter"

  cluster_name          = module.eks.cluster_name
  enable_v1_permissions = true

  # Name needs to match role name passed to the EC2NodeClass
  node_iam_role_use_name_prefix   = false
  node_iam_role_name              = local.name
  create_pod_identity_association = true

  node_iam_role_additional_policies = {
    AmazonSSMManagedInstanceCore = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
  }

  tags = local.tags

  depends_on = [module.eks]
}

################################################################################
# Karpenter Helm chart & manifests
################################################################################
resource "helm_release" "karpenter" {
  namespace           = "kube-system"
  name                = "karpenter"
  repository          = "oci://public.ecr.aws/karpenter"
  repository_username = data.aws_ecrpublic_authorization_token.token.user_name
  repository_password = data.aws_ecrpublic_authorization_token.token.password
  chart               = "karpenter"
  version             = "1.4.0"
  wait                = true

  values = [
    <<-EOT
    nodeSelector:
      karpenter.sh/controller: 'true'
    dnsPolicy: Default
    settings:
      clusterName: ${module.eks.cluster_name}
      clusterEndpoint: ${module.eks.cluster_endpoint}
      interruptionQueue: ${module.karpenter.queue_name}
    webhook:
      enabled: false
    EOT
  ]

  depends_on = [module.karpenter]
}

resource "kubernetes_manifest" "ec2_node_class" {
  manifest = {
    apiVersion = "karpenter.k8s.aws/v1"
    kind       = "EC2NodeClass"
    metadata = {
      name = "default"
    }
    spec = {
      amiSelectorTerms = [
        {
          alias = "bottlerocket@latest"
        }
      ]
      role = "${local.name}"
      subnetSelectorTerms = [
        {
          tags = {
            "karpenter.sh/discovery" = "${local.name}"
          }
        }
      ]
      securityGroupSelectorTerms = [
        {
          tags = {
            "karpenter.sh/discovery" = "${local.name}"
          }
        }
      ]
      tags = {
        "karpenter.sh/discovery" = "${local.name}"
      }
    }
  }
  depends_on = [helm_release.karpenter]
}

resource "kubernetes_manifest" "node_pool" {
  manifest = {
    apiVersion = "karpenter.sh/v1"
    kind       = "NodePool"
    metadata = {
      name = "default"
    }
    spec = {
      template = {
        spec = {
          nodeClassRef = {
            group = "karpenter.k8s.aws"
            kind  = "EC2NodeClass"
            name  = "default"
          }
          requirements = [
            {
              key      = "karpenter.sh/capacity-type"
              operator = "In"
              values   = ["spot"]
            },
            {
              key      = "karpenter.k8s.aws/instance-category"
              operator = "In"
              values   = ["t", "c", "m", "r"]
            },
            {
              key      = "karpenter.k8s.aws/instance-cpu"
              operator = "In"
              values   = ["4", "8", "16", "32"]
            },
            {
              key      = "karpenter.k8s.aws/instance-hypervisor"
              operator = "In"
              values   = ["nitro"]
            },
            {
              key      = "karpenter.k8s.aws/instance-generation"
              operator = "Gt"
              values   = ["2"]
            }
          ]
        }
      }
      limits = {
        cpu = 1000
      }
      disruption = {
        consolidationPolicy = "WhenEmpty"
        consolidateAfter    = "30s"
      }
    }
  }

  depends_on = [helm_release.karpenter]
}

################################################################################
# EBS Role
################################################################################
module "ebs_cni_irsa_role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"

  role_name = format("%s-aws-ebs-csi-driver", local.name)

  attach_ebs_csi_policy = true

  oidc_providers = {
    main = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = ["kube-system:ebs-csi-controller-sa"]
    }
  }

  depends_on = [module.eks]
}

resource "aws_eks_addon" "aws_ebs_csi_driver" {
  cluster_name             = local.name
  addon_name               = "aws-ebs-csi-driver"
  service_account_role_arn = module.ebs_cni_irsa_role.iam_role_arn
  configuration_values = jsonencode({
    defaultStorageClass = {
      enabled = true
    }
  })

  depends_on = [module.ebs_cni_irsa_role]
}

# resource "kubernetes_manifest" "gp2_storage_class" {
#   manifest = {
#     apiVersion = "storage.k8s.io/v1"
#     kind       = "StorageClass"
#     metadata = {
#       name = "gp2"
#       annotations = {
#         "storageclass.kubernetes.io/is-default-class" = "true"
#       }
#     }
#     provisioner = "kubernetes.io/aws-ebs"
#     parameters = {
#       type   = "gp2"
#       fsType = "ext4"
#     }
#     reclaimPolicy       = "Delete"
#     volumeBindingMode   = "WaitForFirstConsumer"
#     allowVolumeExpansion = true
#   }

#   depends_on = [ module.ebs_cni_irsa_role ]
# }

################################################################################
# APP resources Dynamo
################################################################################
resource "aws_dynamodb_table" "products" {
  name         = format("%s-%s", local.name, "products")
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  tags = local.tags
}

resource "aws_dynamodb_table" "orders" {
  name         = format("%s-%s", local.name, "orders")
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  tags = local.tags
}

################################################################################
# APP resources SNS and SQS
################################################################################
resource "aws_sns_topic" "orders" {
  name = format("%s-%s", local.name, "orders-topic")
}

resource "aws_sqs_queue" "products" {
  name = format("%s-%s", local.name, "products-queue")
}

resource "aws_sns_topic_subscription" "products_subscription" {
  topic_arn = aws_sns_topic.orders.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.products.arn
}

resource "aws_sqs_queue_policy" "allow_sns" {
  queue_url = aws_sqs_queue.products.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = "*"
      Action    = "sqs:SendMessage"
      Resource  = aws_sqs_queue.products.arn
      Condition = {
        ArnEquals = {
          "aws:SourceArn" = aws_sns_topic.orders.arn
        }
      }
    }]
  })
}

################################################################################
# Load Balancer
################################################################################
module "alb_irsa_role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"

  role_name                              = format("%s-aws-load-balancer-controller", local.name)
  attach_load_balancer_controller_policy = true

  oidc_providers = {
    main = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = ["kube-system:aws-load-balancer-controller"]
    }
  }
  tags       = local.tags
  depends_on = [module.eks]
}

resource "helm_release" "aws-load-balancer-controller" {
  namespace  = "kube-system"
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  wait       = true

  values = [
    <<-EOT
    clusterName: ${local.name}
    serviceAccount:
      create: true
      name: aws-load-balancer-controller
      annotations:
        eks.amazonaws.com/role-arn: ${module.alb_irsa_role.iam_role_arn}
    region: ${local.region}
    vpcId: ${module.vpc.vpc_id}
    controller:
      extraArgs:
        vpc-id: ${module.vpc.vpc_id}
    EOT
  ]
  depends_on = [module.alb_irsa_role]
}


################################################################################
# Tempo, loki and Thanos S3 configurations.
################################################################################
resource "aws_s3_bucket" "tempo" {
  bucket = "${local.name}-tempo"

  tags = local.tags
}

resource "aws_s3_bucket" "loki" {
  bucket = "${local.name}-loki"

  tags = local.tags
}

resource "aws_s3_bucket" "thanos" {
  bucket = "${local.name}-thanos"

  tags = local.tags
}

resource "aws_iam_policy" "thanos_s3_access" {
  name        = "thanos-s3-access-policy"
  description = "Allow Thanos access to S3 bucket for Get, Put, Delete operations"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = ["${aws_s3_bucket.thanos.arn}/*", "${aws_s3_bucket.thanos.arn}"]
      }
    ]
  })
}

resource "aws_iam_policy" "tempo_s3_access" {
  name        = "tempo-s3-access-policy"
  description = "Allow Tempo access to S3 bucket for Get, Put, Delete operations"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = ["${aws_s3_bucket.tempo.arn}/*", "${aws_s3_bucket.tempo.arn}"]
      }
    ]
  })
}

resource "aws_iam_policy" "loki_s3_access" {
  name        = "loki-s3-access-policy"
  description = "Allow Loki access to S3 bucket for Get, Put, Delete operations"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = ["${aws_s3_bucket.loki.arn}/*", "${aws_s3_bucket.loki.arn}"]
      }
    ]
  })
}

module "loki_irsa_role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"

  role_name = format("%s-loki", local.name)

  role_policy_arns = {
    policy = aws_iam_policy.loki_s3_access.arn
  }

  oidc_providers = {
    main = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = ["loki:loki"]
    }
  }

  tags       = local.tags
  depends_on = [module.eks]

}

module "tempo_irsa_role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"

  role_name = format("%s-tempo", local.name)

  role_policy_arns = {
    policy = aws_iam_policy.tempo_s3_access.arn
  }

  oidc_providers = {
    main = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = ["tempo:tempo"]
    }
  }

  tags       = local.tags
  depends_on = [module.eks]

}

module "thanos_irsa_role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"

  role_name = format("%s-thanos", local.name)

  role_policy_arns = {
    policy = aws_iam_policy.thanos_s3_access.arn
  }

  oidc_providers = {
    main = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = ["prometheus:prometheus-kube-prometheus-prometheus"]
    }
  }

  tags       = local.tags
  depends_on = [module.eks]

}

################################################################################
# Prometheus installation
################################################################################
resource "helm_release" "prometheus_stack" {
  name             = "prometheus"
  namespace        = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  create_namespace = true
  wait             = true
  wait_for_jobs    = true

  values = [
    <<-EOT
    kubeApiServer:
      enabled: false

    kubeEtcd:
      enabled: false

    prometheus:
      serviceAccount:
        annotations:
          eks.amazonaws.com/role-arn: ${module.thanos_irsa_role.iam_role_arn}

      thanosService:
        enabled: true

      thanosServiceMonitor:
        enabled: true

      prometheusSpec:
        thanos:
          image: "quay.io/thanos/thanos:v0.38.0"
          objectStorageConfig:
            secret:
              type: "S3"
              config:
                bucket: "${local.name}-thanos"
                region: "${local.region}"
                endpoint: "s3.${local.region}.amazonaws.com"
        storageSpec:
          volumeClaimTemplate:
            spec:
              resources:
                requests:
                  storage: "5Gi"

    grafana:
      additionalDataSources:
        - name: "Tempo"
          type: "tempo"
          url: "http://tempo.tempo.svc.cluster.local:3100"
          access: "proxy"
        - name: "Loki"
          type: "loki"
          url: "http://loki-gateway.loki.svc.cluster.local"
          access: "proxy"
    EOT
  ]

  depends_on = [module.thanos_irsa_role]

}

################################################################################
# Promtail installation
################################################################################
resource "helm_release" "promtail" {
  name             = "promtail"
  namespace        = "promtail"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "promtail"
  create_namespace = true
  wait             = true

  values = [
    <<-EOT
    config:
      clients:
        - url: http://loki-gateway.loki.svc.cluster.local/loki/api/v1/push
    EOT
  ]

  depends_on = [module.eks]
}

################################################################################
# Loki installation
################################################################################
resource "helm_release" "loki" {
  name             = "loki"
  namespace        = "loki"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "loki"
  create_namespace = true
  wait             = true

  values = [
    <<-EOT
    gateway:
      replicas: 1
    write:
      replicas: 2
    backend:
      replicas: 1
    read:
      replicas: 1
    deploymentMode: SimpleScalable
    tableManager:
      retention_deletes_enabled: true
      retention_period: 30d
    loki:
      auth_enabled: false
      storage:
        type: s3
        s3:
          s3: "s3://${local.region}/${local.name}-loki"
          region: "${var.aws_region}"
        bucketNames:
          chunks: "${local.name}-loki"
          ruler: "${local.name}-loki"
          admin: "${local.name}-loki"
      commonConfig:
        replication_factor: 2
      schemaConfig:
        configs:
        - from: 2022-01-11
          index:
            period: 24h
            prefix: loki_index_
          object_store: s3
          schema: v12
          store: boltdb-shipper
        - from: 2024-10-09
          index:
            period: 24h    
            prefix: loki_index_
          object_store: s3    
          schema: v13
          store: tsdb

    serviceAccount:
      annotations:
        eks.amazonaws.com/role-arn: "${module.loki_irsa_role.iam_role_arn}"

    serviceMonitor:
      enabled: true

    EOT
  ]

  depends_on = [module.loki_irsa_role]

}

################################################################################
# Tempo installation
################################################################################
resource "helm_release" "tempo" {
  name             = "tempo"
  namespace        = "tempo"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "tempo"
  create_namespace = true
  wait             = true

  values = [
    <<-EOT
    serviceAccount:
      annotations:
        eks.amazonaws.com/role-arn: "${module.tempo_irsa_role.iam_role_arn}"

    serviceMonitor:
      enabled: true

    tempo:
      storage:
        trace:
          backend: s3
          s3:
            bucket: "${local.name}-tempo"
            region: "${var.aws_region}"
            endpoint: "s3.${var.aws_region}.amazonaws.com"
    EOT
  ]
  depends_on = [module.tempo_irsa_role]

}


################################################################################
# Applications service accounts
################################################################################
resource "aws_iam_policy" "products_service" {
  name        = "products-service-policy"
  description = "Allow products service access to DynamoDB and SQS"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:PutItem",
          "dynamodb:GetItem",
          "dynamodb:UpdateItem",
        ]
        Resource = [
          aws_dynamodb_table.products.arn
        ]
      }
    ]
  })
}

resource "aws_iam_policy" "products_worker" {
  name        = "products-worker-policy"
  description = "Allow products worker access to DynamoDB and SQS"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "dynamodb:PutItem",
          "dynamodb:GetItem",
          "dynamodb:UpdateItem",
        ]
        Resource = [
          aws_sqs_queue.products.arn,
          aws_dynamodb_table.products.arn
        ]
      }
    ]
  })
}

resource "aws_iam_policy" "orders_service" {
  name        = "orders-service-policy"
  description = "Allow orders service access to DynamoDB and SNS"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:PutItem",
          "dynamodb:GetItem",
          "dynamodb:UpdateItem",
          "sns:Publish",
        ]
        Resource = [
          aws_dynamodb_table.orders.arn,
          aws_sns_topic.orders.arn
        ]
      }
    ]
  })
}


module "products_service_irsa_role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"

  role_name = format("%s-%s", local.name, "products-service")
  role_policy_arns = {
    policy = aws_iam_policy.products_service.arn
  }

  oidc_providers = {
    main = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = ["sample-store:products-service"]
    }
  }

  tags       = local.tags
  depends_on = [module.eks]
}

module "products_worker_irsa_role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"

  role_name = format("%s-%s", local.name, "products-worker")
  role_policy_arns = {
    policy = aws_iam_policy.products_worker.arn
  }

  oidc_providers = {
    main = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = ["sample-store:products-worker"]
    }
  }

  tags       = local.tags
  depends_on = [module.eks]
}

module "orders_service_irsa_role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"

  role_name = format("%s-%s", local.name, "orders-service")
  role_policy_arns = {
    policy = aws_iam_policy.orders_service.arn
  }


  oidc_providers = {
    main = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = ["sample-store:orders-service"]
    }
  }

  tags       = local.tags
  depends_on = [module.eks]
}