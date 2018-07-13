ETCD Recovery of Bad Health Member
---------------------------------


Most relevant documentation about etcd recovery can be found in the [upstream etcd](https://coreos.com/etcd/) documentation.  Specifically, their [Disaster Recovery guide](https://coreos.com/etcd/docs/latest/op-guide/recovery.html)


The workflow for recovery.

1. Go to the node has bad healthy etcd.
2. Shut down atomic-openshift-*  services
3. Shut down etcd
4. Create a backup copy of /var/lib/etcd (this would be a cold backup)
5. Remove all contents of /var/lib/etcd
6. Edit /etc/etcd/etcd.conf and change ETCD_INITIAL_CLUSTER_STATE to existing
7. On one of the remaining nodes, see which member is unhealthy (which should be the one we just destroyed)
```

# etcdctl -C \
    https://${master1}:2379,https://${master2}:2379,https://${master3}:2379 \
    --ca-file=/etc/origin/master/master.etcd-ca.crt \
    --cert-file=/etc/origin/master/master.etcd-client.crt \
    --key-file=/etc/origin/master/master.etcd-client.key cluster-health

member 1eef7c801b8cd921 is healthy: got healthy result from https://172.31.0.121:2379
member 4531fddb8d7f87a0 is healthy: got healthy result from https://172.31.0.122:2379
member 81fd7db1c370aa0a is unhealthy
cluster is healthy

```
8. Delete the unhealthy member
```
# etcdctl -C \
    https://${master1}:2379,https://${master2}:2379,https://${master3}:2379 \
    --ca-file=/etc/origin/master/master.etcd-ca.crt \
    --cert-file=/etc/origin/master/master.etcd-client.crt \
    --key-file=/etc/origin/master/master.etcd-client.key member delete 81fd7db1c370aa0a
```

9. Add a new etcd memeber
```
# etcdctl -C \
    https://${master1}:2379,https://${master2}:2379,https://${master3}:2379 \
    --ca-file=/etc/origin/master/master.etcd-ca.crt \
    --cert-file=/etc/origin/master/master.etcd-client.crt \
    --key-file=/etc/origin/master/master.etcd-client.key member add <name(FQDN)> https://<IP>:2380

eg) etcdctl member add infra3 http://10.0.1.13:2380
```
