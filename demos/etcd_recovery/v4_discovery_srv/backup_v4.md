# Backup ETCD Data for all members

## Gagher information
```
export KUBECONFIG=./aws/auth/kubeconfig

oc login -u kubeadmin -p XXX https://${API_SERVER}:6443
oc project kube-system

RUNNING_ETCD_POD=$(oc get pod -n kube-system --no-headers|grep Running| grep -o -m 1 '\S*etcd\S*')
RUNNING_ETCD_DNS_NAME=$(oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "cat /run/etcd/environment |grep ETCD_DNS|cut -d'=' -f2")


```

## Save ETCD snapshot DB
Actually, we don’t need to have snapshot on each etcd member but I will do it on every etcd member


```
for etcd in $(oc get pod -n kube-system --no-headers|grep -v $(hostname) |grep -o  '\S*etcd\S*' );
do
  echo $etcd

  RUNNING_ETCD_IP=$(oc get pod ${etcd} --template='{{.status.hostIP}}')
  RUNNING_ETCD_DNS_NAME=$(oc exec ${etcd}  -n kube-system  -- /bin/sh -c "cat /run/etcd/environment |grep ETCD_DNS|cut -d'=' -f2")

  oc exec ${etcd}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt  snapshot save /var/lib/etcd/snapshot.db"
  
done
```


## Move ETCD Backup DB

You need to create basion node to access master server. The public DNS is not created by default.

Execute the following command on each master node

```
ssh 
sudo cp /var/lib/etcd/snapshot.db /etc/snapshot.db

```