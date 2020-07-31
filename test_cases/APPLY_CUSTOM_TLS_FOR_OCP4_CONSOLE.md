# Apply a custom cert for ocp4 custom console url

## Configuration
~~~
export cert_dir=/tmp/certs
export work_dir=/tmp/ssl
export subdomain=apps.bell.tamlab.brq.redhat.com

export custom_console_hostname=custom-console.${subdomain}
export custom_console_secret=custom-console-tls

~~~

## [Generate SSL](https://github.com/Jooho/ansible-cheat-sheet/tree/master/ansible-playbooks/ansible-playbook-generate-self-signed-cert)
~~~
# Setup dir & clone git repo
mkdir ${work_dir}
cd ${work_dir}
git clone git@github.com:Jooho/ansible-cheat-sheet.git
cd ansible-cheat-sheet/ansible-playbooks/ansible-playbook-generate-self-signed-cert/
ansible-galaxy install -f -r requirements.yaml

# Generate Certificate
ansible-playbook ./playbook.yaml -e use_intermediate_cert=false -e cert_commonName=${custom_console_hostname} -e cert_base_dir=${cert_dir} -b -vvvv

# Apply a custom cert for OCP4 custom console url
cd ${cert_dir}
sudo mv ./${custom_console_hostname}/${custom_console_hostname}.cert.pem tls.crt
sudo mv ./${custom_console_hostname}/${custom_console_hostname}.key.pem tls.key

# Verify Cert
openssl x509 -in tls.crt -text

# Create a Secret with the cert/key
oc create secret tls ${custom_console_secret} --cert=tls.crt --key=tls.key -n openshift-config

# Update console operator cr
echo "oc patch console.operator cluster --patch '{\"spec\":{\"route\":{\"hostname\":\"${custom_console_hostname}\",\"secret\": {\"name\": \"${custom_console_secret}\"}}}}' --type=merge" |bash -

# If you already have the param
#oc patch console.operator cluster -p="[{'op': 'replace', 'path': '/spec/route/hostname', 'value': ${custom_console_hostname}}]" --type='json'
#oc patch console.operator cluster -p="[{'op': 'replace', 'path': '/spec/route/secret/name', 'value': ${custom_console_secret}}]" --type='json'

# Check console operator cr
oc get console.operator cluster -o json | jq '.spec'
..
spec:
  managementState: Managed
  route:
    hostname: custom-console.apps.bell.tamlab.brq.redhat.com
    secret:
      name: custom-console-tls
..
~~~

## Verify
~~~
oc get route -n openshift-console
~~~

## Clean up
~~~
sudo rm -rf ${cert_dir} ${work_dir}
oc delete secret ${custom_console_secret}  -n openshift-config
oc patch console.operator cluster -p="[{'op': 'remove', 'path': '/spec/route'}]" --type='json'
~~~







~~~
