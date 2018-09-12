
## Add test user (cluster admin)
```
ansible -i /etc/ansible/hosts masters -m command -a "htpasswd -bc /etc/origin/master/htpasswd joe redhat"
ansible -i /etc/ansible/hosts masters[0] -m command -a "oc adm policy add-cluster-role-to-user cluster-admin joe" 


```
