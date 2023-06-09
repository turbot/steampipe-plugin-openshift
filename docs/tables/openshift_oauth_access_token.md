# Table: openshift_oauth_access_token

OpenShift OAuth access token is a security token issued to a user upon successful authentication, which grants authorization to access OpenShift resources. It contains information about the user, their associated client, token expiration time, scopes, and authorization details. The access token is used to authenticate subsequent API requests to OpenShift, ensuring that only authorized users can access and perform actions on the resources within the cluster.

## Examples

### Basic info

```sql
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

```sql
select
  uid,
  name,
  expires_in,
  user_name,
  jsonb_pretty(scopes) as scopes
from
  openshift_oauth_access_token;
```

### List token with admin access

```sql
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

### List expired tokens

```sql
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
