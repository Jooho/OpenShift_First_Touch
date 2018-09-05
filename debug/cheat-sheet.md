
## Add test user (cluster admin)
```
ansible masters -m command -a "htpasswd -bc /etc/origin/master/htpasswd joe redhat"
ansible masters[0]-m command -a "oc adm policy add-cluster-role-to-user cluster-admin joe" 


```
