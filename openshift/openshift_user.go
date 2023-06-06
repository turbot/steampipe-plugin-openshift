package openshift

import (
	"context"

	client_v1 "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func tableOpenShiftUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_user",
		Description: "Retrieve information about your users.",
		List: &plugin.ListConfig{
			Hydrate: listUsers,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getUser,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "full_name",
				Type:        proto.ColumnType_STRING,
				Description: "The full name of the user.",
			},
			{
				Name:        "identities",
				Type:        proto.ColumnType_JSON,
				Description: "Identities are the identities associated with this user.",
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

func listUsers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_user.listUsers", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_user.listUsers", "NewForConfig_error", err)
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
		response, err := client.Users().List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_user.listUsers", "api_error", err)
			return nil, err
		}
		for _, user := range response.Items {
			d.StreamListItem(ctx, user)

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

func getUser(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")

	// Check if name is empty.
	if name == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_user.getUser", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_user.getUser", "NewForConfig_error", err)
		return nil, err
	}

	user, err := client.Users().Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_user.getUser", "api_error", err)
		return nil, err
	}

	return user, nil
}
