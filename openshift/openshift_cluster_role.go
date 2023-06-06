package openshift

import (
	"context"

	client_v1 "github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func tableOpenShiftClusterRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_cluster_role",
		Description: "Retrieve information about your cluster roles.",
		List: &plugin.ListConfig{
			Hydrate: listClusterRoles,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getClusterRole,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "rules",
				Type:        proto.ColumnType_JSON,
				Description: "List of the PolicyRules for this Role.",
			},
			{
				Name:        "aggregation_rule",
				Type:        proto.ColumnType_JSON,
				Description: "An optional field that describes how to build the Rules for this ClusterRole",
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

func listClusterRoles(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_cluster_role.listClusterRoles", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_cluster_role.listClusterRoles", "NewForConfig_error", err)
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
		response, err := client.ClusterRoles().List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_cluster_role.listClusterRoles", "api_error", err)
			return nil, err
		}
		for _, clusterRole := range response.Items {
			d.StreamListItem(ctx, clusterRole)

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

func getClusterRole(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")

	// Check if name is empty.
	if name == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_cluster_role.getClusterRole", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_cluster_role.getClusterRole", "NewForConfig_error", err)
		return nil, err
	}

	clusterRole, err := client.ClusterRoles().Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_cluster_role.getClusterRole", "api_error", err)
		return nil, err
	}

	return clusterRole, nil
}
