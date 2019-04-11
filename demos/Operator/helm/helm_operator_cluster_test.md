# Helm Operator Get Started

This is practical 

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
- Update override values by operator
  Delete most parts because it does not need to override any values now.
  ```
  vi deploy/crds/jhouse_v1alpha1_nfs_cr.yaml 
  
  apiVersion: jhouse.com/v1alpha1
  kind: NFS
  metadata:
    name: nfs-provisioner
  spec:
    ...  

    namespace: nfs-provisioner
    replicaCount: 1    
    
    resources: {}
    ...
    
    tolerations: []
    
    affinity: {}

  ```

- Build operator image
  ```
  operator-sdk build quay.io/jooholee/nfs-provisioner-operator:v0.0.1
  ```

- Login with quay.io
  ```
  Account Settings -> Generate Encrypted Password
  ```

- Push operator image
  ```
  docker push quay.io/jooholee/nfs-provisioner-operator:v0.0.1
  
  ```

- Update yaml files
  ```
  # Operator.yaml
  sed -i 's|REPLACE_IMAGE|quay.io/jooholee/nfs-provisioner-operator:v0.0.1|g'  deploy/operator.yaml


  # role_binding.yaml
  export OPERATOR_NAMESPACE=$(oc config view --minify -o jsonpath='{.contexts[0].context.namespace}')
  sed -i "s|REPLACE_NAMESPACE|$OPERATOR_NAMESPACE|g" deploy/role_binding.yaml
  ```

- Deploy NFS Provisioner Operator
  ```
  oc create -f deploy/crds/jhouse_v1alpha1_nfs_crd.yaml 
  oc create -f deploy/service_account.yaml
  oc create -f deploy/role.yaml
  oc create -f deploy/role_binding.yaml
  oc create -f deploy/operator.yaml
  ```

- Create CR (NFS provisioner)
  ```
  oc create -f deploy/crds/jhouse_v1alpha1_nfs_cr.yaml
  ```
- Test
  ```
  helm template helm-charts/nfs-provisioner -x templates/tests/test-pvc.yaml |oc create -f -

  ```

##Clean up

```
oc project nfs-provisioner-operator
oc delete -f deploy/crds/jhouse_v1alpha1_nfs_cr.yaml 
oc delete -f deploy/operator.yaml
oc delete -f deploy/role_binding.yaml
oc delete -f deploy/role.yaml
oc delete -f deploy/service_account.yaml
oc delete -f deploy/crds/jhouse_v1alpha1_nfs_crd.yaml 
oc delete project nfs-provisioner-operator nfs-provisioner
```

## Reference
- [User Guide](https://github.com/operator-framework/operator-sdk/blob/master/doc/helm/user-guide.md)