Remove ETCD packages on all nodes to break etcd cluster
-------------------------------------------------------

To break all etcd perfectly, I choose package uninstallation.
This is the most critical situation. With this issue, OpenShift will be out of service.

## Video
[![ETCD Recovery All Members](http://img.youtube.com/vi/8mFdXII8NXU/0.jpg)](https://www.youtube.com/embed/8mFdXII8NXU)


## Remove all files under /var/lib/etcd on all ETCD nodes
```
systemctl stop etcd
rm -rf /var/lib/etcd/*
```

## Remove package on all ETCD nodes
```
yum remove etcd -y
```

## Check oc is working
```
oc get pod
```

## [Next](./recover_first_etcd.md)
