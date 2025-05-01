# ============================================
# Deployment - Full Infra
# ============================================

deploy-infra:
	@if [ -z "$$AWS_PROFILE" ]; then \
		read -p "⚠️  AWS_PROFILE is not set. Enter the AWS profile to use: " profile; \
		export AWS_PROFILE=$$profile; \
	fi && \
	echo "Using AWS_PROFILE=$$AWS_PROFILE" && \
	cd deployment/terraform && \
	bash deploy-infra.sh

deploy-app:
	@if [ -z "$$AWS_PROFILE" ]; then \
		read -p "⚠️  AWS_PROFILE is not set. Enter the AWS profile to use: " profile; \
		export AWS_PROFILE=$$profile; \
	fi && \
	echo "Using AWS_PROFILE=$$AWS_PROFILE" && \
	cd deployment/terraform/apps && \
	tofu init && \
	tofu apply -auto-approve

# ============================================
# Help
# ============================================

help:
	@echo "Makefile for Local Development and Infrastructure"
	@echo ""
	@echo "Local Development:"
	@echo "  make compose-up     Start local dev environment with docker-compose"
	@echo "  make compose-down   Stop and remove docker-compose containers"
	@echo ""
	@echo "Infrastructure:"
	@echo "  make deploy-infra     Run full Terraform deployment via deploy-infra.sh"
	@echo "  make deploy-app       Deploy application Helm releases from Terraform in apps/"
	@echo ""
	@echo "Utilities:"
	@echo "  make help           Show this help message"

.PHONY: compose-up compose-down deploy-infra deploy-app help