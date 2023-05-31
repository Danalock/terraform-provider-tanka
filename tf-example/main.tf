resource "tanka_release" "example" {
  name = "tanka_example"
  api_server = var.api_server
  namespace = var.namespace
  config_inline = {
    key_inline: "inline value"
    key_1: "value_1",
    key_2: "value_2",
    key_list: "Not applicable in the inline config"
    key_nested: "Not applicable in the inline config"
  }
  # config_local = "file://tanka_config_override.json"
  config_local = jsonencode({
    key_1 = "local_override",
    key_local = "local value",
    key_list: [
      "a", "few", "list", "items"
    ],
    key_nested: {
      new: {
        deep: "leaf"
      }
      string_value: "value"
    },
  })
}
