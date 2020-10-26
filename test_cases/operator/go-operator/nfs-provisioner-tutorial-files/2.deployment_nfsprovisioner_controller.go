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

	// Check if the deployment already exists, if not create a new one
	deployFound := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: nfsprovisioner.Name, Namespace: nfsprovisioner.Namespace}, deployFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForNFSProvisioner(nfsprovisioner)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)

	}

	// Check if the rbac already exists, if not create a new one
	// Check if the service already exists, if not create a new one
	// Check if the deployment already exists, if not create a new one

	return ctrl.Result{Requeue: true}, nil
}

// deploymentForNFSProvisioner returns a NFSProvisioner Deployment object
func (r *NFSProvisionerReconciler) deploymentForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *appsv1.Deployment {
	ls := labelsForNFSProvisioner(m.Name)

	nodeSelector := m.Spec.NodeSelector
	if nodeSelector == nil {
		nodeSelector = map[string]string{"app": "nfs-provisioner"}
	}

	sa := "nfs-provisioner"

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace, //the namespace that NFSProvisioner requested.
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "quay.io/kubernetes_incubator/nfs-provisioner:latest",
						Name:    "nfs-provisioner",
						Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
						Ports: []corev1.ContainerPort{
							{Name: "nfs",
								ContainerPort: 2049},
							{Name: "nfs-udp",
								ContainerPort: 2049,
								Protocol:      "UDP"},
							{Name: "nlockmgr",
								ContainerPort: 32803},
							{Name: "nlockmgr-udp",
								ContainerPort: 32803,
								Protocol:      "UDP"},
							{Name: "mountd",
								ContainerPort: 20048},
							{Name: "mountd-udp",
								ContainerPort: 20048,
								Protocol:      "UDP"},
							{Name: "rquotad",
								ContainerPort: 875},
							{Name: "rquotad-udp",
								ContainerPort: 875,
								Protocol:      "UDP"},
							{Name: "rpcbind",
								ContainerPort: 111},
							{Name: "rpcbind-udp",
								ContainerPort: 111,
								Protocol:      "UDP"},
							{Name: "statd",
								ContainerPort: 662},
							{Name: "statd-udp",
								ContainerPort: 662,
								Protocol:      "UDP"},
						},
						SecurityContext: &corev1.SecurityContext{
							Capabilities: &corev1.Capabilities{
								Add: []corev1.Capability{"DAC_READ_SEARCH", "SYS_RESOURCE"},
							},
						},
						Args: []string{"-provisioner=example.com/nfs"},
						Env: []corev1.EnvVar{{
							Name: "POD_IP",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									FieldPath: "status.podIP",
								},
							},
						}, {
							Name:  "SERVICE_NAME",
							Value: "nfs-provisioner",
						}, {
							Name: "POD_NAMESPACE",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									FieldPath: "metadata.namespace",
								},
							},
						}},
						ImagePullPolicy: "IfNotPresent",
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "export-volume",
							MountPath: "/export",
						}},
					}},
					NodeSelector:       nodeSelector,
					ServiceAccountName: sa,
					Volumes: []corev1.Volume{{
						Name: "export-volume",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/home/core/exports",
							},
						},
					}},
				},
			},
		},
	}
	// Set Memcached instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
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

