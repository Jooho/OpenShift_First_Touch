Generate PKCS12 Cert for EAP EDGE TLS
--------------------------

# edge

## Export Variable
``
export hostname=abc.cloudapps.exmaple.com
``

## Create keystore (jks file)
keytool -genkey -keyalg RSA -alias selfsigned -keystore keystore.jks -storepass supersecret -validity 360 -keysize 2048 -storepass supersecret  -dname "CN=${hostname}"


## Change keystore type from jks to pkcs12
keytool -importkeystore -srckeystore keystore.jks  -destkeystore keystore.p12 -srcstoretype jks -deststoretype pkcs12 -srcstorepass supersecret -deststorepass supersecret


## Get private key and cert from pkcs12
openssl pkcs12 -in keystore.p12 -nodes -password pass:supersecret
