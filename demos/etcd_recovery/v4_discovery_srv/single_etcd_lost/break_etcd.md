Break one of etcd members for demo
---------------------------------

For demonstration purpose, one of etcd members need to be broken.
This explain how to break it down.

## Login to one of master nodes

~~~
ssh core@etcd-member-ip-X-X-X-X.us-east-2.compute.internal

sudo su -
~~~

## Remove all files under /var/lib/etcd ##
~~~
mkdir /etc/kubernetes/stopped-pod

mv /etc/kubernetes/manifests/etcd-member.yaml /etc/kubernetes/stopped-pod/etcd-member.yaml

rm -rf /var/lib/etcd/*
~~~

## Start ETCD
```
mv /etc/kubernetes/stopped-pod/etcd-member.yaml /etc/kubernetes/manifests/.
```

## Check if the ETCD member is not health ##

```
export KUBECONFIG=./aws/auth/kubeconfig

oc login -u kubeadmin -p XXX https://${API_SERVER}:6443
oc project kube-system

RUNNING_ETCD_POD=$(oc get pod -n kube-system --no-headers|grep Running| grep -o -m 1 '\S*etcd\S*')

RUNNING_ETCD_DNS_NAME=$(oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "cat /run/etcd/environment |grep ETCD_DNS|cut -d'=' -f2")

ETCD_ENDPOINTS=$(oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt member list" | awk '{printf "%s%s",sep,$5; sep=","}')

oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt --endpoints ${ETCD_ENDPOINTS} endpoint health"

```

## [Next](./recover_etcd.md)
