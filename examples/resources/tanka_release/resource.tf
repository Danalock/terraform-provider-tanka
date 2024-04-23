resource "tanka_release" "minimal" {}

resource "tanka_release" "config_override_inline " {
  config = jsonencode({
    key_1 : "value_1"
    key_2 : "value_2"
  })
  config_override = jsonencode({
    key_1 : "overridden_value_1",
  })
}

resource "tanka_release" "config_override_external_file " {
  config = jsonencode({
    key_1 : "value_1"
    key_2 : "value_2"
  })
  config_override = "file://tanka_config_override.json"
}
