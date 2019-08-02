Example Pipeline
--------------------

## Reference

### Install TKN CLI
~~~
curl -L https://github.com/tektoncd/cli/releases/download/v0.2.0/tkn_0.2.0_Linux_x86_64.tar.gz |tar xvz
~~~


### [Install Tekton on OpenShift](https://github.com/tektoncd/pipeline/blob/master/docs/install.md)

```
oc login -u system:admin
oc new-project tekton-pipelines
oc adm policy add-scc-to-user anyuid -z tekton-pipelines-controller
oc apply --filename https://storage.googleapis.com/tekton-releases/latest/release.yaml
oc get pods --namespace tekton-pipelines --watch
```

### [Official Tutorial](https://github.com/tektoncd/pipeline/blob/master/docs/tutorial.md)

### [Red Hat Tutorial Blog](https://blog.openshift.com/cloud-native-ci-cd-with-openshift-pipelines/)
## [Red Hat Tutorial Github](https://github.com/openshift/pipelines-tutorial)


[![Red Hat Tutorial Youtube](http://img.youtube.com/vi/pMDiiW1UqLo/0.jpg)](https://www.youtube.com/embed/pMDiiW1UqLo)





## Demo

### Pre-requisites
~~~
mkdir /sysroot/exports-nfs

sudo chcon -Rt svirt_sandbox_file_t  /sysroot/exports-nfs

oc label node ip-10-0-147-211.us-east-2.compute.internal  app=nfs-provisioner

oc process -f  ./template.yml -e HOST_PATH=/sysroot/exports-nfs -
~~~



### Initialize demo
```
oc new-project pipelines-tutorial
oc create serviceaccount pipeline
oc adm policy add-scc-to-user privileged -z pipeline
oc adm policy add-role-to-user edit -z pipeline

oc create -f https://raw.githubusercontent.com/openshift/pipelines-tutorial/master/resources/petclinic.yaml

```

### Install Tasks
~~~
oc create -f https://raw.githubusercontent.com/tektoncd/catalog/master/openshift-client/openshift-client-task.yaml
oc create -f https://raw.githubusercontent.com/openshift/pipelines-catalog/master/s2i-java-8/s2i-java-8-task.yaml
~~~



### Install Pipeline 
```
oc create -f https://raw.githubusercontent.com/openshift/pipelines-tutorial/master/resources/petclinic-deploy-pipeline.yaml
```


### Install PipelineResource
```
oc create -f https://raw.githubusercontent.com/openshift/pipelines-tutorial/master/resources/petclinic-resources.yaml
```

### Trigger pipeline (Install PipelineRun)
```
https://raw.githubusercontent.com/openshift/pipelines-tutorial/master/resources/petclinic-deploy-pipelinerun.yaml
```