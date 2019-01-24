Executing ETCD with CLI
------------------------

Start etcd process without systemd
```
exec etcd --name "vm125.gsslab.rdu2.redhat.com"
--data-dir "/var/lib/etcd/" 
--initial-advertise-peer-urls "http://10.10.178.125:2380"
--listen-peer-urls "http://10.10.178.125:2380"
--listen-client-urls "http://10.10.178.125:2379"                      
--advertise-client-urls "http://10.10.178.125:2379"
--initial-cluster-token "etcd-cluster-1"
--initial-cluster "vm125.gsslab.rdu2.redhat.com=http://10.10.178.125:2380
--initial-cluster-state "new"
```
