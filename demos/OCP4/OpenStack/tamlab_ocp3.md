# Deploy OpenShift 3.11 on TAM LAB(BRQ)

This doc explains how to deploy openshift 3.11 on BRQ tamlab. 

In order to use this script, you need to know [OpenShift Service Account](https://access.redhat.com/terms-based-registry) so please prepare it at first.

If you know sa_name, you can get the information using the following URL
```
https://access.redhat.com/terms-based-registry/#/token/${sa_name}
```

**NOTE:** This script creates all objects from the scratch and delete them all. At the moment, tamlab did not integrate with corp LDAP so this script will create even user and the user will be removed when it uninstalls.


## Steps

### Pre-requisites

- Install openstack client
```
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
python get-pip.py
pip install python-openstackclient
```
- Clone git repository
```
git clone https://github.com/Jooho/tamlab-ocp3.git
cd tamlab-ocp3
```

### Update openshift.yaml
You will put confidential information here so ansible-vault will be used. 
In this demo, I will use "redhat" for vault password but you must create ansible-password file under `$HOME` folder

- Create vault password
```
cat redhat > ~/ansible-password
```

- Update 'CHANGE_ME' part.
```
vi openshift.yaml

subs_id: 'CHANGE_ME'
subs_pw: 'CHANGE_ME'

oreg_auth_user_name: "CHANGE_ME"    <=== You have to prepare this information from access.redhat.com
oreg_auth_token: 'CHANGE_ME'        <=== You have to prepare this information from access.redhat.com
```

- Encrypt the file
```
ansible-vault encrypt openshift.yaml
Vault password: redhat
```


### Update openstack.yaml
Like openshift.yaml file, this file also will be encrypted using ansible-vault

- Update 'CHANGE_ME'
```
vi openstack.yaml

# openstack cluster admin password
admin_user_pw: CHANGE_ME

# tamlab_dns_server
tamlab_dns_server_ip: CHANGE_ME
tamlab_dns_server_user: CHANGE_ME
tamlab_dns_server_pass: CHANGE_ME

```

- Encrypt the file
```
ansible-vault encrypt openstack.yaml
Vault password: redhat
```


### Update vars/all
This file contains most of the detail information for OpenShift and Openstack but you don't need to change much here.

- Update user information
```
vi vars/all

user_name: CHANGE_ME
user_pw: CHANGE_ME
user_email: CHANGE_ME
user_project: CHANGE_ME

```

### Install
If you follow the above instruction, now you can deploy OCP 3.11 on tamlab.
```
./install.sh
```


### Uninstall
```
./uninstall.sh
```
**NOTE** 
If you destroy OpenShift, you can not create an OCP cluster with the same cluster_name for **1 day**. When you want to recreate OCP, you can change cluster_name from vars/all file. Moreover, this script will delete even **project and user** as well.



## OCP Information

| name            | value                                                      |
| --------------- | ---------------------------------------------------------- |
| web console     | https://ocp-console.{USER_NAME}.{CLUSTER_NAME}.tamlab.brq.redhat.com:8443 |
| default user/pw | joe/redhat                                                 |

