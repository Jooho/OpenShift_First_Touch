# OpenShift v4 ETCD Recovery

This is tutorial for static pod ETCD using discovery-srv 

Static Pod ETCD with Data Schema Version 3

## Test Environment ##
- OpenShift 4.0
- ETCD 3.3.10
- 3 ETCD Nodes

### Diagram ###
```
+-------------------------------------------------------------+  +---------------------------------------------------+ + ---------------------------------------------+
|                          ETCD1                              |  |                      ETCD2                        | |                   ETCD3                      |
|  etcd-member-ip-X-X-X-X.us-east-2.compute.internal|  | etcd-member-ip-Y-Y-Y-Y.us-east-2.compute.internal  | |  etcd-member-ip-Z-Z-Z-Z.us-east-2.compute.internal |   
+-------------------------------------------------------------+  +---------------------------------------------------+ +----------------------------------------------+
```

## Video
[![ETCD Recovery Overview](http://img.youtube.com/vi/JA0vJ5M4I60/0.jpg)](https://www.youtube.com/embed/JA0vJ5M4I60)

This video is from v2 but most of parts are the similar

## Pre-tasks ##
- [Backup ETCD Data](./backup_v4.md)

## Demo Scenarios ##
- **Single ETCD member lost**
  - [Break one ETCD member](./single_etcd_lost/break_etcd.md)
  - [Recover the ETCD member & Synchronizing data from other running ETCD member](./single_etcd_lost/recover_etcd.md)
- **All ETCD nodes lost**
  - [Remove ETCD package on all master nodes](./all_etcd_lost/break_etcd.md)
  - [Recover all ETCD members](./all_etcd_lost/recover_all_etcd.md)

