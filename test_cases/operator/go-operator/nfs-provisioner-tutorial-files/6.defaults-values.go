package defaults

const (

	// SecurityContextContrants is the permission control mechanism in Openshift and it is the same as PodSecurityPolices in Kubenetes
	// For OpenShift, you have to create SCC even though PSP can be created.
	SecurityContextContrants = "nfs-provisioner"
	//HostPathDir is the directory that NFS server will use.
	//NFS server will use PVC by default.
	HostPathDir = "/home/core/nfs"
	//StorageSize is the PVC size for NFS server will use.
	StorageSize = "10G"
	//Pvc is the storage for NFS server will use.
	Pvc = "nfs-server"
	//SCForNFSPvc is the storageClass name to create PVC for NFS server.
	SCForNFSPvc = "local-sc"
	//ServiceAccount is the project level main sa that has power to control NFS provisioners
	ServiceAccount = "nfs-provisioner"
	//ClusterRole is for NFS Provisioner to create SC/PV/PVC
	ClusterRole = "nfs-provisioner-runner"
	//ClusterRoleBinding match ClusterRole and ServiceAccount in the NFS provisioner project
	ClusterRoleBinding = "nfs-provisioner-runner"
	//Role gives the permissions to get endpoints/services for NFS server.
	Role = "leader-locking-nfs-provisioner"
	//RoleBinding gives the Role to the SA
	RoleBinding = "leader-locking-nfs-provisioner"
	//Deployment is for NFS server
	Deployment = "nfs-provisioner"
	//Service is for NFS provisioner to access to NFS Server
	Service = "nfs-provisioner"
	//SCForNFSProvisioner is for NFS Provisioner
	SCForNFSProvisioner = "nfs"
)

var (
	// NodeSelector is for the node where NFS server will be running
	NodeSelector = map[string]string{"app": "nfs-provisioner"}
	 
)
