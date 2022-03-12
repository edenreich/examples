
.PHONY: load
load:
	kubectl --namespace development run --rm -it --image=nginx:alpine temp -- bin/sh -c \
		"while true; do curl -m 5 -sL -X POST http://order:8080/orders; sleep 0.5s; done;"

.PHONY: clean
clean:
	k3d cluster delete

.PHONY: build_services
build_services:
	cd order && docker build -t localhost:5432/order:development --target development .
	docker push localhost:5432/order:development
	cd notification && docker build -t localhost:5432/notification:development --target development .
	docker push localhost:5432/notification:development
	cd shipping && docker build -t localhost:5432/shipping:development --target development .
	docker push localhost:5432/shipping:development
	cd billing && docker build -t localhost:5432/billing:development --target development .
	docker push localhost:5432/billing:development

.PHONY: deploy_services
deploy_services:
	cd order && kubectl apply -k k8s/environments/development/
	cd notification && kubectl apply -k k8s/environments/development/
	cd shipping && kubectl apply -k k8s/environments/development/
	cd billing && kubectl apply -k k8s/environments/development/
	kubectl --namespace development rollout status deployment order
	kubectl --namespace development rollout status deployment notification
	kubectl --namespace development rollout status deployment shipping
	kubectl --namespace development rollout status deployment billing

.PHONY: clean_services
clean_services:
	cd order && kubectl delete -k k8s/environments/development/
	cd notification && kubectl delete -k k8s/environments/development/
	cd shipping && kubectl delete -k k8s/environments/development/
	cd billing && kubectl delete -k k8s/environments/development/
