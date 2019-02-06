Back up ETCD v3 schema data
------------------------------

This doc explains how to back up etcd v3 schema data. 

## Target ETCD member ##
- vm125.gsslab.rdu2.redhat.com(10.10.178.125)

## Create Backup Folders & export variables
```
export ETCD_DATA_PATH=/var/lib/etcd
export ETCD_POD_MANIFEST="/etc/origin/node/pods/etcd.yaml"
export MYBACKUPDIR=/root/backup/etcd/$(date +%Y%m%d)

mkdir -p ${MYBACKUPDIR}/var/lib/etcd/member/snap
mkdir -p /etc/origin/node/pods-stopped

## etcd members (https:// :2379)
export etcd_members=$(etcdctl3 --write-out=fields member list | awk '/ClientURL/{printf "%s%s",sep,$3; sep=","}')
# ex) export etcd_members=https://10.10.182.77:2379,https://10.10.178.126:2379,https://10.10.178.125:2379

oc login -u system:admin
oc project kube-system

```

## Check ETCD health & Data Sync
```
source /etc/etcd/etcd.conf

etcdctl3 --endpoints $etcd_members endpoint health
etcdctl3 --endpoints $etcd_members endpoint status -w table
```

## Backup configuration
```
/bin/cp $ETCD_POD_MANIFEST ${MYBACKUPDIR}/.
```

## Backup Data (snapshot)
~~~
export ETCD_POD=$(oc get pod |grep $(hostname)|  grep -o -m 1 '\S*etcd\S*')

oc exec ${ETCD_POD} -c etcd -- /bin/bash -c "ETCDCTL_API=3 etcdctl \
    --cert /etc/etcd/peer.crt \
    --key /etc/etcd/peer.key \
    --cacert /etc/etcd/ca.crt \
    --endpoints $etcd_members snapshot save /var/lib/etcd/snapshot.db"
~~~

## Copy backup data
```
/bin/cp $ETCD_POD_MANIFEST   ${MYBACKUPDIR}/.
/bin/cp /var/lib/etcd/snapshot.db ${MYBACKUPDIR}/var/lib/etcd/. 

# optional
/bin/cp $ETCD_DATA_PATH/member/snap/db ${MYBACKUPDIR}/var/lib/etcd/member/snap/db

ls ${MYBACKUPDIR}/.
```

### Optional
  #### Backup Data for another etcd member (optional)
  ```

  export ETCD_POD=$(oc get pod --no-headers|grep -v $(hostname) |grep -o -m 1 '\S*etcd\S*' |head -n 1 )

  oc exec ${ETCD_POD} -c etcd -- /bin/bash -c "ETCDCTL_API=3 etcdctl \
    --cert /etc/etcd/peer.crt \
    --key /etc/etcd/peer.key \
    --cacert /etc/etcd/ca.crt \
    --endpoints $etcd_members snapshot save /var/lib/etcd/snapshot.db"

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
