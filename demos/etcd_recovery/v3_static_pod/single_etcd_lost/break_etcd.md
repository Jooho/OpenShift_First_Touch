Break one of etcd members for demo
---------------------------------

For demonstration purpose, one of etcd members need to be broken.
This explain how to break it down.

*Every command should be executed on the target ETCD member node.*

[Variable information](./backup_v3.md)

## Target ETCD member ##
- pvm-fusesource-patches.gsslab.rdu2.redhat.com (10.10.178.126)

## Export target ETCD member ##
```
export target_etcd=pvm-fusesource-patches.gsslab.rdu2.redhat.com
```

## Remove all files under /var/lib/etcd ##
```
mv /etc/origin/node/pods/etcd.yaml /etc/origin/node/pods-stopped/
rm -rf /var/lib/etcd/*
```

## Check if the ETCD member is not health on VM125 ETCD member where ETCD work well##
FYI, `etcd_members` is `https://10.10.178.125:2379,https://10.10.182.77:2379,https://10.10.178.126:2379`
```
export etcd_members=https://10.10.178.125:2379,https://10.10.182.77:2379,https://10.10.178.126:2379
etcdctl3 --endpoints $etcd_members endpoint health
```

## [Next](./recover_etcd.md)
