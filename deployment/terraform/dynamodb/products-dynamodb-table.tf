resource "aws_dynamodb_table" "products" {
  name           = "products"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "id"
  
  attribute {
    name = "id"
    type = "S"
  }

  tags = {
    Environment = "dev"
    Project     = "microservices"
  }
}