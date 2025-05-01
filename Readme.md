# Microservices EKS Demo Project

This project is a microservices-based architecture designed to run entirely on Kubernetes (EKS), including:

- API microservices written in Go using Fiber
- Frontend applications using Next.js
- AWS DynamoDB for data storage
- SNS/SQS for event-driven communication
- OpenTelemetry for observability
- Terraform (OpenTofu) for infrastructure management
- Helm for Kubernetes deployments

## Services

- products-service: Manages the product catalog.
- orders-service: Manages customer orders.
- shipments-service: Manages shipping and order delivery.
- ui-web: Frontend application for customers.
- ui-backoffice: Admin panel for internal management.

## Architecture Overview

- Backend services written in Go using Fiber
- Communication between services through SNS/SQS
- Storage in DynamoDB tables (one table per service)
- Kubernetes cluster (EKS) for orchestration
- Observability through OpenTelemetry, Loki, Tempo
- Continuous deployment via Helm charts
- Infrastructure managed via OpenTofu (Terraform-compatible)

## Project Structure

```
apps/
  products-service/
  orders-service/
  shipments-service/
  ui-web/
  ui-backoffice/
deployment/
  terraform/
  charts/
observability/
  opentelemetry/
  grafana/
docs/
scripts/
```

## Requirements

- Docker
- OpenTofu (tofu)
- AWS CLI configured with appropriate permissions
- kubectl
- Helm

## Quick Start

```bash
make tf-init
make tf-apply
make docker-build
make docker-push
make helm-install
```

## License

This project is licensed under the MIT License.

## Terraform + Helm.
### Terraform:
The EKS cluster (vpc, eks, karpenter)
The SNS, SQS and dynamo resources.
The IAM roles to allow pods use AWS services.
### Helm:
Deploy the application
Deploy the monitoring Stack (prometheus, grafana, tempo, loki, promtail)

## KRO + helm
### KRO
The EKS cluster (vpc, eks, karpenter)
The SNS, SQS and dynamo resources.
The IAM roles to allow pods use AWS services.
Deploy the application
### Helm:
Deploy the monitoring Stack (prometheus, grafana, tempo, loki, promtail)
