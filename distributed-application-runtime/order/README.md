## Order service

This is a service that simulates new client orders.

## Development

Use development namespace in the k3d cluster and run:

```sh
docker build -t localhost:5432/order:development --target development .
docker push localhost:5432/order:development
kubectl apply -k k8s/environments/development/
kubectl --namespace development get services,deployments,pods
```

## Build

To build a container image:

```sh
docker build -t localhost:5432/order:1.0.0 --target production .
docker push localhost:5432/order:1.0.0
```

## Deploy

```sh
kubectl --namespace <environment> apply -f k8s/
```

Note: for the sake of demonstration it's using an in memory database.
