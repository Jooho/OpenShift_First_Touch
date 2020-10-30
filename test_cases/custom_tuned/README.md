# Create a custom tuned

This example show how to apply a custom tuned for infra node

## Pre-requisites
### Create a MC for infra
~~~
oc create -f infra.mc.yaml
~~~

### Create a infra MCP
~~~
oc create -f infra.mcp.yaml
~~~

### Add infra label and remove worker label from a node
~~~
oc label node $NODE_NAME  node-role.kubernetes.io/infra=
oc label node $NODE_NAME  node-role.kubernetes.io/worker-
~~~

Wait for appling a new Rendered MC to the infra node.

### Create a new Tuned for the infra node
~~~
oc create -f custom_tuned_hugepage.yaml
~~~


### Check if the new tuned is applied or not
~~~
oc project openshift-cluster-node-tuning-operator

oc get pod|grep operator
oc logs cluster-node-tuning-operator-5f7d7c74b9-ccw6x     # tuned operator logs
...
I1030 13:27:05.541241       1 controller.go:406] updated Tuned rendered
I1030 13:35:39.907145       1 controller.go:517] updated profile ip-10-0-220-155.us-east-2.compute.internal [openshift-hugepages]


oc get pod -o wide |grep $NODE_NAME
oc logs $tuned_pod
...
2020-10-30 13:35:40,777 INFO     tuned.plugins.base: instance disk: assigning devices dm-0, xvda
2020-10-30 13:35:40,777 INFO     tuned.plugins.base: instance net: assigning devices ens3
2020-10-30 13:35:40,781 INFO     tuned.plugins.plugin_sysctl: reapplying system sysctl
2020-10-30 13:35:40,787 INFO     tuned.daemon.daemon: static tuning from profile 'openshift-hugepages' applied

oc debug node/$NODE_NAME
$ cat /sys/kernel/mm/transparent_hugepage/enabled
always madvise [never]                #<=== changed. it was always

$ sysctl -a|grep vm.nr_hugepages
vm.nr_hugepages = 16                   #<=== changed. it was 0
vm.nr_hugepages_mempolicy = 16         #<=== changed. it was 0
~~~


