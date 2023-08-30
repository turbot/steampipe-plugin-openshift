# Table: openshift_project

In OpenShift, a project is a container for organizing and isolating resources such as applications and services. It provides a boundary for access control and resource allocation, allowing multiple teams or users to work independently within their own project. Projects enable efficient management and collaboration in multi-tenant environments.

## Examples

### Basic info

```sql
select
  uid,
  name,
  resource_version,
  phase,
  creation_timestamp,
  deletion_grace_period_seconds,
  generate_name
from
  openshift_project;
```

### List inactive projects

```sql
select
  uid,
  name,
  resource_version,
  phase,
  creation_timestamp,
  deletion_grace_period_seconds,
  generate_name
from
  openshift_project
where
  phase <> 'Active';
```

### List projects created in the last 30 days

```sql
select
  uid,
  name,
  resource_version,
  phase,
  creation_timestamp,
  deletion_grace_period_seconds,
  generate_name
from
  openshift_project
where
  creation_timestamp >= now() - interval '30' day;
```

### List deleted projects

```sql
select
  uid,
  name,
  resource_version,
  phase,
  creation_timestamp,
  deletion_grace_period_seconds,
  generate_name
from
  openshift_project
where
  deletion_timestamp is not null;
```

### Get project annotations

```sql
select
  uid,
  name,
  phase,
  creation_timestamp,
  jsonb_pretty(annotations) as annotations
from
  openshift_project;
```

### Get project labels

```sql
select
  uid,
  name,
  phase,
  creation_timestamp,
  jsonb_pretty(labels) as labels
from
  openshift_project;
```
