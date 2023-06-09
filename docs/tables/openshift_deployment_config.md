# Table: openshift_deployment_config

DeploymentConfig objects involve one or more replication controller, which contain a point-in-time record of the state of a deployment as a pod template.

## Examples

### Basic info

```sql
select
  uid,
  name,
  namespace,
  status_replicas,
  ready_replicas,
  updated_replicas,
  available_replicas,
  unavailable_replicas,
  creation_timestamp
from
  openshift_deployment_config;
```

### List deployment configs present in default namespace

```sql
select
  uid,
  name,
  namespace,
  status_replicas,
  ready_replicas,
  updated_replicas,
  available_replicas,
  unavailable_replicas,
  creation_timestamp
from
  openshift_deployment_config
where
  namespace = 'default';
```

### List paused deployment configs

```sql
select
  uid,
  name,
  namespace,
  status_replicas,
  ready_replicas,
  updated_replicas,
  available_replicas,
  unavailable_replicas,
  creation_timestamp
from
  openshift_deployment_config
where
  paused;
```

### List deployment configs created in last 30 days

```sql
select
  uid,
  name,
  namespace,
  status_replicas,
  ready_replicas,
  updated_replicas,
  available_replicas,
  unavailable_replicas,
  creation_timestamp
from
  openshift_deployment_config
where
  creation_timestamp >= now() - interval '30' day;
```

### Get container images used in deployment configs

```sql
select
  name,
  namespace,
  c ->> 'name' as container_name,
  c ->> 'image' as image
from
  openshift_deployment_config,
  jsonb_array_elements(template -> 'spec' -> 'containerc') as c
order by
  namespace,
  name;
```

### List pods for a particular deployment config

```sql
select
  pod.namespace,
  d.name as deployment_config_name,
  rc.name as replication_controller_name,
  pod.name as pod_name,
  pod.phase,
  age(current_timestamp, pod.creation_timestamp),
  pod.pod_ip,
  pod.node_name
from
  kubernetes_pod as pod,
  jsonb_array_elements(pod.owner_references) as pod_owner,
  kubernetes_replication_controller as rc,
  jsonb_array_elements(rc.owner_references) as rc_owner,
  openshift_deployment_config as d
where
  pod_owner ->> 'kind' = 'ReplicationController'
  and rc.uid = pod_owner ->> 'uid'
  and rc_owner ->> 'uid' = d.uid
  and d.name = 'sample-deployment'
order by
  pod.namespace,
  d.name,
  rc.name,
  pod.name;
```

### List deployment config with access to the host process ID, IPC, or network

```sql
select
  namespace,
  name,
  template -> 'spec' ->> 'hostPID' as hostPID,
  template -> 'spec' ->> 'hostIPC' as hostIPC,
  template -> 'spec' ->> 'hostNetwork' as hostNetwork
from
  openshift_deployment_config
where
  template -> 'spec' ->> 'hostPID' = 'true' or
  template -> 'spec' ->> 'hostIPC' = 'true' or
  template -> 'spec' ->> 'hostNetwork' = 'true';
```