package openshift

import (
	"context"
	"strings"

	client_v1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//// TABLE DEFINITION
func tableOpenShiftDeploymentConfig(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_deployment_config",
		Description: "Retrieve information about OpenShift deployment configs.",
		List: &plugin.ListConfig{
			Hydrate:    listDeploymentConfigs,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getDeploymentConfig,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "strategy",
				Description: "Strategy describes how a deployment is executed.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.Strategy"),
			},
			{
				Name:        "min_ready_seconds",
				Description: "MinReadySeconds is the minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Spec.MinReadySeconds"),
			},
			{
				Name:        "triggers",
				Description: "Triggers determine how updates to a DeploymentConfig result in new deployments. If no triggers are defined, a new deployment can only occur as a result of an explicit client update to the DeploymentConfig with a new LatestVersion. If null, defaults to having a config change trigger.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.Triggers"),
			},
			{
				Name:        "replicas",
				Description: "Replicas is the number of desired replicas.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Spec.Replicas"),
			},
			{
				Name:        "revision_history_limit",
				Description: "RevisionHistoryLimit is the number of old ReplicationControllers to retain to allow for rollbacks. This field is a pointer to allow for differentiation between an explicit zero and not specified. Defaults to 10.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Spec.RevisionHistoryLimit"),
			},
			{
				Name:        "test",
				Description: "Test ensures that this deployment config will have zero replicas except while a deployment is running.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Spec.Test"),
			},
			{
				Name:        "paused",
				Description: "Paused indicates that the deployment config is paused resulting in no new deployments on template changes or changes in the template caused by other triggers.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Spec.Paused"),
			},
			{
				Name:        "selector",
				Description: "Selector is a label query over pods that should match the Replicas count.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.Selector"),
			},
			{
				Name:        "template",
				Description: "Template is the object that describes the pod that will be created if insufficient replicas are detected.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.Template"),
			},
			{
				Name:        "latest_version",
				Description: "LatestVersion is used to determine whether the current deployment associated with a deployment config is out of sync.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Status.LatestVersion"),
			},
			{
				Name:        "observed_generation",
				Description: "ObservedGeneration is the most recent generation observed by the deployment config controller.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Status.ObservedGeneration"),
			},
			{
				Name:        "replicas",
				Description: "Replicas is the total number of pods targeted by this deployment config.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Status.Replicas"),
			},
			{
				Name:        "updated_replicas",
				Description: "The total number of non-terminated pods targeted by this deployment config that has the desired template spec.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Status.UpdatedReplicas"),
			},
			{
				Name:        "available_replicas",
				Description: "The total number of available pods targeted by this deployment config.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Status.AvailableReplicas"),
			},
			{
				Name:        "unavailable_replicas",
				Description: "The total number of unavailable pods targeted by this deployment config.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Status.UnavailableReplicas"),
			},
			{
				Name:        "Details",
				Description: "Details are the reasons for the update to this deployment config.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Status.Details"),
			},
			{
				Name:        "conditions",
				Description: "Conditions represent the latest available observations of a deployment config's current state.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Status.Conditions"),
			},
			{
				Name:        "ready_replicas",
				Description: "Total number of ready pods targeted by this deployment.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Status.ReadyReplicas"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

// LIST FUNCTION
func listDeploymentConfigs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_deployment_config.listDeploymentConfigs", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_deployment_config.listDeploymentConfigs", "NewForConfig_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int64(1000)
	if d.QueryContext.Limit != nil {
		limit := *d.QueryContext.Limit
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := v1.ListOptions{
		Limit: maxLimit,
	}

	commonFieldSelectorValue := getCommonOptionalKeyQualsValueForFieldSelector(d)

	if len(commonFieldSelectorValue) > 0 {
		input.FieldSelector = strings.Join(commonFieldSelectorValue, ",")
	}

	for {
		response, err := client.DeploymentConfigs("").List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_deployment_config.listDeploymentConfigs", "api_error", err)
			return nil, err
		}
		for _, deploymentConfig := range response.Items {
			d.StreamListItem(ctx, deploymentConfig)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if response.Continue != "" {
			input.Continue = response.Continue
		} else {
			break
		}
	}

	return nil, nil
}

// HYDRATE FUNCTIONS
func getDeploymentConfig(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")
	namespace := d.EqualsQualString("namespace")

	// Check if name or namespace is empty.
	if name == "" || namespace == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_deployment_config.getDeploymentConfig", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_deployment_config.getDeploymentConfig", "NewForConfig_error", err)
		return nil, err
	}

	deploymentConfig, err := client.DeploymentConfigs(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_deployment_config.getDeploymentConfig", "api_error", err)
		return nil, err
	}

	return deploymentConfig, nil
}
