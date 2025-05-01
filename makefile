# ============================================
# Local Development - Docker Compose
# ============================================

COMPOSE_FILE=docker-compose.yaml

compose-up:
	docker-compose -f $(COMPOSE_FILE) up --build

compose-down:
	docker-compose -f $(COMPOSE_FILE) down

# ============================================
# Deployment - Full Infra
# ============================================

deploy-all:
	@if [ -z "$$AWS_PROFILE" ]; then \
		read -p "⚠️  AWS_PROFILE is not set. Enter the AWS profile to use: " profile; \
		export AWS_PROFILE=$$profile; \
	fi && \
	echo "Using AWS_PROFILE=$$AWS_PROFILE" && \
	cd deployment/terraform && \
	bash deploy-all.sh

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
	@echo "  make deploy-all     Run full Terraform deployment via deploy-all.sh"
	@echo ""
	@echo "Utilities:"
	@echo "  make help           Show this help message"

.PHONY: compose-up compose-down deploy-all help