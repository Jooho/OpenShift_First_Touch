Remove ETCD packages on all nodes to break etcd cluster
-------------------------------------------------------

To break all etcd perfectly, I choose package uninstallation.
This is the most critical situation. With this issue, OpenShift will be out of service.


## Check ETCD Backup Data on each master
```
ls  /etc/snapshot.db
```



## Remove all files under /var/lib/etcd on all ETCD nodes
**On Bastion node**
```
export PRIVATE_KEY=~/.ssh/libra.pem

for etcd in $ETCD_IP_LIST
do

  echo “Delete ETCD Data folder on $etcd”
  ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no core@${etcd} -- sudo rm -rf /var/lib/etcd/*

done
```

## Check if oc is working 
   Try this command where oc command worked before
```
oc get pod
```

## [Next](./recover_all_etcd.md)
