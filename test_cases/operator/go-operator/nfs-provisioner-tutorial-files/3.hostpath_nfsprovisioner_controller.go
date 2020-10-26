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
	"k8s.io/apimachinery/pkg/api/resource"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	//Added
	storagev1 "k8s.io/api/storage/v1"
	// secv1client "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"
	// securityv1client "github.com/openshift/client-go/security/clientset/versioned"
	// securityv1typedclient "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"
	// secv1client "github.com/openshift/client-go/security/clientset/versioned/typed/security/v1"

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
// +kubebuilder:rbac:groups=cache.jhouse.com,resources=nfsprovisioners/finalizers,verbs=update
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=storage.k8s.io,resources=storageclasses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=policy,resources=podsecuritypolicies,verbs=use
// +kubebuilder:rbac:groups=core,resources=persistentvolumes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

//Reconcile turn nothing  test
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

	// Check if the PVC already exists, if not create a new one
	pvcFound := &corev1.PersistentVolumeClaim{}
	err = r.Get(ctx, req.NamespacedName, pvcFound)

	if err != nil {
		if errors.IsNotFound(err) {
			pvc := r.pvcForNFSProvisioner(nfsprovisioner)
			log.Info("Creating a new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
			err := r.Create(ctx, pvc)
			if err != nil {
				log.Error(err, "Failed to create a new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
			}
		}

	}

	// Check if the serviceaccount already exists, if not create a new one
	saFound := &corev1.ServiceAccount{}
	err = r.Get(ctx, types.NamespacedName{Name: "nfs-provisioner", Namespace: nfsprovisioner.Namespace}, saFound)

	if err != nil && errors.IsNotFound(err) {
		// Define a new serviceaccount
		sa := r.serviceAccountForNFSProvisioner(nfsprovisioner)
		log.Info("Creating a new Serviceaccount", "Serviceaccount.Namespace", sa.Namespace, "Serviceaccount.Name", sa.Name)
		err := r.Create(ctx, sa)
		if err != nil {
			log.Error(err, "Failed to create a new Serviceaccount", "Serviceaccount.Namespace", sa.Namespace, "Serviceaccount.Name", sa.Name)
		}
	}

	// Check if the rbac(clusterrole/clusterrolebinding/role/rolebinding) already exists, if not create a new one
	//clusterRole
	crFound := &rbacv1.ClusterRole{}
	err = r.Get(ctx, types.NamespacedName{Name: "nfs-provisioner-runner", Namespace: ""}, crFound)
	if err != nil {
		cr := r.clusterRoleForNFSProvisioner(nfsprovisioner)
		log.Info("Creating a new ClusterRole", "ClusterRole.Name", cr.Name)
		err := r.Create(ctx, cr)
		if err != nil {
			log.Error(err, "Failed to create a ClusterRole for NFSProvisioner", "ClusterRole.Namespace", cr.Namespace, "ClusterRole.Name", cr.Name)
		}
	}

	//clusterRoleBinding
	crbFound := &rbacv1.ClusterRoleBinding{}
	err = r.Get(ctx, types.NamespacedName{Name: "run-provisioner-runner", Namespace: ""}, crbFound)
	if err != nil {
		crb := r.clusterRoleBindingForNFSProvisioner(nfsprovisioner)
		log.Info("Creating a new ClusterRoleBinding", "ClusterRoleBinding.Name", crb.Name)
		err := r.Create(ctx, crb)
		if err != nil {
			log.Error(err, "Failed to create a ClusterRoleBinding for NFSProvisioner", "ClusterRoleBinding.Name", crb.Name)
			// SCCClient := securityv1typedclient.SecurityContextConstraintsInterface
			// SCCClient.Get("hostaccess",metav1.GetOptions{})
		}
	}
	//Role
	roleFound := &rbacv1.Role{}
	err = r.Get(ctx, types.NamespacedName{Name: "leader-locking-nfs-provisioner", Namespace: nfsprovisioner.Namespace}, roleFound)
	if err != nil {
		role := r.roleForNFSProvisioner(nfsprovisioner)
		log.Info("Creating a new Role", "Role.Namespace", role.Namespace, "Role.Name", role.Name)
		err := r.Create(ctx, role)
		if err != nil {
			log.Error(err, "Failed to create a Role for NFSProvisioner", "Role.Namespace", role.Namespace, "Role.Name", role.Name)
		}
	}

	//RoleBinding
	roleBindingFound := &rbacv1.RoleBinding{}
	err = r.Get(ctx, types.NamespacedName{Name: "leader-locking-nfs-provisioner", Namespace: nfsprovisioner.Namespace}, roleBindingFound)
	if err != nil {
		roleBinding := r.roleBindingForNFSProvisioner(nfsprovisioner)
		log.Info("Creating a new RoleBinding", "RoleBinding.Namespace", roleBindingFound.Namespace, "roleBinding.Name", roleBindingFound.Name)
		err := r.Create(ctx, roleBinding)
		if err != nil {
			log.Error(err, "Failed to create a RoleBinding for NFSProvisioner", "roleBinding.Namespace", roleBindingFound.Namespace, "roleBinding.Name", roleBindingFound.Name)
		}
	}

	// Check if the scc(optional) already exists, if not create a new one
	//https://github.com/openshift/ocs-operator/blob/f10e2314cac2bc16ed5d73da74a0202d0a4cd392/pkg/controller/ocsinitialization/sccs.go

	// Check if the deployment already exists, if not create a new one
	deployFound := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: nfsprovisioner.Name, Namespace: nfsprovisioner.Namespace}, deployFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForNFSProvisioner(nfsprovisioner)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)

	}

	// To-Do Delete ClusterRole/ClusterRoleBinding

	// if nfsprovisioner.GetDeletionTimestamp().IsZero() {
	// 	log.Info("TEST1")
	// 	err = r.Delete(ctx, r.clusterRoleForNFSProvisioner(nfsprovisioner))
	// 	err = r.Delete(ctx, r.clusterRoleBindingForNFSProvisioner(nfsprovisioner))
	// }

	// Check if the service already exists, if not create a new one
	svcFound := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: "nfs-provisioner", Namespace: nfsprovisioner.Namespace}, svcFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		svc := r.serviceForNFSProvisioner(nfsprovisioner)
		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create a Service for NFSProvisioner", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		}
	}
	// Check if the storageclass already exists, if not create a new one

	scFound := &storagev1.StorageClass{}
	err = r.Get(ctx, types.NamespacedName{Name: "example-nfs", Namespace: ""}, scFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		sc := r.storageclassForNFSProvisioner(nfsprovisioner)
		log.Info("Creating a new Storageclass", "Storageclass.Name", sc.Name)
		err = r.Create(ctx, sc)
		if err != nil {
			log.Error(err, "Failed to create a Storageclass for NFSProvisioner", "Storageclass.Name", sc.Name)
		}
	}

	return ctrl.Result{Requeue: true}, nil
}

// deploymentForNFSProvisioner returns a NFSProvisioner Deployment object
func (r *NFSProvisionerReconciler) deploymentForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *appsv1.Deployment {
	ls := labelsForNFSProvisioner(m.Name)

	nodeSelector := m.Spec.NodeSelector
	if nodeSelector == nil {
		nodeSelector = map[string]string{"app": "nfs-provisioner"}
	}
	var privi bool = true
	// sa := "nfs-provisioner"

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
							Privileged:               &privi,
							AllowPrivilegeEscalation: &privi,
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
					NodeSelector: nodeSelector,
					Volumes: []corev1.Volume{{
						Name: "export-volume",
						VolumeSource: corev1.VolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/home/core/exports-nfs",
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

func (r *NFSProvisionerReconciler) pvcForNFSProvisioner(m *cachev1alpha1.NFSProvisioner) *corev1.PersistentVolumeClaim {

	pvc := &corev1.PersistentVolumeClaim{

		ObjectMeta: metav1.ObjectMeta{
			Name:      "nfs-provisioner",
			Namespace: m.Namespace, //the namespace that NFSProvisioner requested.
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteMany"},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
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

	// ctrl.SetControllerReference(m, cr, r.Scheme)
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

	// ctrl.SetControllerReference(m, crb, r.Scheme)
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

	sc := &storagev1.StorageClass{

		ObjectMeta: metav1.ObjectMeta{
			Name: "example-nfs",
		},
		Provisioner: "example.com/nfs",
		Parameters:  map[string]string{"mountOptions": "vers=4.1"},
	}

	ctrl.SetControllerReference(m, sc, r.Scheme)
	return sc
}

// func (r *NFSProvisionerReconciler) addUserToSCC(scc *secv1.SecurityContextConstraints, m *cachev1alpha1.NFSProvisioner) *secv1.SecurityContextConstraints {

// 	scc.Users = []string{
// 		fmt.Sprintf("system:serviceaccount:%s:nfs-provisioner", m.Namespace),
// 	}

// 	return scc

// }

// labelsForNFSProvisioner returns the labels for selecting the resources
// belonging to the given NFSProvisioner CR name.
func labelsForNFSProvisioner(name string) map[string]string {
	return map[string]string{"app": "nfs-provisioner", "nfsprovisioner_cr": name}
}

//SetupWithManager return error
func (r *NFSProvisionerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.NFSProvisioner{}).
		Complete(r)
}

