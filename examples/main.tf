// Runnable example used in development

variable "endpoint" {
  type = string
}

variable "cluster_ca_certificate" {
  type = string
}

variable "token" {
  type = string
}

terraform {
  required_providers {
    tanka = {
      source = "Danalock/tanka"
    }
  }
}

provider "tanka" {
  endpoint               = var.endpoint
  cluster_ca_certificate = var.cluster_ca_certificate
  token                  = var.token
}

resource "tanka_release" "example" {
  version = "456"
  config = jsonencode({
    key_1 : "value_1"
    key_2 : "value_2"
    key_list : [
      "a", "few", "list", "items"
    ],
    key_nested : {
      new : {
        deep : "leaf"
      }
      string_value : "value"
    },
  })
  # config_override = "file://tanka_config_override.json"
  config_override = jsonencode({
    key_1 : "overridden_value_1"
    key_override : "value_only_existing_in_override"
  })
}
