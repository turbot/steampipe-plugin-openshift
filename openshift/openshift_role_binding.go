package openshift

import (
	"context"

	client_v1 "github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func tableOpenShiftRoleBinding(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_role_binding",
		Description: "Retrieve information about your role bindings.",
		List: &plugin.ListConfig{
			Hydrate: listRoleBindings,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getRoleBinding,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "user_names",
				Type:        proto.ColumnType_JSON,
				Description: "Holds all the usernames directly bound to the role.",
			},
			{
				Name:        "group_names",
				Type:        proto.ColumnType_JSON,
				Description: "Holds all the groups directly bound to the role.",
			},
			{
				Name:        "subjects",
				Type:        proto.ColumnType_JSON,
				Description: "Subjects hold object references to authorize with this rule. This field is ignored if UserNames or GroupNames are specified to support legacy clients and servers. Thus newer clients that do not need to support backwards compatibility should send only fully qualified Subjects and should omit the UserNames and GroupNames fields. Clients that need to support backwards compatibility can use this field to build the UserNames and GroupNames.",
			},
			{
				Name:        "role_name",
				Type:        proto.ColumnType_STRING,
				Description: "Name of the referent.",
				Transform:   transform.FromField("RoleRef.Name"),
			},
			{
				Name:        "role_namespace",
				Type:        proto.ColumnType_STRING,
				Description: "Namespace of the referent.",
				Transform:   transform.FromField("RoleRef.Namespace"),
			},
			{
				Name:        "role_kind",
				Type:        proto.ColumnType_STRING,
				Description: "Kind of the referent.",
				Transform:   transform.FromField("RoleRef.Kind"),
			},
			{
				Name:        "role_uid",
				Type:        proto.ColumnType_STRING,
				Description: "UID of the referent.",
				Transform:   transform.FromField("RoleRef.UID"),
			},
			{
				Name:        "role_api_version",
				Type:        proto.ColumnType_STRING,
				Description: "API version of the referent.",
				Transform:   transform.FromField("RoleRef.APIVersion"),
			},
			{
				Name:        "role_resource_version",
				Type:        proto.ColumnType_STRING,
				Description: "Specific resourceVersion to which this reference is made, if any.",
				Transform:   transform.FromField("RoleRef.ResourceVersion"),
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

func listRoleBindings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_role_binding.listRoleBindings", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_role_binding.listRoleBindings", "NewForConfig_error", err)
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
		response, err := client.RoleBindings("").List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift__role_binding.listRoleBindings", "api_error", err)
			return nil, err
		}
		for _, RoleBinding := range response.Items {
			d.StreamListItem(ctx, RoleBinding)

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

func getRoleBinding(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")
	namespace := d.EqualsQualString("namespace")

	// Check if name or namespace is empty.
	if name == "" || namespace == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_role_binding.getRoleBinding", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_role_binding.getRoleBinding", "NewForConfig_error", err)
		return nil, err
	}

	RoleBinding, err := client.RoleBindings(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_role_binding.getRoleBinding", "api_error", err)
		return nil, err
	}

	return RoleBinding, nil
}
