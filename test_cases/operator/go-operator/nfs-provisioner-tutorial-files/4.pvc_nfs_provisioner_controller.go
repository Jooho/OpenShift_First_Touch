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
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	
	securityv1 "github.com/openshift/api/security/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	cachev1alpha1 "github.com/jooho/nfs-provisioner-operator/api/v1alpha1"
	"github.com/jooho/nfs-provisioner-operator/controllers/defaults"
)

// NFSProvisionerReconciler reconciles a NFSProvisioner object
type NFSProvisionerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.jhouse.com,resources=nfsprovisioners,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.jhouse.com,resources=nfsprovisioners/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.jhouse.com,resources=nfsprovisioners/finalizers,verbs=update
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=security.openshift.io,resources=securitycontextconstraints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=storage.k8s.io,resources=storageclasses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=endpoints,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch
// +kubebuilder:rbac:groups=policy,resources=podsecuritypolicies,verbs=use
// +kubebuilder:rbac:groups=core,resources=persistentvolumes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is main method for operator
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
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if SCC already exists, if not create a new one
	//https://github.com/openshift/ocs-operator/blob/f10e2314cac2bc16ed5d73da74a0202d0a4cd392/pkg/controller/ocsinitialization/sccs.go

	sccFound := &securityv1.SecurityContextConstraints{}
	err = r.Get(ctx, types.NamespacedName{Name: "nfs-provisioner", Namespace: ""}, sccFound)
	if err != nil {
		if errors.IsNotFound(err) {
			scc := r.sccForNFSProvisioner(nfsprovisioner)
			
			log.Info("Creating a new SecurityContextConstraints", "SecurityContextConstraints.Name", scc.Name)
			
			err := r.Create(ctx, scc)

			if err != nil {
				log.Error(err, "Failed to create a new SecurityContextConstraints", "SecurityContextConstraints.Name", scc.Name)
			}
		}

	}

	// Check if the PVC already exists, if not create a new one
	pvcName := defaults.Pvc

	if nfsprovisioner.Spec.Pvc != "" {
		pvcName = nfsprovisioner.Spec.Pvc
	}

	pvcFound := &corev1.PersistentVolumeClaim{}
	err = r.Get(ctx, types.NamespacedName{Name: pvcName, Namespace: nfsprovisioner.Namespace}, pvcFound)

	if err != nil && nfsprovisioner.Spec.Pvc != "" {
		if errors.IsNotFound(err) {
			pvc := r.pvcForNFSProvisioner(nfsprovisioner)
			
			log.Info("Creating a new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
			
			if err := r.Create(ctx, pvc); err != nil {
				log.Error(err, "Failed to create a new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
			}
		}

	}

	// Check if the serviceaccount already exists, if not create a new one
	saFound := &corev1.ServiceAccount{}
	err = r.Get(ctx, types.NamespacedName{Name: defaults.ServiceAccount, Namespace: nfsprovisioner.Namespace}, saFound)

	if err != nil && errors.IsNotFound(err) {
		// Define a new serviceaccount
		sa := r.serviceAccountForNFSProvisioner(nfsprovisioner)
	
		log.Info("Creating a new Serviceaccount", "Serviceaccount.Namespace", sa.Namespace, "Serviceaccount.Name", sa.Name)
	
		if err := r.Create(ctx, sa); err != nil {
			log.Error(err, "Failed to create a new Serviceaccount", "Serviceaccount.Namespace", sa.Namespace, "Serviceaccount.Name", sa.Name)
		}
	}

	// Check if the rbac(clusterrole/clusterrolebinding/role/rolebinding) already exists, if not create a new one
	// clusterRole
	crFound := &rbacv1.ClusterRole{}
	err = r.Get(ctx, types.NamespacedName{Name: defaults.ClusterRole, Namespace: ""}, crFound)
	if err != nil {
		cr := r.clusterRoleForNFSProvisioner(nfsprovisioner)
		
		log.Info("Creating a new ClusterRole", "ClusterRole.Name", cr.Name)
		
		if err := r.Create(ctx, cr); err != nil {
			log.Error(err, "Failed to create a ClusterRole for NFSProvisioner", "ClusterRole.Namespace", cr.Namespace, "ClusterRole.Name", cr.Name)
		}
	}

	// clusterRoleBinding
	crbFound := &rbacv1.ClusterRoleBinding{}
	err = r.Get(ctx, types.NamespacedName{Name: defaults.ClusterRoleBinding, Namespace: ""}, crbFound)
	if err != nil {
		crb := r.clusterRoleBindingForNFSProvisioner(nfsprovisioner)
		
		log.Info("Creating a new ClusterRoleBinding", "ClusterRoleBinding.Name", crb.Name)
		
		if err := r.Create(ctx, crb); err != nil {
			log.Error(err, "Failed to create a ClusterRoleBinding for NFSProvisioner", "ClusterRoleBinding.Name", crb.Name)
		}
	}
	// Role
	roleFound := &rbacv1.Role{}
	err = r.Get(ctx, types.NamespacedName{Name: defaults.Role, Namespace: nfsprovisioner.Namespace}, roleFound)
	if err != nil {
		role := r.roleForNFSProvisioner(nfsprovisioner)
		log.Info("Creating a new Role", "Role.Namespace", role.Namespace, "Role.Name", role.Name)
		if err := r.Create(ctx, role); err != nil {
			log.Error(err, "Failed to create a Role for NFSProvisioner", "Role.Namespace", role.Namespace, "Role.Name", role.Name)
		}
	}

	// RoleBinding
	roleBindingFound := &rbacv1.RoleBinding{}
	err = r.Get(ctx, types.NamespacedName{Name: defaults.RoleBinding, Namespace: nfsprovisioner.Namespace}, roleBindingFound)
	if err != nil {
		roleBinding := r.roleBindingForNFSProvisioner(nfsprovisioner)
		
		log.Info("Creating a new RoleBinding", "RoleBinding.Namespace", roleBindingFound.Namespace, "roleBinding.Name", roleBindingFound.Name)
		
		if err := r.Create(ctx, roleBinding); err != nil {
			log.Error(err, "Failed to create a RoleBinding for NFSProvisioner", "roleBinding.Namespace", roleBindingFound.Namespace, "roleBinding.Name", roleBindingFound.Name)
		}
	}

	// Check if the deployment already exists, if not create a new one
	deployFound := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: defaults.Deployment, Namespace: nfsprovisioner.Namespace}, deployFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForNFSProvisioner(nfsprovisioner)
		
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		if err = r.Create(ctx, dep); err != nil {
			log.Error(err, "Failed to create a Deployment for NFSProvisioner", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		}
	}

	// Check if the service already exists, if not create a new one
	svcFound := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: defaults.Service, Namespace: nfsprovisioner.Namespace}, svcFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		svc := r.serviceForNFSProvisioner(nfsprovisioner)
		
		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		
		if err = r.Create(ctx, svc); err != nil {
			log.Error(err, "Failed to create a Service for NFSProvisioner", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		}
	}

	// Check if the storageclass already exists, if not create a new one
	scName := defaults.SCForNFSProvisioner

	if nfsprovisioner.Spec.SCForNFSProvisioner != "" {
		scName = nfsprovisioner.Spec.SCForNFSProvisioner
	}

	scFound := &storagev1.StorageClass{}
	err = r.Get(ctx, types.NamespacedName{Name: scName, Namespace: ""}, scFound)
	if err != nil && errors.IsNotFound(err) && nfsprovisioner.Spec.SCForNFSProvisioner == "" {
		// Define a new deployment
		sc := r.storageclassForNFSProvisioner(nfsprovisioner)
		
		log.Info("Creating a new Storageclass", "Storageclass.Name", sc.Name)
		
		if err = r.Create(ctx, sc); err != nil {
			log.Error(err, "Failed to create a Storageclass for NFSProvisioner", "Storageclass.Name", sc.Name)
		}
	}

	// Delete Logic
	// name of our custom finalizer
	const finalizerName = "nfsprovisioner.finalizers.jhouse.io"

	// examine DeletionTimestamp to determine if object is under deletion
	isNFSProvisionerdMarkedToBeDeleted := nfsprovisioner.GetDeletionTimestamp() != nil

	// if nfsprovisioner.ObjectMeta.DeletionTimestamp.IsZero() {
	if isNFSProvisionerdMarkedToBeDeleted {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		
		if !containsString(nfsprovisioner.GetFinalizers(), finalizerName) {
			log.Info("Adding Finalizer for the NFSProvisioner")

			controllerutil.AddFinalizer(nfsprovisioner, finalizerName)

			if err := r.Update(context.Background(), nfsprovisioner); err != nil {
				log.Error(err, "Failed to update CR NFSProvisioner to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(nfsprovisioner.ObjectMeta.Finalizers, finalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.deleteExternalResources(nfsprovisioner); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried

				log.Error(err, "Failed to delete external resoureces")
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(nfsprovisioner, finalizerName)
			
			if err := r.Update(context.Background(), nfsprovisioner); err != nil {
				log.Error(err, "Failed to update CR NFSProvisioner with finalizer to remove finalizer")
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	return ctrl.Result{Requeue: true}, nil
}

// Delete any external resources associated with the nfs server
func (r *NFSProvisionerReconciler) deleteExternalResources(m *cachev1alpha1.NFSProvisioner) error {
	
	// To-Do Delete all pvc that pawned by NFS SC
	ctx := context.Background()

	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaults.ClusterRole,
		}}
	err := r.Get(ctx, types.NamespacedName{Name: defaults.ClusterRole, Namespace: ""}, clusterRole)
	if err == nil {
		err = r.Delete(ctx, clusterRole, &client.DeleteOptions{})
		if err != nil {
			log.Error(err, "Failed to delete ClusterRole for NFSProvisioner", "ClusterRole.Name", defaults.ClusterRole)
			return err
		}
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaults.ClusterRoleBinding,
		}}

	err = r.Get(ctx, types.NamespacedName{Name: defaults.ClusterRoleBinding, Namespace: ""}, clusterRoleBinding)
	if err == nil {
		
		err = r.Delete(ctx, clusterRoleBinding, &client.DeleteOptions{})

		if err != nil {
			log.Error(err, "Failed to delete ClusterRoleBinding for NFSProvisioner", "ClusterRoleBinding.Name", defaults.ClusterRoleBinding)
			return err
		}
	}

	return nil
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

// deploymentForNFSProvisioner returns a NFSProvisioner Deployment object
func (r *NFSProvisionerReconciler) deploymentForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *appsv1.Deployment {
	ls := labelsForNFSProvisioner(m.Name)

	nodeSelector := defaults.NodeSelector

	if m.Spec.NodeSelector != nil {
		nodeSelector = m.Spec.NodeSelector
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
						Image: "quay.io/kubernetes_incubator/nfs-provisioner:latest",
						Name:  "nfs-provisioner",
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
								Add:  []corev1.Capability{"DAC_READ_SEARCH", "SYS_RESOURCE"},
								Drop: []corev1.Capability{"KILL", "MKNOD", "SYS_CHROOT"},
							},
						},

						Args: []string{"'-provisioner=example.com/nfs'"},
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
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: defaults.Pvc,
							},
						},
					}},
				},
			},
		},
	}
	// Set NFSProvisioner instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

//https://github.com/openshift/origin/blob/master/docs/proposals/security-context-constraints.md

func (r *NFSProvisionerReconciler) sccForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *securityv1.SecurityContextConstraints {
	sccName := "nfs-provisioner"
	scc := &securityv1.SecurityContextConstraints{

		ObjectMeta: metav1.ObjectMeta{
			Name:      sccName,
			Namespace: m.Namespace, //the namespace that NFSProvisioner requested.
		},
		AllowHostDirVolumePlugin: true,
		AllowHostIPC:             false,
		AllowHostNetwork:         false,
		AllowHostPID:             false,
		AllowHostPorts:           false,
		AllowPrivilegedContainer: false,
		AllowedCapabilities:      []corev1.Capability{"DAC_READ_SEARCH", "SYS_RESOURCE"},
		DefaultAddCapabilities:   nil,
		Priority:                 nil,
		ReadOnlyRootFilesystem:   false,
		RequiredDropCapabilities: []corev1.Capability{"KILL", "MKNOD", "SYS_CHROOT"},
		RunAsUser: securityv1.RunAsUserStrategyOptions{
			Type: securityv1.RunAsUserStrategyRunAsAny,
		},
		SELinuxContext: securityv1.SELinuxContextStrategyOptions{
			Type: securityv1.SELinuxStrategyMustRunAs,
		},
		Users: []string{"system:serviceaccount:nfs-provisioner:nfs-provisioner"},
		SupplementalGroups: securityv1.SupplementalGroupsStrategyOptions{
			Type: securityv1.SupplementalGroupsStrategyRunAsAny,
		},
		Volumes: []securityv1.FSType{securityv1.FSTypeConfigMap, securityv1.FSTypeDownwardAPI, securityv1.FSTypeEmptyDir, securityv1.FSTypePersistentVolumeClaim, securityv1.FSTypeHostPath, securityv1.FSTypeSecret},
		FSGroup: securityv1.FSGroupStrategyOptions{
			Type: securityv1.FSGroupStrategyMustRunAs,
		},
	}

	ctrl.SetControllerReference(m, scc, r.Scheme)
	return scc
}

func (r *NFSProvisionerReconciler) pvcForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *corev1.PersistentVolumeClaim {
	scName := defaults.SCForNFSPvc

	if m.Spec.SCForNFSPvc != "" {
		scName = m.Spec.SCForNFSPvc
	}

	pvc := &corev1.PersistentVolumeClaim{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "nfs",
			Namespace: m.Namespace, //the namespace that NFSProvisioner requested.
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("10Gi"),
				},
			},
			StorageClassName: &scName,
		},
	}

	ctrl.SetControllerReference(m, pvc, r.Scheme)
	return pvc
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

func (r *NFSProvisionerReconciler) clusterRoleForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *rbacv1.ClusterRole {

	cr := &rbacv1.ClusterRole{

		ObjectMeta: metav1.ObjectMeta{
			Name: "nfs-provisioner-runner",
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"persistentvolumes"},
				Verbs:     []string{"get", "list", "watch", "create", "delete"},
			}, {
				APIGroups: []string{""},
				Resources: []string{"persistentvolumeclaims"},
				Verbs:     []string{"get", "list", "watch", "update"},
			}, {
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"get", "list", "watch"},
			}, {
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"create", "update", "patch"},
			}, {
				APIGroups:     []string{"policy"},
				Resources:     []string{"podsecuritypolicies"},
				ResourceNames: []string{"nfs-provisioner"},
				Verbs:         []string{"use"},
			},
		},
	}

	ctrl.SetControllerReference(m, cr, r.Scheme)
	return cr
}

func (r *NFSProvisionerReconciler) clusterRoleBindingForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *rbacv1.ClusterRoleBinding {

	crb := &rbacv1.ClusterRoleBinding{

		ObjectMeta: metav1.ObjectMeta{
			Name: "run-provisioner-runner",
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      "nfs-provisioner",
			Namespace: m.Namespace,
		}},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     "nfs-provisioner-runner",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	ctrl.SetControllerReference(m, crb, r.Scheme)
	return crb
}

func (r *NFSProvisionerReconciler) roleForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *rbacv1.Role {

	role := &rbacv1.Role{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "leader-locking-nfs-provisioner",
			Namespace: m.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"endpoints"},
				Verbs:     []string{"get", "list", "watch", "create", "update", "delete"},
			}, {
				APIGroups: []string{""},
				Resources: []string{"services"},
				Verbs:     []string{"get"},
			},
		},
	}

	ctrl.SetControllerReference(m, role, r.Scheme)
	return role
}

func (r *NFSProvisionerReconciler) roleBindingForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *rbacv1.RoleBinding {

	rolebinding := &rbacv1.RoleBinding{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "leader-locking-nfs-provisioner",
			Namespace: m.Namespace,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      "nfs-provisioner",
			Namespace: m.Namespace,
		}},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     "leader-locking-nfs-provisioner",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	ctrl.SetControllerReference(m, rolebinding, r.Scheme)
	return rolebinding
}

func (r *NFSProvisionerReconciler) serviceForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *corev1.Service {

	ls := labelsForNFSProvisioner(m.Name)
	svc := &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "nfs-provisioner",
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Name: "nfs",
					Port: 2049},
				{Name: "nfs-udp",
					Port:     2049,
					Protocol: "UDP"},
				{Name: "nlockmgr",
					Port: 32803},
				{Name: "nlockmgr-udp",
					Port:     32803,
					Protocol: "UDP"},
				{Name: "mountd",
					Port: 20048},
				{Name: "mountd-udp",
					Port:     20048,
					Protocol: "UDP"},
				{Name: "rquotad",
					Port: 875},
				{Name: "rquotad-udp",
					Port:     875,
					Protocol: "UDP"},
				{Name: "rpcbind",
					Port: 111},
				{Name: "rpcbind-udp",
					Port:     111,
					Protocol: "UDP"},
				{Name: "statd",
					Port: 662},
				{Name: "statd-udp",
					Port:     662,
					Protocol: "UDP"},
			},
			Selector: ls,
		},
	}

	ctrl.SetControllerReference(m, svc, r.Scheme)
	return svc
}

func (r *NFSProvisionerReconciler) storageclassForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *storagev1.StorageClass {

	scName := defaults.SCForNFSProvisioner

	if m.Spec.SCForNFSProvisioner != "" {
		scName = m.Spec.SCForNFSProvisioner
	}
	sc := &storagev1.StorageClass{

		ObjectMeta: metav1.ObjectMeta{
			Name: scName,
		},
		Provisioner: "example.com/nfs",
		Parameters:  map[string]string{"mountOptions": "vers=4.1"},
	}

	ctrl.SetControllerReference(m, sc, r.Scheme)
	return sc
}

// labelsForNFSProvisioner returns the labels for selecting the resources
// belonging to the given NFSProvisioner CR name.
func labelsForNFSProvisioner(name string) map[string]string {
	return map[string]string{"app": "nfs-provisioner", "nfsprovisioner_cr": name}
}

// SetupWithManager return error
func (r *NFSProvisionerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.NFSProvisioner{}).
		Complete(r)
}
