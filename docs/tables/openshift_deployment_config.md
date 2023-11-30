---
title: "Steampipe Table: openshift_deployment_config - Query OpenShift Deployment Configs using SQL"
description: "Allows users to query OpenShift Deployment Configs, specifically the deployment configurations in an OpenShift cluster, providing insights into deployment strategies, triggers, and details about the latest deployments."
---

# Table: openshift_deployment_config - Query OpenShift Deployment Configs using SQL

OpenShift Deployment Configs are a feature of OpenShift, a Kubernetes distribution from Red Hat. Deployment Configs define the template for a pod and manage the deployment of new pod versions when the template changes. They provide detailed control over the lifecycle of deployments, including strategies for rolling updates, manual intervention, and rollback.

## Table Usage Guide

The `openshift_deployment_config` table provides insights into Deployment Configs within OpenShift. As a DevOps engineer or system administrator, explore deployment-specific details through this table, including triggers, strategies, and information about the latest deployments. Utilize it to manage and monitor the deployment of new pod versions in an OpenShift cluster.

## Examples

### Basic info
Explore the status and details of different deployment configurations within an OpenShift environment. This allows for efficient monitoring and management of resources, ensuring optimal application performance and availability.

```sql
select
  uid,
  name,
  namespace,
  spec_replicas
  ready_replicas,
  updated_replicas,
  available_replicas,
  unavailable_replicas,
  creation_timestamp
from
  openshift_deployment_config;
```

### List deployment configs present in the default namespace
Explore the deployment configurations within the default namespace to understand the status and availability of replicas. This can be useful for managing resources and identifying potential issues in your OpenShift environment.

```sql
select
  uid,
  name,
  namespace,
  spec_replicas
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
Discover the segments that have paused deployment configurations. This can help in identifying instances where updates or changes have been halted, allowing for quick resolution of issues that may be causing the pause.

```sql
select
  uid,
  name,
  namespace,
  spec_replicas
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

### List deployment configs created in the last 30 days
Identify recent deployment configurations within the past month to understand their status and performance. This can be beneficial in monitoring system health and identifying potential issues early.

```sql
select
  uid,
  name,
  namespace,
  spec_replicas
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
Discover the specific container images used in your deployment configurations. This allows you to assess and manage your resource usage and maintain an inventory of your deployed images.

```sql
select
  name,
  namespace,
  c ->> 'name' as container_name,
  c ->> 'image' as image
from
  openshift_deployment_config,
  jsonb_array_elements(template -> 'Spec' -> 'Containers') as c
order by
  namespace,
  name;
```

### List pods for a particular deployment config
This query helps you to monitor and manage your Kubernetes deployment by listing all the pods associated with a specific deployment configuration. It's particularly useful when you need to track the status and location of pods within a particular deployment, such as troubleshooting issues or optimizing resource allocation.

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
Discover the deployment configurations that have access to the host process ID, IPC, or network. This is useful for identifying potential security risks in your Openshift environment.

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