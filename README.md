# Tanka Provider for Terraform

This is an unofficial [Tanka](https://github.com/grafana/tanka) provider for [Terraform](https://www.terraform.io/).

This provider allows you to install and manage Tanka ressources in a Kubernetes cluster using Terraform.

While configuration objects can be passed to the tanka package, secrets should be handled appropriately by some other means, as to not be stored in state.

## Getting Started

You'll find a small usage example in the tf-example directory.

You'll need to have [terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli) and [kubectl](https://kubernetes.io/docs/tasks/tools/) installed.
