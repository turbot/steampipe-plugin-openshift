connection "openshift" {
  plugin = "openshift"

  # By default, the plugin will use credentials in "~/.kube/config" with the current context.
  # The kubeconfig path and context can also be specified with the following config arguments:

  # Specify the file path to the kubeconfig.
  # Can also be set with the `KUBE_CONFIG` or `KUBECONFIG` environment variables.
  # config_path = "~/.kube/config"

  # Specify a context other than the current one. Optional.
  # config_context = "default/api-openshift-test-dq1i-p2-openshiftapps-com:6443/test"
}
