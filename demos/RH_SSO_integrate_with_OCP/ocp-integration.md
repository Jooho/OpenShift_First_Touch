Integration with OCP
--------------------

## SSO Side

Access to SSO Admin Console "https://sso.${subdomain}"
### New Realm for OCP 

- Mouse over `Master` on left menu top
- Clieck `Add realm`
- Type `OpenShift` on name field
- Click `Create`

### Create a new Client

- Click `Clients` on left menu
- Click `Create` on the table
- Type `ocp-auth` on client ID field
- Set `confidential` from Access Type field
- Type `Valid Redirect URIs` to "https://${API_SERVER}/*"
- Click "Save"


### Gather SSO Information

- Click `Credentials`
- Copy Secret text
- Export the text 
```
export sso_secret_text=$test
```

- Click `Installation`
- Click `Format` and set `Keyloak OIDC JSON`
- Copy `auth-server-url`
- Export the url 
```
export auth_server_url=$test

(ex)
export auth_server_url='https://sso.cloudapps-37-0911.gsslab.rdu2.redhat.com/auth'
```

### New User 

- Create a user
  - Click `Users` on left menu
  - Click `Add user` on table
  - Type `jhouse` on Username field
  - Click `Save`

- Set password of the user
  - Click `Credentials` tab
  - Type `redhat` on "New password" and "Password Confirmation" field.
  - Disable `Temporary` slide
  - Click `Reset Password`
  - Click `Change password`
  - Disable `Temporary` slide again  <== Seems like a bug


### Get information
```
curl -k  https://sso.cloudapps-37-0911.gsslab.rdu2.redhat.com/auth/realms/$REALM_NAME/.well-known/openid-configuration | python -m json.tool

(ex)
curl -k  https://sso.cloudapps-37-0911.gsslab.rdu2.redhat.com/auth/realms/OpenShift/.well-known/openid-configuration | python -m json.tool
```

## Master Server

### Prepare configuration for openID Identity Provider
```
cat <<EOF> ./openid.txt 
  - name: rh_sso
    challenge: false
    login: true
    mappingInfo: add
    provider:
      apiVersion: v1
      kind: OpenIDIdentityProvider
      clientID: ocp-auth
      clientSecret: ${sso_secret_text}
      ca: xpaas.crt
      urls:
        authorize: ${auth_server_url}/realms/Openshift/protocol/openid-connect/auth
        token: ${auth_server_url}/realms/Openshift/protocol/openid-connect/token
        userInfo: ${auth_server_url}/realms/Openshift/protocol/openid-connect/userinfo
      claims:
        id:
        - sub
        preferredUsername:
        - preferred_username
        name:
        - name
        email:
        - email
EOF
```





### Update All Masters

####  Prior OCP 3.10

### Copy xpaas.crt  (CA file)
```
cp  ${cert_folder}/xpaas.crt  /etc/origin/master/
```
```
# Back up /etc/origin/master/master
cp  /etc/origin/master/master-config.yaml /etc/origin/master/master-config.yaml.orig

# Insert openid.txt to master config file after `identityProvider:`
vi /etc/origin/master/master-config.yaml

...
oauthConfig:
  assetPublicURL: https://dhcp181-240.gsslab.rdu2.redhat.com:8443/console/
  grantConfig:
    method: auto
  identityProviders:
  - name: rh_sso
    challenge: false
    login: true
    mappingInfo: add
    provider:
    ....

# Restart Master API Server/ Controller
systemctl restart atomic-openshfit-master-api  atomic-openshift-master-controller
```


Tip. Using Ansible
```
ansible -i /etc/ansible/hosts masters -m blockinfile -a "block={{lookup('file','openid.txt')}} dest=/etc/origin/master/master-config.yaml backup=yes marker='### RH SSO By Ansible ###' insertafter='identityProviders'"

ansible -i /etc/ansible/hosts masters -m copy -a "src=${cert_folder}/xpaas.crt dest=/etc/origin/master/xpaas.crt"
ansible -i /etc/ansible/hosts masters -m systemd -a "name=atomic-openshift-master-api state=restarted"
ansible -i /etc/ansible/hosts masters -m systemd -a "name=atomic-openshift-master-controllers state=restarted"

```






Update master-config ConfigMap




## Test

- Go to Master API Server on browser
- Login with `rh-sso` identityProvider using `jhouse/redhat`



