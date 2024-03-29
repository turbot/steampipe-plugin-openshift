package openshift

import (
	"context"

	client_v1 "github.com/openshift/client-go/oauth/clientset/versioned/typed/oauth/v1"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//// TABLE DEFINITION
func tableOpenShiftOAuthAccessToken(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_oauth_access_token",
		Description: "Retrieve information about OpenShift OAuth access tokens.",
		List: &plugin.ListConfig{
			Hydrate: listOAuthAccessTokens,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getOAuthAccessToken,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "client_name",
				Description: "ClientName references the client that created this token.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "expires_in",
				Description: "ExpiresIn is the seconds from CreationTime before this token expires.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "scopes",
				Description: "Scopes is an array of the requested scopes.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "redirect_uri",
				Description: "RedirectURI is the redirection associated with the token.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RedirectURI"),
			},
			{
				Name:        "user_name",
				Description: "The user name associated with this token.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_uid",
				Description: "UserUID is the unique UID associated with this token.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserUID"),
			},
			{
				Name:        "authorize_token",
				Description: "AuthorizeToken contains the token that authorized this token.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "refresh_token",
				Description: "RefreshToken is the value by which this token can be renewed. Can be blank.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "inactivity_timeout_seconds",
				Description: "InactivityTimeoutSeconds is the value in seconds, from the CreationTimestamp, after which this token can no longer be used. The value is automatically incremented when the token is used.",
				Type:        proto.ColumnType_INT,
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
func listOAuthAccessTokens(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_oauth_access_token.listOAuthAccessTokens", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_oauth_access_token.listOAuthAccessTokens", "NewForNetwork_error", err)
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

	for {
		response, err := client.OAuthAccessTokens().List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_oauth_access_token.listOAuthAccessTokens", "api_error", err)
			return nil, err
		}
		for _, token := range response.Items {
			d.StreamListItem(ctx, token)

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
func getOAuthAccessToken(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")

	// Check if name is empty.
	if name == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_oauth_access_token.getOAuthAccessToken", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_oauth_access_token.getOAuthAccessToken", "NewForNetwork_error", err)
		return nil, err
	}

	clusterNetwork, err := client.OAuthAccessTokens().Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_oauth_access_token.getOAuthAccessToken", "api_error", err)
		return nil, err
	}

	return clusterNetwork, nil
}
