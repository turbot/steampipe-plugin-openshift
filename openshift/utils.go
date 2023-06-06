package openshift

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"gopkg.in/yaml.v2"
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
	var configPath = "~/.kube/config"

	if openshiftConfig.ConfigPath != nil {
		configPath = *openshiftConfig.ConfigPath
	} else if v := os.Getenv("KUBE_CONFIG"); v != "" {
		configPath = v
	} else if v := os.Getenv("KUBERNETES_MASTER"); v != "" {
		configPath = v
	}
	path, err := homedir.Expand(configPath)
	if err != nil {
		return nil, err
	}

	// by default plugin will consider first available openshift cluster context
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// take the first context available
	var flag = 0
	for k, v := range config {
		if k == "current-context" && strings.Contains(v.(string), "openshift") {
			overrides.CurrentContext = v.(string)
			flag = 1
			break
		}
	}

	// check if there is any openshift config available
	if flag == 0 {
		return nil, errors.New("openshift cluster details is unavailable in: " + path)
	}

	loader.ExplicitPath = path

	// override context if provided in the connection config file
	if openshiftConfig.ConfigContext != nil {
		overrides.CurrentContext = *openshiftConfig.ConfigContext
		overrides.Context = clientcmdapi.Context{}
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
