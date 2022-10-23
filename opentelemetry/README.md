## OpenTelemetry Auto Instrumentation 

This is a simple example of how to use OpenTelemetry to automatically instrument applications in Kubernetes.

### Architecture

There are three components to this example:

1. A frontend that makes requests to a backend order service
2. A backend order service that makes requests to an email service (bottleneck)
3. A backend email service that simulating sending emails

In this example we'll create 3 nodes cluster with 1 frontend pod, 1 order service pod, and 1 email service pod. The frontend will make requests to the order service, which will make requests to the email service. The order service will be the bottleneck in this example.

### Prerequisites

1. Docker
2. K3D or KinD (for installing a local kubernetes cluster)
3. Kubernetes client - kubectl
4. OpenTelemetry Collector

### Setup

1. Install a local kubernetes cluster

```bash
# Install K3D
curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash

# Create a cluster
make cluster
```

2. Install OpenTelemetry Operator

```bash
make deploy-opentelemetry-operator
```

3. Install the services

```bash
make build
make import-container-images
make deploy-services
```

4. Setup automatic instrumentation

```bash
make deploy-auto-instrumentation
```

5. Finally cleanup:

```bash
make clean
```
