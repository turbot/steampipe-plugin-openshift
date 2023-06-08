---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/openshift.svg"
brand_color: ""
display_name: "OpenShift"
short_name: "openshift"
description: "Steampipe plugin to query projects, routes, builds and more from OpenShift."
og_description: "Query OpenShift with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/openshift-social-graphic.png"
---

# OpenShift + Steampipe

[OpenShift](https://docs.openshift.com/) is a container platform that provides a scalable and secure environment for deploying, managing, and scaling applications based on Kubernetes, enabling organizations to develop and run applications more efficiently and reliably.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

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

| Item        | Description                                                                                                                                                                                       |
| ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Credentials | OpenShift requires an `config_path` or `config_path` and `config_context` for all requests. |
| Permissions | The permission scope of Secret IDs is set by the Admin at the creation time of the [ACL tokens](https://developer.hashicorp.com/openshift/tutorials/web-ui/web-ui-access).                        |
| Radius      | Each connection represents a single OpenShift Installation.                                                                                                                                           |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/openshift.spc`)<br />2. Credentials specified in environment variables, e.g., `KUBE_CONFIG` and `KUBECONFIG`.      |

### Configuration

Installing the latest openshift plugin will create a config file (`~/.steampipe/config/openshift.spc`) with a single connection named `openshift`:

Configure your account details in `~/.steampipe/config/openshift.spc`:

```hcl
connection "openshift" {
  plugin = "openshift"

  # By default, the plugin will use credentials in "~/.kube/config" with the current context.
  # The kubeconfig path and context can also be specified with the following config arguments:

  # Specify the file path to the kubeconfig.
  # Can also be set with the "KUBE_CONFIG" or "KUBECONFIG" environment variables.
  # config_path = "~/.kube/config"

  # Specify a context other than the current one.
  # config_context = "default/api-openshift-test-dq1i-p2-openshiftapps-com:6443/test"
}
```

Alternatively, you can also use the standard OpenShift environment variable to obtain credentials **only if other arguments (`address`, `token`, and `namespace`) are not specified** in the connection:

```sh
export NOMAD_ADDR=http://18.118.144.168:4646
export NOMAD_TOKEN=c178b810-8b18-6f38-016f-725ddec5d58
export NOMAD_NAMESPACE=*
```

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-openshift
- Community: [Slack Channel](https://steampipe.io/community/join)
