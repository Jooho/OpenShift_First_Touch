Recover the all ETCD members at one time
----------------------------------------

## Gather ETCD_IP_LIST
Suppose all etcd nodes are down and openshift cluster is also down. Hence, you need to get etcd ip list by AWS cli or AWS console because oc cli is not working.

**On the node** that has *openshift cluster metadata.json* file

```
export INFRA_ID=$(cat ./metadata.json |jq -r .infraID)

export ETCD_IP_LIST=""

for i in {0..2}     # If you have more than 3 masters, you should change the “2” 
do 
  etcd_ip=$(aws ec2 describe-instances --filters "Name=tag:Name,Values=${INFRA_ID}-master-${i}"  --query 'Reservations[*].Instances[*].PrivateIpAddress')
  ETCD_IP_LIST="${ETCD_IP_LIST} $etcd_ip"
done

# Copy this output
echo ${ETCD_IP_LIST}    


##Example output
## 10.0.142.250 10.0.150.243 10.0.163.223
```
*You should copy the output IP list and will be used in the next script*



## Gather ETCD_INITIAL_CLUSTER
Execute this script on the **bastion node**
You MUST change `export ETCD_IP_LIST="$CHANGE_ME"` that is in the script
```
export PRIVATE_KEY=~/.ssh/id_rsa.pem

export ETCD_IP_LIST="$CHANGE_ME"  (From above output)
export ETCD_INITIAL_CLUSTER=""
for etcd in $ETCD_IP_LIST
do
  etcd_name=$(ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no core@${etcd} -- cat /run/etcd/environment |grep NAME|awk -F= '{print $2}'|cut -d'.' -f1|cut -d'-' -f2)
  etcd_ip=$(ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no core@${etcd} -- cat /run/etcd/environment |grep ADDRESS|awk -F= '{print $2}')

  if [[ z${ETCD_INITIAL_CLUSTER} != z ]];then
    ETCD_INITIAL_CLUSTER=("${ETCD_INITIAL_CLUSTER},${etcd_name}=https://${etcd_ip}:2380")
  else
    ETCD_INITIAL_CLUSTER=("${etcd_name}=https://${etcd_ip}:2380")
  fi
done

echo ${ETCD_INITIAL_CLUSTER}

```

## Stop all ETCD and delete ETCD data folder
Execute this script on the **bastion node**
```
export PRIVATE_KEY=~/.ssh/libra.pem

for etcd in $ETCD_IP_LIST
do
  echo "Create Stopped pods folder(/etc/kubernetes/stopped-pods) on $etcd"
  ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no core@${etcd} -- sudo mkdir /etc/kubernetes/stopped-pods

  echo "Stop ETCD on $etcd"
  ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no core@${etcd} -- sudo mv /etc/kubernetes/manifests/etcd-member.yaml /etc/kubernetes/stopped-pods/.

  echo "Clean ETCD data folder(/var/lib/etcd) on $etcd"
  ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no core@${etcd} -- sudo rm -rf /var/lib/etcd

  echo "Check ETCD data folder(/var/lib/etcd) is empty on $etcd"
  ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no core@${etcd} -- sudo ls -al /var/lib/etcd

done
```

## Restore ETCD data
Execute this script on the **bastion node**
```
cat <<EOF> recover_etcd.sh
curl -L https://github.com/etcd-io/etcd/releases/download/v3.3.12/etcd-v3.3.12-linux-amd64.tar.gz|tar xvz

## (TBD) enhance curl command to specify etcd home
cd ./etcd-v3.3.12-linux-amd64

source  /run/etcd/environment
export \$(cut -d= -f1 /run/etcd/environment)

## (TBD) need to check if the files exist. (looks it is a bug)
mkdir /etc/ssl/etcd
ln -s /etc/kubernetes/static-pod-resources/etcd-member/ca.crt /etc/ssl/etcd/.
ln -s /etc/kubernetes/static-pod-resources/etcd-member/system* /etc/ssl/etcd/.



export ETCD_NAME=\$(cat /run/etcd/environment |grep NAME|awk -F= '{print $2}'|cut -d'.' -f1|cut -d'-' -f2)
export ETCD_INITIAL_CLUSTER="${ETCD_INITIAL_CLUSTER}"
export ETCD_INITIAL_CLUSTER_TOKEN=etcd-cluster-1 
export ETCD_INITIAL_ADVERTISE_PEER_URLS=https://\${ETCD_IPV4_ADDRESS}:2380

sudo ETCDCTL_API=3 ./etcdctl snapshot restore /etc/snapshot.db \
  --name \${ETCD_NAME} \
  --initial-cluster \${ETCD_INITIAL_CLUSTER} \
  --initial-cluster-token \${ETCD_INITIAL_CLUSTER_TOKEN} \
  --initial-advertise-peer-urls \${ETCD_INITIAL_ADVERTISE_PEER_URLS} \
  --data-dir /var/lib/etcd 

sudo  restorecon -Rv /var/lib/etcd
EOF


for etcd in $ETCD_IP_LIST
do
 echo "Copy script(recover_etcd.sh) to $etcd"
 scp -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no ./recover_etcd.sh core@${etcd}:/home/core

 echo "Change script file permission on $etcd"
 ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no core@${etcd} -- chmod 777 /home/core/recover.etcd.sh

 echo "Execute recover_etcd.sh on $etcd"
 ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no core@${etcd} -- /home/core/recover_etcd.sh

 echo "Change selinux context for ETCD data folder on $etcd"
 ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no core@${etcd} -- sudo  restorecon -Rv /var/lib/etcd

done
```
## Restart CRIO and start ETCD pods
```
for etcd in $ETCD_IP_LIST
do
 
  echo "Restart crio on $etcd"
  ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no core@${etcd} -- sudo  systemctl restart crio

  echo "Start ETCD on $etcd"
  ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no core@${etcd} -- sudo  mv /etc/kubernetes/stopped-pods/etcd-member.yaml /etc/kubernetes/manifests/.
done
```
## Check ETCD 

Execuete this script on where can be accessible to the OpenShift cluster.
```
# Login OpenShift Cluster
## If it is possible to login, you can assume ETCD nodes are recovered

oc login -u kubeadmin -p xxxxx

# Gather data to check ETCD health
 RUNNING_ETCD_POD=$(oc get pod -n kube-system -l k8s-app=etcd -o  jsonpath='{range .items[?(@.status.phase=="Running")]}{.metadata.name}{"\n"}{end}' |head -n 1)
RUNNING_ETCD_DNS_NAME=$(oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "cat /run/etcd/environment |grep ETCD_DNS|cut -d'=' -f2")


# Member list
oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt member list"

# Endpoint IP list
ETCD_ENDPOINTS=$(oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt member list" | awk '{printf "%s%s",sep,$5; sep=","}')

# Endpoints health
oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt --endpoints ${ETCD_ENDPOINTS} endpoint health"

# ETCD data size
oc exec ${RUNNING_ETCD_POD}  -n kube-system  -- /bin/sh -c "ETCDCTL_API=3 etcdctl --cert /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.crt --key /etc/ssl/etcd/system:etcd-peer:${RUNNING_ETCD_DNS_NAME}.key --cacert /etc/ssl/etcd/ca.crt --endpoints ${ETCD_ENDPOINTS} endpoint status --write-out table"

```