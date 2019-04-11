# Helm Operator Local Test


- Create a project for nfs-provisioner-operator
  ```
  oc new-project nfs-provisioner-operator
  ```

- Create Helm Operator for nfs provisioner 
  if you don't have `operator-sdk`, please refer [this](../operator-sdk.md)

  ```
  operator-sdk new nfs-provisioner-operator --api-version=jhouse.com/v1alpha1 --kind=NFS --type=helm
  
  cd nfs-provisioner-operator
  ```

- Remove helm chart
  ```
  rm -rf helm-charts/nfs
  ```

- Copy NFS provisioner helm chart
  ```
  cp -R ../helm-charts/nfs-provisioner ./helm-charts/
  ```

- Download nfs-provisioner helm chart
  ```
  sed -i "s/nfs/nfs-provisioner/g" ./watches.yaml 
  ```

- Create symbolic link for local test
  ```
  sudo mkdir -p /opt/helm/helm-charts
  sudo ln -s $(pwd)/helm-charts/nfs-provisioner  /opt/helm/helm-charts/nfs-provisioner
  ```

- Create CRD/Role/RoleBinding/ServiceAccount
  ```
  oc create -f deploy/crds/jhouse_v1alpha1_nfs_crd.yaml
  oc create -f deploy/service_account.yaml
  oc create -f deploy/role.yaml
  oc create -f deploy/role_binding.yaml
  ```

- Deploy operator
  ```
  operator-sdk up local
  ```

- Deploy CR
  ```
  oc create -f deploy/crds/jhouse_v1alpha1_nfs_cr.yaml
  ```

- Clean up
  ```
  oc delete -f deploy/crds/jhouse_v1alpha1_nfs_cr.yaml -n nfs-provisioner-operator
  oc delete -f deploy/role_binding.yaml
  oc delete -f deploy/role.yaml
  oc delete -f deploy/service_account.yaml -n nfs-provisioner-operator
  oc delete -f deploy/crds/jhouse_v1alpha1_nfs_crd.yaml 
  oc delete pvc --all -n nfs-provisioner
  oc delete project nfs-provisioner-operator nfs-provisioner
  
  ```