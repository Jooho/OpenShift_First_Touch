Recover the all ETCD members at one time
----------------------------------------

If you encounter this issue, this is really critical. Dislike v2 api, using snapshot way with v3 api, we can recover all etcd members at one time.

This doc will help you restore all ETCD members.

*Every command should be executed on the target ETCD member node unless there is not NOTICE for node change*

[Variable information](../backup_v3.md)

## Target ETCD member ##
- vm49.gsslab.rdu2.redhat.com(10.10.178.49)


**Note: Execute the following commands on each ETCD node**

## Export snapshot db folder
~~~
export MYBACKUPDIR=/root/backup/etcd/20190301    #UPDATE
~~~

## Stop docker/atomic-openshift-node
```
systemctl stop docker atomic-openshift-node 
rm -rf /var/lib/etcd
```

## Install ETCD on all ETCD nodes ##
```
yum install -y etcd

mv /etc/etcd/etcd.conf.rpmsave /etc/etcd/etcd.conf
```

**NOTE: Execute the following commands on target ETCD node(vm49)**
## Restore Data on a etcd node
```
export ETCDCTL_API=3
rm -rf /var/lib/etcd

source /etc/etcd/etcd.conf

etcdctl snapshot restore ${MYBACKUPDIR}/var/lib/etcd/snapshot.db \
  --name $ETCD_NAME \
  --initial-cluster $ETCD_INITIAL_CLUSTER \
  --initial-cluster-token $ETCD_INITIAL_CLUSTER_TOKEN \
  --initial-advertise-peer-urls $ETCD_INITIAL_ADVERTISE_PEER_URLS \
  --data-dir /var/lib/etcd

chown -R etcd:etcd /var/lib/etcd
restorecon -Rv /var/lib/etcd
```

**Note: Execute the following commands on each ETCD node**
## Start docker/atomic-openshift-node for all nodes
```
mkdir /var/lib/etcd
chown -R etcd:etcd /var/lib/etcd
restorecon -Rv /var/lib/etcd
systemctl start docker atomic-openshift-node 
```

## Check if backup data is recovered ##
```
etcdctl3 --endpoints $etcd_members endpoint status -w table
```

## Try oc command ##
```
oc get pod
```

