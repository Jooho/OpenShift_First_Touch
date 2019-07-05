# OpenShift 4 UPI Installation 

This script is for the demonstration how you can install Openshift 4 by UPI for bare metal with KVM. 

The script is created by ansible, terraform.

Environment:
- KVM ([Config_KVM](../Config_KVM/README.md))

## Overview

Please refer this [slide](https://www.slideshare.net/jooholee81/openshift4-installation-by-upi-on-kvm)

## Video
[![OpenShift 4 installation using UPI(Bare metal) with KVM Overview ](http://img.youtube.com/vi/qMyyZpoJ6o4/0.jpg)](https://www.youtube.com/watch?v=qMyyZpoJ6o4)

[![OpenShift 4 installation using UPI(Bare metal) with KVM Demo ](http://img.youtube.com/vi/UySkD8iWoSU/0.jpg)](https://www.youtube.com/watch?v=UySkD8iWoSU)

## Steps
**1. Init  (using ansible)**
   - Create:
     -  terraform
        - terraform.tfvars 
        - connection.tf
     - ansible
        - inventory
     - cloud-init
        - user-data
        - meta-data

**2. Prep (using ansible + terraform)**
- install_packages.yml (ansible)
    - docker
    - httpd
    - download_files.yml (ansible)
      - terraform_binary 
      - terraform_provider_libvirt
      - terraform_provider_matchbox
      - rhcos_kernel
      - rhcos_initramfs
      - rhcos_bios
      - openshift_installer
      - matchbox_git_repo
 
    - dns_config.yml (ansible)
      - DSNMasq

    - matchbox_config.yml (ansible)
      - Generate certs
      - Deploy matchbox server
  
    - ocp_vm_config.yml (ansible)
      - Create ocp_module.tf file

    - Network (terrraform)
  
    - MatchBox config (terraform)
        - config matchbox
  
    - Deploy LB 
      - Deploy VM (terraform)
      - Config haproxy (ansible)  

**3. OCP**
- Deploy OCP VM

**4. Post**
- Patch registry storage to emptyDir
- Wait for bootstrap-complete
- Remove bootstrap from lb pool
- Wait for install-complete
  


## jkit commands

jkit is python script to provide easier way to config each steps.

**1. Init**
    ```
    ./jkit.py init
    ```

**2. Prep**
    ```
   ./jkit.py prep -t [apply,dtr]  # dtr = destory
    ```

**3. OCP**
    ```
    ./jkit.py ocp -t [apply,dtr]   # dtr = destory
    ```
**4. Post**
    ```
    ./jkit.py post
    ```

- Additional cmds
    - **Clean**
        ```
        jkit clean
        ```

    - **Update**
        ```
        jkit.py update -t [inventory(default),ocp,ocp_module]
        ```
  
    - **Oneshot**
        ```
        ./jkit.py oneshot
        ```
     

## OCP debugging cmds

*certificate* 
```
oc --config ${INSTALL_DIR}/auth/kubeconfig get nodes
oc --config ${INSTALL_DIR}/auth/kubeconfig get csr
oc --config ${INSTALL_DIR}/auth/kubeconfig get csr -o json | jq -r '.items[] | select(.status == {}) | .metadata.name' | xargs oc --config ${INSTALL_DIR}/auth/kubeconfig adm certificate approve
```

*cluster operator*
```
oc --config ${INSTALL_DIR}/auth/kubeconfig get clusteroperators
```
