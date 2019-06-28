
ansible-playbook -i config prep/ansible/tasks/generate_config_files.yml -vvv



ansible-playbook -i config prep/ansible/tasks/download_files.yml -vvv
ansible-playbook -i ./ansible/inventory ./ansible/tasks/dns_config.yml
ansible-playbook -i ./ansible/inventory ./ansible/tasks/cloud_init.yml




 TF_LOG=trace terraform apply -auto-approve



Destroy
```
ssh known_host
prep/ansible.log
prep/meta-data
prep/user-data
prep/terraform.tfstate.backup
prep/terraform.tfstate
prep/terraform.tfvars
prep/{cluster_name}
prep/matchbox-master


/usr/bin/openshift-installer
/usr/bin/README.md

~/.terraform.d/plugins/terraform-provider-libvirt
~/.terraform.d/plugins/terraform-provider-matchbox

docker kill matchbox_server

ansible-playbook -i prep/ansible/inventory prep/ansible/tasks/matchbox_config.yml -vvv -e @prep/ansible/defaults/main.yml

/etc/matchbox
/home/jooho/.matchbox

Process

0. Init  (ansible)
    -  terraform
        - terraform.tfvars 
        - connection.tf
    - ansible
        - inventory

0.1. cloud_init.yml
    - user-data
    - meta-data

```
command: jkit.py init
```



1. Prep - 
1.1. install_packages.yml
    - docker
    - httpd

1.2. download_files.yml
    - terraform_binary 
    - terraform_provider_libvirt
    - terraform_provider_matchbox
    - rhcos_kernel
    - rhcos_initramfs
    - rhcos_bios
    - openshift_installer
    - matchbox_git_repo
 
1.3. dns_config.yml
        - DSNMasq


1.4. Network (terrraform)
1.5. LB 
    - Create VM
    - Config haproxy
    
        - matchbox
        - matchbox_module.tf
        -  terraform

TF_LOG=trace terraform apply -auto-approve


verify 
```
 openssl s_client -connect matchbox.example.com:8081 -CAfile /etc/matchbox/ca.crt -cert ~/.matchbox/client.crt -key ~/.matchbox/client.key
 ```

 ```
 $ dig matchbox.example.com
 $ curl http://matchbox.example.com:8080
matchbox
```



cmd

```
jkit.py init
jkit.py update --type ivn
jkit.py update --type ocp

.split("__")

#        'ansible-playbook %s  -i config prep/ansible/tasks/generate_config_files.yml --flush-cache; \
