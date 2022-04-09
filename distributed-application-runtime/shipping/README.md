## Shipping service

This is a service that simulates shipping the customer order.

## Development

Use development namespace in the k3d cluster and run:

```sh
docker build -t localhost:5432/shipping:development --target development .
docker push localhost:5432/shipping:development
kubectl apply -k k8s/environments/development/
kubectl --namespace development get deployments,pods
```

## Build

To build a container image:

```sh
docker build -t localhost:5432/shipping:1.0.0 --target production .
docker push localhost:5432/shipping:1.0.0
```

## Deploy

```sh
kubectl --namespace <environment> apply -f k8s/
```
