Recover ETCD member & Synchronizing Data
-----------------------------------------

OCP 4 ETCD use `discovery-srv` so the cluster information is gathered by DNS automatically.

In this demo, you need to open 2 terminals
- T1: A node where can access to cluster(oc login)
- T2: The staled ETCD member node



```
oc rsh ${RUNNING_ETCD_POD} 

sh-4.2# dig +noall  +answer SRV _etcd-server-ssl._tcp.ocp4.jlee.rhcee.support
_etcd-server-ssl._tcp.ocp4.jlee.rhcee.support. 52 IN SRV 0 10 2380 etcd-2.ocp4.jlee.rhcee.support.
_etcd-server-ssl._tcp.ocp4.jlee.rhcee.support. 52 IN SRV 0 10 2380 etcd-0.ocp4.jlee.rhcee.support.
_etcd-server-ssl._tcp.ocp4.jlee.rhcee.support. 52 IN SRV 0 10 2380 etcd-1.ocp4.jlee.rhcee.support.
```
So we need to tweak the ETCD yaml file.



## Gather information
*T1 Terminal*
**Execute the following commands on anywhere except the stale ETCD member node.**
```
export KUBECONFIG=./aws/auth/kubeconfig

oc login -u kubeadmin -p XXX https://${API_SERVER}:6443
oc project kube-system

RUNNING_ETCD_POD=$(oc get pod -n kube-system --no-headers|grep Running| grep -o -m 1 '\S*etcd\S*')
RUNNING_ETCD_DNS_NAME=$(oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "cat /run/etcd/environment |grep ETCD_DNS|cut -d'=' -f2")
STALE_ETCD_POD=$(oc get pod -n kube-system  --no-headers|grep -v Running| awk '{print $1}')
STALE_ETCD_ID=$(oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt member list|grep ${STALE_ETCD_POD}|cut -d',' -f1")
STALE_ETCD_DNS_NAME=$(oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt member list|grep ${STALE_ETCD_POD}|cut -d',' -f4|cut -d/ -f3|cut -d: -f1")

# Check
echo RUNNING_ETCD_POD=${RUNNING_ETCD_POD}
echo RUNNING_ETCD_DNS_NAME=${RUNNING_ETCD_DNS_NAME}
echo STALE_ETCD_POD=${STALE_ETCD_POD}
echo STALE_ETCD_ID=${STALE_ETCD_ID}
echo STALE_ETCD_DNS_NAME=${STALE_ETCD_DNS_NAME}

```


## Stop ETCD Server
*T2 Terminal*

```
ssh core@etcd-member-ip-X-X-X-X.us-east-2.compute.internal

sudo su -

mkdir /etc/kubernetes/stopped-pod

mv /etc/kubernetes/manifests/etcd-member.yaml /etc/kubernetes/stopped-pod/.
```

## Clean etcd data
*T2 Terminal*
```
rm -rf /var/lib/etcd/*
```

## Remove the stale ETCD member from the cluster
*T1 Terminal*

```
oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt member remove ${STALE_ETCD_ID}"

oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt member list"
```

## Add the stale ETCD member to the cluster
*T1 Terminal*
```
oc exec ${RUNNING_ETCD_POD} -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt member add  ${STALE_ETCD_POD} --peer-urls=https://${STALE_ETCD_DNS_NAME}:2380"


ETCD_NAME="new_etcd-member-ip-X-X-X-X.us-east-2.compute.internal"
ETCD_INITIAL_CLUSTER="etcd-member-ip-Z-Z-Z-Z.us-east-2.compute.internal=https://etcd-2.ocp4.jlee.rhcee.support:2380,etcd-member-ip-Y-Y-Y-Y.us-east-2.compute.internal=https://etcd-1.ocp4.jlee.rhcee.support:2380,etcd-member-ip-X-X-X-X.us-east-2.compute.internal=https://etcd-0.ocp4.jlee.rhcee.support:2380"
ETCD_INITIAL_ADVERTISE_PEER_URLS="https://etcd-0.ocp4.jlee.rhcee.support:2380"
ETCD_INITIAL_CLUSTER_STATE="existing"
```


## Tweak ETCD yaml file
*T2 Terminal*

```

cp /etc/kubernetes/stopped-pod/etcd-member.yaml /etc/kubernetes/stopped-pod/etcd-member-recover.yaml

## Add the following information
vi /etc/kubernetes/stopped-pod/etcd-member-recover.yaml

93         --listen-peer-urls=https://0.0.0.0:2380 \
94         --initial-cluster-state=existing \            # <====== ADD
95         --initial-cluster-token=etcd-cluster \        # <====== ADD
```

## Start ETCD with the tweaked yaml file
*T2 Terminal*
```
mv /etc/kubernetes/stopped-pod/etcd-member-recover.yaml /etc/kubernetes/manifests/.
```

## Check ETCD member healths
*T1 Terminal*
```
# Member list
oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt member list"

ETCD_ENDPOINTS=$(oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt member list" | awk '{printf "%s%s",sep,$5; sep=","}')

# Member health
oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt --endpoints ${ETCD_ENDPOINTS} endpoint health" 

# ETCD data size
oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt --endpoints ${ETCD_ENDPOINTS} endpoint status --write-out table" 

```

