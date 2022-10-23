
help:
	@printf "Usage commands:\n"
	@printf "  make <command>\n\n"
	@printf "Available commands:\n"
	@awk '/^[a-zA-Z\-0-9%:\\]+/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = $$1; \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			gsub(":+$$", "", helpCommand); \
			printf "  \x1b[32;01m%-35s\x1b[0m %s\n", helpCommand, helpMessage; \
		} \
	} { lastLine = $$0 }' $(MAKEFILE_LIST) | sort -u
	@printf "\n"

default: help

docker:
	@which docker || (echo "docker is not installed" && exit 1)

k3d:
	@which k3d || (echo "k3d is not installed" && exit 1)

kubectl:
	@which kubectl || (echo "kubectl is not installed" && exit 1)

## build the order service
build-order: docker
	cd services/backend/order; docker build -t order --target production .

## build the email service
build-email: docker
	cd services/backend/email; docker build -t email --target production .

## build the frontend service
build-frontend: docker
	cd services/frontend; docker build -t frontend --target production .

## build all services
build: build-order build-email build-frontend

## create a k3d cluster
cluster: k3d
	k3d cluster create --servers 1 --agents 3 --port 80:80@loadbalancer --port 443:443@loadbalancer

## delete the k3d cluster
delete-cluster: k3d
	k3d cluster delete

## import container images into the k3d cluster
import-container-images: k3d
	k3d image import frontend
	k3d image import order
	k3d image import email

## deploy opentelemetry operator
deploy-opentelemetry-operator: kubectl
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.0/cert-manager.yaml
	kubectl -n cert-manager rollout status deployment cert-manager --timeout=2m
	kubectl -n cert-manager rollout status deployment cert-manager-webhook --timeout=2m
	kubectl -n cert-manager rollout status deployment cert-manager-cainjector --timeout=2m
	kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml
	kubectl -n opentelemetry-operator-system rollout status deployment opentelemetry-operator-controller-manager --timeout=2m
	kubectl -n opentelemetry-operator-system get all

## deploy-opentelemetry collector
deploy-opentelemetry-collector: kubectl
	kubectl apply -f opentelemetry/collector.yaml

## deploy the services
deploy-services: kubectl
	kubectl apply -R -f k8s/

## synonym for delete-cluster
clean: delete-cluster

.PHONY: help default docker k3d kubectl build-order build-email build-frontend build cluster delete-cluster import-container-images deploy-opentelemetry-operator deploy-opentelemetry-collector deploy-services clean
