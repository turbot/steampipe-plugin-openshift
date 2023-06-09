package openshift

import (
	"context"

	client_v1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func tableOpenShiftProject(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openshift_project",
		Description: "Retrieve information about OpenShift Projects.",
		List: &plugin.ListConfig{
			Hydrate: listProjects,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    getProject,
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "phase",
				Type:        proto.ColumnType_STRING,
				Description: "Phase is the current lifecycle phase of the project.",
				Transform:   transform.FromField("Status.Phase"),
			},
			{
				Name:        "spec_finalizers",
				Type:        proto.ColumnType_JSON,
				Description: "Finalizers is an opaque list of values that must be empty to permanently remove object from storage.",
				Transform:   transform.FromField("Spec.Finalizers"),
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

func listProjects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_project.listProjects", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_project.listProjects", "NewForConfig_error", err)
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
		response, err := client.Projects().List(ctx, input)
		if err != nil {
			plugin.Logger(ctx).Error("openshift_project.listProjects", "api_error", err)
			return nil, err
		}
		for _, project := range response.Items {
			d.StreamListItem(ctx, project)

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

func getProject(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQualString("name")

	// Check if name is empty.
	if name == "" {
		return nil, nil
	}

	config, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_project.getProject", "connection_error", err)
		return nil, err
	}
	client, err := client_v1.NewForConfig(config)
	if err != nil {
		plugin.Logger(ctx).Error("openshift_project.getProject", "NewForConfig_error", err)
		return nil, err
	}

	project, err := client.Projects().Get(ctx, name, v1.GetOptions{})
	if err != nil {
		plugin.Logger(ctx).Error("openshift_project.getProject", "api_error", err)
		return nil, err
	}

	return project, nil
}
