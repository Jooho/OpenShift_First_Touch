Recover second ETCD member & Synchronizing data from the first ETCD member
---------------------------------------------------------------------------

Dislike the first ETCD member recovery, you don't need to use backup data for synchronization.

However, still you need the backup data for `/etc/etcd`. 

The way to recover the second ETCD member is similar to single ETCD lost situation.

*Every command should be executed on the target ETCD member node unless there is not NOTICE for node change*

[Variable information](./backup_v2.md)

## Target ETCD member ##
- pvm-fusesource-patches.gsslab.rdu2.redhat.com(10.10.178.126)

## Export target ETCD member hostname ##
```
export target_etcd=pvm-fusesource-patches.gsslab.rdu2.redhat.com
```

## Install ETCD
```
yum install -y etcd
```

## Restore Data ##
### Restore /etc/etcd ###
```
/bin/cp -Rf $MYBACKUPDIR/etc/etcd/* /etc/etcd/  
chown -R etcd:etcd /etc/etcd/*
```
## Make sure to clean /var/lib/etcd/* ##
```
rm -Rf /var/lib/etcd/*
```

## Make sure /var/lib/etcd exist ##
```
mkdir -p /var/lib/etcd
chown -R etcd:etcd $ETCD_DATA_PATH
restorecon -Rv $ETCD_DATA_PATH
```

## Edit etcd.conf like a new ETCD ###
```
vi /etc/etcd/etcd.conf
ETCD_INITIAL_CLUSTER="pvm-fusesource-patches.gsslab.rdu2.redhat.com=https://10.10.178.126:2380"
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
# a15c096752ae0cd4, started, pvm-fusesource-patches.gsslab.rdu2.redhat.com, https://10.10.178.126:2380, https://10.10.178.126:2379
  

# Update member info
etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers https://$(hostname):2379 member update  $(etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers https://$(hostname):2379 member list|awk -F: '{print $1}') https://$(dig +short $(hostname)):2380

# Recheck member info
etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers https://$(hostname):2379 member list
# ex)
# a15c096752ae0cd4, started, pvm-fusesource-patches.gsslab.rdu2.redhat.com, https://10.10.178.126:2380, https://10.10.178.126:2379
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
## Stop etcd for joining to existing cluster ##
```
systemctl stop etcd
```

**NOTE: You must execute this command on a ETCD node where ETCD is running well in the cluster!!**
In this case, `vm125.gsslab.rdu2.redhat.com` is the ETCD member where this command will be execueted.

## Add the recovered ETCD member to the cluster again ##
```
export ETCD_CA_HOST="vm125.gsslab.rdu2.redhat.com"                 #<==== update
export NEW_ETCD="pvm-fusesource-patches.gsslab.rdu2.redhat.com"    #<==== update 
export NEW_ETCD_IP="10.10.178.126"                                 #<==== update

etcdctl -C https://${ETCD_CA_HOST}:2379 \
  --ca-file=/etc/etcd/ca.crt     \
  --cert-file=/etc/etcd/peer.crt     \
  --key-file=/etc/etcd/peer.key member add ${NEW_ETCD} https://${NEW_ETCD_IP}:2380
```

**NOTE: You must execute this command on a recovered ETCD node!!**
Go back to `pvm-fusesource-patches.gsslab.rdu2.redhat.com`


## Update etcd.conf with output after adding the recovered ETCD member to the cluster ##
### Specify only running etcd members ###
```
ETCD_INITIAL_CLUSTER="vm125.gsslab.rdu2.redhat.com=https://10.10.178.125:2380,pvm-fusesource-patches.gsslab.rdu2.redhat.com=https://10.10.178.126:2380"
## from new to existing
ETCD_INITIAL_CLUSTER_STATE="existing"
```

## Delete member data & change owner for /var/lib/etcd to rejoin to the cluster ##
```
rm -rf /var/lib/etcd/member
chown -R etcd:etcd /var/lib/etcd
```


## Start ETCD and join to the cluster ###
```
systemctl start etcd
```

## Check cluster health & data sync ###
```
export etcd_members=https://10.10.178.126:2379,https://10.10.178.125:2379
etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers $etcd_members member list
etcdctl3 --endpoints $etcd_members endpoint health
etcdctl3 --endpoints $etcd_members endpoint status -w table
```

# The Second ETCD is recovered. Let's finish the last one # 

## [Next](./recover_third_etcd.md)
