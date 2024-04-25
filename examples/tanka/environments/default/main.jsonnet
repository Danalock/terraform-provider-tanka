local k = import 'github.com/grafana/jsonnet-libs/ksonnet-util/kausal.libsonnet';

local configMap = k.core.v1.configMap;

function(apiServer, namespace, tf_config={}, tf_config_override={}) {

  local default_config = {
    key_default: 'value only existing in default',
  },
  local tf_config_values = std.mergePatch(tf_config, tf_config_override),
  local config = std.mergePatch(default_config, tf_config_values),

  trace: std.trace('obj config: %s' % [config], config),

  apiVersion: 'tanka.dev/v1alpha1',
  kind: 'Environment',
  metadata: {
    name: 'default',
  },
  spec: {
    namespace: namespace,
    apiServer: apiServer,
    injectLabels: true,
  },

  data: {
    terraformprovidertest: configMap.new('terraformprovidertest', {
      "key_default": config.key_default,
      "key_1": config.key_1,
      "key_2": config.key_2,
      "key_list_item": config.key_list[2],
      "key_nested_item": config.key_nested.new.deep,
      "key_override": config.key_override,
    }),
  },
}
