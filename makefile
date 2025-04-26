
TF_DIR=deployment/terraform
AWS_REGION=us-west-2
AWS_ACCOUNT_ID=$(shell aws sts get-caller-identity --query "Account" --output text)
ECR_REPOSITORY=products-service
ECR_URI=$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
IMAGE_TAG=latest
DOCKER_IMAGE=$(ECR_URI)/$(ECR_REPOSITORY):$(IMAGE_TAG)
HELM_DIR=charts/products-service
HELM_RELEASE=products-service
NAMESPACE=microservices-system

tf-init:
	cd $(TF_DIR) && tofu init

tf-plan:
	cd $(TF_DIR) && tofu plan

tf-apply:
	cd $(TF_DIR) && tofu apply -auto-approve

tf-destroy:
	cd $(TF_DIR) && tofu destroy -auto-approve

tf-validate:
	cd $(TF_DIR) && tofu validate

tf-fmt:
	cd $(TF_DIR) && tofu fmt

# ============================================
# Docker
# ============================================

docker-build:
	docker build -t $(DOCKER_IMAGE) apps/products-service/

docker-login:
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR_URI)

docker-push: docker-login docker-build
	docker push $(DOCKER_IMAGE)

# ============================================
# Helm
# ============================================

helm-install:
	helm install $(HELM_RELEASE) $(HELM_DIR) --namespace $(NAMESPACE) --create-namespace

helm-upgrade:
	helm upgrade $(HELM_RELEASE) $(HELM_DIR) --namespace $(NAMESPACE)

helm-uninstall:
	helm uninstall $(HELM_RELEASE) --namespace $(NAMESPACE)

# ============================================
# Utilities
# ============================================

clean:
	rm -rf $(TF_DIR)/.terraform $(TF_DIR)/.terraform.lock.hcl

help:
	@echo "Global Makefile - Usage:"
	@echo ""
	@echo "Infraestructura:"
	@echo "  make tf-init        - Initialize OpenTofu"
	@echo "  make tf-plan        - Plan OpenTofu deployment"
	@echo "  make tf-apply       - Apply OpenTofu changes"
	@echo "  make tf-destroy     - Destroy OpenTofu resources"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-login   - Login to AWS ECR"
	@echo "  make docker-push    - Push Docker image to ECR"
	@echo ""
	@echo "Helm:"
	@echo "  make helm-install   - Install Helm release"
	@echo "  make helm-upgrade   - Upgrade Helm release"
	@echo "  make helm-uninstall - Uninstall Helm release"
	@echo ""
	@echo "Utilidades:"
	@echo "  make clean          - Clean .terraform and locks"
	@echo "  make help           - Show this help message"

.PHONY: tf-init tf-plan tf-apply tf-destroy tf-validate tf-fmt docker-build docker-login docker-push helm-install helm-upgrade helm-uninstall clean help