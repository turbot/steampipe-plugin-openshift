---
title: "Steampipe Table: openshift_project - Query OpenShift Projects using SQL"
description: "Allows users to query OpenShift Projects, specifically retrieving details about project metadata, status, and specifications."
---

# Table: openshift_project - Query OpenShift Projects using SQL

OpenShift Projects are a top-level scope for managing and organizing resources in an OpenShift cluster. They provide a logical, hierarchical organization for a set of resources and users. Projects are essentially Kubernetes namespaces with additional annotations, providing a unique scope for objects such as pods, services, and replication controllers.

## Table Usage Guide

The `openshift_project` table provides insights into Projects within OpenShift. As a DevOps engineer or system administrator, explore project-specific details through this table, including metadata, status, and specifications. Utilize it to uncover information about projects, such as those with specific resource quotas, role bindings, and service accounts, aiding in the management and organization of your OpenShift cluster.

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
