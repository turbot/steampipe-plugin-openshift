# Table: openshift_build_config

A build configuration, or BuildConfig, is characterized by a build strategy and one or more sources. The strategy determines the aforementioned process, while the sources provide its input.

## Examples

### Basic info

```sql
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

### List build configs present in default namespace

```sql
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

```sql
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

### List build configs created in last 30 days

```sql
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

### List common specs of the build configs

```sql
select
  uid,
  name,
  namespace,
  jsonb_pretty(common_spec) as common_spec
from
  openshift_build_config;
```

### Get triggers of the build configs

```sql
select
  uid,
  name,
  namespace,
  jsonb_pretty(triggers) as triggers
from
  openshift_build_config;
```

### List builds associated with a particular build config

```sql
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
