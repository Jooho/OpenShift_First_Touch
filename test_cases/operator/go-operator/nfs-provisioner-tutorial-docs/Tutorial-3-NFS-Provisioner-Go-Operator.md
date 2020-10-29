# NFS Provisioner Go Operator 

This tutorial show how to create a new NFS Provisioner go operator with operator sdk. 


## Pre-requisites
- [Setup Go Operator Dev Environment](./1.Setup.md)
- [HostPath/Deploy Local Storage Operator](./3.LocalStorage.md)

## Tutorial Flows
- Set environment variables
- Create Go Operator
- Update Makefile for test
- Add API for nfsprovisioner Kind
- Update Resource Type (CRD)
- Generate CRD 
- Build Operator
- Push Operator
- Create CRD
- Deploy Operator 
  - On Local
  - On Cluster
- Create CR
  - On Local
  - On Cluster
- Update CR
- Clean up


## Steps
### 1. Set Environment variables for a new operator

```
export NEW_OP_NAME=test-nfs-provisioner-operator
export NEW_OP_HOME=${ROOT_HOME}/operator-projects/${NEW_OP_NAME}
export NAMESPACE=${NEW_OP_NAME}
export VERSION=0.0.1
export IMG=quay.io/jooholee/${NEW_OP_NAME}:${VERSION}
```

### 2. Create Go Operator
  ```
  mkdir -p ${NEW_OP_HOME}
  cd ${NEW_OP_HOME}

  operator-sdk init --domain=jhouse.com --repo=github.com/jooho/${NEW_OP_NAME}
  ```

### 3. [Update Makefile for test](2.Update_Makefile.md)

### 4. Add API for nfsprovisioner Kind
```
operator-sdk create api --group=cache --version=v1alpha1 --kind=NFSProvisioner
Create Resource [y/n]
y
Create Controller [y/n]
y
```

### 5. Update Resource Type (CRD)
- Update `api/v1alpha1/nfsprovisioner_types.go`
  ~~~
  vi api/v1alpha1/nfsprovisioner_types.go 

  type NFSProvisionerSpec struct {
        	// HostPathDir is the folder that NFS server will use.
          HostPathDir string `json:"hostPathdir,omitempty"`

          // PVC resource will be attached to NFS server
          // If there is already pvc created, you can use this param
          Pvc string `json:"pvc,omitempty"`

          // NFS server will be running on a specific node by NodeSeletor
          NodeSelector map[string]string `json:"nodeSelector,omitempty"`

          // StorageClass for NFS PVC
          // If you want to delegate this operator to create a pvc for NFS Server from StorageClass, you set this param for the storageclass name
          SCForNFSPvc string `json:"scForNFSPvc,omitempty"` //https://golang.org/pkg/encoding/json/
          
          // StorageClass name gor NFS provisioner
          // This operator will create a StorageClass for NFS provisioner named "nfs". If you want to change it, you set this param.
          SCForNFSProvisioner string `json:"scForNFS,omitempty"` //https://golang.org/pkg/encoding/json/

  }

  // NFSProvisionerStatus defines the observed state of NFSProvisioner
  type NFSProvisionerStatus struct {

          // Nodes are the node names of the NFS server are running.
          Nodes []string `json:"nodes"`
          // Error show error messages briefly
          Error string `json:"error"`
  }
  ~~~
- Update `api/v1alpha1/zz_generated.deepcopy.go file` based on type.go
  ~~~
  make generate
  ~~~
- Generate CRD Manifests
  ~~~
  make manifests
  ~~~

### 6. Controller Update 
From this tutorial, we will use [this file](
../nfs-provisioner-tutorial/5.full_finished.go)
but Tutorial 4 will explain this controller in detail.

- Controller
  - ServiceAccount
  - SecurityContextContraints(psp)
  - ClusterRole
  - ClusterroleBinding
  - Rolebinding
  - Role
  - Deployment
  - Service
  - Pvc

  ~~~
  cp ../nfs-provisioner-tutorial/5.full_finished.go ${NEW_OP_HOME}/controllers/nfsprovisioner_controller.go 
  ~~~

- Copy [Default vaules](../nfs-provisioner-tutorial/6.defaults-values.md) 
  ~~~
  mkdir ${NEW_OP_HOME}/controllers/defaults
  cp ../nfs-provisioner-tutorial/6.defaults-values.md ${NEW_OP_HOME}/controllers/defaults/default.go
  ~~~

- Add securityv1 schema
  ~~~
  import (
    ...
    securityv1 "github.com/openshift/api/security/v1"
  )
  func init(){
  ..
	//Add 3rd API Scheme
	utilruntime.Must(securityv1.AddToScheme(scheme))

  }
  ~~~

### 7. Build Operator Image

- Small refactoring
  - Move main.go
    ~~~
    cd ${NEW_OP_HOME}
    mkdir ./cmd

    mv ./main.go  ./cmd/main.go
    ~~~
  - Update Makefile
    ~~~
    vi Makefil
    s/main.go /cmd\/main.go/g
    ~~~

- Bduil
  ```
  make podman-build IMG=${IMG}
  ```

### 8. Push Operator Image
```
make podman-push IMG=${IMG}
```

### 9. Deploy Operator 

#### 9.1 On Local
```
oc new-project ${NAMESPACE}
make run ENABLE_WEBHOOKS=false
```

#### 9.2 On Cluster

```
# Update namespace
cd config/default/; kustomize edit set namespace "${NAMESPACE}" ;cd ../..

# Deploy Operator
make deploy IMG=${IMG}

oc project ${NAMESPACE}

oc get pod
``` 
### 10. Create CR for test

#### 10.0 Update Sample CR
- Default(LocalStorage)
  ~~~
  vi config/samples/cache_v1alpha1_nfsprovisioner.yaml

  apiVersion: cache.jhouse.com/v1alpha1
  kind: NFSProvisioner
  metadata:
    name: nfsprovisioner-sample
  spec:
    scForNFSPvc: local-sc   #<==  this tutorial will use local sc first
  ~~~
- Create a new CR for existed PVC
  ~~~
  vi config/samples/cache_v1alpha1_nfsprovisioner_pvc.yaml

  apiVersion: cache.jhouse.com/v1alpha1
  kind: NFSProvisioner
  metadata:
    name: nfsprovisioner-sample
  spec:
    Pvc: nfs-server
  ~~~
- Create a new CR for HostPath
  ~~~
  vi config/samples/cache_v1alpha1_nfsprovisioner_hostPath.yaml

  apiVersion: cache.jhouse.com/v1alpha1
  kind: NFSProvisioner
  metadata:
    name: nfsprovisioner-sample
  spec:
    nodeSelector:
      kubernetes.io/hostname: worker-0.bell.tamlab.brq.redhat.com
    hostPathDir: "/home/core/test"
  ~~~


#### 10.1 On Local
~~~
// Open a new terminal
oc apply -f config/samples/cache_v1alpha1_nfsprovisioner.yaml 

oc get pod
~~~

#### 10.2 On Cluster
~~~
oc apply -f config/samples/cache_v1alpha1_nfsprovisioner.yaml 
oc logs deployment.apps/nfs-provisioner-operator-controller-manager  -c manager
~~~


### 11. Create NFS PVC
~~~
oc create -f $TEST_HOME/test-pvc.yaml
~~~

### 12. Clean up

```
# Delete PVC
oc delete -f $TEST_HOME/test-pvc.yaml

# Delete CR
oc delete -f config/samples/cache_v1alpha1_nfsprovisioner.yaml

# Delete All Objects
kustomize build config/default | kubectl delete -f -

```

### Tip
- go build
  ~~~
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/manager ./cmd/main.go
  ~~~

- go run
  ~~~
  go run ./cmd/main.go
  ~~~
