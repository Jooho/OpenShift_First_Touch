ETCD Recovery 
-------------

This is tutorial for ETCD recovery on OpenShift.

There are many documentations and blogs regarding this. This doc try to show you practical commands from A to Z under specific scenarios.

During demonstration, it uses v2 or v3 cli version. It does not matter if you are using one specific version way. I just want to show you each version cli is working fine.

[Standalone ETCD with Data Schema Version 2/3](./v2_standalone/README.md)
- Test Environment 
  - OpenShift 3.6
  - ETCD 3.2.22 with scheme v2/v3 data (**No snapshot**)
  - 3 ETCD Nodes

[Static Pod ETCD with Data Schema Version 3](./v3_static_pod/README.md)
- Test Environment 
  - OpenShift 3.11
  - ETCD 3.2.22 with scheme v3 data (**Snapshot way**)
  - 3 ETCD Nodes

