package openshift

import (
	"context"
	"strings"

	client_v1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//// TABLE DEFINITION
func tableOpenShiftRoute(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_route",
		Description: "Retrieve information about OpenShift routes.",
		List: &plugin.ListConfig{
			Hydrate:    listRoutes,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getRoute,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "host",
				Description: "Host is an alias/DNS that points to the service. Optional. If not specified a route name will typically be automatically chosen. Must follow DNS952 subdomain conventions.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Host"),
			},
			{
				Name:        "path",
				Description: "Path that the router watches for, to route traffic to the service.",
				Transform:   transform.FromField("Spec.Path"),
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "spec_to",
				Description: "To is an object the route should use as the primary backend. Only the Service kind is allowed, and it will be defaulted to Service. If the weight field (0-256 default 1) is set to zero, no traffic will be sent to this backend.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.To"),
			},
			{
				Name:        "alternate_backends",
				Description: "Alternate backends allow up to 3 additional backends to be assigned to the route. Only the Service kind is allowed, and it will be defaulted to Service. Use the weight field in RouteTargetReference object to specify relative preference.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.AlternateBackends"),
			},
			{
				Name:        "port",
				Description: "If specified, the port to be used by the router. Most routers will use all endpoints exposed by the service by default - set this value to instruct routers which port to use.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.Port"),
			},
			{
				Name:        "subdomain",
				Description: "Subdomain is a DNS subdomain that is requested within the ingress controller's domain (as a subdomain). If the host is set this field is ignored. An ingress controller may choose to ignore this suggested name, in which case the controller will report the assigned name in the status.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.Subdomain"),
			},
			{
				Name:        "tls",
				Description: "The tls field provides the ability to configure certificates and termination for the route.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.TLS"),
			},
			{
				Name:        "wildcard_policy",
				Description: "Wildcard policy if any for the route. Currently only 'Subdomain' or 'None' is allowed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.WildcardPolicy"),
			},
			{
				Name:        "ingress",
				Description: "Ingress describes the places where the route may be exposed. The list of ingress points may contain duplicate Host or RouterName values. Routes are considered live once they are `Ready`.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Status.Ingress"),
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
func listRoutes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_route.listRoutes", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_route.listRoutes", "NewForConfig_error", err)
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
		response, err := client.Routes("").List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_route.listRoutes", "api_error", err)
			return nil, err
		}
		for _, route := range response.Items {
			d.StreamListItem(ctx, route)

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
func getRoute(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")
	namespace := d.EqualsQualString("namespace")

	// Check if name or namespace is empty.
	if name == "" || namespace == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_route.getRoute", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_route.getRoute", "NewForConfig_error", err)
		return nil, err
	}

	route, err := client.Routes(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_route.getRoute", "api_error", err)
		return nil, err
	}

	return route, nil
}
