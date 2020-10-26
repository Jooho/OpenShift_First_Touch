# Go Operator Summary

### Operator Scope

- Cluster Level
  - Watching resources in all namesapces

- Namespaced Level
  - Watching resources in a specific namespace

- MultiNamespaced
  - Watching resources in specific namespaces
    ~~~
    ...
    namespaces := []string{"foo", "bar"} // List of Namespaces
    ...
    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme:             scheme,
        MetricsBindAddress: metricsAddr,
        Port:               9443,
        LeaderElection:     enableLeaderElection, 
        LeaderElectionID:   "f1c5ece8.example.com",
        NewCache:           cache.MultiNamespacedCacheBuilder(namespaces), <===
    })
    ...
    ~~~

### How to change `cluster-wide` to a `namespace wide`?
- Default way (A namespace is hard-coded)
  - Add a namespace into `main.go`
  - Add the namespace into RBAC marker in `controller.go` 
  - `make manifests`
- Better way (Set a namespace using Env Variable)
  - How to get `ENV` ?
    ~~~
    // getWatchNamespace returns the Namespace the operator should be watching for changes
    func getWatchNamespace() (string, error) {
        // WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
        // which specifies the Namespace to watch.
        // An empty value means the operator is running with cluster scope.
        var watchNamespaceEnvVar = "WATCH_NAMESPACE"

        ns, found := os.LookupEnv(watchNamespaceEnvVar)
        if !found {
            return "", fmt.Errorf("%s must be set", watchNamespaceEnvVar)
        }
        return ns, nil
    }
    ~~~
  - In manager, it sets a `Namespace` dynamically
    ~~~
    watchNamespace, err := getWatchNamespace()
    if err != nil {
        setupLog.Error(err, "unable to get WatchNamespace, " +
        "the manager will watch and manage resources in all namespaces")
    }

    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
        Scheme:             scheme,
        MetricsBindAddress: metricsAddr,
        Port:               9443,
        LeaderElection:     enableLeaderElection,
        LeaderElectionID:   "f1c5ece8.example.com",
        Namespace:          watchNamespace, // namespaced-scope when the value is not an empty string
    })
    ~~~
  - Set `ENV` for a namespace in the config/manager/manager.yaml:
    ~~~
      spec:
        containers:
        - command:
            - /manager
            ...
            env:
            - name: WATCH_NAMESPACE
                valueFrom:
                fieldRef:
                    fieldPath: metadata.namespace  <====
    ~~~

  - Multi namespaces
    ~~~
    watchNamespace, err := getWatchNamespace()
    if err != nil {
        setupLog.Error(err, "unable to get WatchNamespace, " +
            "the manager will watch and manage resources in all Namespaces")
    }

    // Add support for MultiNamespace set in WATCH_NAMESPACE (e.g ns1,ns2)
    if strings.Contains(watchNamespace, ",") {
        setupLog.Infof("manager will be watching namespace %q", watchNamespace)
        // configure cluster-scoped with MultiNamespacedCacheBuilder
        options.Namespace = ""
        options.NewCache = cache.MultiNamespacedCacheBuilder(strings.Split(watchNamespace, ","))
    }
    ~~~
    - Set `ENV` for a namespace in the config/manager/manager.yaml:
      ~~~
        spec:
          containers:
          - command:
              - /manager
              ...
              env:
              - name: WATCH_NAMESPACE
                value: "ns1,ns2"  <====
      ~~~




- OLM integration
  - Cluster Scope is default
  - Namespace Scope
  
    - If the operator can watch its own namespace, set the following in your spec.installModes list:
        ~~~
        - type: OwnNamespace
        supported: true
        ~~~
    - If the operator can watch a single namespace that is not its own, set the following in your spec.installModes list:
        ~~~
        - type: SingleNamespace
        supported: true
        ~~~
  - Multi Namespaces Scope
    - If the operator can watch multiple namespaces, set the following in your spec.installModes list:
        ~~~
        - type: MultiNamespace
        supported: true
        ~~~

### CRD Scope
Basically CRD Scope is for Cluster.

- Set `Cluster` Scope explicitly 
  ~~~
  vi api/v1alpha1/memcached_types.go

    // Memcached is the Schema for the memcacheds API
    // +kubebuilder:resource:path=memcacheds,scope=Cluster  <== 
    type Memcached struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   MemcachedSpec   `json:"spec,omitempty"`
    Status MemcachedStatus `json:"status,omitempty"`
    }
  ~~~
  - **Note** The resouce of `// +kubebuilder:resource:path=<resource>,scope=Cluster` is `plural` value in `/config/crd/bases/cache.example.com_memcacheds.yaml`


## Testing Operator project(TODO)
- https://book.kubebuilder.io/cronjob-tutorial/writing-tests.html


### Add 3rd API

- If the 3rd party resource have AddToScheme()
    ~~~
    vi controller.go

    import(
        securityv1 "github.com/openshift/api/security/v1"
        // +kubebuilder:scaffold:imports
    )

    var (
        scheme   = apiruntime.NewScheme()
        setupLog = ctrl.Log.WithName("setup")
    )

    func init() {
        utilruntime.Must(clientgoscheme.AddToScheme(scheme))

        utilruntime.Must(cachev1alpha1.AddToScheme(scheme))  <=====
    }
    ~~~
- If the 3rd party resource does not have AddToScheme()
  ~~~
    import (
        ...
        "k8s.io/apimachinery/pkg/runtime/schema"
        "sigs.k8s.io/controller-runtime/pkg/scheme"
        ...
        // DNSEndoints
        externaldns "github.com/kubernetes-incubator/external-dns/endpoint"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    )

    func init() {
        ...

        log.Info("Registering Components.")

        schemeBuilder := &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "externaldns.k8s.io", Version: "v1alpha1"}}
        schemeBuilder.Register(&externaldns.DNSEndpoint{}, &externaldns.DNSEndpointList{})
        if err := schemeBuilder.AddToScheme(mgr.GetScheme()); err != nil {
            log.Error(err, "")
            os.Exit(1)
        }
        ...
    }
  ~~~
- **NOTE**
  - After adding new import paths to your operator project, run `go mod vendor` if a vendor/ directory is present in the root of your project directory to fulfill these dependencies.
  - Your 3rd party resource needs to be added before add the controller in "Setup all Controllers".


### Delete Reconcile 
- NFS Provisioner updated

### Advanced Topics
- Scorecard [TODO](https://sdk.operatorframework.io/docs/advanced-topics/scorecard/)
- Operator Capability levels([TODO](https://sdk.operatorframework.io/docs/advanced-topics/operator-capabilities/operator-capabilities/))
