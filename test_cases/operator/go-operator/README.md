# Go Operator 

[Go operator doc](https://sdk.operatorframework.io/docs/building-operators/golang)

## Tutorials
- [Go-Operator-without-logic](docs/Tutorial-1-Go-Operator-without-logic.md)
- [Go-Operator-with-logic](docs/Tutorial-2-Memcached-Go-Operator.md)

## Commands
- Create Operator
  ```
  mkdir -p $HOME/projects/memcached-operator
  cd $HOME/projects/memcached-operator

  operator-sdk init --domain=example.com --repo=github.com/example-inc/memcached-operator
  ```

- Create a new API and Controller
  ```
  operator-sdk create api --group=cache --version=v1alpha1 --kind=Memcached
  ```

- Update resource type
  After modifying the `*_types.go` file always run the following command
  `api/v1alpha1/memcached_types.go`
  ```
  make generate
  ```

- Generate CRD
  [validation marker](https://book.kubebuilder.io/reference/markers/crd-validation.html)
  `CRD path: config/crd/bases/cache.example.com_memcacheds.yaml`
  ```
  make manifests
  ```

- Build/Push image
  ```
   make podman-build podman-push IMG=quay.io/jooholee/memcached-operator:latest
  ``` 

- Run Operator
  ```
  make install
  make deploy IMG=quay.io/jooholee/memcached-operator:latest
  ```

- Deploy CR
  ```
   oc apply -f config/samples/cache_v1alpha1_memcached.yaml 
   oc logs deployment.apps/memcached-operator2-controller-manager  -c manager
  ```

- Clean up
  ```
  oc delete -f config/samples/cache_v1alpha1_memcached.yaml


  kustomize build config/default | kubectl delete -f -
  ```

## Test Operator
- Deploy locally
  ```
   make run ENABLE_WEBHOOKS=false
  ```

- Deploy on cluster
  ```
   make podman-build podman-push IMG=quay.io/jooholee/memcached-operator:latest
   #update namespace
   cd config/default/; kustomize edit set namespace "test-mem" ;cd ../..
 
   make deploy IMG=quay.io/jooholee/memcached-operator:latest
  ``` 

- Create sample CR
  ```
  oc create -f config/samples/cache_v1alpha1_memcached.yaml
  ```



## Issues
- Finalizer issue
  - [Helm operator Issue](https://github.com/operator-framework/operator-sdk/issues/3767)
  - [Go operator Issue](https://github.com/kubernetes-sigs/kubebuilder/issues/1654) 
  - How to solve?
    - Add this rbac marker 
      ```
      // +kubebuilder:rbac:groups=<group>.<domain>,resources=<resource>/finalizers,verbs=get;update;patch
      example) // +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/finalizers,verbs=update
      ```

## TODO
- Deleted all PVC that created by nfs StorageClass

## Reference
- [API Group](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#persistentvolumeclaim-v1-core)
- [Go doc for API](https://godoc.org/k8s.io/api/core)
- [Testing Method Example](https://github.com/digitalocean/clusterlint/pull/46/files)
- [client test Example](https://github.com/kubernetes-sigs/controller-runtime/blob/master/pkg/client/client_test.go)
- [client Source Example](https://github.com/kubernetes-sigs/controller-runtime/tree/master/pkg/client)
- [How to put Namespace/Resource for client](https://github.com/kubernetes/apimachinery/blob/master/pkg/types/namespacedname.go)
- [Prometheus operator](https://github.com/prometheus-operator/prometheus-operator/blob/master/pkg/apis/monitoring/v1/types.go)
- [How to check reconcile](https://developers.redhat.com/blog/2019/10/04/getting-started-with-golang-operators-by-using-operator-sdk/)
- [rbac helper](https://github.com/soltysh/oc/blob/aff24a5b00adecd97fbc06ad64a8246fe8482a18/pkg/cli/admin/policy/modify_scc.go)
- [How to uninstall out of Reconcile](https://github.com/openshift/ocs-operator/blob/d085c469806539ee10a46ffb738e10d4963e583a/pkg/controller/storagecluster/reconcile.go)
- [SCC Information](https://www.openshift.com/blog/managing-sccs-in-openshift)
- [Golang Function Option Pattern](https://github.com/tmrts/go-patterns/blob/master/idiom/functional-options.md)
- [Golang Function Option Pattern2](https://www.sohamkamani.com/golang/options-pattern)
- [Golang Function Option Pattern Example](https://github.com/kubernetes-sigs/controller-runtime/blob/e63948755ba14b3dd8506fafbef4065fb9457817/pkg/client/options.go)
- [Error Handle](https://github.com/kubernetes/apimachinery/blob/master/pkg/api/errors/errors.go)

