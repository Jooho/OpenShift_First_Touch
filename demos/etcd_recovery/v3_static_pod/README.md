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
|  pvm-fusesource-patches.gsslab.rdu2.redhat.com(10.10.182.77)|  | dhcp182-77.gsslab.rdu2.redhat.com(10.10.178.126)  | |  vm125.gsslab.rdu2.redhat.com(10.10.178.125) |   
+-------------------------------------------------------------+  +---------------------------------------------------+ +----------------------------------------------+
```

## Video
[![ETCD Recovery Overview](http://img.youtube.com/vi/JA0vJ5M4I60/0.jpg)](https://www.youtube.com/embed/JA0vJ5M4I60)

## Pre-tasks ##
- [Backup ETCD Data](./backup_v3.md)

## Demo Scenarios ##
- **Single ETCD member lost**
  - [Break one ETCD member](./single_etcd_lost/break_etcd.md)
  - [Recover the ETCD member & Synchronizing data from other running ETCD member](./single_etcd_lost/recover_etcd.md)
- **All ETCD nodes lost**
  - [Remove ETCD package on all master nodes](./all_etcd_lost/break_etcd.md)
  - [Recover all ETCD members](./all_etcd_lost/recover_all_etcd.md)



