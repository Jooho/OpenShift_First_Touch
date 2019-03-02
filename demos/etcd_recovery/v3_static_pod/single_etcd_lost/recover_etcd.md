Recover ETCD member & Synchronizing Data
-----------------------------------------

Now, the ETCD is broken so it is not starting. The first thing we should do is making it start up.

This doc explains how to make the etcd member start up even it will have different id. Then, it will show you how to put the recovered etcd member in the cluster and data synchronization.

[![Static POD Single ETCD Recovery](http://img.youtube.com/vi/1-1WzcvDRLY/0.jpg)](https://www.youtube.com/embed/1-1WzcvDRLY)


**Note: Execute commands on vm49 node**

## Export target ETCD member ##
```
export target_etcd=dhcp181-165.gsslab.rdu2.redhat.com   #UPDATE
```

## Remove broken etcd member from the cluster ##
```
export etcd_members=$(etcdctl3 --write-out=fields member list | awk '/ClientURL/{printf "%s%s",sep,$3; sep=","}')

# Check the broken etcd member health
etcdctl3 --endpoints $etcd_members member list

# Remover the broker etcd member from the cluster
etcdctl3 --endpoints $etcd_members member  remove $(etcdctl3 --endpoints $etcd_members member list|grep ${target_etcd} |awk -F, '{print $1}')
# Format
# etcdctl3 --endpoints $etcd_members member  remove $FAILED_ETCD_ID)

# Check if the broken etcd member is removed from the cluster
etcdctl3 --endpoints $etcd_members member list
````

**Note: Execute commands on target etcd node(dhcp181-165.gsslab.rdu2.redhat.com) node**

## Stop ETCD 
```
mv /etc/origin/node/pods/etcd.yaml /etc/origin/node/pods-stopped/
```

## Remove /var/lib/etcd/* for clean up
```
export ETCD_DATA_PATH=/var/lib/etcd
rm -rf $ETCD_DATA_PATH/*
```

## Change selinux in etcd data directory
```
chown -R etcd:etcd $ETCD_DATA_PATH
restorecon -Rv $ETCD_DATA_PATH
```

## Update etcd.conf like new etcd

```
vi /etc/etcd/etcd.conf

ETCD_INITIAL_CLUSTER="dhcp181-165.gsslab.rdu2.redhat.com=https://10.10.181.165:2380"   #UPDATE
# change existing to new
ETCD_INITIAL_CLUSTER_STATE=new
```

## Add force new cluster option to systemd file
```
/bin/cp /etc/etcd/etcd.conf /etc/etcd/etcd.conf.bak
echo "ETCD_FORCE_NEW_CLUSTER=true" >> /etc/etcd/etcd.conf
mv /etc/origin/node/pods-stopped/etcd.yaml /etc/origin/node/pods/.
```
## Check if a new etcd start up
```
export etcd_members="https://10.10.181.165:2379"  #UPDATE
etcdctl3 --endpoints $etcd_members endpoint health

# Check logs
docker logs $(docker ps|grep etcd|grep -v pod|awk '{print $1}')
```

## Remove force new cluster option from systemd file
```
mv /etc/origin/node/pods/etcd.yaml  /etc/origin/node/pods-stopped/.
rm  /etc/etcd/etcd.conf
/bin/mv /etc/etcd/etcd.conf.bak /etc/etcd/etcd.conf
mv /etc/origin/node/pods-stopped/etcd.yaml /etc/origin/node/pods/.
```

## Check if a new etcd start up properly without --force-new-cluster option
```
etcdctl3 --endpoints $etcd_members endpoint health

# Check logs
docker logs $(docker ps|grep etcd|grep -v pod|awk '{print $1}')
```

## Stop etcd for joining to existing cluster
```
mv /etc/origin/node/pods/etcd.yaml /etc/origin/node/pods-stopped/
```

**NOTE: Go to one of ETCD node `vm49`**

## Add the recovered ETCD member to the cluster again
```
export RUNNING_ETCD=$(oc get pod -n kube-system --no-headers | grep $(hostname) |grep -o -m 1 '\S*etcd\S*' )
export ETCD_CA_HOST="$(echo ${RUNNING_ETCD}|sed 's/master-etcd-//g')" 
export NEW_ETCD="dhcp181-165.gsslab.rdu2.redhat.com"            # <===== Update
export NEW_ETCD_IP="10.10.181.165"                              # <===== Update

export ETCD_EP=$(dig +short $ETCD_CA_HOST)

# If you don't have etcdctl cli
oc exec ${RUNNING_ETCD} -c etcd -- /bin/bash -c "ETCDCTL_API=3 etcdctl \
    --cert /etc/etcd/peer.crt \
    --key /etc/etcd/peer.key \
    --cacert /etc/etcd/ca.crt \
    --endpoints ${ETCD_EP}:2379 member add  ${NEW_ETCD} --peer-urls=https://${NEW_ETCD_IP}:2380"

# If you have etcdctl cli
etcdctl3 --endpoints=https://${ETCD_CA_HOST}:2379 member add  ${NEW_ETCD} --peer-urls="https://${NEW_ETCD_IP}:2380" 
```

**NOTE: You must execute this command on a recovered ETCD node!!**
Go back to `dhcp181-165.gsslab.rdu2.redhat.com` 

## Update etcd.conf with output after adding the recovered ETCD member to the cluster
### Specify all etcd nodes
```
vi /etc/etcd/etcd.conf

ETCD_INITIAL_CLUSTER="vm49.gsslab.rdu2.redhat.com=https://10.10.178.49:2380,dhcp181-165.gsslab.rdu2.redhat.com=https://10.10.181.165:2380,dhcp179-170.gsslab.rdu2.redhat.com=https://10.10.179.170:2380"
## from new to existing
ETCD_INITIAL_CLUSTER_STATE="existing"
```

## Delete member data & change owner for /var/lib/etcd
```
rm -rf /var/lib/etcd/member
chown -R etcd:etcd /var/lib/etcd
```

## Start ETCD and join to the cluster
```
mv /etc/origin/node/pods-stopped/etcd.yaml /etc/origin/node/pods/.
```

## Check cluster health & data sync
```
export etcd_members=$(etcdctl3 --write-out=fields member list | awk '/ClientURL/{printf "%s%s",sep,$3; sep=","}')
etcdctl3 --endpoints $etcd_members endpoint health
etcdctl3 --endpoints $etcd_members endpoint status -w table
```

## [Next](../all_etcd_lost/break_etcd.md)

