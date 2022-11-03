# Gloo Portal Demo

This portal demo is based on Book Store sample application.

## Prerequisites

1. Create required env vars

    ```
    export CLUSTER_OWNER="kasunt"
    export CLUSTER_REGION=australia-southeast1
    export PROJECT="gloo-portal-demo"

    export DOMAIN_NAME=bookstore.development.internal

    export GLOO_EDGE_HELM_VERSION=1.12.33
    export GLOO_EDGE_VERSION=v${GLOO_EDGE_HELM_VERSION}

    export GLOO_PORTAL_HELM_VERSION=1.3.0-beta9
    export GLOO_PORTAL_VERSION=v${GLOO_PORTAL_HELM_VERSION}
    ```

2. Provisioned cluster

    ```
    ./cluster-provision/scripts/provision-gke-cluster.sh create -n $PROJECT -o $CLUSTER_OWNER -r $CLUSTER_REGION
    ```

## Instructions

### Prerequisite

Deploy Gloo Edge (with license)

```
helm repo add gloo-ee https://storage.googleapis.com/gloo-ee-helm
helm repo update

helm install gloo-ee gloo-ee/gloo-ee -n gloo-system \
  --version ${GLOO_EDGE_HELM_VERSION} \
  --create-namespace \
  --set-string license_key=${GLOO_EDGE_LICENSE_KEY} \
  -f k8s/gloo-edge-helm-values.yaml
```

### Install Dev Portal

```
helm repo add dev-portal https://storage.googleapis.com/dev-portal-helm
helm repo update

helm install gloo-portal gloo-portal/gloo-portal -n gloo-portal \
  --version ${GLOO_PORTAL_HELM_VERSION} \
  --create-namespace \
  -f k8s/gloo-portal-helm-values.yaml
```

### Install Application

```
kubectl create ns apps
kubectl create ns apps-configuration

./k8s/bookstrap-demo-portal-environment.sh prov
```