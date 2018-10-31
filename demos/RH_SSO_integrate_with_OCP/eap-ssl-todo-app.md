Deploy EAP SSO
--------------


## Description

Red Hat Single Sign-on is an integrated sign-on solution. This help deploying RH-SSO image on openshift

## Key Parameters

```
SSO_ADMIN_USERNAME=ssoadmin    # User for RH-SSO Admin Console
SSO_ADMIN_PASSWORD=redhat      # PW for RH-SSO Admin Console
SSO_REALM=ocp-auth             # Default Realm
SSO_SERVICE_USERNAME=svcuser   # Test user for the default Realm
SSO_SERVICE_PASSWORD=redhat    # Test user pw for the default Realm
```
To varify default test user:
 - Access RH-SSO Admin Console
 - Click `Clients` memu
 - Click Base URL of `account`
 - Then, login with test user/pw



## Required Certs
The RH-SSO template requires an SSL keystore and a JGroups keystore.



## Images (10.30.2018)

`For RH-SSO 7.2`

- sso72-https: RH-SSO 7.2 backed by internal H2 database on the same pod.
- sso72-mysql: RH-SSO 7.2 backed by ephemeral MySQL database on a separate pod.
- sso72-mysql-persistent: RH-SSO 7.2 backed by persistent MySQL database on a separate pod.
- sso72-postgresql: RH-SSO 7.2 backed by ephemeral PostgreSQL database on a separate pod.
- sso72-postgresql-persistent: RH-SSO 7.2 backed by persistent PostgreSQL database on a separate pod.


`For RH-SSO 7.1`

- sso71-https: RH-SSO 7.1 backed by internal H2 database on the same pod.
- sso71-mysql: RH-SSO 7.1 backed by ephemeral MySQL database on a separate pod.
- sso71-mysql-persistent: RH-SSO 7.1 backed by persistent MySQL database on a separate pod.
- sso71-postgresql: RH-SSO 7.1 backed by ephemeral PostgreSQL database on a separate pod.
- sso71-postgresql-persistent: RH-SSO 7.1 backed by persistent PostgreSQL database on a separate pod.

`Other templates that integrate with RH-SSO are also available`

- eap64-sso-s2i: RH-SSO-enabled Red Hat JBoss Enterprise Application Platform 6.4.
- eap70-sso-s2i: RH-SSO-enabled Red Hat JBoss Enterprise Application Platform 7.0.
- eap71-sso-s2i: RH-SSO enabled Red Hat JBoss Enterprise Application Platform 7.1.
- datavirt63-secure-s2i: RH-SSO-enabled Red Hat JBoss Data Virtualization 6.3.


## Deployment


### Using eap71-sso-s2i

```
export sso_hostname=secure-todo.cloudapps-37-0911.gsslab.rdu2.redhat.com

export SECRETS_KEYSTORE_PASSWORD=$(openssl rand -base64 512 | tr -dc A-Z-a-z-0-9 | head -c 17)

keytool -genkey -keyalg RSA -alias eapdemo-selfsigned -keystore keystore.jks -validity 360 -keysize 2048 -storepass redhat -dname "CN=Red Hat"

==> `Enter`

oc create secret generic eap7-app-secret --from-file=keystore.jks 

 oc new-app --template=eap71-mysql-persistent-s2i \
-p MYSQL_IMAGE_STREAM_TAG=5.7 \
-p APPLICATION_NAME=todo-app \
-p DB_PASSWORD=redhat \
-p DB_USERNAME=student \
-p DB_DATABASE=eap71db \
-p VOLUME_CAPACITY=1Gi \
-p HTTPS_NAME=eapdemo-selfsigned \
-p HTTPS_PASSWORD=redhat \
-p HTTPS_SECRET=eap7-app-secret \
--name=todo






## Reference 
[Tutorial](https://access.redhat.com/documentation/en-us/red_hat_jboss_middleware_for_openshift/3/html/red_hat_single_sign-on_for_openshift/tutorials)
  
