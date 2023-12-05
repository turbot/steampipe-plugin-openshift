---
title: "Steampipe Table: openshift_build_config - Query OpenShift Build Configurations using SQL"
description: "Allows users to query OpenShift Build Configurations, specifically providing details about the build strategy, source, output, and triggers."
---

# Table: openshift_build_config - Query OpenShift Build Configurations using SQL

OpenShift Build Configurations are a key resource in the OpenShift platform that defines the parameters for building an application. They specify the source code repository, the build strategy (Docker, Source-to-Image, or Custom), and the output image. Build Configurations also specify triggers, which cause a new build to be run automatically when certain events occur.

## Table Usage Guide

The `openshift_build_config` table provides insights into the configuration details of builds within the OpenShift platform. As a DevOps engineer or system administrator, explore build-specific details through this table, including the source code repository, build strategy, output image, and triggers. Utilize it to manage and monitor the build process, ensuring that builds are configured correctly and running as expected.

## Examples

### Basic info
Explore which OpenShift build configurations have been created and their respective run policies. This can help in understanding the frequency of successful builds, thereby aiding in optimizing the build process.

```sql+postgres
select
  uid,
  name,
  namespace,
  run_policy,
  creation_timestamp,
  successful_builds_history_limit
from
  openshift_build_config;
```

```sql+sqlite
select
  uid,
  name,
  namespace,
  run_policy,
  creation_timestamp,
  successful_builds_history_limit
from
  openshift_build_config;
```

### List build configs present in the default namespace
Explore the build configurations present in the default workspace to understand the build policies, creation times, and history limits of successful builds. This can be useful in managing and optimizing your build processes.

```sql+postgres
select
  uid,
  name,
  namespace,
  run_policy,
  creation_timestamp,
  successful_builds_history_limit
from
  openshift_build_config
where
  namespace = 'default';
```

```sql+sqlite
select
  uid,
  name,
  namespace,
  run_policy,
  creation_timestamp,
  successful_builds_history_limit
from
  openshift_build_config
where
  namespace = 'default';
```

### List build configs with default run policy
Explore which build configurations are using the default 'Serial' run policy. This can help in managing and optimizing build processes in an OpenShift environment.

```sql+postgres
select
  uid,
  name,
  namespace,
  run_policy,
  creation_timestamp,
  successful_builds_history_limit
from
  openshift_build_config
where
  run_policy = 'Serial';
```

```sql+sqlite
select
  uid,
  name,
  namespace,
  run_policy,
  creation_timestamp,
  successful_builds_history_limit
from
  openshift_build_config
where
  run_policy = 'Serial';
```

### List build configs created in the last 30 days
Explore the recently created build configurations to understand their run policies and success rates. This can help assess the efficiency of your configurations and identify areas for improvement.

```sql+postgres
select
  uid,
  name,
  namespace,
  run_policy,
  creation_timestamp,
  successful_builds_history_limit
from
  openshift_build_config
where
  creation_timestamp >= now() - interval '30' day;
```

```sql+sqlite
select
  uid,
  name,
  namespace,
  run_policy,
  creation_timestamp,
  successful_builds_history_limit
from
  openshift_build_config
where
  creation_timestamp >= datetime('now','-30 day');
```

### List common specs of the build configs
Discover the common specifications of your build configurations in OpenShift to better understand and manage your application's build process. This can be particularly useful for identifying potential areas of optimization or troubleshooting issues.

```sql+postgres
select
  uid,
  name,
  namespace,
  jsonb_pretty(common_spec) as common_spec
from
  openshift_build_config;
```

```sql+sqlite
select
  uid,
  name,
  namespace,
  common_spec
from
  openshift_build_config;
```

### Get triggers of the build configs
Explore which build configurations in your Openshift environment have triggers set up. This can help you understand and manage automated build processes in your system.

```sql+postgres
select
  uid,
  name,
  namespace,
  jsonb_pretty(triggers) as triggers
from
  openshift_build_config;
```

```sql+sqlite
select
  uid,
  name,
  namespace,
  triggers
from
  openshift_build_config;
```

### List builds associated with a particular build config
Explore which builds are linked to a specific configuration to understand their performance and status. This is beneficial for assessing the efficiency of different configurations and identifying any potential issues.

```sql+postgres
select
  b.uid,
  b.name,
  b.namespace,
  b.start_timestamp,
  b.reason,
  b.phase,
  b.cancelled,
  b.duration,
  b.completion_timestamp
from
  openshift_build as b,
  jsonb_array_elements(owner_references) as ref,
  openshift_build_config as c
where
  ref ->> 'uid' = c.uid
  and c.name = 'config_name';
```

```sql+sqlite
select
  b.uid,
  b.name,
  b.namespace,
  b.start_timestamp,
  b.reason,
  b.phase,
  b.cancelled,
  b.duration,
  b.completion_timestamp
from
  openshift_build as b,
  json_each(b.owner_references) as ref,
  openshift_build_config as c
where
  json_extract(ref.value, '$.uid') = c.uid
  and c.name = 'config_name';
```