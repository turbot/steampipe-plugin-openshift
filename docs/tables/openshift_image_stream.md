---
title: "Steampipe Table: openshift_image_stream - Query OpenShift Image Streams using SQL"
description: "Allows users to query OpenShift Image Streams, offering detailed information about the streams and related metadata."
---

# Table: openshift_image_stream - Query OpenShift Image Streams using SQL

An OpenShift Image Stream is a resource in OpenShift that provides an abstraction over related images. It allows users to track, tag, import and reference images based on tags, without explicit knowledge of the image repository or the image's specific registry URL. Image Streams also enable automatic updates whenever a new image is pushed to the registry.

## Table Usage Guide

The `openshift_image_stream` table provides insights into Image Streams within OpenShift. If you are a DevOps engineer or system administrator, you can explore details about image streams, including tags, annotations, and associated metadata through this table. Utilize it to manage and monitor image streams effectively, ensuring smooth application deployments and updates.

## Examples

### Basic info

```sql
select
  uid,
  name,
  namespace,
  resource_version,
  generation,
  docker_image_repository,
  creation_timestamp
from
  openshift_image_stream;
```

### List image streams present in the default namespace

```sql
select
  uid,
  name,
  namespace,
  resource_version,
  generation,
  docker_image_repository,
  creation_timestamp
from
  openshift_image_stream
where
  namespace = 'default';
```

### List deleted image streams

```sql
select
  uid,
  name,
  namespace,
  resource_version,
  generation,
  docker_image_repository,
  creation_timestamp
from
  openshift_image_stream
where
  deletion_timestamp is not null;
```

### Get image stream annotations

```sql
select
  uid,
  name,
  namespace,
  jsonb_pretty(annotations) as annotations
from
  openshift_image_stream;
```

### Get image stream spec tags

```sql
select
  uid,
  name,
  namespace,
  jsonb_pretty(spec_tags) as spec_tags
from
  openshift_image_stream;
```

### Get image stream status tags

```sql
select
  uid,
  name,
  namespace,
  jsonb_pretty(status_tags) as status_tags
from
  openshift_image_stream;
```
