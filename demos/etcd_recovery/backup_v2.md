Back up ETCD v2 schema data
------------------------------

This doc explains how to back up etcd v2 schema data. For schema v3 data, you can use snapshot (refer - TBD)


## Video
[![ETCD Data Backup](http://img.youtube.com/vi/TkkT92Y4nho/0.jpg)](https://www.youtube.com/embed/TkkT92Y4nho)


## Create Backup Folders & export variables
```
export ETCD_DATA_PATH=/var/lib/etcd
export MYBACKUPDIR=/root/backup/etcd/$(date +%Y%m%d)
mkdir -p ${MYBACKUPDIR}/etc
mkdir -p ${MYBACKUPDIR}/var/lib

## etcd members (https:// :2379)
export etcd_members=https://10.10.182.77:2379,https://10.10.178.126:2379,https://10.10.178.125:2379
```


## Check ETCD health & Data Sync
```
source /etc/etcd/etcd.conf

etcdctl3 --endpoints $etcd_members endpoint health
etcdctl3 --endpoints $etcd_members endpoint status -w table
```

## Backup Data (One by One) 
**Note: All etcd member should be backuped repectively**
```
# Remove temp backup directory for new backup
rm -rf  ${ETCD_DATA_PATH}_bak 

systemctl stop etcd

etcdctl backup \
  --data-dir $ETCD_DATA_PATH \
   --backup-dir ${ETCD_DATA_PATH}_bak

cp $ETCD_DATA_PATH/member/snap/db ${ETCD_DATA_PATH}_bak/member/snap/db
cp -R /etc/etcd   ${MYBACKUPDIR}/etc/.
cp -R ${ETCD_DATA_PATH}_bak  ${MYBACKUPDIR}/var/lib/.

ls ${MYBACKUPDIR}/.
```

*Tip*
You should backup main master folder as well just in case, you lost the main master.
Main master is the master that is the first node in the ansible hosts file.
```
[masters]
a.example.com  <==== Main Master
b.example.com
c.example.com
```
Backup
```
cp -R /etc/origin/master ${MYBACKUPDIR}/etc/.
```



Backup is done.
