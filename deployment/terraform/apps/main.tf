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
  name  = "products-service"
  chart = "./charts/products-service"

  set {
    name  = "image.repository"
    value = aws_ecr_repository.products-service.repository_url
  }

  set {
    name  = "AWS_REGION"
    value = local.region
  }

  set {
    name  = "TEMPO_ENDPOINT"
    value = "http://tempo.tempo.svc.cluster.local:4317"
  }

  set {
    name  = "PRODUCTS_TABLE"
    value = data.terraform_remote_state.eks.outputs.products_table_name
  }
}

resource "helm_release" "products-worker" {
  name  = "products-worker"
  chart = "./charts/products-worker"

  set {
    name  = "image.repository"
    value = aws_ecr_repository.products-worker.repository_url
  }

  set {
    name  = "AWS_REGION"
    value = local.region
  }

  set {
    name  = "TEMPO_ENDPOINT"
    value = "http://tempo.tempo.svc.cluster.local:4317"
  }

  set {
    name  = "SQS_QUEUE_URL"
    value = data.terraform_remote_state.eks.outputs.products_sqs_url
  }

  set {
    name  = "PRODUCTS_TABLE"
    value = data.terraform_remote_state.eks.outputs.products_table_name
  }
}

resource "helm_release" "orders-service" {
  name  = "orders-service"
  chart = "./charts/orders-service"

  set {
    name  = "image.repository"
    value = aws_ecr_repository.orders-service.repository_url
  }

  set {
    name  = "AWS_REGION"
    value = local.region
  }

  set {
    name  = "TEMPO_ENDPOINT"
    value = "http://tempo.tempo.svc.cluster.local:4317"
  }

  set {
    name  = "ORDERS_TOPIC_ARN"
    value = data.terraform_remote_state.eks.outputs.orders_sns_arn
  }

  set {
    name  = "ORDERS_TABLE"
    value = data.terraform_remote_state.eks.outputs.orders_table_name
  }
}


