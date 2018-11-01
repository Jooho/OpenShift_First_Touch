
## Add test user (cluster admin)
```
ansible -i /etc/ansible/hosts masters -m command -a "htpasswd -bc /etc/origin/master/htpasswd joe redhat"
ansible -i /etc/ansible/hosts masters[0] -m command -a "oc adm policy add-cluster-role-to-user cluster-admin joe" 

```

## Template

### Clean error pods

*Project level*
```
oc get pods -o jsonpath='{.items[?(@.status.phase!="Running")].metadata.name}'|xargs oc delete pod
```

*Cluster level*
```
oc get pod --all-namespaces --template='{{ range $pod := .items}} {{if ne $pod.status.phase "Running"}}  -n {{$pod.metadata.namespace}} {{$pod.metadata.name}} {{end}} {{end}}' |bash -

```

*Orphan docker layers*
```
for im in $(docker images|grep '\<none' |awk '{print $3}'); do docker rmi --force $im;done
```



### [Tips] 

- Change Terminal Text Color to Normal
```
tput sgr0
```
