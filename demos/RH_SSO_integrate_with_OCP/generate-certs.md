Certificates for RH SSO
-----------------------

## Required Certs
The RH-SSO template requires an SSL keystore and a JGroups keystore.

## Generate Certs

- Create Cert Folder
```
mkdir ${cert_folder};cd ${cert_folder}
```

- Generate a CA certificate
```
openssl req -new -newkey rsa:4096 -x509 -keyout xpaas.key -out xpaas.crt -days 365 -subj "/CN=xpaas-sso-demo.ca" -passin pass:${SECRETS_KEYSTORE_PASSWORD} -passout pass:${SECRETS_KEYSTORE_PASSWORD}
```

- Generate a Certificate for the SSL keystore:
```
keytool -genkeypair -keyalg RSA -keysize 2048 -dname "CN=sso.${subdomain}" -alias eap-ssl-key -keystore eap-keystore.jks -storepass ${SECRETS_KEYSTORE_PASSWORD}

Enter
```

- Generate a Certificate Sign Request for the SSL keystore
```
keytool -certreq -keyalg rsa -alias eap-ssl-key -keystore eap-keystore.jks -file sso.csr -storepass ${SECRETS_KEYSTORE_PASSWORD}
```

- Sign the Certificate Sign Request with the CA certificate
```
openssl x509 -req -CA xpaas.crt -CAkey xpaas.key -in sso.csr -out sso.crt -days 365 -CAcreateserial -passin pass:${SECRETS_KEYSTORE_PASSWORD}
```

- Import the CA into the SSL keystore
```
keytool -import -file xpaas.crt -alias xpaas.ca -keystore eap-keystore.jks -storepass ${SECRETS_KEYSTORE_PASSWORD}

Type "yes" and Enter
```

- Import the signed Certificate Sign Request(cert) into the SSL keystore
```
keytool -import -file sso.crt -alias eap-ssl-key -keystore eap-keystore.jks -storepass ${SECRETS_KEYSTORE_PASSWORD}
```

- Import the CA into a new truststore keystore
```
keytool -import -file xpaas.crt -alias xpaas.ca -keystore truststore.jks -storepass ${SECRETS_KEYSTORE_PASSWORD}

Type "yes" and Enter
```

- Generate a secure key for the JGroups keystore
```
keytool -genseckey -alias jgroups -storetype JCEKS -keystore jgroups.jceks -storepass ${SECRETS_KEYSTORE_PASSWORD}

Enter
```

## Create Secrets with the generated Certs
```
oc create secret generic sso-ssl-secret --from-file=truststore.jks
oc create secret generic eap-ssl-secret --from-file=eap-keystore.jks
oc create secret generic jgroups-secret --from-file=jgroups.jceks
```
## Link the secrets to the default service account
```
oc secrets link default jgroups-secret sso-ssl-secret eap-ssl-secret
```

