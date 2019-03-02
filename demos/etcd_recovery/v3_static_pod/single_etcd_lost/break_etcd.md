Break one of etcd members for demo
---------------------------------

For demonstration purpose, one of etcd members need to be broken.
This explain how to break it down.

*Every command should be executed on the target ETCD member node.*

[Variable information](./backup_v3.md)

## Target ETCD member ##
- dhcp181-165.gsslab.rdu2.redhat.com (10.10.181.165)

## Export target ETCD member ##
```
export target_etcd=dhcp181-165.gsslab.rdu2.redhat.com
```

## Remove all files under /var/lib/etcd ##
```
mv /etc/origin/node/pods/etcd.yaml /etc/origin/node/pods-stopped/
rm -rf /var/lib/etcd/*
```

## Check if the ETCD member is not health on VM125 ETCD member where ETCD work well##

```
export etcd_members=https://$(hostname):2379
etcdctl3 --endpoints $etcd_members endpoint health
```

## [Next](./recover_etcd.md)
