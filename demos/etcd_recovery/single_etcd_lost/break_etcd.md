Break one of etcd members for demo
---------------------------------

For demonstration purpose, one of etcd members need to be broken.
This explain how to break it down.

*Every command should be executed on the target ETCD member node.*

##Target ETCD member##
- pvm-fusesource-patches.gsslab.rdu2.redhat.com (10.10.178.126)

##Export target ETCD member##
```
export target_etcd=pvm-fusesource-patches.gsslab.rdu2.redhat.com
```

##Remove all files under /var/lib/etcd##
```
systemctl stop etcd
rm -rf /var/lib/etcd/*
```

##Test ETCD is broken##
```
systemctl start etcd
```

**Errors: ETCD is not starting anymore**


##Check ETCD cluster lost the etcd node##
```
etcdctl3 --endpoints $etcd_members endpoint health
etcdctl3 --endpoints $etcd_members endpoint status -w table
```

*Just in case*
etcd_members is etcd_members (refer [ETCD Backup](../backup_v2.md))
