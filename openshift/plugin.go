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
			"openshift_project":              tableOpenShiftProject(ctx),
			"openshift_user":                 tableOpenShiftUser(ctx),
			"openshift_deployment_config":    tableOpenShiftDeploymentConfig(ctx),
			"openshift_cluster_role":         tableOpenShiftClusterRole(ctx),
			"openshift_cluster_role_binding": tableOpenShiftClusterRoleBinding(ctx),
			"openshift_role_binding":         tableOpenShiftRoleBinding(ctx),
			"openshift_role":                 tableOpenShiftRole(ctx),
			//"openshift_build":                tableOpenShiftBuild(ctx),
			"openshift_secret":     tableOpenShiftSecret(ctx),
			"openshift_pod":        tableOpenShiftPod(ctx),
			"openshift_service":    tableOpenShiftService(ctx),
			"openshift_deployment": tableOpenShiftDeployment(ctx),
			"openshift_node":       tableOpenShiftNode(ctx),
			"openshift_ingress":    tableOpenShiftIngress(ctx),
		},
	}
	return p
}
