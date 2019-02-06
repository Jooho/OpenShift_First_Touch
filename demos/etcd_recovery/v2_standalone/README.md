Standalone ETCD with Data Schema Version v2/v3 
---------------------------------------------

This is tutorial for standalone ETCD with data schema v2/v3 on OpenShift.


## Test Environment ##
- OpenShift 3.6
- ETCD 3.2.22 with scheme v2/v3 data (No snapshot)
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
- [Backup ETCD Data](./backup_v2.md)

## Demo Scenarios ##
- **Single ETCD member lost**
  - [Break one ETCD member](./single_etcd_lost/break_etcd.md)
  - [Recover the ETCD member & Synchronizing data from other running ETCD member](./single_etcd_lost/recover_etcd.md)
- **All ETCD nodes lost**
  - [Remove ETCD package on all master nodes](./all_etcd_lost/break_etcd.md)
  - [Recover the first ETCD member & Restore data from backup](./all_etcd_lost/recover_first_etcd.md)
  - [Recover the second ETCD member & Synchronizing data from the first ETCD member](./all_etcd_lost/recover_second_etcd.md)
  - [Recover the third ETCD member & Synchronizing data from other running ETCD member](./all_etcd_lost/recover_third_etcd.md)


*Tip*
- [Creating ETCD certificates manually](./create_etcd_certs_manally.md)
- [Execute ETCD using cli](./execute_etcd_using_cli.md)
