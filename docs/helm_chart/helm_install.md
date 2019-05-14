# Helm Install

## Components
- HELM client
- Tiller
  
> Tiller is the in-cluster component of Helm. It interacts directly with the Kubernetes API server to install, upgrade, query, and remove Kubernetes resources. It also stores the objects that represent releases.


## Installation

### Install HELM

There are many ways to install helm client. please check [here](https://helm.sh/docs/using_helm/#installing-helm)

Here, I will use script way
```
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get > get_helm.sh
chmod 700 get_helm.sh
./get_helm.sh
```

### Install Tiller

[offical doc](https://helm.sh/docs/using_helm/#role-based-access-control)

- Create RBAC 
```
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tiller
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: tiller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: tiller
    namespace: kube-system
```

- Deploy cluster-wide tiller
~~~
helm init
~~~


Take the permission to access tiller server
```
oc policy remove-role-from-group edit system:authenticated -n tiller
```

You should give the `edit` or `admin` role to only reliable users 
```
oc policy add-role-to-user edit peter -n tiller
```

```
oc policy add-role-to-user admin peter -n tiller
```


https://engineering.bitnami.com/articles/helm-security.html

```
helm version --host=10.128.2.202:44134
```

```
kubectl  patch deployment tiller-deploy --patch '
spec:
  template:
    spec:
      containers:
        - name: tiller
          ports: []
          command: ["/tiller"]
          args: ["--listen=localhost:44134"]
'
```
