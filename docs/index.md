---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/openshift.svg"
brand_color: "#DB212E"
display_name: "OpenShift"
short_name: "openshift"
description: "Steampipe plugin to query projects, routes, builds and more from OpenShift."
og_description: "Query OpenShift with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/openshift-social-graphic.png"
---

# OpenShift + Steampipe

[OpenShift](https://docs.openshift.com/) is a container platform that provides a scalable and secure environment for deploying, managing, and scaling applications based on Kubernetes, enabling organizations to develop and run applications more efficiently and reliably.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

List your OpenShift projects:

```sql
select
  name,
  uid,
  phase,
  creation_timestamp,
  resource_version
from
  openshift_project;
```

```
+--------------------------------------------------+--------------------------------------+--------+---------------------------+------------------+
| name                                             | uid                                  | phase  | creation_timestamp        | resource_version |
+--------------------------------------------------+--------------------------------------+--------+---------------------------+------------------+
| openshift-authentication                         | cf62541e-1ad0-45f6-b023-cf695f32bffd | Active | 2023-06-05T17:45:53+05:30 | 6659             |
| openshift-apiserver                              | fb619658-fbb2-4735-87fe-386fa2897816 | Active | 2023-06-05T17:45:49+05:30 | 6650             |
| openshift-backplane-csa                          | 693ce132-49a7-4e00-b295-29de6be2fca7 | Active | 2023-06-05T18:07:01+05:30 | 30294            |
+--------------------------------------------------+--------------------------------------+--------+---------------------------+------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/openshift/tables)**

## Quick start

### Install

Download and install the latest OpenShift plugin:

```sh
steampipe plugin install openshift
```

### Credentials

No credentials are required.

### Configuration

Installing the latest openshift plugin will create a config file (`~/.steampipe/config/openshift.spc`) with a single connection named `openshift`:

Configure your account details in `~/.steampipe/config/openshift.spc`:

```hcl
connection "openshift" {
  plugin = "openshift"

  # By default, the plugin will use credentials in "~/.kube/config" with the current context.
  # The kubeconfig path and context can also be specified with the following config arguments:

  # Specify the file path to the kubeconfig. If not set, the plugin will check `~/.kube/config`.
  # Can also be set with the "KUBE_CONFIG" or "KUBECONFIG" environment variables.
  # config_path = "~/.kube/config"

  # Specify a context other than the current one. If not set, the current context will be used. Optional.
  # config_context = "default/api-openshift-test-dq1i-p2-openshiftapps-com:6443/test"
}
```

## Configuring OpenShift Credentials

By default, the plugin will use the kubeconfig in `~/.kube/config` with the current context. If using the default oc CLI configurations, the kubeconfig will be in this location and the OpenShift plugin connections will work by default.

You can also set the kubeconfig file path and context with the `config_path` and `config_context` config arguments respectively.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-openshift
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)