Break one of etcd members for demo
---------------------------------

For demonstration purpose, one of etcd members need to be broken.
This explain how to break it down.

## Video
[![ETCD Recovery Single Member](http://img.youtube.com/vi/yx3Pjpkasl4/0.jpg)](https://www.youtube.com/embed/yx3Pjpkasl4)


*Every command should be executed on the target ETCD member node.*

[Variable information](../backup_v2.md)

## Target ETCD member ##
- pvm-fusesource-patches.gsslab.rdu2.redhat.com (10.10.178.126)

## Export target ETCD member ##
```
export target_etcd=pvm-fusesource-patches.gsslab.rdu2.redhat.com
```

## Remove all files under /var/lib/etcd ##
```
systemctl stop etcd
rm -rf /var/lib/etcd/*
```

## Test ETCD is broken ##
```
systemctl start etcd
```

**Errors: ETCD is not starting anymore**

## Check ETCD cluster lost the etcd node ##
FYI, `etcd_members` is `https://10.10.178.126:2379,https://10.10.178.125:2379,https://10.10.182.77:2379`
```
etcdctl3 --endpoints $etcd_members endpoint health
etcdctl3 --endpoints $etcd_members endpoint status -w table
```

## [Next](./recover_etcd.md)
