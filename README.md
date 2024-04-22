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

## Development

The provider is written in [Go](https://golang.org/doc/install).

Activate the developer override by exporting the dev environment vars (run `source ./set_debug_env_vars.sh`). If `terrafom init` is necessary, it should be run before activation of the dev override.

Build the provider with `go build -o terraform-provider-tanka`

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.
