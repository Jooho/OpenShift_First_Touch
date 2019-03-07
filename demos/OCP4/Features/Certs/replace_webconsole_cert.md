# How to replace Web console certificate

In OpenShift 4, you can replace the certificate for web console recreating secret `router-certs-default` in openshift-ingress project


## Pre-requisites
- OpenShift Container Platform 4
- Test cert/key 

## Generate self certificate for demo

Refer [this playbook](https://github.com/Jooho/ansible-cheat-sheet/tree/master/ansible-playbooks/ansible-playbook-generate-self-signed-cert)

```
git clone https://github.com/Jooho/ansible-cheat-sheet.git

cd ./ansible-cheat-sheet

ansible-galaxy install -f -r requirements.yaml

ansible-playbook ./playbook.yaml -e use_intermediate_cert=false -e cert_commonName=*.apps.ocp4.jlee.rhcee.support  -e cert_base_dir=/tmp/cert_base -b -vvvv
```
The result cert/key are in the `/tmp/cert_base/wild.apps.ocp4.jlee.rhcee.support/`
- *wild.apps.ocp4.jlee.rhcee.support.cert.pem*
- *wild.apps.ocp4.jlee.rhcee.support.key.pem*


## Replace cert

```
oc project openshift-ingress

oc delete secret router-certs-default

oc create secret tls router-certs-default  --cert=/tmp/cert_base/wild.apps.ocp4.jlee.rhcee.support/wild.apps.ocp4.jlee.rhcee.support.cert.pem --key=/tmp/cert_base/wild.apps.ocp4.jlee.rhcee.support/wild.apps.ocp4.jlee.rhcee.support.key.pem

oc rollout status deployment/router-default

```

## Check new cert
```
 openssl s_client -connect console-openshift-console.apps.ocp4.jlee.rhcee.support:443 -servername console-openshift-console.apps.ocp4.jlee.rhcee.support

CONNECTED(00000004)
depth=0 C = CA, ST = True, O = REDHAT, OU = SCE, CN = *.apps.ocp4.jlee.rhcee.support, emailAddress = test@test.com
verify error:num=20:unable to get local issuer certificate
verify return:1
depth=0 C = CA, ST = True, O = REDHAT, OU = SCE, CN = *.apps.ocp4.jlee.rhcee.support, emailAddress = test@test.com
verify error:num=21:unable to verify the first certificate
verify return:1
---
Certificate chain
 0 s:/C=CA/ST=True/O=REDHAT/OU=SCE/CN=*.apps.ocp4.jlee.rhcee.support/emailAddress=test@test.com
   i:/C=CA/ST=True/L=MILTON/O=REDHAT/OU=SCE/CN=Root CA/emailAddress=test@test.com

```




