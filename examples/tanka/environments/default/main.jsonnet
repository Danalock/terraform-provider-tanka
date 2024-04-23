local k = import 'github.com/grafana/jsonnet-libs/ksonnet-util/kausal.libsonnet';

local configMap = k.core.v1.configMap;

function(apiServer, namespace, config={}, config_override={}) {

  local default = {
    key_default: 'value only existing in default',
  },
  local override = std.mergePatch(config, config_override),
  local conf = std.mergePatch(default, override),

  trace: std.trace('obj config: %s' % [conf], conf),

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
      "key_default": conf.key_default,
      "key_1": conf.key_1,
      "key_2": conf.key_2,
      "key_list_item": conf.key_list[2],
      "key_nested_item": conf.key_nested.new.deep,
      "key_override": conf.key_override,
    }),
  },
}
