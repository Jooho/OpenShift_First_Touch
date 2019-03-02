Back up ETCD v3 schema data
------------------------------

This doc explains how to back up etcd v3 schema data. 

[![ETCD Backup using Snapshot](http://img.youtube.com/vi/T1AzaG3NapA/0.jpg)](https://www.youtube.com/embed/T1AzaG3NapA)


## Gather etcd meber list

- Go to `one of masters`

- Copy the output of the following command
```
## etcd members (https:// :2379)
etcdctl3 --write-out=fields member list | awk '/ClientURL/{printf "%s%s",sep,$3; sep=","}'
# ex) export etcd_members=https://10.10.182.77:2379,https://10.10.178.126:2379,https://10.10.178.125:2379
```

## Backup Steps
- Go to `Ansible Controller`


## Create Backup Folders & export variables
```
export etcd_members=$PASTE_THE_OUTPUT       #UPDATE

# cluster-admin user
oc login 
oc project kube-system

export ETCD_DATA_PATH=/var/lib/etcd
export ETCD_POD_MANIFEST="/etc/origin/node/pods/etcd.yaml"
export MYBACKUPDIR=/root/backup/etcd/$(date +%Y%m%d)
export RUNNING_ETCD=$(oc get pod -n kube-system --no-headers | grep -o -m 1 '\S*etcd\S*' )

mkdir -p ${MYBACKUPDIR}/var/lib/etcd/member/snapshot
mkdir -p /etc/origin/node/pods-stopped



```

## Check ETCD health & Data Sync
```
oc exec ${RUNNING_ETCD} -c etcd -- /bin/bash -c "ETCDCTL_API=3 etcdctl \
    --cert /etc/etcd/peer.crt \
    --key /etc/etcd/peer.key \
    --cacert /etc/etcd/ca.crt \
    --endpoints $etcd_members endpoint health"


oc exec ${RUNNING_ETCD} -c etcd -- /bin/bash -c "ETCDCTL_API=3 etcdctl \
    --cert /etc/etcd/peer.crt \
    --key /etc/etcd/peer.key \
    --cacert /etc/etcd/ca.crt \
    --endpoints $etcd_members endpoint status -w table"
```

## Backup 
```
for etcd in $(oc get pod -n kube-system --no-headers|grep -v $(hostname) |grep -o  '\S*etcd\S*' );
do
  etcd_node=$(echo ${etcd}|sed 's/master-etcd-//g')
  echo "Backup Node: ${etcd_node}"

  oc exec ${etcd} -c etcd -- /bin/bash -c "ETCDCTL_API=3 etcdctl \
    --cert /etc/etcd/peer.crt \
    --key /etc/etcd/peer.key \
    --cacert /etc/etcd/ca.crt \
    --endpoints $etcd_members snapshot save /var/lib/etcd/snapshot.db"


  ssh ${etcd_node} -- mkdir -p ${MYBACKUPDIR}/var/lib/etcd/member/snap
  ssh ${etcd_node} -- mkdir -p /etc/origin/node/pods-stopped/
  ssh ${etcd_node} -- /bin/cp $ETCD_POD_MANIFEST   ${MYBACKUPDIR}/.
  ssh ${etcd_node} -- /bin/cp /var/lib/etcd/snapshot.db ${MYBACKUPDIR}/var/lib/etcd/. 
  ssh ${etcd_node} -- /bin/cp $ETCD_DATA_PATH/member/snap/db ${MYBACKUPDIR}/var/lib/etcd/member/snap/db
  echo "Backup Data==>"
  ssh ${etcd_node} -- ls ${MYBACKUPDIR}
  echo ""
done 
```


  #### Backup Data (Optional - V2/V3 copy way) 
  **Note: All etcd member should be backuped repectively**
  ```
  # Remove temp backup directory for new backup
  rm -rf  ${ETCD_DATA_PATH}_bak 
  mv /etc/origin/node/pods/etcd.yaml  /etc/origin/node/pods-stopped/.

  etcdctl2 backup \
  --data-dir $ETCD_DATA_PATH \
   --backup-dir ${ETCD_DATA_PATH}_bak

  mkdir -p ${ETCD_DATA_PATH}_bak/member/snap
  /bin/cp $ETCD_DATA_PATH/member/snap/db ${ETCD_DATA_PATH}_bak/member/snap/db
  /bin/cp -R ${ETCD_DATA_PATH}_bak  ${MYBACKUPDIR}/var/lib/.
  ```


*Tip*
Backup origin folder on main Master.
```
[masters]
a.example.com  <==== Main Master
b.example.com
c.example.com
```
Backup
```
cp -R /etc/origin ${MYBACKUPDIR}/etc/.
```



Backup is done.
