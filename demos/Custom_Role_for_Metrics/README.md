# How to get data of `oc adm top pod` for normal user or service account


`oc adm top pod` show you the resource utilization but it requires enough permission.
Using custom cluster role, you can give the permission to user/group or service account.

From this tutorial, I will use Service Account.

## Test Envionment

- OKD 
  - 3.11

### Pre-requisites
  - Create user `joe` who has cluster-admin role
  - Create user `sue` & project `sue-prj` & service account `sue-sa`
  - Deploy test applications
  - Test `oc adm top pod`
```
ansible -i /etc/ansible/hosts masters -m command -a "htpasswd -bc /etc/origin/master/htpasswd joe redhat"
ansible -i /etc/ansible/hosts masters[0] -m command -a "oc adm policy add-cluster-role-to-user cluster-admin joe" 


ansible -i /etc/ansible/hosts masters -m command -a "htpasswd -b /etc/origin/master/htpasswd sue redhat"
oc login --username=sue --password=redhat

oc new-project sue-prj
oc create sa sue-sa
oc new-app --template=cakephp-mysql-example

oc adm top pod --heapster-namespace='openshift-infra' --heapster-scheme="https" 
oc policy who-can get pods.metrics.k8s.io
```


### Solution

- Creating a cluster role
  Note: Only cluster-admin can create this role.
  ```
  oc login --username=joe --password=redhat

  cat <<EOF> custom-role.yaml
  kind: ClusterRole
  apiVersion: rbac.authorization.k8s.io/v1
  metadata:
    name: system:aggregated-metrics-reader
    labels:
      rbac.authorization.k8s.io/aggregate-to-view: "true"
      rbac.authorization.k8s.io/aggregate-to-edit: "true"
      rbac.authorization.k8s.io/aggregate-to-admin: "true"
  rules:
  - apiGroups: ["metrics.k8s.io"]
    resources: ["pods"]
    verbs: ["get", "list", "watch"]
  EOF

  oc create -f ./custom-role.yaml

  ```

- Add the *service account* to role
  ```
  oc adm policy add-cluster-role-to-user system:aggregated-metrics-reader system:serviceaccount:sue-prj:sue-sa
  ```
  From now on, the service account have power to see metrics of all namespaces


  - Login with the service account
    ```
    TOKEN=$(oc sa get-token sue-sa -n sue-prj) 
    oc login --token=$TOKEN 
    ```

  - Check if the `oc adm top pod` is working
    ```
    oc adm top pod --heapster-namespace='openshift-infra' --heapster-scheme="https" 
    ```

- Add the *user* to role
  ```
  oc adm policy add-cluster-role-to-user system:aggregated-metrics-reader sue
  ```
  With the new clusterrole, user already have the power to see the metric of the project `sue-prj` only.

  If you want to give the user more power can see metrics of all projects, you can also give the clusterrole `system:aggregated-metrics-reader` to the user
  
  - Login with the user
    ```
    oc login --username=sue --password=redhat
    ```

  - Check if the `oc adm top pod` is working for default project
    ```
    oc adm top pod --heapster-namespace='openshift-infra' --heapster-scheme="https" -n default
    ```
    
    
## Clean up

```
oc login --username=joe --password=redhat
oc adm policy remove-cluster-role-from-user system:aggregated-metrics-reader sue
oc adm policy remove-cluster-role-from-user system:aggregated-metrics-reader system:serviceaccount:sue-prj:sue-sa
oc delete clusterrole system:aggregated-metrics-reader
oc delete all --all -n sue-prj
oc delete project sue-prj
oc delete user sue
oc delete identity htpasswd_auth:sue
```
  

