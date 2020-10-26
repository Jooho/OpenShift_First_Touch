# NFS Provisioner Go Operator - Controller Analysis

This tutorial show the typical controller develop flow


## Reconcile
Controller main function but logic is simple.

- Main logic
  - Check CR. If not, skip. 
  - Check resources. If not, create them
  - Check if CR deleted. If yes, delete some resources that are not deleted by garbage collector

## Create resources
- Find a scheme 
  - Using [API Group](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19)
  - Using [Go doc for API](https://godoc.org/k8s.io/api/core)

- Use [Kubernete Client] ("sigs.k8s.io/controller-runtime/pkg/client") for Get/Delete/Update/Create resources
  
- Use [Kubenetes Error](https://github.com/kubernetes/apimachinery/blob/master/pkg/api/errors/errors.go) for error handling

- In order to add a 3rd API scheme, use `utilruntime`
  ~~~
  vi  main.go

  //Add 3rd API Scheme
	utilruntime.Must(securityv1.AddToScheme(scheme))
  ~~~


## Handle Log
- Setup Log
  ~~~
  // main.go
  var setupLog = ctrl.Log.WithName("setup")

  setupLog.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
  ~~~
- Controller Log
  ~~~
  Log    logr.Logger
  log.Info("NFSProvisioner resource not found. Ignoring since object must be deleted")
  ~~~
- [logr](https://github.com/go-logr/logr) (log interface) - [Zap](https://github.com/uber-go/zap) (logging framework)


### Must-gather
https://github.com/openshift/ocs-operator/tree/master/must-gather