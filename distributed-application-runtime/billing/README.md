## Billing service

This is a service that simulates billing the customer order, it simulates the communication with external API's for charging the amount of the ordered line items.

## Development

Use development namespace in the k3d cluster and run:

```sh
docker build -t localhost:5432/billing:development --target development .
docker push localhost:5432/billing:development
kubectl apply -k k8s/environments/development/
kubectl --namespace development get deployments,pods
```

## Build

To build a container image:

```sh
docker build -t localhost:5432/billing:1.0.0 --target production .
docker push localhost:5432/billing:1.0.0
```

## Deploy

```sh
kubectl --namespace <environment> apply -f k8s/
```
