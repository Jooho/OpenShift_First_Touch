
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
 oc get pod --all-namespaces --template='{{ range $pod := .items}}{{if ne $pod.status.phase "Running"}} oc delete pod -n {{$pod.metadata.namespace}} {{$pod.metadata.name}} {{"\n"}}{{end}}{{end}}'  |bash -

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

### [Storages]
*storage speed check*
```
fio --rw=write --ioengine=sync --fdatasync=1 --directory=./test-data --size=22m --bs=2300 --name=mytest

dd if=/dev/zero of=/var/lib/etcd/test/abc.img bs=8k count=10k oflag=dsync

```
