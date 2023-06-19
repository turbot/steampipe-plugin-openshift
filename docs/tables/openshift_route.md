# Table: openshift_route

A route allows you to host your application at a public URL. It can either be secure or unsecured, depending on the network security configuration of your application. An HTTP-based route is an unsecured route that uses the basic HTTP routing protocol and exposes a service on an unsecured application port.

## Examples

### Basic info

```sql
select
  uid,
  name,
  path,
  host,
  creation_timestamp,
  resource_version,
  namespace
from
  openshift_route;
```

### List routes which are present in default namespace

```sql
select
  uid,
  name,
  path,
  host,
  creation_timestamp,
  resource_version,
  namespace
from
  openshift_route
where
  namespace = 'default';
```

### List deleted routes

```sql
select
  uid,
  name,
  path,
  host,
  creation_timestamp,
  resource_version,
  namespace
from
  openshift_route
where
  deletion_timestamp is not null;
```

### List routes ingresses

```sql
select
  uid,
  name,
  namespace,
  jsonb_pretty(ingress) as ingress
from
  openshift_route;
```

### List routes associated with a particular service

```sql
select
  uid,
  name,
  path,
  host,
  creation_timestamp,
  resource_version,
  namespace
from
  openshift_route
where
  spec_to ->> 'kind' = 'Service'
  and spec_to ->> 'name' = 'console';
```

### List routes associated with a particular daemonset

```sql
select
  uid,
  name,
  path,
  host,
  creation_timestamp,
  resource_version,
  namespace
from
  openshift_route,
  jsonb_array_elements(owner_references) owner
where
  owner ->> 'kind' = 'daemonset'
  and owner ->> 'name' = 'ingress-canary';
```
