local k = import 'github.com/grafana/jsonnet-libs/ksonnet-util/kausal.libsonnet';

local configMap = k.core.v1.configMap;

function(apiServer, namespace, config_inline = {}, config_local = {}) {

  local config = {
    "key_jsonnet_default": "testvalue"
  } + config_inline + config_local,

  trace: std.trace("obj config: %s" % [config], config),

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
      'key_jsonnet_default': config.key_jsonnet_default,
      'key_inline': config.key_inline,
      'key_local': config.key_local,
      'key_1': config.key_1,
      'key_2': config.key_2,
    }),
  },
}
