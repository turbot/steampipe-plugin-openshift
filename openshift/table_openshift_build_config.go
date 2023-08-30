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

//// TABLE DEFINITION
func tableOpenShiftBuildConfig(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_build_config",
		Description: "Retrieve information about OpenShift build configs.",
		List: &plugin.ListConfig{
			Hydrate:    listBuildConfigs,
			KeyColumns: getCommonOptionalKeyQuals(),
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "namespace"}),
			Hydrate:    getBuildConfig,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "common_spec",
				Description: "CommonSpec is the desired build specification.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.CommonSpec"),
			},
			{
				Name:        "triggers",
				Description: "Triggers determine how new Builds can be launched from a BuildConfig. If no triggers are defined, a new build can only occur as a result of an explicit client build creation.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Spec.Triggers"),
			},
			{
				Name:        "run_policy",
				Description: "RunPolicy describes how the new build created from this build configuration will be scheduled for execution. This is optional, if not specified we default to Serial.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Spec.RunPolicy"),
			},
			{
				Name:        "successful_builds_history_limit",
				Description: "It is the number of old successful builds to retain. If not specified, all successful builds are retained.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Spec.SuccessfulBuildsHistoryLimit"),
			},
			{
				Name:        "failed_builds_history_limit",
				Description: "It is the number of old failed builds to retain. If not specified, all failed builds are retained.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Spec.FailedBuildsHistoryLimit"),
			},
			{
				Name:        "last_version",
				Description: "Last version is used to inform about the number of last triggered builds.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Status.LastVersion"),
			},
			{
				Name:        "image_change_triggers",
				Description: "ImageChangeTriggers captures the runtime state of any ImageChangeTrigger specified in the BuildConfigSpec, including the value reconciled by the OpenShift APIServer for the lastTriggeredImageID.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Status.ImageChangeTriggers"),
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

func listBuildConfigs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_build_config.listBuildConfigs", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_build_config.listBuildConfigs", "NewForConfig_error", err)
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
		response, err := client.BuildConfigs("").List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_build_config.listBuildConfigs", "api_error", err)
			return nil, err
		}
		for _, buildConfig := range response.Items {
			d.StreamListItem(ctx, buildConfig)

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

func getBuildConfig(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")
	namespace := d.EqualsQualString("namespace")

	// Check if name or namespace is empty.
	if name == "" || namespace == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_build_config.getBuildConfig", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_build_config.getBuildConfig", "NewForConfig_error", err)
		return nil, err
	}

	buildConfig, err := client.BuildConfigs(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_build_config.getBuildConfig", "api_error", err)
		return nil, err
	}

	return buildConfig, nil
}
