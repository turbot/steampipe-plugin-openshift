package openshift

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client_v1 "k8s.io/client-go/kubernetes"
)

func tableOpenShiftNode(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_node",
		Description: "Retrieve information about your nodes.",
		List: &plugin.ListConfig{
			Hydrate: listNodes,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getNode,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "pod_cidr",
				Type:        proto.ColumnType_CIDR,
				Description: "Pod IP range assigned to the node.",
				Transform:   transform.FromField("Spec.PodCIDR"),
			},
			{
				Name:        "pod_cidrs",
				Type:        proto.ColumnType_JSON,
				Description: "List of the IP ranges assigned to the node for usage by Pods.",
				Transform:   transform.FromField("Spec.PodCIDRs"),
			},
			{
				Name:        "provider_id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the node assigned by the cloud provider in the format: <ProviderName>://<ProviderSpecificNodeID>.",
				Transform:   transform.FromField("Spec.ProviderID"),
			},
			{
				Name:        "unschedulable",
				Type:        proto.ColumnType_BOOL,
				Description: "Unschedulable controls node schedulability of new pods. By default, node is schedulable.",
				Transform:   transform.FromField("Spec.Unschedulable"),
			},
			{
				Name:        "taints",
				Type:        proto.ColumnType_JSON,
				Description: "List of the taints attached to the node to has the \"effect\" on pod that does not tolerate the Taint",
				Transform:   transform.FromField("Spec.Taints"),
			},
			{
				Name:        "config_source",
				Type:        proto.ColumnType_JSON,
				Description: "The source to get node configuration from.",
				Transform:   transform.FromField("Spec.ConfigSource"),
			},
			{
				Name:        "capacity",
				Type:        proto.ColumnType_JSON,
				Description: "Capacity represents the total resources of a node.",
				Transform:   transform.FromField("Status.Capacity"),
			},
			{
				Name:        "allocatable",
				Type:        proto.ColumnType_JSON,
				Description: "Allocatable represents the resources of a node that are available for scheduling. Defaults to capacity.",
				Transform:   transform.FromField("Status.Allocatable"),
			},
			{
				Name:        "phase",
				Type:        proto.ColumnType_STRING,
				Description: "The recently observed lifecycle phase of the node.",
				Transform:   transform.FromField("Status.Phase"),
			},
			{
				Name:        "conditions",
				Type:        proto.ColumnType_JSON,
				Description: "List of current observed node conditions.",
				Transform:   transform.FromField("Status.Conditions"),
			},
			{
				Name:        "addresses",
				Type:        proto.ColumnType_JSON,
				Description: "Endpoints of daemons running on the Node.",
				Transform:   transform.FromField("Status.Addresses"),
			},
			{
				Name:        "daemon_endpoints",
				Type:        proto.ColumnType_JSON,
				Description: "Set of ids/uuids to uniquely identify the node.",
				Transform:   transform.FromField("Status.DaemonEndpoints"),
			},
			{
				Name:        "node_info",
				Type:        proto.ColumnType_JSON,
				Description: "List of container images on this node.",
				Transform:   transform.FromField("Status.NodeInfo"),
			},
			{
				Name:        "images",
				Type:        proto.ColumnType_JSON,
				Description: "List of container images on this node.",
				Transform:   transform.FromField("Status.Images"),
			},
			{
				Name:        "volumes_in_use",
				Type:        proto.ColumnType_JSON,
				Description: "List of attachable volumes in use (mounted) by the node.",
				Transform:   transform.FromField("Status.VolumesInUse"),
			},
			{
				Name:        "volumes_attached",
				Type:        proto.ColumnType_JSON,
				Description: "List of volumes that are attached to the node.",
				Transform:   transform.FromField("Status.VolumesAttached"),
			},
			{
				Name:        "config",
				Type:        proto.ColumnType_JSON,
				Description: "Status of the config assigned to the node via the dynamic Kubelet config feature.",
				Transform:   transform.FromField("Status.Config"),
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

func listNodes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_node.listNodes", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_node.listNodes", "NewForConfig_error", err)
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
		response, err := client.CoreV1().Nodes().List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_node.listNodes", "api_error", err)
			return nil, err
		}
		for _, node := range response.Items {
			d.StreamListItem(ctx, node)

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

func getNode(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")

	// Check if name is empty.
	if name == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_node.getNode", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_node.getNode", "NewForConfig_error", err)
		return nil, err
	}

	node, err := client.CoreV1().Nodes().Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_node.getNode", "api_error", err)
		return nil, err
	}

	return node, nil
}
