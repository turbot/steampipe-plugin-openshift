package openshift

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-openshift",
		DefaultTransform: transform.FromCamel(),
		DefaultIgnoreConfig: &plugin.IgnoreConfig{
			ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"404"}),
		},
		DefaultRetryConfig: &plugin.RetryConfig{
			ShouldRetryErrorFunc: shouldRetryError([]string{"429"})},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		TableMap: map[string]*plugin.Table{
			"openshift_build_config":       tableOpenShiftBuildConfig(ctx),
			"openshift_build":              tableOpenShiftBuild(ctx),
			"openshift_deployment_config":  tableOpenShiftDeploymentConfig(ctx),
			"openshift_image_stream":       tableOpenShiftImageStream(ctx),
			"openshift_oauth_access_token": tableOpenShiftOAuthAccessToken(ctx),
			"openshift_project":            tableOpenShiftProject(ctx),
			"openshift_route":              tableOpenShiftRoute(ctx),
			"openshift_user":               tableOpenShiftUser(ctx),
		},
	}
	return p
}
