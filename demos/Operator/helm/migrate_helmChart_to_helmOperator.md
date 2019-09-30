# Migrate Helm Chart to Helm Operator

## Check List
- namespace
  - check if values.yaml has parameter to override namespace or not
    - if not, you need to add it
    - if so, you can set the namespace for a CR.
- `Release.Name` might be changed

## Test NFS Provisioner Helm Chart
Refer [helm chart doc](./helm-charts/nfs-provisioner/REAEME.md)

## NFS Provisioner Helm Operator
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
  cp -R ~/jhouse_openshift/demos/Operator/helm/helm-charts/nfs-provisioner ./helm-charts/nfs-provisioner
  ```

- Update watch path for nfs-provisioner helm chart
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
    name: nfs-provisioner  # <== Update
  spec:
    ...  

    namespace: nfs-provisioner    # <== Add
    replicaCount: 1    
    
    resources: {}
    ...
    
    tolerations: []
    
    affinity: {}

  ```

- Create namespace.yaml for new provisioner project
  ```
  vi helm-charts/nfs-provisioner/templates/namespace.yaml

  apiVersion: v1
  kind: Namespace
  metadata:
    annotations:
      app.kubernetes.io/name: {{ include "nfs-provisioner.name" . }}
      helm.sh/chart: {{ include "nfs-provisioner.chart" . }}
      app.kubernetes.io/instance: {{ .Release.Name }}
      app.kubernetes.io/managed-by: {{ .Release.Service }}
    name: {{ .Values.namespace | default ( .Release.Namespace ) }}
  spec:
    finalizers:
    - kubernetes
  ```

## Local Test
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
- Open new terminal
  
- Deploy CR
  ```
  oc create -f deploy/crds/jhouse_v1alpha1_nfs_cr.yaml
  ```

- Check NFS provisioner is running well
  ```
  oc project nfs-provisioner
  oc get pods
  oc get pvc
  ```

- Clean up
  ```
  oc delete -f deploy/crds/jhouse_v1alpha1_nfs_cr.yaml -n nfs-provisioner-operator
  oc delete -f deploy/role_binding.yaml
  oc delete -f deploy/role.yaml
  oc delete -f deploy/service_account.yaml -n nfs-provisioner-operator
  oc delete -f deploy/crds/jhouse_v1alpha1_nfs_crd.yaml 
  oc delete pv --all
  oc delete pvc --all -n nfs-provisioner
  oc delete project nfs-provisioner-operator nfs-provisioner
  rm -rf /opt/helm/helm-charts
  
  ```




## Cluster Test 

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


- Create operator project
  ```
  oc new-project nfs-provisioner-operator
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

- Check NFS provisioner is running well
  ```
  oc project nfs-provisioner
  oc get pods
  oc get pvc
  ```

- Test
  ```
  helm template helm-charts/nfs-provisioner -x templates/tests/test-pod.yaml |oc create -f -
  ```


- Clean up
  ```
  oc delete pvc --all -n nfs-provisioner
  oc project nfs-provisioner-operator
  oc delete -f deploy/crds/jhouse_v1alpha1_nfs_cr.yaml 
  oc delete -f deploy/role_binding.yaml
  oc delete -f deploy/role.yaml
  oc delete -f deploy/service_account.yaml -n nfs-provisioner-operator
  oc delete -f deploy/crds/jhouse_v1alpha1_nfs_crd.yaml 
  oc delete pv --all  
  oc delete project nfs-provisioner-operator nfs-provisioner  
  ```
