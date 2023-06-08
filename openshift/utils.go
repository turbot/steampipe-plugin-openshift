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

	// default kube config path
	var configPath = "~/.kube/config"

	if openshiftConfig.ConfigPath != nil {
		configPath = *openshiftConfig.ConfigPath
	} else if v := os.Getenv("KUBE_CONFIG"); v != "" {
		configPath = v
	} else if v := os.Getenv("KUBECONFIG"); v != "" {
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

	// by default plugin will consider the current context available
	var flag = 0
	for k, v := range config {
		if k == "current-context" && strings.Contains(v.(string), "openshift") {
			overrides.CurrentContext = v.(string)
			overrides.Context = clientcmdapi.Context{}
			flag = 1
			break
		}
	}

	// return err if no openshift config available
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

func getCommonOptionalKeyQuals() []*plugin.KeyColumn {
	return []*plugin.KeyColumn{
		{Name: "name", Require: plugin.Optional},
		{Name: "namespace", Require: plugin.Optional},
	}
}

func getCommonOptionalKeyQualsValueForFieldSelector(d *plugin.QueryData) []string {
	fieldSelectors := []string{}

	if d.EqualsQualString("name") != "" {
		fieldSelectors = append(fieldSelectors, fmt.Sprintf("metadata.name=%v", d.EqualsQualString("name")))
	}

	if d.EqualsQualString("namespace") != "" {
		fieldSelectors = append(fieldSelectors, fmt.Sprintf("metadata.namespace=%v", d.EqualsQualString("namespace")))
	}

	return fieldSelectors
}
