
## Add test user (cluster admin)
```
ansible masters -m command -a "htpasswd -bc /etc/origin/master/htpasswd joe redhat"
ansible tag_Name_36-1020_master_node_vms -i ./inventory/rhev/hosts/ovirt.py -m command -a "htpasswd -bc /etc/origin/master/htpasswd joe redhat"
```
