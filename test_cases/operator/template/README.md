# NFS Provioner Template

Refer [Original doc](https://github.com/Jooho/openshift-first-touch/blob/master/docs/storage/nfs/nfs-provisioner.md)

## Pre-requisites



- Set environmental variables & Git clone
  ```
  export ROOT_HOME=/tmp    #Update
  export REPO_HOME=${ROOT_HOME}/jhouse_openshift    
  export DEMO_HOME=${REPO_HOME}/test_cases/operator/template
  export UTIL_HOME=${DEMO_HOME}/utils    
  export TEST_HOME=${REPO_HOME}/test

  export NODE_NAME=ip-10-0-158-235.us-east-2.compute.internal #Update
  
  git clone https://github.com/Jooho/jhouse_openshift.git 

  ```

- Create a `/home/core/exports-nfs` folder on a node
  ~~~
  cd ${DEMO_HOME}/utils
  ./login_node.sh $NODE_NAME

  mkdir /home/core/exports-nfs
  
  exit
  ~~~


## Optional pre-requisites
- Attach a storage to a node
  - By default, it will try to use `/home/core/exports-nfs folder` as a NFS server volume point.
  
- Steps
  After you add volume to a node, execute the following commands:
  - Update params according to your environment
  ~~~
  sudo su -
  export disk_name=vdc
  export mount_folder=/home/core/exports-nfs
  ~~~

  - Create a LVM
  ~~~
  fdisk /dev/${disk_name} 
  # (type n p enter enter enter t 8e w enter)

  pvcreate /dev/${disk_name}1
  vgcreate nfs-vg /dev/${disk_name}1
  lvcreate -l 50%VG -n nfs-lv nfs-vg
  mkfs.xfs /dev/nfs-vg/nfs-lv

  mkdir -p ${mount_folder}
  echo "/dev/nfs-vg/nfs-lv ${mount_folder} xfs defaults 0 0" >> /etc/fstab

  mount -a
  chcon -Rt svirt_sandbox_file_t ${mount_folder}
  ls -alZ ${mount_folder}

  ~~~

## Test Flows
- Add node label to the node that volume attached to
  ~~~ 
  oc label node $NODE_NAME app=nfs-provisioner
  ~~~

- Update parameters of template.yaml
  ~~~
  ...
  name: NAMESPACE
  value: nfs-provisioner

  name: NFS_PATH
  value: "/export"

  name: HOST_PATH
  required: true
  value: "/home/core/exports-nfs"
  ~~~

## Deploy steps
  ~~~
  cd $DEMO_HOME/template
  oc new-project nfs-provisioner
  oc process -f  ./template.yaml|oc create -f -

  oc adm policy add-scc-to-user nfs-provisioner system:serviceaccount:nfs-provisioner:nfs-provisioner
  oc adm policy add-cluster-role-to-user nfs-provisioner-runner system:serviceaccount:nfs-provisioner:nfs-provisioner

  oc create -f storageclass.yaml
  ~~~

## Test Steps
  ~~~
  oc create -f ${DEMO_HOME}/test/test-pvc.yaml
  oc get pvc
  ~~~


## Clean up steps
  ~~~
  oc project nfs-provisioner
  oc delete po/test-pod
  oc delete pvc --all
  oc delete sa nfs-provisioner
  oc delete clusterrole nfs-provisioner-runner
  oc delete clusterrolebindings "run-nfs-provisioner"
  oc delete roles "leader-locking-nfs-provisioner"
  oc delete rolebindings "leader-locking-nfs-provisioner"

  oc delete scc nfs-provisioner
  oc delete sc example-nfs
  oc delete project nfs-provisioner 
  ~~~

## Troubleshooting
  - Mount point failed to write files
    ~~~
    Error setting up NFS server: error writing ganesha config /export/vfs.conf: open /export/vfs.conf: permission denied
    ~~~
    - Check the nfs-provisioner pod scc
      ~~~
      oc get pod nfs-XXX -o yaml|grep scc
      openshift.io/scc: nfs-provisioner
      ~~~ 
  
    - Check the selinux label of the folder
      ~~~
      oc exec nfs-provisioner-68f7b994d5-d4klq -- ls -alZ /export
      drwxrwxrwx. 4 root root system_u:object_r:container_file_t:s0         112 Sep 22 20:13 .
      ~~~
      - if not, execute this command
        ~~~
        chcon -Rvt svirt_sandbox_file_t /export
        ~~~


## Utils
- Deploy bastion pod to login to node 
  - [deploy.sh](./utils/deploy.sh)
  - This file can be download from internet.
  - `login_node.sh` check if the file exist
- SSH to a node(same as `oc debug node/$NODE` )
  - [login_node.sh](./utils/login_node.sh)
  - Usage `login_node.sh $NODE_NAME`
