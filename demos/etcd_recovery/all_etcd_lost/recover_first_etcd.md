Recover the first ETCD member
-----------------------------

The way to recover the first ETCD member is similar to creating a new ETCD member.

The first ETCD will restore data from backup data and this task is the first thing you should do.

*Every command should be executed on the target ETCD member node unless there is not NOTICE for node change*

[Variable information](../backup_v2.md)

## Target ETCD member ##
- vm125.gsslab.rdu2.redhat.com(10.10.178.125)

## Export target ETCD member hostname ##
```
export target_etcd=vm125.gsslab.rdu2.redhat.com
```

## Install ETCD ##
```
yum install -y etcd
```

## Restore Data ##
FYI, ETCD_DATA_PATH is `/var/lib/etcd`

### Restore /var/lib/etcd ###
```
rm -rf $ETCD_DATA_PATH
cp -Rp  $MYBACKUPDIR/var/lib/etcd_bak $ETCD_DATA_PATH
chcon -R --reference $MYBACKUPDIR/var/lib/etcd_bak $ETCD_DATA_PATH
chown -R etcd:etcd $ETCD_DATA_PATH
```

FYI, MYBACKUPDIR is `/root/backup/etcd/$(date +%Y%m%d)`

### Restore /etc/etcd ###
```
/bin/cp -Rf $MYBACKUPDIR/etc/etcd/* /etc/etcd/  
chown -R etcd:etcd /etc/etcd/*
```

## Edit etcd.conf like a new ETCD member ##
```
ETCD_INITIAL_CLUSTER=vm125.gsslab.rdu2.redhat.com=https://10.10.178.125:2380  
ETCD_INITIAL_CLUSTER_STATE=new   
```

## ADD force new cluster option to systemd file ##

```
sed -i '/ExecStart/s/"$/  --force-new-cluster"/' /usr/lib/systemd/system/etcd.service
systemctl daemon-reload
systemctl show etcd.service --property ExecStart --no-pager
systemctl restart etcd
```

## Check if a new etcd start up ##
```
systemctl status etcd
```

## If peerURL does not contain localhost, you can skip this steps ##
### Update peerURL from localhost to ip ###

```
# Check member info
etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers https://$(hostname):2379 member list
# ex)
# 687c3648cc01: name=vm125.gsslab.rdu2.redhat.com peerURLs=http://localhost:2380 clientURLs=https://10.10.178.125:2379 isLeader=true   

# Update member info
etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers https://$(hostname):2379 member update  $(etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers https://$(hostname):2379 member list|awk -F: '{print $1}') https://$(dig +short $(hostname)):2380

# Recheck member info
etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers https://$(hostname):2379 member list
# ex)
# 687c3648cc01: name=vm125.gsslab.rdu2.redhat.com peerURLs=https://10.10.178.125:2380 clientURLs=https://10.10.178.125:2379 isLeader=true 
```

## Remove force new cluster option from systemd file ##
```
sed -i '/ExecStart/s/ --force-new-cluster//' /usr/lib/systemd/system/etcd.service
systemctl daemon-reload
systemctl show etcd.service --property ExecStart --no-pager
systemctl restart etcd
```
## Check if a new etcd start up ##
```
systemctl status etcd
```

## Check if backup data is recovered ##
```
etcdctl3 --endpoints $etcd_members endpoint status -w table
```

## Try oc command ##
```
oc get pod
```

# WoW, the first ETCD is recovered, now move on the next etcd node #

## [Next](./recover_second_etcd.md)
