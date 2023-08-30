# Table: openshift_build

A build in OpenShift Container Platform is the process of transforming input parameters into a resulting object. Most often, builds are used to transform source code into a runnable container image.

## Examples

### Basic info

```sql
select
  uid,
  name,
  namespace,
  start_timestamp,
  reason,
  phase,
  cancelled,
  duration,
  completion_timestamp
from
  openshift_build;
```

### List incomplete builds

```sql
select
  uid,
  name,
  namespace,
  start_timestamp,
  reason,
  phase,
  cancelled,
  duration,
  completion_timestamp
from
  openshift_build
where
  phase <> 'Complete';
```

### List cancelled builds

```sql
select
  uid,
  name,
  namespace,
  start_timestamp,
  reason,
  phase,
  cancelled,
  duration,
  completion_timestamp
from
  openshift_build
where
  cancelled;
```

### List common specs of the builds

```sql
select
  uid,
  name,
  namespace,
  phase,
  jsonb_pretty(common_spec) as common_spec
from
  openshift_build;
```

### Get trigger details of the builds

```sql
select
  uid,
  name,
  namespace,
  phase,
  jsonb_pretty(triggered_by) as triggered_by
from
  openshift_build;
```

### Get stage details of the builds

```sql
select
  uid,
  name,
  namespace,
  phase,
  jsonb_pretty(stages) as stages
from
  openshift_build;
```