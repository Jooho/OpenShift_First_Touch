Red Hat Sigle Sign On
---------------------

There are several ways to use RH SSO on OpenShift. This docs explain how to deploy RH SSO for demo.


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


## Prerequisites
```
export SECRETS_KEYSTORE_PASSWORD=redhat

export subdomain="example.redhat.com"

export API_SERVER="https://master.example.com:8443"

export cert_folder=/tmp/cert

oc new-project rh-sso
```

Tip. If you want to use random value for SECRETS_KEYSTORE_PASSWORD
```
export SECRETS_KEYSTORE_PASSWORD=$(openssl rand -base64 512 | tr -dc A-Z-a-z-0-9 | head -c 17)
```

## Steps 

- [Generate Certs and Create Secret.](./generate-certs.md)

- [Deploy RH-SSO](./deploy-sso.md)

- [Integrate RH SSO with OpenShift](./ocp-integration.md)

## Referecne

[Tutorial](https://access.redhat.com/documentation/en-us/red_hat_jboss_middleware_for_openshift/3/html/red_hat_single_sign-on_for_openshift/tutorials)

