# Essential Variables
export SECRETS_KEYSTORE_PASSWORD=redhat
export subdomain="example.redhat.com"
export API_SERVER="https://master.example.com:8443"
export cert_folder=~/oft_cert
oc new-project rh-sso

# Generate Certs
mkdir ${cert_folder};cd ${cert_folder}

openssl req -new -newkey rsa:4096 -x509 -keyout xpaas.key -out xpaas.crt -days 365 -subj "/CN=xpaas-sso-demo.ca" -passin pass:${SECRETS_KEYSTORE_PASSWORD} -passout pass:${SECRETS_KEYSTORE_PASSWORD}

keytool -genkeypair -keyalg RSA -keysize 2048 -dname "CN=sso.${subdomain}" -alias eap-ssl-key -keystore eap-keystore.jks -storepass ${SECRETS_KEYSTORE_PASSWORD}
keytool -certreq -keyalg rsa -alias eap-ssl-key -keystore eap-keystore.jks -file sso.csr -storepass ${SECRETS_KEYSTORE_PASSWORD}

openssl x509 -req -CA xpaas.crt -CAkey xpaas.key -in sso.csr -out sso.crt -days 365 -CAcreateserial -passin pass:${SECRETS_KEYSTORE_PASSWORD}

keytool -import -file xpaas.crt -alias xpaas.ca -keystore eap-keystore.jks -storepass ${SECRETS_KEYSTORE_PASSWORD}
keytool -import -file sso.crt -alias eap-ssl-key -keystore eap-keystore.jks -storepass ${SECRETS_KEYSTORE_PASSWORD}
keytool -import -file xpaas.crt -alias xpaas.ca -keystore truststore.jks -storepass ${SECRETS_KEYSTORE_PASSWORD}
keytool -genseckey -alias jgroups -storetype JCEKS -keystore jgroups.jceks -storepass ${SECRETS_KEYSTORE_PASSWORD}


oc create secret generic sso-ssl-secret --from-file=truststore.jks
oc create secret generic eap-ssl-secret --from-file=eap-keystore.jks
oc create secret generic jgroups-secret --from-file=jgroups.jceks

oc secrets link default jgroups-secret sso-ssl-secret eap-ssl-secret


## Deploy
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


# Integrate with OCP
echo "Try to interate RH SSO with OCP: https://github.com/Jooho/jhouse_openshift/blob/master/demos/RH_SSO_integrate_with_OCP/ocp-integration.md"

