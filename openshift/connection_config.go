package openshift

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/schema"
)

type openshiftConfig struct {
	ConfigPath    *string `cty:"config_path"`
	ConfigContext *string `cty:"config_context"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"config_path": {
		Type: schema.TypeString,
	},
	"config_context": {
		Type: schema.TypeString,
	},
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
