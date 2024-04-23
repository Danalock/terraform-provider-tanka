# Terraform Provider for Tanka

This is an unofficial [Terraform](https://www.terraform.io/) provider for [Tanka](https://github.com/grafana/tanka).

This provider allows you to install and manage Tanka ressources in a Kubernetes cluster using Terraform.

# Usage

Refer to the docs at [the terraform registry](https://registry.terraform.io/providers/Danalock/tanka/latest/docs).

*Note:* While configuration objects can be passed to the tanka package, secrets should be handled appropriately by some other means, as to not be stored in state.

## Getting Started

The example directory contains a small usage example.

You'll need to have [terraform](https://developer.hashicorp.com/terraform/downloads) and [kubectl](https://kubernetes.io/docs/tasks/tools/) installed.

The example creates a configmap demonstrating the override hierarchy of the different config vars.

### Minikube example

Run the example against Minikube by starting Minikube with a user token.

*Note:* The following creates blanket permissions with an unsecure password - do not do something like this in a publicly accessible cluster (including production environments).

Create `.minikube/files/etc/ca-certificates/token.csv` with the following content: `123,system:serviceaccount:default:default,1` to create a token for the default serviceaccount with the value "123". Start the cluster with `minikube start --extra-config=apiserver.token-auth-file=/etc/ca-certificates/token.csv`.

Add permission to everything for this service account by applying the following kubernetes yaml with `kubectl apply -f rbac.yml`:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole

metadata:
  name: default

rules:
  - apiGroups: ["*"]
    resources: ["*"]
    verbs: ["*"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding

metadata:
  name: default-binding

roleRef:
  kind: ClusterRole
  name: default
  apiGroup: rbac.authorization.k8s.io

subjects:
  - kind: ServiceAccount
    name: default
    namespace: default
```

Create the file `config.auto.tfvars` from the example, filling in the CA certificate and the correct local address for the cluster. The token value will be "123" after the above procedure.

It should now be possible to run `terrafom apply` inside the example directory (all the tf code is in `main.tf`).

## Development

The provider is written in [Go](https://golang.org/doc/install).

Activate the developer override by exporting the dev environment vars (run `source ./set_debug_env_vars.sh`). If `terrafom init` is necessary, it should be run before activation of the dev override.

Build the provider with `go build -o terraform-provider-tanka`

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.
