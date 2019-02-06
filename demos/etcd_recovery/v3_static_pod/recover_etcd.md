Recover ETCD member & Synchronizing Data
-----------------------------------------

Now, the ETCD is broken so it is not starting. The first thing we should do is making it start up.

This doc explains how to make the etcd member start up even it will have different id. Then, it will show you how to put the recovered etcd member in the cluster and data synchronization.


## Remove broken etcd member from the cluster ##
```
# Check the broken etcd member health
etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers $etcd_members member list

# Remover the broker etcd member from the cluster
etcdctl3 --endpoints $etcd_members member  remove $(etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers $etcd_members member list|grep ${target_etcd} |awk -F: '{print $1}')
# Format
# etcdctl3 --endpoints $etcd_members member  remove $FAILED_ETCD_ID)

# Check if the broken etcd member is removed from the cluster
etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers $etcd_members member list
````

## Remove /var/lib/etcd/* for clean up ##
```
rm -rf $ETCD_DATA_PATH/*
```

## Change selinux in etcd data directory ##
```
chown -R etcd:etcd $ETCD_DATA_PATH
restorecon -Rv $ETCD_DATA_PATH
```

## Update etcd.conf like new etcd ##

```
vi /etc/etcd/etcd.conf

ETCD_INITIAL_CLUSTER="pvm-fusesource-patches.gsslab.rdu2.redhat.com=https://10.10.178.126:2380"
# change existing to new
ETCD_INITIAL_CLUSTER_STATE=new
```

## Add force new cluster option to systemd file ##
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

## Remove force new cluster option from systemd file ##
```
sed -i '/ExecStart/s/ --force-new-cluster//' /usr/lib/systemd/system/etcd.service
systemctl daemon-reload
systemctl show etcd.service --property ExecStart --no-pager
systemctl restart etcd
```

## Check if a new etcd start up properly without --force-new-cluster option ##
```
systemctl status etcd
```

## Stop etcd for joining to existing cluster
```
systemctl stop etcd
```

**NOTE: You must execute this command on a ETCD node where ETCD is running well in the cluster!!**
In this case, `vm125.gsslab.rdu2.redhat.com` is the ETCD member where this command will be execueted.

## Add the recovered ETCD member to the cluster again ##
```
export ETCD_CA_HOST="vm125.gsslab.rdu2.redhat.com"              # <===== Update
export NEW_ETCD="pvm-fusesource-patches.gsslab.rdu2.redhat.com" # <===== Update
export NEW_ETCD_IP="10.10.178.126"                              # <===== Update

etcdctl -C https://${ETCD_CA_HOST}:2379 \
  --ca-file=/etc/etcd/ca.crt     \
  --cert-file=/etc/etcd/peer.crt     \
  --key-file=/etc/etcd/peer.key member add ${NEW_ETCD} https://${NEW_ETCD_IP}:2380
```

**NOTE: You must execute this command on a recovered ETCD node!!**
Go back to `pvm-fusesource-patches.gsslab.rdu2.redhat.com` 

## Update etcd.conf with output after adding the recovered ETCD member to the cluster ##
### Specify all etcd nodes ###
```
vi /etc/etcd/etcd.conf

ETCD_INITIAL_CLUSTER="dhcp182-77.gsslab.rdu2.redhat.com=https://10.10.182.77:2380,pvm-fusesource-patches.gsslab.rdu2.redhat.com=https://10.10.178.126:2380,vm125.gsslab.rdu2.redhat.com=https://10.10.178.125:2380"
## from new to existing
ETCD_INITIAL_CLUSTER_STATE="existing"
```

## Delete member data & change owner for /var/lib/etcd ##
```
rm -rf /var/lib/etcd/member
chown -R etcd:etcd /var/lib/etcd
```

## Start ETCD and join to the cluster##
```
systemctl start etcd
```

## Check cluster health & data sync ##
```
etcdctl3 --endpoints $etcd_members endpoint health
etcdctl3 --endpoints $etcd_members endpoint status -w table
```

## [Next](../all_etcd_lost/break_etcd.md)

