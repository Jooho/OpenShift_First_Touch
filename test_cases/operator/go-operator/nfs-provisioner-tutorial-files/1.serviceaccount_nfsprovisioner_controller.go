/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	cachev1alpha1 "github.com/jooho/nfs-provisioner-operator/api/v1alpha1"
)

// NFSProvisionerReconciler reconciles a NFSProvisioner object
type NFSProvisionerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.jhouse.com,resources=nfsprovisioners,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.jhouse.com,resources=nfsprovisioners/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.jhouse.com,resources=nfsprovisioners/finalize,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

func (r *NFSProvisionerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("nfsprovisioner", req.NamespacedName)

	// Fetch the NFSProvisioner instance
	nfsprovisioner := &cachev1alpha1.NFSProvisioner{}
	err := r.Get(ctx, req.NamespacedName, nfsprovisioner)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("NFSProvisioner resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get NFSProvisioner")
		return ctrl.Result{}, err
	}

	// Check if the serviceaccount already exists, if not create a new one
	saFound := &corev1.ServiceAccount{}
	err = r.Get(ctx, types.NamespacedName{Name: "nfs-provisioner", Namespace: nfsprovisioner.Namespace}, saFound)

	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		sa := r.serviceAccountForNFSProvisioner(nfsprovisioner)
		log.Info("Creating a new Serviceaccount", "Serviceaccount.Namespace", sa.Namespace, "Serviceaccount.Name", sa.Name)
		err := r.Create(ctx, sa)
		if err != nil {
			log.Info("Failed to create a new Serviceaccount", "Serviceaccount.Namespace", sa.Namespace, "Serviceaccount.Name", sa.Name, "Error", err)
		}
	}

	// Check if the rbac(clusterrole/clusterrolebinding/role/rolebinding) already exists, if not create a new one
	// Check if the scc(optional) already exists, if not create a new one
	// Check if the deployment already exists, if not create a new one
	// Check if the service already exists, if not create a new one

	return ctrl.Result{Requeue: true}, nil
}


func (r *NFSProvisionerReconciler) serviceAccountForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *corev1.ServiceAccount {

	sa := &corev1.ServiceAccount{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "nfs-provisioner",
			Namespace: m.Namespace, //the namespace that NFSProvisioner requested.
		},
	}

	ctrl.SetControllerReference(m, sa, r.Scheme)
	return sa
}

// labelsForMemcached returns the labels for selecting the resources
// belonging to the given memcached CR name.
func labelsForNFSProvisioner(name string) map[string]string {
	return map[string]string{"app": "nfs-provisioner", "nfsprovisioner_cr": name}
}

func (r *NFSProvisionerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.NFSProvisioner{}).
		Complete(r)
}

