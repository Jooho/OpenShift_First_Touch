# NFS Provisioner Operator OLM(Operator Lifecyle Management) - Bundle - 1
[OLM Summary](Tutorial-6-Operator-OLM-Summary.md#Bundle)

## Objects
- Install OLM or check OLM status on the cluster
- Create a bundle
- Try to use useful commands for Bundle
  - build/push/validate/run/cleanup



## Tutorial Flows
- Set environment variables for a new operator bundle
- Check OLM 
- If there is no OLM installed, install OLM
- Create a new operator bundle
- Build the new operator bundle
- Validate the new operator bundle
- Push the new operator bundle


## Pre-requisites

- Add the following command into `Makefile` for bundle push
  ~~~
  vi Makefile
  ...
  # Push the bundle image.
  .PHONY: bundle-push
  bundle-push:
          podman push $(BUNDLE_IMG)
  ~~~

## Steps

### 1. Set environment variables for a new bundle

~~~
export NEW_OP_NAME=test-nfs-provisioner-operator
export NEW_OP_HOME=${ROOT_HOME}/operator-projects/${NEW_OP_NAME}
export NAMESPACE=${NEW_OP_NAME}

export VERSION=0.0.1
export IMG=quay.io/jooholee/${NEW_OP_NAME}:${VERSION}
export BUNDLE_IMG=quay.io/jooholee/${NEW_OP_NAME}-bundle:${VERSION}
export CHANNELS=alpha
export DEFAULT_CHANNEL=alpha
~~~


### 2. Check OLM 
Openshift OLM namespace is `openshift-operator-lifecycle-manager`

~~~
operator-sdk olm status --olm-namespace="openshift-operator-lifecycle-manager"
~~~

### 3. If there is no OLM installed, install OLM
You can **NOT** install multiple OLM in a cluster
~~~
operator-sdk olm install  # failed
~~~

### 4. Create a new operator bundle
~~~
make bundle
...

Display name for the operator (required): 
> NFS Provisioner Operator

Description for the operator (required): 
> This operator deploy NFS server with local storage and also provide provisioner for storageClass.

Provider's name for the operator (required): 
> Jooho Lee

Any relevant URL for the provider name (optional): 
> github.com/nfs-provisioner-operator

Comma-separated list of keywords for your operator (required): 
> nfs,storage,provisioner                            

Comma-separated list of maintainers and their emails (e.g. 'name1:email1, name2:email2') (required): 
> jooho:ljhiyh@gmail.com
~~~

### 5. Build the new operator bundle
~~~
make bundle-build 
~~~

### 6. Push the new operator bundle
~~~
make bundle-push
~~~

### 7. Validate the new operator bundle(optional)
This validation command is one of `make bundle` commands so you can skip this

- Basic
~~~
operator-sdk bundle validate $BUNDLE_IMG -b podman
~~~
You can change build provider using `-b`

Since 1.1.0, you can set optional validator.
- New version
  ~~~
  operator-sdk bundle validate --list-optional
  operator-sdk bundle validate ./bundle --select-optional name=operatorhub
  ~~~


### 8. Deploy bundle and Check recources
~~~
operator-sdk run bundle ${BUNDLE_IMG}
~~~

- Check resources
  ~~~
  oc get catalogsource,operatorgroup
  
  //operator registry
  oc get pod
  ~~~



### 9. Clean up
For next tutorial, do not delete this.

~~~
operator-sdk cleanup ${NEW_OP_NAME}
~~~




