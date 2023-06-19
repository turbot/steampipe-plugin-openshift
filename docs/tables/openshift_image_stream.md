# Table: openshift_image_stream

Image streams provide a means of creating and updating container images in an on-going way. As improvements are made to an image, tags can be used to assign new version numbers and keep track of changes.

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

### List streams present in default namespace

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

### List deleted streams

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

### Get stream annotations

```sql
select
  uid,
  name,
  namespace,
  jsonb_pretty(annotations) as annotations
from
  openshift_image_stream;
```

### Get stream spec tags

```sql
select
  uid,
  name,
  namespace,
  jsonb_pretty(spec_tags) as spec_tags
from
  openshift_image_stream;
```

### Get stream status tags

```sql
select
  uid,
  name,
  namespace,
  jsonb_pretty(status_tags) as status_tags
from
  openshift_image_stream;
```