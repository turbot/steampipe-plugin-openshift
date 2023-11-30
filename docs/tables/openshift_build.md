---
title: "Steampipe Table: openshift_build - Query OpenShift Builds using SQL"
description: "Allows users to query OpenShift Builds, specifically providing details about the build configuration, status, source, output, and strategy."
---

# Table: openshift_build - Query OpenShift Builds using SQL

OpenShift Builds are a part of the OpenShift Container Platform build system, which provides a consistent method of turning application source code into a running containerized application. The build system captures source code from a version control system, processes it into a new container image, and pushes it to a container image registry. It supports several strategies to produce images.

## Table Usage Guide

The `openshift_build` table provides insights into the build configurations within OpenShift Container Platform. As a DevOps engineer, explore build-specific details through this table, including the build strategy, source, status, and output. Utilize it to track the progress of builds, understand the build strategies being used, and verify the output of completed builds.

## Examples

### Basic info
Analyze the settings to understand the overall performance and status of your OpenShift builds. This query helps to pinpoint specific instances where builds may have failed or been cancelled, providing insights to improve system efficiency and stability.

```sql
select
  uid,
  name,
  namespace,
  start_timestamp,
  reason,
  phase,
  cancelled,
  duration,
  completion_timestamp
from
  openshift_build;
```

### List incomplete builds
Discover the segments that are still in progress within your OpenShift environment. This query is useful for tracking ongoing processes and identifying any builds that may be stalled or delayed.

```sql
select
  uid,
  name,
  namespace,
  start_timestamp,
  reason,
  phase,
  cancelled,
  duration,
  completion_timestamp
from
  openshift_build
where
  phase <> 'Complete';
```

### List cancelled builds
Explore which builds have been cancelled in OpenShift to understand the reasons behind it and analyze the duration and time of cancellation. This helps in identifying any recurring issues and taking proactive measures to improve build success rates.

```sql
select
  uid,
  name,
  namespace,
  start_timestamp,
  reason,
  phase,
  cancelled,
  duration,
  completion_timestamp
from
  openshift_build
where
  cancelled;
```

### List common specs of the builds
Explore the common specifications of different builds to understand their unique configurations and phases. This can aid in identifying patterns or irregularities in build setups, thereby facilitating improved management and optimization of resources.

```sql
select
  uid,
  name,
  namespace,
  phase,
  jsonb_pretty(common_spec) as common_spec
from
  openshift_build;
```

### Get trigger details of the builds
Analyze the settings to understand the stages and conditions that initiate specific builds in a system. This can help in identifying the triggers and resolving any issues related to the build process.

```sql
select
  uid,
  name,
  namespace,
  phase,
  jsonb_pretty(triggered_by) as triggered_by
from
  openshift_build;
```

### Get stage details of the builds
Analyze the stages of various builds to understand their progress and status, which can be useful for managing and optimizing build processes within an OpenShift environment.

```sql
select
  uid,
  name,
  namespace,
  phase,
  jsonb_pretty(stages) as stages
from
  openshift_build;
```