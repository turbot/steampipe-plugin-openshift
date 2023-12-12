package openshift

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type openshiftConfig struct {
	ConfigPath    *string `hcl:"config_path"`
	ConfigContext *string `hcl:"config_context"`
}

func ConfigInstance() interface{} {
	return &openshiftConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) openshiftConfig {
	if connection == nil || connection.Config == nil {
		return openshiftConfig{}
	}
	config, _ := connection.Config.(openshiftConfig)
	return config
}
