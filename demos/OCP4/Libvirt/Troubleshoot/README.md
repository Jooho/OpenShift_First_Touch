# Troubleshooting 


## Worker is not created

### Check List

- Check Machine-controller log on master node
```
crictl ps|grep machine-controller

crictl logs ${machine-controller-id}
...
I0312 19:20:54.861763       1 client.go:332] Create a libvirt volume with name ocp4-pxcsb-worker-0-qrsrp for pool default from the base volume /var/lib/libvirt/images/default/ocp4-pxcsb-base
E0312 19:20:54.866011       1 actuator.go:106] Machine error: error creating volume storage volume 'ocp4-pxcsb-worker-0-qrsrp' already exists
E0312 19:20:54.866048       1 actuator.go:50] ocp4-pxcsb-worker-0-qrsrp: error creating libvirt machine: error creating volume storage volume 'ocp4-pxcsb-worker-0-qrsrp' already exists
I0312 19:20:54.866064       1 client.go:140] Closing libvirt connection: 0xc420524020
W0312 19:20:54.866093       1 controller.go:186] unable to create machine ocp4-pxcsb-worker-0-qrsrp: ocp4-pxcsb-worker-0-qrsrp: error creating libvirt machine: error creating volume storage volume 'ocp4-pxcsb-worker-0-qrsrp' already exists
...
```
You should check storage-pool `default` directory selinux/permission/disk path (`/var/lib/libvirt/images`).

Disk path of `default` storage pool is critical because installer assume the base image under this folder `/var/lib/libvirt/images`

```
virsh pool-list
virsh pool-dumpxml default
<pool type='dir'>
  <name>default</name>
  <uuid>5f39b6ea-e6cb-4524-945c-e3e93f9c0421</uuid>
  <capacity unit='bytes'>52576092160</capacity>
  <allocation unit='bytes'>40304312320</allocation>
  <available unit='bytes'>12271779840</available>
  <source>
  </source>
  <target>
    <path>/var/lib/libvirt/images</path>      <============= Check
    <permissions>
      <mode>0755</mode>
      <owner>1000</owner>
      <group>1000</group>
      <label>system_u:object_r:svirt_image_t:s0</label>
    </permissions>
  </target>
</pool>

```


## Wildecard DNS entry is not added

Console pod try to check oauth via router `https://openshift-authentication-openshift-authentication.apps.ocp4.tt.testing/oauth/token` and it failed.

Because of the following reasons:
- the router does not listen with host IP:PORT
- the wildcard dns entry is not added
- the domain `apps.ocp4.tt.testing` is resolved by libvirt dnsmasq always

Solution
- Add wildcard dns entry to `/etc/NetworkManager/dnsmasq.d/tt.testing.conf
```
server=/tt.testing/192.168.126.1
address=/.apps.tt.testing/192.168.126.51     <==== Add
```
- If you don't install cluster yet, you customize ingress 
```
./openshift-install create manifests --dir=./ocp4


cat ocp4/manifests/cluster-ingress-02-config.yml

apiVersion: config.openshift.io/v1
kind: Ingress
metadata:
  creationTimestamp: null
  name: cluster
spec:
  domain: apps.ocp4.tt.testing  <============>  apps.tt.testing   (Change)
status: {}
```
- If you already install the cluster, you need to update route in openshift-console
~~~
TBD
~~~


## Storage Pool selinux

- If you recreate the default folder, sometimes the selinux context is changed.

```
chcon -R system_u:object_r:svirt_image_t:s0 /var/lib/libvirt/images
restorecon -Rv  system_u:object_r:svirt_image_t:s0
```


