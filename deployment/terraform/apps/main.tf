data "terraform_remote_state" "eks" {
  backend = "local"
  config = {
    path = "../terraform.tfstate"
  }
}

locals {
  name   = format("%s-%s", var.environment, var.project_name)
  region = var.aws_region

  tags = {
    CreatedBy   = "csepulveda"
    Owner       = "cesar.sepulveda.b@gmail.com"
    Project     = var.project_name
    Environment = var.environment
    OpenTofu    = "true"
  }
}

resource "helm_release" "products-service" {
  name             = "products-service"
  chart            = "./charts/products-service"
  namespace        = "sample-store"
  create_namespace = true

  set {
    name  = "image.repository"
    value = "ghcr.io/csepulveda/products-service"
  }

  set {
    name  = "AWS_REGION"
    value = local.region
  }

  set {
    name  = "TEMPO_ENDPOINT"
    value = "tempo.tempo.svc.cluster.local:4318"
  }

  set {
    name  = "PRODUCTS_TABLE"
    value = data.terraform_remote_state.eks.outputs.products_table_name
  }

  set {
    name  = "serviceAccountAnnotations.eks\\.amazonaws\\.com/role-arn"
    value = data.terraform_remote_state.eks.outputs.products_service_service_account_role_arn
  }
}

resource "helm_release" "products-worker" {
  name             = "products-worker"
  chart            = "./charts/products-worker"
  namespace        = "sample-store"
  create_namespace = true

  set {
    name  = "image.repository"
    value = "ghcr.io/csepulveda/products-worker"
  }

  set {
    name  = "AWS_REGION"
    value = local.region
  }

  set {
    name  = "TEMPO_ENDPOINT"
    value = "tempo.tempo.svc.cluster.local:4318"
  }

  set {
    name  = "SQS_QUEUE_URL"
    value = data.terraform_remote_state.eks.outputs.products_sqs_url
  }

  set {
    name  = "PRODUCTS_TABLE"
    value = data.terraform_remote_state.eks.outputs.products_table_name
  }

  set {
    name  = "serviceAccountAnnotations.eks\\.amazonaws\\.com/role-arn"
    value = data.terraform_remote_state.eks.outputs.products_worker_service_account_role_arn
  }
}

resource "helm_release" "orders-service" {
  name             = "orders-service"
  chart            = "./charts/orders-service"
  namespace        = "sample-store"
  create_namespace = true

  set {
    name  = "image.repository"
    value = "ghcr.io/csepulveda/orders-service"
  }

  set {
    name  = "AWS_REGION"
    value = local.region
  }

  set {
    name  = "TEMPO_ENDPOINT"
    value = "tempo.tempo.svc.cluster.local:4318"
  }

  set {
    name  = "ORDERS_TOPIC_ARN"
    value = data.terraform_remote_state.eks.outputs.orders_sns_arn
  }

  set {
    name  = "ORDERS_TABLE"
    value = data.terraform_remote_state.eks.outputs.orders_table_name
  }

  set {
    name  = "serviceAccountAnnotations.eks\\.amazonaws\\.com/role-arn"
    value = data.terraform_remote_state.eks.outputs.orders_service_service_account_role_arn
  }

}


resource "helm_release" "ui-web" {
  name             = "ui-web"
  chart            = "./charts/ui-web"
  namespace        = "sample-store"
  create_namespace = true


  set {
    name  = "image.repository"
    value = "ghcr.io/csepulveda/ui-web"
  }

  set {
    name  = "OTEL_EXPORTER_OTLP_ENDPOINT"
    value = "http://tempo.tempo.svc.cluster.local:4318/v1/traces"
  }

  set {
    name  = "ORDER_API_BASE_URL"
    value = "http://orders-service.sample-store.svc.cluster.local:8080"
  }

  set {
    name  = "PRODUCT_API_BASE_URL"
    value = "http://products-service.sample-store.svc.cluster.local:8080"
  }
}


resource "helm_release" "ui-backoffice" {
  name             = "ui-backoffice"
  chart            = "./charts/ui-backoffice"
  namespace        = "sample-store"
  create_namespace = true

  set {
    name  = "image.repository"
    value = "ghcr.io/csepulveda/ui-backoffice"
  }

  set {
    name  = "OTEL_EXPORTER_OTLP_ENDPOINT"
    value = "http://tempo.tempo.svc.cluster.local:4318/v1/traces"
  }

  set {
    name  = "ORDER_API_BASE_URL"
    value = "http://orders-service.sample-store.svc.cluster.local:8080"
  }

  set {
    name  = "PRODUCT_API_BASE_URL"
    value = "http://products-service.sample-store.svc.cluster.local:8080"
  }

}
