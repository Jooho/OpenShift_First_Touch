# Ansible Operator based on existing ansible role

 This tutorial show how to create ansible operator with developed ansible role   
    

## Operator Creation Steps
- Create ansible operator framework 
  ```
  operator-sdk new nfs-provisioner-ansible-operator --api-version=jhouse.com/v1alpha1 --kind=NFS --type=ansible
  cd nfs-provisioner-ansible-operator
  rm -rf roles/nfs
  ```

- Download existing ansible role from ansible-galaxy
  ```
  ansible-galaxy install Jooho.nfs_provisioner -p ./roles
  sed -i "s/nfs/Jooho.nfs_provisioner/g" ./watches.yaml 
  ```

- Update CR(custom Resource)
  ```
  sed -i "s/example-nfs/nfs-provisioner/g" deploy/crds/jhouse_v1alpha1_nfs_cr.yaml
  sed -i "s/size: 3/replicaCount: 1/g" deploy/crds/jhouse_v1alpha1_nfs_cr.yaml  
  ```
- Update role_binding
  This operator is not cluster wide scope so you need to add a namespace
  ```
  oc new-project nfs-provisioner-operator
  export OPERATOR_NAMESPACE=$(oc config view --minify -o jsonpath='{.contexts[0].context.namespace}')
  sed -i "s|REPLACE_NAMESPACE|$OPERATOR_NAMESPACE|g" deploy/role_binding.yaml
  ```

-  Create CRD
   ```
   oc create -f deploy/crds/jhouse_v1alpha1_nfs_crd.yaml 
   ```

## Local Test Steps

-  Install necessary packages
   -  [PIP installation](../../Ansible_Molecule/tools/pip.md)
   ```
   sudo yum install gcc python-devel 
   pip install ansible-runner ansible-runner-http openshift
   
   # CentOS or RHEL
   # yum install python-openshift
   ```
- Copy ansible role to local folder(where watchs.yaml point out)
   ```
   sudo mkdir -p /opt/ansible/roles
   sudo cp -R $(pwd)/roles/Jooho.nfs_provisioner  /opt/ansible/roles
   ```

- Deploy Operator locally
   ```  
   operator-sdk up local --namespace=nfs-provisioner-operator
   ```
- Create CR for test
   ```
   oc new-project nfs-provisioner-operator
   oc apply -f deploy/crds/jhouse_v1alpha1_nfs_cr.yaml 
   ```

- Fix Error!!
  host/api_key are not needed anymore
  ```
  vi /opt/ansible/roles/Jooho.nfs_provisioner/tasks/main.yml
  ...
    api_key: "{{api_key}}"    # Remove
    host: "{{host}}"          # Remove
  ```
  **NOTE** You must change *$(pwd)/roles/Jooho.nfs_provisioner/tasks/main.yml* as well

- Stop/Deploy Operator locally
   ```  
   operator-sdk up local --namespace=nfs-provisioner-operator
   ```

- Create test pvc
  ```
  oc create -f roles/Jooho.nfs_provisioner/tests/test-pvc.yaml
  ``` 

- Clean up
  ```
  oc delete pvc test-pvc
  oc delete nfs --all
  oc delete project nfs-provisioner nfs-provisioner-operator
  ```


## Cluster Test Steps
- Build Operator image
  ```
  operator-sdk build quay.io/jooholee/nfs-provisioner-ansible-operator:v0.0.1
  docker push quay.io/jooholee/nfs-provisioner-ansible-operator:v0.0.1
  ```

- Update operator.yaml 
  ```
  sed -i 's|{{ REPLACE_IMAGE }}|quay.io/jooholee/nfs-provisioner-ansible-operator:v0.0.1|g' deploy/operator.yaml
  sed -i 's|{{ pull_policy\|default('\''Always'\'') }}|Always|g' deploy/operator.yaml
  ```

- Deploy objects for Operator
  ```
  oc new-project nfs-provisioner-operator
  oc create -f deploy/service_account.yaml
  oc create -f deploy/role.yaml
  oc create -f deploy/role_binding.yaml
  oc create -f deploy/operator.yaml
  ```

- Create CR for test
   ```
   oc apply -f  deploy/crds/jhouse_v1alpha1_nfs_cr.yaml 
   ```
- Create test pvc
  ```
  oc create -f roles/Jooho.nfs_provisioner/tests/test-pvc.yaml
  ``` 

### Steps (Clean up)   
   ```
   oc delete pvc --all -n nfs-provisioner
   oc project nfs-provisioner-operator 
   oc delete -f deploy/crds/jhouse_v1alpha1_nfs_cr.yaml   
   oc delete -f deploy/service_account.yaml
   oc delete -f deploy/role.yaml
   oc delete -f deploy/role_binding.yaml
   oc delete -f deploy/operator.yaml
   oc delete -f deploy/crds/jhouse_v1alpha1_nfs_crd.yaml
   oc delete project nfs-provisioner-operator 
   ```
