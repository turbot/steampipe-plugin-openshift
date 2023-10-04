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
