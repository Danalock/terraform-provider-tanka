# Terraform Provider for Tanka

This is an unofficial [Terraform](https://www.terraform.io/) provider for [Tanka](https://github.com/grafana/tanka).

This provider allows you to install and manage Tanka ressources in a Kubernetes cluster using Terraform.

While configuration objects can be passed to the tanka package, secrets should be handled appropriately by some other means, as to not be stored in state.

## Getting Started

You'll find a small usage example in the tf-example directory.

You'll need to have [terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli) and [kubectl](https://kubernetes.io/docs/tasks/tools/) installed.

## Development

The `tf-example` dir contains a small terraform example that creates a configmap demonstrating the override hierarchy of the different config vars. Activate the developer override by exporting the dev environment vars (run `. set_debug_env_vars.sh`). If `terrafom init` is necessary, it should be run before activation the dev override.
