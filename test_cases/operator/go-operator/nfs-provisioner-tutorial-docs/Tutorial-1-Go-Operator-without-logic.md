# Tutorial - Your first Go operator(Memcached) without local

This tutorial show how to create a new go operator with operator sdk.
The new operator controller does not have any logic and it does not create a sample CR.


## Pre-requisites
- [Setup Go Operator Dev Environment](./1.Setup.md)

## Tutorial Flows
- Set environment variables
- Create Go Operator
- Update Makefile
- Create a sample API
- Build Operator
- Push Operator
- Deploy Operator on the cluster


## Steps
### 1. Set Environment variables for a new operator

```
export NEW_OP_NAME=memcached-operator-without-logic
export NEW_OP_HOME=${ROOT_HOME}/operator-projects/${NEW_OP_NAME}
export NAMESPACE=${NEW_OP_NAME}
export IMG=quay.io/jooholee/${NEW_OP_NAME}:latest

```

### 2. Create Go Operator
  ```
  mkdir -p ${NEW_OP_HOME}
  cd ${NEW_OP_HOME}

  operator-sdk init --domain=example.com --repo=github.com/example-inc/${NEW_OP_NAME}
  ```

### 3. [Update Makefile for test](2.Update_Makefile.md)

### 4. Create a sample API
```
operator-sdk create api --group=cache --version=v1alpha1 --kind=Memcached
Create Resource [y/n]
y
Create Controller [y/n]
y

```

### 5. Build Operator Image
```
make podman-build IMG=quay.io/jooholee/${NEW_OP_NAME}:latest
```

### 6. Push Operator Image
```
make podman-push IMG=quay.io/jooholee/${NEW_OP_NAME}:latest
```

### 7. Deploy Operator on the cluster
```
#update namespace
cd config/default/; kustomize edit set namespace "${NAMESPACE}" ;cd ../..

make deploy IMG=quay.io/jooholee/memcached-operator:latest

oc project ${NAMESPACE}

oc get pod
NAME                                                          READY   STATUS    RESTARTS   AGE
memcached-operator-test-controller-manager-6c888b589b-cphjk   2/2     Running   0          20s
``` 

### 8. Clean up
```
# Delete All Objects
kustomize build config/default | kubectl delete -f -
```


  