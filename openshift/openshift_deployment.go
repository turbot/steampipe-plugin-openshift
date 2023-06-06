package openshift

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client_v1 "k8s.io/client-go/kubernetes"
)

func tableOpenShiftDeployment(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_deployment",
		Description: "Retrieve information about your deployments.",
		List: &plugin.ListConfig{
			Hydrate: listDeployments,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getDeployment,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "replicas",
				Type:        proto.ColumnType_INT,
				Description: "Number of desired pods. Defaults to 1.",
				Transform:   transform.FromField("Spec.Replicas"),
			},
			{
				Name:        "selector_query",
				Type:        proto.ColumnType_STRING,
				Description: "A query string representation of the selector.",
				Transform:   transform.FromField("Spec.Selector").Transform(labelSelectorToString),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: " Label selector for pods. A label selector is a label query over a set of resources.",
				Transform:   transform.FromField("Spec.Selector"),
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Template describes the pods that will be created.",
				Transform:   transform.FromField("Spec.Template"),
			},
			{
				Name:        "strategy",
				Type:        proto.ColumnType_JSON,
				Description: "The deployment strategy to use to replace existing pods with new ones.",
				Transform:   transform.FromField("Spec.Strategy"),
			},
			{
				Name:        "min_ready_seconds",
				Type:        proto.ColumnType_INT,
				Description: "Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0.",
				Transform:   transform.FromField("Spec.MinReadySeconds"),
			},
			{
				Name:        "revision_history_limit",
				Type:        proto.ColumnType_INT,
				Description: "The number of old ReplicaSets to retain to allow rollback.",
				Transform:   transform.FromField("Spec.RevisionHistoryLimit"),
			},
			{
				Name:        "paused",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates that the deployment is paused.",
				Transform:   transform.FromField("Spec.Paused"),
			},
			{
				Name:        "progress_deadline_seconds",
				Type:        proto.ColumnType_INT,
				Description: "The maximum time in seconds for a deployment to make progress before it is considered to be failed.",
				Transform:   transform.FromField("Spec.ProgressDeadlineSeconds"),
			},

			//// Status Columns
			{
				Name:        "observed_generation",
				Type:        proto.ColumnType_INT,
				Description: "The generation observed by the deployment controller.",
				Transform:   transform.FromField("Status.ObservedGeneration"),
			},
			{
				Name:        "status_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of non-terminated pods targeted by this deployment (their labels match the selector).",
				Transform:   transform.FromField("Status.Replicas"),
			},
			{
				Name:        "updated_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of non-terminated pods targeted by this deployment that have the desired template spec.",
				Transform:   transform.FromField("Status.UpdatedReplicas"),
			},
			{
				Name:        "ready_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of ready pods targeted by this deployment.",
				Transform:   transform.FromField("Status.ReadyReplicas"),
			},
			{
				Name:        "available_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.",
				Transform:   transform.FromField("Status.AvailableReplicas"),
			},
			{
				Name:        "unavailable_replicas",
				Type:        proto.ColumnType_INT,
				Description: "Total number of unavailable pods targeted by this deployment.",
				Transform:   transform.FromField("Status.UnavailableReplicas"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "Represents the latest available observations of a deployment's current state.",
				Transform:   transform.FromField("Status.Conditions"),
			},
			{
				Name:        "collision_count",
				Type:        proto.ColumnType_INT,
				Description: "Count of hash collisions for the Deployment. The Deployment controller uses this field as a collision avoidance mechanism when it needs to create the name for the newest ReplicaSet.",
				Transform:   transform.FromField("Status.CollisionCount"),
			},

			/// Steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

func listDeployments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_deployment.listDeployments", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_deployment.listDeployments", "NewForConfig_error", err)
		return nil, err
	}

	// Limiting the results
	maxLimit := int64(500)
	if d.QueryContext.Limit != nil {
		limit := *d.QueryContext.Limit
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := v1.ListOptions{
		Limit: maxLimit,
	}

	for {
		response, err := client.AppsV1().Deployments("").List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_deployment.listDeployments", "api_error", err)
			return nil, err
		}
		for _, deployment := range response.Items {
			d.StreamListItem(ctx, deployment)

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

func getDeployment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")
	namespace := d.EqualsQualString("namespace")

	// Check if name or namespace is empty.
	if name == "" || namespace == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_deployment.getDeployment", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_deployment.getDeployment", "NewForConfig_error", err)
		return nil, err
	}

	deployment, err := client.AppsV1().Deployments(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_deployment.getDeployment", "api_error", err)
		return nil, err
	}

	return deployment, nil
}
