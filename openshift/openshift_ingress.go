package openshift

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client_v1 "k8s.io/client-go/kubernetes"
)

func tableOpenShiftIngress(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_ingress",
		Description: "Retrieve information about your ingresses.",
		List: &plugin.ListConfig{
			Hydrate: listIngresses,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getIngress,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "ingress_class_name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the IngressClass cluster resource. The associated IngressClass defines which controller will implement the resource.",
				Transform:   transform.FromField("Spec.IngressClassName"),
			},
			{
				Name:        "default_backend",
				Type:        proto.ColumnType_JSON,
				Description: "A default backend capable of servicing requests that don't match any rule. At least one of 'backend' or 'rules' must be specified.",
				Transform:   transform.FromField("Spec.DefaultBackend"),
			},
			{
				Name:        "tls",
				Type:        proto.ColumnType_JSON,
				Description: "TLS configuration.",
				Transform:   transform.FromField("Spec.TLS"),
			},
			{
				Name:        "rules",
				Type:        proto.ColumnType_JSON,
				Description: "A list of host rules used to configure the Ingress.",
				Transform:   transform.FromField("Spec.Rules"),
			},
			{
				Name:        "load_balancer",
				Type:        proto.ColumnType_JSON,
				Description: "a list containing ingress points for the load-balancer. Traffic intended for the service should be sent to these ingress points.",
				Transform:   transform.FromField("Status.LoadBalancer.Ingress"),
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

func listIngresses(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_ingress.listIngresses", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_ingress.listIngresses", "NewForConfig_error", err)
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
		response, err := client.NetworkingV1().Ingresses("").List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_ingress.listIngresses", "api_error", err)
			return nil, err
		}
		for _, ingress := range response.Items {
			d.StreamListItem(ctx, ingress)

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

func getIngress(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")
	namespace := d.EqualsQualString("namespace")

	// Check if name or namespace is empty.
	if name == "" || namespace == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_ingress.getIngress", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_ingress.getIngress", "NewForConfig_error", err)
		return nil, err
	}

	ingress, err := client.NetworkingV1().Ingresses(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_ingress.getIngress", "api_error", err)
		return nil, err
	}

	return ingress, nil
}
