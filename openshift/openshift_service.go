package openshift

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client_v1 "k8s.io/client-go/kubernetes"
)

func tableOpenShiftService(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_service",
		Description: "Retrieve information about your services.",
		List: &plugin.ListConfig{
			Hydrate: listServices,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getService,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "Type determines how the Service is exposed.",
				Transform:   transform.FromField("Spec.Type").Transform(transform.ToString),
			},
			{
				Name:        "allocate_load_balancer_node_ports",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether NodePorts will be automatically allocated for services with type LoadBalancer, or not.",
				Transform:   transform.FromField("Spec.AllocateLoadBalancerNodePorts"),
			},
			{
				Name:        "cluster_ip",
				Type:        proto.ColumnType_STRING,
				Description: "IP address of the service and is usually assigned randomly.",
				Transform:   transform.FromField("Spec.ClusterIP"),
			},
			{
				Name:        "external_name",
				Type:        proto.ColumnType_STRING,
				Description: "The external reference that discovery mechanisms will return as an alias for this service (e.g. a DNS CNAME record).",
				Transform:   transform.FromField("Spec.ExternalName"),
			},
			{
				Name:        "external_traffic_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Denotes whether the service desires to route external traffic to node-local or cluster-wide endpoints.",
				Transform:   transform.FromField("Spec.ExternalTrafficPolicy").Transform(transform.ToString),
			},
			{
				Name:        "health_check_node_port",
				Type:        proto.ColumnType_INT,
				Description: "Specifies the healthcheck nodePort for the service.",
				Transform:   transform.FromField("Spec.HealthCheckNodePort"),
			},
			{
				Name:        "ip_family_policy",
				Type:        proto.ColumnType_STRING,
				Description: "Specifies the dual-stack-ness requested or required by this service, and is gated by the 'IPv6DualStack' feature gate.",
				Transform:   transform.FromField("Spec.IPFamilyPolicy").Transform(transform.ToString),
			},
			{
				Name:        "load_balancer_ip",
				Type:        proto.ColumnType_IPADDR,
				Description: "The IP specified when the load balancer was created.",
				Transform:   transform.FromField("Spec.LoadBalancerIP"),
			},
			{
				Name:        "publish_not_ready_addresses",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates that any agent which deals with endpoints for this service should disregard any indications of ready/not-ready.",
				Transform:   transform.FromField("Spec.PublishNotReadyAddresses"),
			},
			{
				Name:        "session_affinity",
				Type:        proto.ColumnType_STRING,
				Description: "Supports 'ClientIP' and 'None'. Used to maintain session affinity.",
				Transform:   transform.FromField("Spec.SessionAffinity").Transform(transform.ToString),
			},
			{
				Name:        "session_affinity_client_ip_timeout",
				Type:        proto.ColumnType_INT,
				Description: "Specifies the ClientIP type session sticky time in seconds.",
				Transform:   transform.FromField("Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds"),
			},
			{
				Name:        "cluster_ips",
				Type:        proto.ColumnType_JSON,
				Description: "A list of IP addresses assigned to this service, and are usually assigned randomly.",
				Transform:   transform.FromField("Spec.ClusterIPs"),
			},
			{
				Name:        "external_ips",
				Type:        proto.ColumnType_JSON,
				Description: "A list of IP addresses for which nodes in the cluster will also accept traffic for this service.",
				Transform:   transform.FromField("Spec.ExternalIPs"),
			},
			{
				Name:        "ip_families",
				Type:        proto.ColumnType_JSON,
				Description: "A list of IP families (e.g. IPv4, IPv6) assigned to this service, and is gated by the 'IPv6DualStack' feature gate.",
				Transform:   transform.FromField("Spec.IPFamilies"),
			},
			{
				Name:        "load_balancer_ingress",
				Type:        proto.ColumnType_JSON,
				Description: "A list containing ingress points for the load-balancer.",
				Transform:   transform.FromField("Status.LoadBalancer.Ingress"),
			},
			{
				Name:        "load_balancer_source_ranges",
				Type:        proto.ColumnType_JSON,
				Description: "A list of source ranges that will restrict traffic through the cloud-provider load-balancer will be restricted to the specified client IPs.",
				Transform:   transform.FromField("Spec.LoadBalancerSourceRanges"),
			},
			{
				Name:        "ports",
				Type:        proto.ColumnType_JSON,
				Description: "A list of ports that are exposed by this service.",
				Transform:   transform.FromField("Spec.Ports"),
			},
			{
				Name:        "selector_query",
				Type:        proto.ColumnType_STRING,
				Description: "A query string representation of the selector.",
				Transform:   transform.FromField("Spec.Selector").Transform(selectorMapToString),
			},
			{
				Name:        "selector",
				Type:        proto.ColumnType_JSON,
				Description: "Route service traffic to pods with label keys and values matching this selector.",
				Transform:   transform.FromField("Spec.Selector"),
			},
			{
				Name:        "topology_keys",
				Type:        proto.ColumnType_JSON,
				Description: "A preference-order list of topology keys which implementations of services should use to preferentially sort endpoints when accessing this Service, it can not be used at the same time as externalTrafficPolicy=Local.",
				Transform:   transform.FromField("Spec.TopologyKeys"),
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

func listServices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_service.listServices", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_service.listServices", "NewForConfig_error", err)
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
		response, err := client.CoreV1().Services("").List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_service.listServices", "api_error", err)
			return nil, err
		}
		for _, service := range response.Items {
			d.StreamListItem(ctx, service)

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

func getService(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")
	namespace := d.EqualsQualString("namespace")

	// Check if name or namespace is empty.
	if name == "" || namespace == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_service.getService", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_service.getService", "NewForConfig_error", err)
		return nil, err
	}

	service, err := client.CoreV1().Services(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_service.getService", "api_error", err)
		return nil, err
	}

	return service, nil
}
