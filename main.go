package main

import (
	"github.com/turbot/steampipe-plugin-openshift/openshift"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: openshift.Plugin})
}
