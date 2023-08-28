# Table: openshift_user

A user is an entity that interacts with the OpenShift Container Platform API. These can be a developer for developing applications or an administrator for managing the cluster. Users can be assigned to groups, which sets the permissions applied to all the group’s members. For example, you can give API access to a group, which gives all members of the group API access.

## Examples

### Basic info

```sql
select
  uid,
  name,
  full_name,
  resource_version,
  creation_timestamp,
  generation,
  deletion_grace_period_seconds
from
  openshift_user;
```

### List users who are not associated with any identities

```sql
select
  uid,
  name,
  full_name,
  resource_version,
  creation_timestamp,
  generation,
  deletion_grace_period_seconds
from
  openshift_user
where
  identities is null;
```

### List users created in the last 30 days

```sql
select
  uid,
  name,
  full_name,
  resource_version,
  creation_timestamp,
  generation,
  deletion_grace_period_seconds
from
  openshift_user
where
  creation_timestamp >= now() - interval '30' day;
```

### List users who have admin access

```sql
select
  distinct u.uid,
  u.name,
  full_name,
  u.resource_version,
  u.creation_timestamp,
  u.generation,
  u.deletion_grace_period_seconds
from
  openshift_user as u,
  openshift_oauth_access_token as t,
  jsonb_array_elements_text(scopes) as scope
where
  u.uid = t.user_uid
  and scope = 'user:full';
```

### List rules associated with a particular user

```sql
select
  uid,
  name,
  namespace,
  jsonb_pretty(rules) as rules
from
  kubernetes_role
where
  name in
  (
    select
      role_name
    from
      kubernetes_role_binding,
      jsonb_array_elements(subjects) as s
    where
      s ->> 'kind' = 'User'
      and s ->> 'name' = 'openshift_user'
  );
```
