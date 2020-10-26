# OLM (Operator Lifecyle Management) Summary

https://sdk.operatorframework.io/docs/olm-integration/

## [OLM Commands](https://sdk.operatorframework.io/docs/olm-integration/cli-overview/)
- Status
    ~~~
    operator-sdk olm status --olm-namespace="openshift-operator-lifecycle-manager"
    ~~~

- Install     
    ~~~
    operator-sdk olm install
    ~~~
    **(NOTE)** You can *NOT* install multiple OLM in a cluster

- Uninstall     
    ~~~
    operator-sdk olm install
    ~~~

## Bundle
**(NOTE)** OpenShift OLM project is `openshift-operator-lifecycle-manager`


### Create Bundle
- Pre-requisites
    ~~~
    export VERSION=0.0.1
    export IMG=quay.io/jooholee/nfs-provisioner-operator:${VERSION}
    export BUNDLE_IMG=quay.io/jooholee/nfs-provisioner-operator-bundle:${VERSION}
    ~~~

- Create bundle using make
    ~~~
    make bundle
    ....

    Display name for the operator (required): 
    > NFS Provisioner Operator

    Description for the operator (required): 
    > This operator deploy NFS server with local storage and also provide provisioner for storageClass.

    Provider's name for the operator (required): 
    > Jooho Lee

    Any relevant URL for the provider name (optional): 
    > 

    Comma-separated list of keywords for your operator (required): 
    > nfs,storage,pv provisioner                            

    Comma-separated list of maintainers and their emails (e.g. 'name1:email1, name2:email2') (required): 
    > jooho:ljhiyh@gmail.com
    ...
    ~~~

### Build/Push Bundle
- (Optional) Add a Make target for bundle push
    ~~~
    vi Makefile
    ...
    # Push the bundle image.
    .PHONY: bundle-push
    bundle-push:
            podman push $(BUNDLE_IMG)
    ~~~
~~~
make bundle-build bundle-push
~~~




### Bundle Commands
  - bundle build (standard)
    ~~~ 
    make bundle-build quay.io/jooholee/nfs-provisioner-operator-bundle
    ~~~
  - bundle push (standard)
    ~~~
    make bundle-push quay.io/jooholee/nfs-provisioner-operator-bundle
    ~~~
  - bundle validate (optional)
    This validation command is one of `make bundle` commands 
    ~~~
    operator-sdk bundle validate quay.io/jooholee/nfs-provisioner-operator-bundle -b podman
    ~~~
    - bundle validation list
      ~~~
      operator-sdk bundle validate --list-optional
      ~~~
    - Select bundle validation
      ~~~
      operator-sdk bundle validate ./bundle --select-optional name=operatorhub
      ~~~
    - Trouble-shooting
      - Icon
        ~~~
        ERRO[0002] Error: Value : (nfs-provisioner-operator.v0.0.3) csv.Spec.Icon elements should contain both data and mediatype 
        ~~~
        - Solution
          - Add icon and mediatype :)

### Package Manifests
packageManifest is deprecated



### Generating Manifests and Metadata

- CSV
  - Update type file for CSV update using [ClusterServiceVersion Markers](https://sdk.operatorframework.io/docs/building-operators/golang/references/markers/)

  - Generate manifests and CSV
    ~~~
    make manifests
    make bundle
    ~~~

  - Check CSV file: `bundle/manifests/nfs-provisioner-operator.clusterserviceversion.yaml`

- Channel
  - If you want to change channel, you can set the following:
    ~~~
    export CHANNELS="stable"  //default "alpha"

    export BUNDLE_DEFAULT_CHANNEL   // default "alpha"
    ~~~


### Update operator version
Let’s say you added a new API App with group app and version v1alpha1 to your Operator project, and added a port to your manager Deployment in config/manager/manager.yaml.
~~~
export VERSION=0.0.4
export IMG=quay.io/jooholee/nfs-provisioner-operator:${VERSION}
export BUNDLE_IMG=quay.io/jooholee/nfs-provisioner-operator-bundle:${VERSION}

make podman-build podman-push  # for operator image

make bundle 

make bundle-build bundle-push   # for bundle image
~~~

- `bundle/manifests/nfs-provisioner-operator.clusterserviceversion.yaml` is updated.
  ~~~
  ...
  replaces: nfs-provisioner-operator.v0.0.3
  version: 0.0.4
  ~~~


- [User fields should not be overwritten](https://sdk.operatorframework.io/docs/olm-integration/generation/#csv-fields)
  - icon data was deleted. (bug)



### Upgrade operator

~~~
export VERSION=0.0.5
export IMG=quay.io/jooholee/nfs-provisioner-operator:${VERSION}
export BUNDLE_IMG=quay.io/jooholee/nfs-provisioner-operator-bundle:${VERSION}
export CHANNELS=beta
export DEFAULT_CHANNEL=beta

make podman-build podman-push  # for operator image

make bundle 

make bundle-build bundle-push  # for bundle image

~~~

- `bundle/manifests/nfs-provisioner-operator.clusterserviceversion.yaml` is updated.
  ~~~
  ...
  replaces: nfs-provisioner-operator.v0.0.4
  version: 0.0.5
  ~~~

- `./bundle/metadata/annotations.yaml` is updated
  ~~~
  operators.operatorframework.io.bundle.channel.default.v1: beta
  operators.operatorframework.io.bundle.channels.v1: beta
  ~~~
  
### Test Operator by OLM

[Official doc](https://sdk.operatorframework.io/docs/olm-integration/testing-deployment/) does not explain about bundle way yet(2020/10/20)

- Create a test project
  ~~~
  oc new-project nfs-provisioner-operator
  ~~~
- CSV - ` bundle/manifests/nfs-provisioner-operator.clusterserviceversion.yaml`
  - [InstallMode](https://github.com/operator-framework/operator-lifecycle-manager/blob/4197455/Documentation/design/building-your-csv.md#operator-metadata) is checking the supportability for OperatorGroup. 

  - OperatorGroup can use :
    - `OwnNamespace`
      - The operator can be a member of an OperatorGroup that selects its own namespace
      - Present namespace
    - `SingleNamespace`
      - The operator can be a member of an OperatorGroup that selects one namespace
      - Choose 1 namespace
    - `MultiNamespace`
      - The operator can be a member of an OperatorGroup that selects more than one namespace
      - Choose more than 1 namespaces
    - `AllNamespaces`
      - The operator can be a member of an OperatorGroup that selects all namespaces (target namespace set is the empty string "")
      - All namespaces
  ~~~
  installModes:
  - supported: false
    type: OwnNamespace
  - supported: false
    type: SingleNamespace
  - supported: false
    type: MultiNamespace
  - supported: true
    type: AllNamespaces
  ~~~
  - How to update `installModes`
    - Update `./config/manifests/bases/nfs-provisioner-operator.clusterserviceversion.yaml`
      ~~~
      installModes:
      - supported: true
        type: OwnNamespace
      - supported: false
        type: SingleNamespace
      - supported: false
        type: MultiNamespace
      - supported: true
        type: AllNamespaces
      ~~~
    - Build a new version
      ~~~
      export VERSION=0.0.6
      export IMG=quay.io/jooholee/nfs-provisioner-operator:${VERSION}
      export BUNDLE_IMG=quay.io/jooholee/nfs-provisioner-operator-bundle:${VERSION}
      export CHANNELS=beta
      export DEFAULT_CHANNEL=beta

      make podman-build podman-push  # for operator image

      make bundle 

      make bundle-build bundle-push  # for bundle image
      ~~~
    - Without changing version, CSV changes are not applied correctly
      - It might be a bug
  
- Deploy
  ~~~
  operator-sdk run bundle ${BUNDLE_IMG}

  // Test install-mode
  operator-sdk run bundle ${BUNDLE_IMG}  --verbose --install-mode OwnNamespace

  // this will fail if csv set false for SingleNamespace
  operator-sdk run bundle ${BUNDLE_IMG}  --verbose --install-mode SingleNamespace=test-op2
  ~~~

- Cleanup
  ~~~
  operator-sdk cleanup nfs-provisioner-operator --verbose
  ~~~


## Operator Registry Tooling
[Official Doc](https://github.com/operator-framework/operator-registry/blob/master/docs/design/opm-tooling.md)


### Concept
OPM tool mainly has 2 commands(registry/index).

- registry
  - create a db file (`test-registry.db`)
  - push a bundle (`0.0.1`) to the db file(`test-registry.db`)
    - push a new bundle (`0.0.2`) to the db file(`test-registry.db`)
  - deploy registry 
- index
  - create an tag image that contains db file.
  - push a bundle (`0.0.1`) to a db file and create an new tag image (storage-index:0.0.1)
    - push a new bundle (`0.0.2`)
      - download the db file from previous tag image(`storage-index:0.0.1`) 
      - push the new bundle (`0.0.2`) to the previous db file 
      - create an new tag image (`storage-index:0.0.2`) with the updated db file
  - `run an tag image` == `deploy registry`


### OPM Commands

- Registry
  - Add a bundle / Add a new bundle
      ~~~
      # First
      opm registry add -b quay.io/jooholee/nfs-provisioner-operator-bundle:0.0.1 -d "test-registry.db"

      # Second
      opm registry add -b quay.io/jooholee/nfs-provisioner-operator-bundle:0.0.2 -d "test-registry.db"
      ~~~
  - Delete bundles 
    - Delete entire bundles (not specify specific bundle at the moment-2020/10/21)
      ~~~
      opm registry rm -o "nfs-provisioner-operator" -d "test-registry.db"
      ~~~
  - Prune bundles 
    - Prune all bundles except specified bundle. (opposite of rm)
      ~~~
      opm registry prune -p "nfs-provisioner-operator" -d "test-registry.db"
      ~~~
  - Start Operator Registy 
    ~~~
    opm registry serve -d "test-registry.db" -p 50051
    ~~~

- Index
  - Add a bundle
    ~~~
    opm index add --bundles nfs-provisioner-operator-bundle:0.0.1 --tag quay.io/jooholee/storage-index:0.0.1
    ~~~
  - Add a new bundle
    ~~~
    opm index add --bundles nfs-provisioner-operator-bundle:0.0.2 --from-index quay.io/jooholee/storage-index:0.0.1 --tag quay.io/jooholee/storage-index:0.0.2
    ~~~
  - Add a new custom bundle using generated dockerfile
    ~~~
    opm index add --bundles nfs-provisioner-operator-bundle:0.0.3 --generate --out-dockerfile "my.Dockerfile" -p podman
    
    podman build -t nfs-provisioner-operator-bundle:0.0.3 my.Dockerfile
    ~~~
  - Delete bundles 
    - Delete entire bundles (not specify specific bundle at the moment-2020/10/21)
      ~~~
      opm index rm --operators nfs-provisioner-operator --tag quay.io/jooholee/storage-index:0.0.3
      ~~~
  - Prune bundles
    - Prune all bundles except specified bundle. (opposite of rm)
      ~~~
      opm index prune -p "nfs-provisioner-operator" --from-index quay.io/jooholee/storage-index:0.0.3--tag quay.io/jooholee/storage-index:$0.0.4
      ~~~
  - Export bundles
    ~~~
    opm index export --index=quay.io/jooholee/storage-index:0.0.4 --package="nfs-provisioner-operator" -c "podman"
    ~~~
  - Specify opm version
    ~~~
    opm index add --bundles nfs-provisioner-operator-bundle:0.0.2 --from-index quay.io/jooholee/storage-index:0.0.1 --tag quay.io/jooholee/storage-index:0.0.2 --binary-image quay.io/$user/my-opm-source
    ~~~
  - Start Operator Registy 
    ~~~
    podman run -d --name operator-registry quay.io/jooholee/storage-index:0.0.4 
    ~~~





