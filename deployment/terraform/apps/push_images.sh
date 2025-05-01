#!/bin/bash

## setup buildx
aws_region=$(aws configure get region)
aws_account_id=$(aws sts get-caller-identity --query Account --output text)
export AWS_REGION=${aws_region}
export AWS_ACCOUNT_ID=${aws_account_id}

docker buildx create --use --name multiarch

aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com
# This script is used to push images to the ECR repository
productsService=$(echo aws_ecr_repository.products-service.repository_url | tofu console | tr -d '"' )
productsWorker=$(echo aws_ecr_repository.products-worker.repository_url | tofu console | tr -d '"' )
ordersService=$(echo aws_ecr_repository.orders-service.repository_url | tofu console | tr -d '"' )


docker buildx build --platform=linux/amd64,linux/arm64 -t ${productsService} ../../../apps/products-service/ --push
docker buildx build --platform=linux/amd64,linux/arm64 -t ${productsWorker} ../../../apps/products-worker/ --push
docker buildx build --platform=linux/amd64,linux/arm64 -t ${ordersService} ../../../apps/orders-service/ --push