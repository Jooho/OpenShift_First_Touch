# Custom Role for Metrics with Service Account


`oc adm top pod` show you the resource utils but it requires enough permission.
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

- Add the service account to role
```
oc adm policy add-cluster-role-to-user system:aggregated-metrics-reader system:serviceaccount:sue-prj:sue-sa
```


- Login with the service account
```
TOKEN=$(oc sa get-token sue-sa -n sue-prj) 
oc login --token=$TOKEN 
```

- Check if the `oc adm top pod` is working
```
oc adm pod pod --heapster-namespace='openshift-infra' --heapster-scheme="https" 
```

  


