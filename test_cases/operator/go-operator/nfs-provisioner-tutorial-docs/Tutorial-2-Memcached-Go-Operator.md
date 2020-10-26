# Tutorial 2 - Memcached Go operator

This tutorial show how to create a new memcached go operator with operator sdk. You can find Official doc from [here](https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/)


## Pre-requisites
- [Setup Go Operator Dev Environment](./1.Setup.md)

## Tutorial Flows
- Set environment variables
- Create Go Operator
- Update Makefile for test
- Add API for memcached Kind
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
export NEW_OP_NAME=memcached-operator-with-logic
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

### 4. Add API for memcached Kind
```
operator-sdk create api --group=cache --version=v1alpha1 --kind=Memcached
Create Resource [y/n]
y
Create Controller [y/n]
y

```
- What this mean, operator sdk will:
  - Create a controller `controllers/memcached_controller.go`
    - You need to add logic in the Reconcile
    - You should add [RBAC Marker](https://book.kubebuilder.io/reference/markers/rbac.html)
    - You can refer [this example](https://github.com/operator-framework/operator-sdk/blob/e4635fa2eb8d8e07229723575a4f7a5ac89e79b4/example/memcached-operator/memcached_controller.go.tmpl). This tutorial will use this file.
  - Create memcached type `api/v1alpha1/memcached_types.go`
    - You need to add spec and [marker](https://book.kubebuilder.io/reference/markers.html) regarding validation, generation and processing.


### 5. Update Resource Type (CRD)
~~~
vi api/v1alpha1/memcached_types.go
..

type MemcachedSpec struct {
    // INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
    // Important: Run "make" to regenerate code after modifying this file
    // Foo is an example field of Memcached. Edit Memcached_types.go to remove/update
    Foo string `json:"foo,omitempty"`
     
    ## ADDED ## 
    // +kubebuilder:validation:Minimum=1
    Size int `json:"size"`

}

// MemcachedStatus defines the observed state of Memcached
type MemcachedStatus struct {
        // INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
        // Important: Run "make" to regenerate code after modifying this file

        ## ADDED ##
        // Nodes are the names of the memcached pods
        Nodes []string `json:"nodes"`
}

 make generate
~~~
 
- After modifying the `*_types.go` file always run the following command
  `api/v1alpha1/memcached_types.go`
  ```
  make generate
  ```

### 6.Generate CRD 
  ```  
  make manifests

  vi config/crd/bases/cache.example.com_memcacheds.yaml
  ..
  size:
    minimum: 1
    type: integer

  ```

  Check CRD: `config/crd/bases/cache.example.com_memcacheds.yaml`
- You have to use [marker](https://book.kubebuilder.io/reference/markers.html) properly.
  - [validation marker](https://book.kubebuilder.io/reference/markers/crd-validation.html)


### 7. Build Operator Image
  ```
   make podman-build IMG=${IMG}
  ```

### 8. Push Operator Image
  ```
  make podman-push IMG=${IMG}
  ```

### 9. Create CRD
  ```
  make install
  ```

### 10. Deploy Operator 

#### 10.1 On Local
```
oc new-project ${NAMESPACE}
make run ENABLE_WEBHOOKS=false
```

#### 10.2 On Cluster

```
# Update namespace
cd config/default/; kustomize edit set namespace "${NAMESPACE}" ;cd ../..

# Deploy Operator
make deploy IMG=${IMG}

oc project ${NAMESPACE}

oc get pod
``` 

### 11. Create CR for test

#### 11.0 Update CR
~~~
vi config/samples/cache_v1alpha1_memcached.yaml
..
spec:
  # Add fields here
  foo: bar
  size: 2   #<=== ADDED
~~~

#### 11.1 On Local
~~~
# Open a new terminal
oc apply -f config/samples/cache_v1alpha1_memcached.yaml 

oc get pod
~~~

#### 11.2 On Cluster
```
oc apply -f config/samples/cache_v1alpha1_memcached.yaml 
oc logs deployment.apps/memcached-operator-controller-manager  -c manager
```

### 12. Update CR
~~~
oc patch memcached memcached-sample -p '{"spec":{"size": 3}}' --type=merge
~~~

### 13. Clean up

```
# Delete CR
oc delete -f config/samples/cache_v1alpha1_memcached.yaml

# Delete All Objects
kustomize build config/default | kubectl delete -f -
```

## Issues
- Finalizer issue
  - [Helm operator Issue](https://github.com/operator-framework/operator-sdk/issues/3767)
  - [Go operator Issue](https://github.com/kubernetes-sigs/kubebuilder/issues/1654) 
  - How to solve?
    - Add this rbac marker 
      ```
      // +kubebuilder:rbac:groups=<group>.<domain>,resources=<resource>/finalizers,verbs=get;update;patch
      example) 
      // +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/finalizers,verbs=update
      ```
- Failed to watch *v1.Pod: unknown (get pods) issue
  - How to solve?
    - Add watch permission in `controllers/memcached_controller.go `
    ```
    vi controllers/memcached_controller.go 
    ..
    // +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

    make generate
    oc delete -f ./config/crd/bases/cache.example.com_memcacheds.yaml 
    make install
    ```

