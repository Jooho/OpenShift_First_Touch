Deploying RH SSO
----------------

## Important Parameters

### EAP Parameters
------------------------------------------------
|Variable|      Demo Vaule        | Description|
|--------|------------------------|------------|
|APPLICATION_NAME| sso | |
|HTTPS_KEYSTORE| eap-keystore.jks | This keystore is for eap application that SSO is deployed on|
|HTTPS_PASSWORD| redhat | keytool -storepass value|
|HTTPS_SECRET| eap-ssl-secret |   |
|JGROUPS_ENCRYPT_KEYSTORE| jgroups.jceks | ?|
|JGROUPS_ENCRYPT_PASSWORD| redhat| |
|JGROUPS_ENCRYPT_SECRET| jgroups-secret| |
-----------------------------------------

### SSO Admin Center Parameters
-----------------------------------------------
|Variable|      Demo Vaule        |Description|
|--------|------------------------|-----------|
|SSO_REALM  | OpenShift| Test Realm for SSO that Openshift will use|
|SSO_SERVICE_USERNAME|testUser| User for SSO login|
|SSO_SERVICE_PASSWORD|testPw||
|SSO_ADMIN_USERNAME|ssoAdmin| User for SSO Admin Console login|
|SSO_ADMIN_USERNAME | redhat||
-------------------------------------

### SSO Certificate Parameters
-----------------------------------------------
|Variable|      Demo Vaule        |Description|
|--------|------------------------|-----------|
|SSO_TRUSTSTORE| truststore.jks| KeyStore between Client and SSO|
|SSO_TRUSTSTORE_SECRET| sso-ssl-secret||
|SSO_TRUSTSTORE_PASSWORD| redhat  ||
---------------------------------------

## Deploy

```
oc new-app  --template=sso72-mysql-persistent \
  -p APPLICATION_NAME=sso \
  -p HOSTNAME_HTTPS=sso.${subdomain} \
  -p HOSTNAME_HTTP=sso-http.${subdomain} \
  -p DB_JNDI=java:jboss/datasources/KeycloakDS \
  -p DB_DATABASE=sso72db \
  -p HTTPS_SECRET=eap-ssl-secret \
  -p HTTPS_KEYSTORE=eap-keystore.jks \
  -p HTTPS_NAME=eap-ssl-key \
  -p HTTPS_PASSWORD=${SECRETS_KEYSTORE_PASSWORD} \
  -p DB_USERNAME=jhouse \
  -p DB_PASSWORD=redhat \
  -p VOLUME_CAPACITY=2Gi \
  -p JGROUPS_ENCRYPT_SECRET=jgroups-secret \
  -p JGROUPS_ENCRYPT_KEYSTORE=jgroups.jceks \
  -p JGROUPS_ENCRYPT_NAME=jgroups \
  -p JGROUPS_ENCRYPT_PASSWORD=${SECRETS_KEYSTORE_PASSWORD} \
  -p JGROUPS_CLUSTER_PASSWORD=${SECRETS_KEYSTORE_PASSWORD} \
  -p IMAGE_STREAM_NAMESPACE=openshift \
  -p SSO_ADMIN_USERNAME=ssoadmin \
  -p SSO_ADMIN_PASSWORD=redhat \
  -p SSO_REALM=testRealm \
  -p SSO_SERVICE_USERNAME=testUser \
  -p SSO_SERVICE_PASSWORD=testPw \
  -p SSO_TRUSTSTORE_PASSWORD=${SECRETS_KEYSTORE_PASSWORD} \
  -p SSO_TRUSTSTORE=truststore.jks \
  -p MYSQL_IMAGE_STREAM_TAG=5.7 \
  -p SSO_TRUSTSTORE_SECRET=sso-ssl-secret \
  -p MEMORY_LIMIT=1Gi  \
  --name=sso
```

## Issues
- sso application will be restarted because database connection sometimes at deploying stage.
  So wait for the sso deploying properly.
 
## Test

- Login SSO Admin Console
  - Go to "https://sso.${subdomain} on browser
  - Click "Administration Console"
  - Login with "ssoadmin/redhat"

- Verify SSO with Test User
  - Click "Clients" on left Menu
  - Click "Base URL" of `account`(Client ID) on table
  - Login with "testUser/testPw"




