---
title: "Steampipe Table: openshift_user - Query OpenShift Users using SQL"
description: "Allows users to query OpenShift Users, specifically the user profiles and their associated metadata, providing insights into user management and access control within OpenShift."
---

# Table: openshift_user - Query OpenShift Users using SQL

OpenShift Users are the fundamental identity elements within OpenShift for authentication and authorization. They represent individual end users who may interact with the OpenShift API, and are associated with specific roles and permissions. User management in OpenShift is critical for controlling access, ensuring security, and maintaining operational efficiency.

## Table Usage Guide

The `openshift_user` table provides insights into user profiles within OpenShift. As a system administrator, explore user-specific details through this table, including user names, identities, and associated metadata. Utilize it to uncover information about users, such as their roles, permissions, and the overall user management landscape within your OpenShift environment.

## Examples

### Basic info

```sql+postgres
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

```sql+sqlite
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

```sql+postgres
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

```sql+sqlite
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

```sql+postgres
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

```sql+sqlite
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
  creation_timestamp >= datetime('now', '-30 day');
```

### List users who have admin access

```sql+postgres
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

```sql+sqlite
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
  json_each(t.scopes) as scope
where
  u.uid = t.user_uid
  and scope.value = 'user:full';
```

### List rules associated with a particular user

```sql+postgres
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

```sql+sqlite
select
  uid,
  name,
  namespace,
  rules
from
  kubernetes_role
where
  name in
  (
    select
      role_name
    from
      kubernetes_role_binding,
      json_each(subjects) as s
    where
      json_extract(s.value, '$.kind') = 'User'
      and json_extract(s.value, '$.name') = 'openshift_user'
  );
```
