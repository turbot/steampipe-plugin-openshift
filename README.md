![image](https://hub.steampipe.io/images/plugins/turbot/openshift-social-graphic.png)

# OpenShift Plugin for Steampipe

Use SQL to query projects, routes, builds and more from OpenShift.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/openshift)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/openshift/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-openshift/issues)

## Quick start

### Install

Download and install the latest OpenShift plugin:

```bash
steampipe plugin install openshift
```

Configure your [config file](https://hub.steampipe.io/plugins/turbot/openshift#configuration).

Add your configuration details in `~/.steampipe/config/openshift.spc`:

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

Run steampipe:

```shell
steampipe query
```

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
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-openshift.git
cd steampipe-plugin-openshift
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/openshift.spc
```

Try it!

```
steampipe query
> .inspect openshift
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-openshift/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [OpenShift Plugin](https://github.com/turbot/steampipe-plugin-openshift/labels/help%20wanted)