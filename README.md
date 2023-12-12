![image](https://hub.steampipe.io/images/plugins/turbot/openshift-social-graphic.png)

# OpenShift Plugin for Steampipe

Use SQL to query projects, routes, builds and more from OpenShift.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/openshift)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/openshift/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-openshift/issues)

## Quick start

### Install

Download and install the latest OpenShift plugin:

```bash
steampipe plugin install openshift
```

Configure your [config file](https://hub.steampipe.io/plugins/turbot/openshift#configuration).

Configure your account details in `~/.steampipe/config/openshift.spc`:

```hcl
connection "openshift" {
  plugin = "openshift"

  # Authentication information
  config_path = "~/.kube/config"
  config_context = "default/api-openshift-test-dq1i-p2-openshiftapps-com:6443/test"
}
```

Or through environment variables:

```sh
export KUBE_CONFIG=~/.kube/config
```

or

```sh
export KUBECONFIG=~/.kube/config
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
+--------------------------------------------------+--------------------------------------+--------+---------------------------+------------------+
```

## Engines

This plugin is available for the following engines:

| Engine        | Description
|---------------|------------------------------------------
| [Steampipe](https://steampipe.io/docs) | The Steampipe CLI exposes APIs and services as a high-performance relational database, giving you the ability to write SQL-based queries to explore dynamic data. Mods extend Steampipe's capabilities with dashboards, reports, and controls built with simple HCL. The Steampipe CLI is a turnkey solution that includes its own Postgres database, plugin management, and mod support.
| [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/index) | Steampipe Postgres FDWs are native Postgres Foreign Data Wrappers that translate APIs to foreign tables. Unlike Steampipe CLI, which ships with its own Postgres server instance, the Steampipe Postgres FDWs can be installed in any supported Postgres database version.
| [SQLite Extension](https://steampipe.io/docs//steampipe_sqlite/index) | Steampipe SQLite Extensions provide SQLite virtual tables that translate your queries into API calls, transparently fetching information from your API or service as you request it.
| [Export](https://steampipe.io/docs/steampipe_export/index) | Steampipe Plugin Exporters provide a flexible mechanism for exporting information from cloud services and APIs. Each exporter is a stand-alone binary that allows you to extract data using Steampipe plugins without a database.
| [Turbot Pipes](https://turbot.com/pipes/docs) | Turbot Pipes is the only intelligence, automation & security platform built specifically for DevOps. Pipes provide hosted Steampipe database instances, shared dashboards, snapshots, and more.

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

## Open Source & Contributing

This repository is published under the [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0) (source code) and [CC BY-NC-ND](https://creativecommons.org/licenses/by-nc-nd/2.0/) (docs) licenses. Please see our [code of conduct](https://github.com/turbot/.github/blob/main/CODE_OF_CONDUCT.md). We look forward to collaborating with you!

[Steampipe](https://steampipe.io) is a product produced from this open source software, exclusively by [Turbot HQ, Inc](https://turbot.com). It is distributed under our commercial terms. Others are allowed to make their own distribution of the software, but cannot use any of the Turbot trademarks, cloud services, etc. You can learn more in our [Open Source FAQ](https://turbot.com/open-source).

## Get Involved

**[Join #steampipe on Slack →](https://turbot.com/community/join)**

Want to help but don't know where to start? Pick up one of the `help wanted` issues:

- [OpenShift Plugin](https://github.com/turbot/steampipe-plugin-openshift/labels/help%20wanted)
- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
