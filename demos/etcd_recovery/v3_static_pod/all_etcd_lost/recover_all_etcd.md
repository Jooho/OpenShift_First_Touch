Recover the all ETCD members at one time
----------------------------------------

If you encounter this issue, this is really critical. Dislike v2 api, using snapshot way with v3 api, we can recover all etcd members at one time.

This doc will help you restore all ETCD members.

*Every command should be executed on the target ETCD member node unless there is not NOTICE for node change*

[Variable information](../backup_v3.md)

## Target ETCD member ##
- vm125.gsslab.rdu2.redhat.com(10.10.178.125)


**Note: Execute the following commands on each ETCD node**

## Stop docker/atomic-openshift-node
```
systemctl stop docker atomic-openshift-node 
```

## Install ETCD on all ETCD nodes ##
```
yum install -y etcd

mv /etc/etcd/etcd.conf.rpmsave /etc/etcd/etcd.conf
```

## Restore Data
```
export ETCDCTL_API=3
export ETCD_DATA_PATH=/var/lib/etcd
rm -rf $ETCD_DATA_PATH

source /etc/etcd/etcd.conf

etcdctl3 snapshot restore ${MYBACKUPDIR}/var/lib/etcd/snapshot.db \
  --name $ETCD_NAME \
  --initial-cluster $ETCD_INITIAL_CLUSTER \
  --initial-cluster-token $ETCD_INITIAL_CLUSTER_TOKEN \
  --initial-advertise-peer-urls $ETCD_INITIAL_ADVERTISE_PEER_URLS \
  --data-dir /var/lib/etcd

chown -R etcd:etcd $ETCD_DATA_PATH
restorecon -Rv $ETCD_DATA_PATH
```

## Start docker/atomic-openshift-node
```
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

