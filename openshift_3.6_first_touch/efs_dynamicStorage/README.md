# EFS (Amazon Elastic File System) Dynamic Storage

## Description
AWS provide [EFS service](https://aws.amazon.com/efs/?nc1=h_ls) that is scalable file storage that could be sharable. OpenShift support `Deploying External Persistent Volume Provisioners`(An external provisioner is an application that enables dynamic provisioning for a particular storage provider). **Note: EFS can only attach to EC2 Instance in AWS** 

## Demo Scenario
This show the way to deploy external pv provisioner pod and storageclass. Moreover, it will deploy sample application to assign the efs pvc.

### Demo Env
OCP version: 3.6
CNS installed: registry, app 
Architecture:
  - 1 x Master
  - 3 x Infra (glusterfs)
  - 3 x App   (glusterfs)

### Demo Steps
- Deploy external provioner pod
- Create StorageClass
- Create Sample PVC
- Deploy test application.
  For persistent volume, storageClass name will be used.

Note: If python version is 3, ansible script might complain with "iteritems()" 
      Refer [this doc](http://docs.ansible.com/ansible/latest/playbooks_python_version.html)

*1. Deploy External Provionser*
```
ansible-playbook  -vvv   /usr/share/ansible/openshift-ansible/playbooks/byo/openshift-cluster/openshift-provisioners.yml \
-e openshift_provisioners_install_provisioners=True \
-e openshift_provisioners_efs=True \
-e openshift_provisioners_efs_fsid=fs-%UPDATE% \     
-e openshift_provisioners_efs_region=%UPDATE% \
-e openshift_provisioners_efs_aws_access_key_id=%UPDATE% \
-e openshift_provisioners_efs_aws_secret_access_key=%UPDATE% \
-e openshift_provisioners_project=openshif-infra \
-e '{openshift_provisioners_efs_nodeselector: {role: "infra"}}'
-e openshift_provisioners_image_prefix=openshift3/ \
-e openshift_provisioners_image_version=v3.6 \
-e openshift_provisioners_efs_path=/
```

*2. Create Storage Class*

```
echo "
kind: StorageClass
apiVersion: storage.k8s.io/v1beta1
metadata:
  name: aws-efs
provisioner: openshift.org/aws-efs
parameters:
  gidMin: '40000'
  gidMax: '50000'"| oc create -f -

```

*3. Create Sample PVC*
```
echo "
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
 name: efs-pvc
spec:
 accessModes:
  - ReadWriteMany
 resources:
   requests:
     storage: 10Gi
 storageClassName: aws-efs"|oc create -f -


$ oc get pv
NAME                                       CAPACITY   ACCESSMODES   RECLAIMPOLICY   STATUS    CLAIM                                 STORAGECLASS   REASON    AGE
pvc-XXXXXX                                 1Gi        RWX           Delete          Bound     efs-test-project/efs-pvc               aws-efs                 10M


$ oc get pvc
NAME         STATUS    VOLUME                                     CAPACITY   ACCESSMODES   STORAGECLASS   AGE
efs-pvc      Bound     pvc-XXXXX                                  1Gi        RWX           aws-efs        10M

```

*4. Deploy Test App *
```
# By default, there is no storageClassName in template
$ oc get template  django-psql-persistent  -n openshift -o jsonpath='{.objects[?(@.kind=="PersistentVolumeClaim")].spec.storageClassName}'

# Create template with storageClassName 'aws-efs'
$ oc export template  django-psql-persistent  -n openshift   -o yaml|sed  '/storage:/a \    storageClassName: aws-efs' |oc create -f -

# Add storageClassName into test-template
$ oc new-app --template=django-psql-persistent

$ oc get pv
NAME                                       CAPACITY   ACCESSMODES   RECLAIMPOLICY   STATUS    CLAIM                                 STORAGECLASS   REASON    AGE
provisioners-efs                           1Mi        RWX           Retain          Bound     openshift-infra/provisioners-efs                               1d
pvc-a8b5f508-c97c-11e7-a691-0e66b4f8265c   1Gi        RWO           Delete          Bound     efs-test-project/postgresql           aws-efs                  1d


$ oc get pvc
NAME         STATUS    VOLUME                                     CAPACITY   ACCESSMODES   STORAGECLASS   AGE
postgresql   Bound     pvc-a8b5f508-c97c-11e7-a691-0e66b4f8265c   1Gi        RWO           aws-efs        1d


```

## Tip ##
- Manual mount efs file system to host
```
$ mount -t nfs -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2 fs-XXXX.efs.$Region.amazonaws.com:/ ./test -vvvv
```

- All Security Group for ocp nodes(Master/App/Infra) have to be assigned to the EFS.

## Reference
- [Deploying External Persistent Volume Provisioners](https://docs.openshift.com/container-platform/3.6/install_config/provisioners.html)
