Remove ETCD packages on all nodes to break etcd cluster
-------------------------------------------------------

To break all etcd perfectly, I choose package uninstallation.
This is the most critical situation. With this issue, OpenShift will be out of service.

[![ETCD Recovery Overview](http://img.youtube.com/vi/J6VYlY5PlsE/0.jpg)](https://www.youtube.com/embed/J6VYlY5PlsE)

## Check ETCD Backup Data
```
ls ${MYBACKUPDIR}

## export MYBACKUPDIR=/root/backup/etcd/$(date +%Y%m%d) 
```

## Remove all files under /var/lib/etcd on all ETCD nodes
```
rm -rf /var/lib/etcd/*
```

## Remove package on all ETCD nodes
```
yum remove etcd -y
```

## Check if oc is working
```
oc get pod
```

## [Next](./recover_all_etcd.md)
