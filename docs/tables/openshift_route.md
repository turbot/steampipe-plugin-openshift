---
title: "Steampipe Table: openshift_route - Query OpenShift Routes using SQL"
description: "Allows users to query OpenShift Routes, providing insights into the route objects that define the desired host for externally-reachable services."
---

# Table: openshift_route - Query OpenShift Routes using SQL

OpenShift Routes is a resource within Red Hat OpenShift that helps expose a service at a host name, like www.example.com, so that external clients can reach it by name. It provides a way to aggregate multiple services under the same IP address and differentiate them with the host name. OpenShift Routes makes it easy to expose services to the internet and manage traffic to your applications.

## Table Usage Guide

The `openshift_route` table provides insights into the route objects within Red Hat OpenShift. As a DevOps engineer, you can explore route-specific details through this table, including the host, path, and the associated services. Utilize it to manage and monitor the accessibility of your applications, ensuring they are reachable and functioning as expected.

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

### List routes that are present in the default namespace

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
