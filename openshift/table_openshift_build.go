package openshift

import (
	"context"
	"strings"

	client_v1 "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func tableOpenShiftBuild(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_build",
		Description: "Retrieve information about OpenShift builds.",
		List: &plugin.ListConfig{
			Hydrate:    listBuilds,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getBuild,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "common_spec",
				Type:        proto.ColumnType_JSON,
				Description: "CommonSpec is the information that represents a build.",
				Transform:   transform.FromField("Spec.CommonSpec"),
			},
			{
				Name:        "triggered_by",
				Type:        proto.ColumnType_JSON,
				Description: "It describes which triggers started the most recent update to the build configuration and contains information about those triggers.",
				Transform:   transform.FromField("Spec.TriggeredBy"),
			},
			{
				Name:        "phase",
				Type:        proto.ColumnType_STRING,
				Description: "phase is the point in the build lifecycle. Possible values are New, Pending, Running, Complete, Failed, Error, and Cancelled.",
				Transform:   transform.FromField("Status.Phase"),
			},
			{
				Name:        "cancelled",
				Type:        proto.ColumnType_BOOL,
				Description: "Cancelled describes if a cancel event was triggered for the build.",
				Transform:   transform.FromField("Status.Cancelled"),
			},
			{
				Name:        "reason",
				Type:        proto.ColumnType_STRING,
				Description: "Reason is a brief CamelCase string that describes any failure and is meant for machine parsing and tidy display in the CLI.",
				Transform:   transform.FromField("Status.Reason"),
			},
			{
				Name:        "message",
				Type:        proto.ColumnType_STRING,
				Description: "Message is a human-readable message indicating details about why the build has this status.",
				Transform:   transform.FromField("Status.Message"),
			},
			{
				Name:        "start_timestamp",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "It is a timestamp representing the server time when this Build started running in a Pod.",
				Transform:   transform.FromField("Status.StartTimestamp").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "completion_timestamp",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "It is a timestamp representing the server time when this Build was finished, whether that build failed or succeeded. It reflects the time at which the Pod running the Build terminated.",
				Transform:   transform.FromField("Status.CompletionTimestamp").Transform(v1TimeToRFC3339),
			},
			{
				Name:        "duration",
				Type:        proto.ColumnType_INT,
				Description: "Duration contains time.Duration object describing build time.",
				Transform:   transform.FromField("Status.Duration"),
			},
			{
				Name:        "output_docker_image_reference",
				Type:        proto.ColumnType_STRING,
				Description: "It contains a reference to the Docker image that will be built by this build. Its value is computed from Build.Spec.Output.To, and should include the registry address, so that it can be used to push and pull the image.",
				Transform:   transform.FromField("Status.OutputDockerImageReference"),
			},
			{
				Name:        "config",
				Type:        proto.ColumnType_JSON,
				Description: "It is an ObjectReference to the BuildConfig this Build is based on.",
				Transform:   transform.FromField("Status.Config"),
			},
			{
				Name:        "output",
				Type:        proto.ColumnType_JSON,
				Description: "Output describes the Docker image the build has produced.",
				Transform:   transform.FromField("Status.Output"),
			},
			{
				Name:        "stages",
				Type:        proto.ColumnType_JSON,
				Description: "Stages contains details about each stage that occurs during the build including start time, duration (in milliseconds), and the steps that occured within each stage.",
				Transform:   transform.FromField("Status.Stages"),
			},
			{
				Name:        "log_snippet",
				Type:        proto.ColumnType_STRING,
				Description: "It is the last few lines of the build log. This value is only set for builds that failed.",
				Transform:   transform.FromField("Status.LogSnippet"),
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

func listBuilds(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_build.listBuilds", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_build.listBuilds", "NewForConfig_error", err)
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
		response, err := client.Builds("").List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_build.listBuilds", "api_error", err)
			return nil, err
		}
		for _, build := range response.Items {
			d.StreamListItem(ctx, build)

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

func getBuild(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")
	namespace := d.EqualsQualString("namespace")

	// Check if name or namespace is empty.
	if name == "" || namespace == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_build.getBuild", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_build.getBuild", "NewForConfig_error", err)
		return nil, err
	}

	build, err := client.Builds(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_build.getBuild", "api_error", err)
		return nil, err
	}

	return build, nil
}
