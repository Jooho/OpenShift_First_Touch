# Install NFS provisioner by Helm Chart

## Create export folder on specific node

Here, I will use app1.example.com as a target node where nfs-provisioner will run.

default nodeSelect is *{"app":"nfs-provisioner"}*
```
ssh centos@app1.example.com -- sudo mkdir /exports-nfs
ssh centos@app1.example.com -- sudo chcon -Rt svirt_sandbox_file_t /exports-nfs/

oc label node app1.example.com app=nfs-provisioner 
oc get node -l app=nfs-provisioner

```

## Deploy Helm Chart
```
helm install ../nfs-provisioner/ --name=nfs-provisioner --namespace=nfs-provioner
```

## Test Helm Chart
```
helm test nfs-provisioner --cleanup
```

## Delete pvc object( that is not controlled by Helm )
```
oc delete pvc test-pvc
```

## Clean Up
```
helm del nfs-provisioner --purge
```
