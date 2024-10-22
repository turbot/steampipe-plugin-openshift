## v1.0.0 [2024-10-22]

There are no significant changes in this plugin version; it has been released to align with [Steampipe's v1.0.0](https://steampipe.io/changelog/steampipe-cli-v1-0-0) release. This plugin adheres to [semantic versioning](https://semver.org/#semantic-versioning-specification-semver), ensuring backward compatibility within each major version.

_Dependencies_

- Recompiled plugin with Go version `1.22`. ([#24](https://github.com/turbot/steampipe-plugin-openshift/pull/24))
- Recompiled plugin with [steampipe-plugin-sdk v5.10.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5104-2024-08-29) that fixes logging in the plugin export tool. ([#24](https://github.com/turbot/steampipe-plugin-openshift/pull/24))

## v0.2.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/docs), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#21](https://github.com/turbot/steampipe-plugin-openshift/pull/21))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#21](https://github.com/turbot/steampipe-plugin-openshift/pull/21))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-openshift/blob/main/docs/LICENSE). ([#21](https://github.com/turbot/steampipe-plugin-openshift/pull/21))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#20](https://github.com/turbot/steampipe-plugin-openshift/pull/20))

## v0.1.1 [2023-10-04]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#9](https://github.com/turbot/steampipe-plugin-openshift/pull/9))

## v0.1.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#6](https://github.com/turbot/steampipe-plugin-openshift/pull/6))
- Recompiled plugin with Go version `1.21`. ([#6](https://github.com/turbot/steampipe-plugin-openshift/pull/6))

## v0.0.1 [2023-08-30]

_What's new?_

- New tables added
  - [openshift_build_config](https://hub.steampipe.io/plugins/turbot/openshift/tables/openshift_build_config)
  - [openshift_build](https://hub.steampipe.io/plugins/turbot/openshift/tables/openshift_build)
  - [openshift_deployment_config](https://hub.steampipe.io/plugins/turbot/openshift/tables/openshift_deployment_config)
  - [openshift_image_stream](https://hub.steampipe.io/plugins/turbot/openshift/tables/openshift_image_stream)
  - [openshift_oauth_access_token](https://hub.steampipe.io/plugins/turbot/openshift/tables/openshift_oauth_access_token)
  - [openshift_project](https://hub.steampipe.io/plugins/turbot/openshift/tables/openshift_project)
  - [openshift_route](https://hub.steampipe.io/plugins/turbot/openshift/tables/openshift_route)
  - [openshift_user](https://hub.steampipe.io/plugins/turbot/openshift/tables/openshift_user)
