---
title: "Steampipe Table: openshift_oauth_access_token - Query OpenShift OAuth Access Tokens using SQL"
description: "Allows users to query OpenShift OAuth Access Tokens, specifically the unique identifiers, user names, client names, and expiration timestamps, providing insights into access control and potential security risks."
---

# Table: openshift_oauth_access_token - Query OpenShift OAuth Access Tokens using SQL

OpenShift OAuth Access Tokens are part of OpenShift's OAuth server which is responsible for managing user login and token-based authentication. These tokens are used to authenticate users and clients in the OpenShift cluster and provide them with the necessary permissions to access resources. The tokens consist of unique identifiers, user names, client names, and expiration timestamps, which can be used to monitor and manage access control and potential security risks.

## Table Usage Guide

The `openshift_oauth_access_token` table provides insights into OAuth Access Tokens within OpenShift. As a system administrator, explore token-specific details through this table, including unique identifiers, user names, client names, and expiration timestamps. Utilize it to uncover information about tokens, such as those associated with specific users or clients, the tokens' expiration timestamps, and potential security risks associated with expired or misused tokens.

## Examples

### Basic info
Gain insights into the authorization patterns by analyzing the validity and usage of access tokens in your Openshift environment. This is particularly useful for assessing security measures and identifying potential vulnerabilities.

```sql+postgres
select
  uid,
  name,
  authorize_token,
  refresh_token,
  expires_in,
  user_name,
  user_uid
from
  openshift_oauth_access_token;
```

```sql+sqlite
select
  uid,
  name,
  authorize_token,
  refresh_token,
  expires_in,
  user_name,
  user_uid
from
  openshift_oauth_access_token;
```

### Get token scopes
Assess the elements within your system to understand which user tokens are nearing expiration. This is beneficial in maintaining security and ensuring uninterrupted user access.

```sql+postgres
select
  uid,
  name,
  expires_in,
  user_name,
  jsonb_pretty(scopes) as scopes
from
  openshift_oauth_access_token;
```

```sql+sqlite
select
  uid,
  name,
  expires_in,
  user_name,
  scopes
from
  openshift_oauth_access_token;
```

### List tokens with admin access
Explore which tokens have full administrative access. This is beneficial in analyzing your security configuration and identifying potential vulnerabilities.

```sql+postgres
select
  uid,
  name,
  authorize_token,
  refresh_token,
  expires_in,
  user_name,
  user_uid
from
  openshift_oauth_access_token,
  jsonb_array_elements_text(scopes) as scope
where
  scope like '%full%';
```

```sql+sqlite
select
  uid,
  name,
  authorize_token,
  refresh_token,
  expires_in,
  user_name,
  user_uid
from
  openshift_oauth_access_token,
  json_each(scopes) as scope
where
  scope.value like '%full%';
```

### List expired tokens
Discover the segments that have expired tokens in the OpenShift OAuth access to identify potential security risks or unauthorized access. This is particularly useful for maintaining system integrity and ensuring user account safety.

```sql+postgres
select
  uid,
  name,
  authorize_token,
  refresh_token,
  expires_in,
  user_name,
  user_uid
from
  openshift_oauth_access_token
where
  extract(epoch from age(current_timestamp,creation_timestamp))::int > expires_in;
```

```sql+sqlite
select
  uid,
  name,
  authorize_token,
  refresh_token,
  expires_in,
  user_name,
  user_uid
from
  openshift_oauth_access_token
where
  cast((julianday('now') - julianday(creation_timestamp)) * 86400 as integer) > expires_in;
```