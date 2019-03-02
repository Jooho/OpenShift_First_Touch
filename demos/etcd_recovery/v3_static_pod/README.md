Static Pod ETCD with Data Schema Version 3
------------------------------------------

This is tutorial for static pod ETCD with data schema version 3 using snapshot.

Static Pod ETCD with Data Schema Version 3

## Test Environment ##
- OpenShift 3.11
- ETCD 3.2.22 with data scheme v3 data (**Using snapshot**)
- 3 ETCD Nodes

### Diagram ###
```
+-------------------------------------------------------------+  +---------------------------------------------------+ + ---------------------------------------------+
|                          ETCD1                              |  |                      ETCD2                        | |                   ETCD3                      |
|  dhcp181-165.gsslab.rdu2.redhat.com(10.10.181.165)|  | dhcp179-170.gsslab.rdu2.redhat.com(10.10.179.170)  | |  vm49.gsslab.rdu2.redhat.com(10.10.178.49) |   
+-------------------------------------------------------------+  +---------------------------------------------------+ +----------------------------------------------+
```

## Video
[![ETCD Recovery Overview](http://img.youtube.com/vi/JA0vJ5M4I60/0.jpg)](https://www.youtube.com/embed/JA0vJ5M4I60)

This video is from v2 but most of parts are the similar

## Pre-tasks ##
- [Backup ETCD Data](./backup_v3.md)

## Demo Scenarios ##
- **Single ETCD member lost**
  - [Break one ETCD member](./single_etcd_lost/break_etcd.md)
  - [Recover the ETCD member & Synchronizing data from other running ETCD member](./single_etcd_lost/recover_etcd.md)
- **All ETCD nodes lost**
  - [Remove ETCD package on all master nodes](./all_etcd_lost/break_etcd.md)
  - [Recover all ETCD members](./all_etcd_lost/recover_all_etcd.md)



