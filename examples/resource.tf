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
    key_1        = "overridden_value_1",
    key_override = "value_only_existing_in_override",
  })
}
