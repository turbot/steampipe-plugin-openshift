package openshift

import (
	"context"
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func getClient(ctx context.Context, d *plugin.QueryData) (*rest.Config, error) {
	conn, err := GetNewClientCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	return conn.(*rest.Config), nil
}

var GetNewClientCached = plugin.HydrateFunc(GetNewClientUncached).Memoize()

// GetNewClientUncached :: gets client for querying openshift apis for the provided context
func GetNewClientUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	// get openshift config info
	openshiftConfig := GetConfig(d.Connection)

	// Set default loader and overriding rules
	loader := &clientcmd.ClientConfigLoadingRules{}
	overrides := &clientcmd.ConfigOverrides{}

	// variable to store paths for kubernetes config
	// default kube config path
	var configPaths = []string{"~/.kube/config"}

	// if openshiftConfig.ConfigPath != nil {
	// 	configPaths = []string{*openshiftConfig.ConfigPath}
	// } else if openshiftConfig.ConfigPaths != nil && len(openshiftConfig.ConfigPaths) > 0 {
	// 	configPaths = openshiftConfig.ConfigPaths
	// } else if v := os.Getenv("KUBE_CONFIG_PATHS"); v != "" {
	// 	configPaths = filepath.SplitList(v)
	// } else if v := os.Getenv("KUBERNETES_MASTER"); v != "" {
	// 	configPaths = []string{v}
	// }

	if len(configPaths) > 0 {
		expandedPaths := []string{}
		for _, p := range configPaths {
			path, err := homedir.Expand(p)
			if err != nil {
				return nil, err
			}
			expandedPaths = append(expandedPaths, path)
		}

		if len(expandedPaths) == 1 {
			loader.ExplicitPath = expandedPaths[0]
		} else {
			loader.Precedence = expandedPaths
		}

		if openshiftConfig.ConfigContext != nil {
			overrides.CurrentContext = *openshiftConfig.ConfigContext
			overrides.Context = clientcmdapi.Context{}
		}
	}

	osConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, overrides)

	// Get a rest.Config from the osConfig file.
	restconfig, err := osConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return restconfig, err
}

func v1TimeToRFC3339(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	switch v := d.Value.(type) {
	case v1.Time:
		return v.ToUnstructured(), nil
	case *v1.Time:
		if v == nil {
			return nil, nil
		}
		return v.ToUnstructured(), nil
	default:
		return nil, fmt.Errorf("invalid time format %T! ", v)
	}
}

func selectorMapToString(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	selector_map := d.Value.(map[string]string)

	if len(selector_map) == 0 {
		return nil, nil
	}

	selector_string := labels.SelectorFromSet(selector_map).String()

	return selector_string, nil
}

func labelSelectorToString(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	selector := d.Value.(*v1.LabelSelector)

	ss, err := v1.LabelSelectorAsSelector(selector)
	if err != nil {
		return nil, err
	}

	return ss.String(), nil
}
