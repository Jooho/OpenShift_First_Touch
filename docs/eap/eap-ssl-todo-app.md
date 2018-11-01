Deploy EAP SSL Sample Application using Paththrough
---------------------------------------------------

## Description

Red Hat EAP SSL Sample Application TODO App

## Deployment

```
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
```

