Create User (HTPASSWD)
```
 oc create secret generic joe -n kube-system --from-literal=joe="$(htpasswd -bnB joe redhat | cut -d':' -f2)"
```

Delete kubeadmin
```
oc delete secret kubeadmin -n kube-system

```
