output "products_table_name" {
  description = "Name of the DynamoDB products table"
  value       = aws_dynamodb_table.products.name
}

output "orders_table_name" {
  description = "Name of the DynamoDB orders table"
  value       = aws_dynamodb_table.orders.name
}

output "orders_sns_arn" {
  description = "ARN of the SNS topic for orders"
  value       = aws_sns_topic.orders.arn
}

output "products_sqs_url" {
  description = "URL of the SQS queue for products"
  value       = aws_sqs_queue.products.url
}

output "region" {
  description = "AWS region"
  value       = var.aws_region
}

output "cluster_endpoint" {
  description = "EKS cluster endpoint"
  value       = module.eks.cluster_endpoint
}

output "cluster_certificate_authority_data" {
  description = "EKS cluster certificate authority data"
  value       = module.eks.cluster_certificate_authority_data
}

output "cluster_name" {
  description = "EKS cluster name"
  value       = module.eks.cluster_name
}

output "products_service_service_account_role_arn" {
  description = "Service account ARN for the products service"
  value       = module.products_service_irsa_role.iam_role_arn
}

output "orders_service_service_account_role_arn" {
  description = "Service account ARN for the orders service"
  value       = module.orders_service_irsa_role.iam_role_arn
}

output "products_worker_service_account_role_arn" {
  description = "Service account ARN for the products worker"
  value       = module.products_worker_irsa_role.iam_role_arn
}
