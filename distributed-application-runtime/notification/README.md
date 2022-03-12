## Notification service

This is a service that simulates notifying clients about their order.

## Development

Use development namespace in the k3d cluster and run:

```sh
docker build -t localhost:5432/notification:development --target development .
docker push localhost:5432/notification:development
kubectl apply -k k8s/environments/development/
kubectl --namespace development get deployments,pods
```

## Build

To build a container image:

```sh
docker build -t localhost:5432/notification:1.0.0 --target production .
docker push localhost:5432/notification:1.0.0
```

## Deploy

```sh
kubectl --namespace <environment> apply -f k8s/
```
