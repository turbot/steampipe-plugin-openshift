package openshift

import (
	"context"
	"strings"

	client_v1 "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func tableOpenShiftImageStream(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_image_stream",
		Description: "Retrieve information about your image streams.",
		List: &plugin.ListConfig{
			Hydrate:    listImageStreams,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getImageStream,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "lookup_policy",
				Type:        proto.ColumnType_JSON,
				Description: "Lookup policy controls how other resources reference images within this namespace.",
				Transform:   transform.FromField("Spec.LookupPolicy"),
			},
			{
				Name:        "spec_tags",
				Type:        proto.ColumnType_JSON,
				Description: "Tags map arbitrary string values to specific image locators.",
				Transform:   transform.FromField("Spec.Tags"),
			},
			{
				Name:        "docker_image_repository",
				Type:        proto.ColumnType_STRING,
				Description: "It represents the effective location this stream may be accessed at. May be empty until the server determines where the repository is located.",
				Transform:   transform.FromField("Status.DockerImageRepository"),
			},
			{
				Name:        "public_docker_image_repository",
				Type:        proto.ColumnType_STRING,
				Description: "It represents the public location from where the image can be pulled outside the cluster. This field may be empty if the administrator has not exposed the integrated registry externally.",
				Transform:   transform.FromField("Status.PublicDockerImageRepository"),
			},
			{
				Name:        "status_tags",
				Type:        proto.ColumnType_JSON,
				Description: "Tags are a historical record of images associated with each tag. The first entry in the TagEvent array is the currently tagged image.",
				Transform:   transform.FromField("Status.Tags"),
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

func listImageStreams(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_image_stream.listImageStreams", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_image_stream.listImageStreams", "NewForStream_error", err)
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

	commonFieldSelectorValue := getCommonOptionalKeyQualsValueForFieldSelector(d)

	if len(commonFieldSelectorValue) > 0 {
		input.FieldSelector = strings.Join(commonFieldSelectorValue, ",")
	}

	for {
		response, err := client.ImageStreams("").List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_image_stream.listImageStreams", "api_error", err)
			return nil, err
		}
		for _, imageStream := range response.Items {
			d.StreamListItem(ctx, imageStream)

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

func getImageStream(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")
	namespace := d.EqualsQualString("namespace")

	// Check if name or namespace is empty.
	if name == "" || namespace == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_image_stream.getImageStream", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_image_stream.getImageStream", "NewForStream_error", err)
		return nil, err
	}

	imageStream, err := client.ImageStreams(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_image_stream.getImageStream", "api_error", err)
		return nil, err
	}

	return imageStream, nil
}
