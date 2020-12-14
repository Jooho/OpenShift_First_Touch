# NFS Provisioner Operator OLM(Operator Lifecyle Management) - Bundle - 2
[OLM Summary](Tutorial-6-Operator-OLM-Summary.md#Bundle)

## Objects
- Update manifests and metadata of CSV 
  - using markers for CRD
  - using user part for Operator itself.
- Update Bundle version, channel for upgrading an Operator
  - Update your Operator and Bundle
  - Set Channel and DefaultChannel
  - Upgrade your Operator


## Tutorial Flows
- Set environment variables for a new operator bundle



## Steps

### 1. Set environment variables for a new bundle

~~~
export NEW_OP_NAME=test-nfs-provisioner-operator
export NEW_OP_HOME=${ROOT_HOME}/operator-projects/${NEW_OP_NAME}
export NAMESPACE=${NEW_OP_NAME}

export VERSION=0.0.2
export IMG=quay.io/jooholee/${NEW_OP_NAME}:${VERSION}
export BUNDLE_IMG=quay.io/jooholee/${NEW_OP_NAME}-bundle:${VERSION}
~~~

### 2. Update manifests and metadata of CSV 
- Using  [ClusterServerVersion Markers](https://sdk.operatorframework.io/cs/building-operators/golang/references/markers/) for CRD
- This marker help customers to understand the operators
- File - `api/v1alpha1/nfsprovisioner_types.go`
  - Add default information (Operator name/Resources/Spec/Status)
    ~~~
    // +operator-sdk:csv:customresourcedefinitions:displayName="NFS Provisioner App",resources={{ServiceAccount,v1,nfs-provisioner},{SecurityContextConstraints,v1,nfs-provisioner},{Deployment,v1,nfs-provisioner},{PersistentVolumeClaim,v1,nfs-server},{ClusterRole,v1,nfs-provisioner-runner},{ClusterRoleBinding,v1,nfs-provisioner-runner},{Role,v1,leader-locking-nfs-provisioner},{RoleBinding,v1,leader-locking-nfs-provisioner},{Service,v1,nfs-provisioner},{StorageClass,v1,nfs}}

    type NFSProvisioner struct {
    ..
    }


    type NFSProvisionerSpec struct {
    	// HostPathDir is the direcotry where NFS server will use.
      // +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="HostPath directory",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:string", "urn:alm:descriptor:io.kubernetes:custom"}
      HostPathDir string `json:"hostPathDir,omitempty"`

      // PVC Name is the PVC resource that already created for NFS server.
      // Do not set StorageClass name with this param. Then, operator will fail to deploy NFS Server.
      // +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="PVC Name",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:string", "urn:alm:descriptor:io.kubernetes:custom"}
      Pvc string `json:"pvc,omitempty"`

      // StorageSize is the PVC size for NFS server.
      // By default, it sets 10G.
      // +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Storage Size",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:string", "urn:alm:descriptor:io.kubernetes:custom"}
      StorageSize string `json:"storageSize,omitempty"`

      // torageClass Name for NFS server will provide a PVC for NFS server.
      // Do not set PVC name with this param. Then, operator will fail to deploy NFS Server
      // +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="StorageClass Name for NFS server",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:string","urn:alm:descriptor:io.kubernetes:custom"}
      SCForNFSPvc string `json:"scForNFSPvc,omitempty"` //https://golang.org/pkg/encoding/json/

      // NFS server will be running on a specific node by NodeSeletor
      // +operator-sdk:csv:customresourcedefinitions:type=spec
      NodeSelector map[string]string `json:"nodeSelector,omitempty"`

      // StorageClass Name for NFS Provisioner is the StorageClass name that NFS Provisioner will use. Default value is `nfs`
      // +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="StorageClass Name for NFS Provisioner",xDescriptors={"urn:alm:descriptor:com.tectonic.ui:string","urn:alm:descriptor:io.kubernetes:custom"}
      SCForNFSProvisioner string `json:"scForNFS,omitempty"` //https://golang.org/pkg/encoding/json/
     ~~~

- Build a new operator image and push it
  ~~~
  make generate 
  make manifests
  make podman-build podman-push
  ~~~
  

- Using user part for Operator itself.
  - File - `./config/manifests/bases/${NEW_OP_NAME}.clusterserviceversion.yaml`
    ~~~
    description: This operator deploy NFS server with local storage and also provide provisioner for storageClass.
    displayName: NFS Provisioner Operator
    icon:
    - base64data: ""
      mediatype: ""

    installModes:
      - supported: false
        type: OwnNamespace
      - supported: false
        type: SingleNamespace
      - supported: false
        type: MultiNamespace
      - supported: true
        type: AllNamespaces
      keywords:
      - nfs
      - storage
      - provisioner
      links:
      - name: Nfs Provisioner Operator
        url: https://github.com/jooho/nfs-provisioner-operator
      maintainers:
      - email: ljhiyh@gmail.com
        name: jooho
      maturity: alpha
      provider:
        name: Jooho Lee
      version: 0.0.1
    ~~~

### 3. Update Bundle 
- Update your Operator and Bundle 
  - Here, we will change supported installMode 
  - File - `./config/manifests/bases/nfs-provisioner-operator.clusterserviceversion.yaml`
    ~~~
    installModes:
    - supported: true     <== changed 
      type: OwnNamespace
    - supported: false
      type: SingleNamespace
    - supported: false
      type: MultiNamespace
    - supported: true
      type: AllNamespaces
    ~~~

- Set Channel and DefaultChannel
  ~~~
  export CHANNELS=beta
  export DEFAULT_CHANNEL=beta
  ~~~

- Update bundle 
  ~~~
  make bundle 
  ~~~

- Check if `bundle/manifests/nfs-provisioner-operator.clusterserviceversion.yaml` is updated.
  ~~~
  ...
  replaces: test-nfs-provisioner-operator.v0.0.1
  version: 0.0.2
  ~~~

- Check if `./bundle/metadata/annotations.yaml` is updated
  ~~~
  operators.operatorframework.io.bundle.channel.default.v1: beta
  operators.operatorframework.io.bundle.channels.v1: beta
  ~~~

- Update bundle 
  ~~~
  make bundle-build 
  ~~~

- Push the bundle image
  ~~~
  make bundle-push  
  ~~~

- Upgrade your Operator
  ~~~
  No commands to test this now (2020.10.26)
  ~~~

### 4. Clean up
~~~
operator-sdk cleanup ${NEW_OP_NAME}
~~~


### 5. Test [InstallMode](https://github.com/operator-framework/operator-lifecycle-manager/blob/4197455/Documentation/design/building-your-csv.md#operator-metadata)
~~~
operator-sdk run bundle ${BUNDLE_IMG}

// Test install-mode
operator-sdk run bundle ${BUNDLE_IMG}  --verbose --install-mode OwnNamespace

// this will fail if csv set false for SingleNamespace
oc new-project test-op2
operator-sdk run bundle ${BUNDLE_IMG}  --verbose --install-mode SingleNamespace=test-op2
FATA[0002] Failed to run bundle: install mode type "SingleNamespace" not supported in CSV "nfs-provisioner-operator.v0.0.2"
~~~


### 5. Clean up
~~~
operator-sdk cleanup ${NEW_OP_NAME}
~~~




### Tip
image to base64
~~~
base64 -w 0 ${DEMO_HOME}/nfs-provisioner-tutorial-docs/images/128px-Human-folder-remote-nfs.svg.png
~~~
