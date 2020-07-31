# Apply a custom cert for ocp4 custom console url and a custom cert for default grafana url using SAN 

## Scenario
- Create SAN cert that has 2 DNS records
  - custom-console.apps.bell.tamlab.brq.redhat.com
  - grafana-openshift-monitoring.apps.bell.tamlab.brq.redhat.com

- Apply the cert to console route
- Apply the cert to grafana route



## Configuration
~~~
export cert_dir=/tmp/certs
export work_dir=/tmp/ssl
export subdomain=apps.bell.jlee.rhcee.support

export custom_console_hostname=custom-console.${subdomain}
export custom_console_secret=custom-console-tls

export grafana_hostname=grafana-openshift-monitoring.${subdomain}
~~~

## [Generate SSL](./GENERATE_SSL.md)
~~~
# Generate Certificate
mkdir ${work_dir}
cd ${work_dir}
git clone git@github.com:Jooho/ansible-cheat-sheet.git
cd ansible-cheat-sheet/ansible-playbooks/ansible-playbook-generate-self-signed-cert/
ansible-galaxy install -f -r requirements.yaml

ansible-playbook ./playbook.yaml -e use_intermediate_cert=false -e cert_commonName=${custom_console_hostname} -e cert_base_dir=${cert_dir} -b -e "{san_dns: [{ index: 1, dns: ${grafana_hostname}}]}" -vvvv

openssl x509 -in ${cert_dir}/${custom_console_hostname}/${custom_console_hostname}.cert.pem -text




# Apply a custom cert for OCP4 custom console url
cd ${cert_dir}
sudo mv ./${custom_console_hostname}/${custom_console_hostname}.cert.pem tls.crt
sudo mv ./${custom_console_hostname}/${custom_console_hostname}.key.pem tls.key

oc create secret tls ${custom_console_secret} --cert=tls.crt --key=tls.key -n openshift-config

echo "oc patch console.operator cluster --patch '{\"spec\":{\"route\":{\"hostname\":\"${custom_console_hostname}\",\"secret\": {\"name\": \"${custom_console_secret}\"}}}}' --type=merge" |bash -


# If you already have the param
#oc patch console.operator cluster -p="[{'op': 'replace', 'path': '/spec/route/hostname', 'value': ${custom_console_hostname}}]" --type='json'
#oc patch console.operator cluster -p="[{'op': 'replace', 'path': '/spec/route/secret/name', 'value': ${custom_console_secret}}]" --type='json'

oc get console.operator cluster -o json | jq '.spec'
..
spec:
  managementState: Managed
  route:
    hostname: custom-console.apps.bell.tamlab.brq.redhat.com
    secret:
      name: custom-console-tls
..

# Apply the custom cert for OCP4 grafana route
export grafana_route_crt=$(awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' tls.crt)
export grafana_route_key=$(awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;}' tls.key)

echo "oc patch route/grafana -n openshift-monitoring --patch='{\"spec\":{\"tls\":{\"key\":\"${grafana_route_key}\", \"certificate\": \"${grafana_route_crt}\"}}}' --type=merge" |bash -
~~~

## Verify
~~~
oc get route -n openshift-console
~~~

## Clean up
~~~
rm -rf ${cert_dir} ${work_dir}
oc patch console.operator cluster -p="[{'op': 'remove', 'path': '/spec/route'}]" --type='json'

oc delete secret ${custom_console_secret} -n openshift-config
~~~

