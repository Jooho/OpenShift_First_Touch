#Init Container

## Description
You run init containers in the same pod as your application container to create the environment your application requires or to satisfy any preconditions the application might have. You can run utilities that you would otherwise need to place into your application image. You can run them in different file system namespaces (view of the same file system) and offer them different secrets than your application container.

Init containers run to completion and each container must finish before the next one starts. The init containers will honor the restart policy. Leverage initContainers in the podspec.

## Sample
```
$ cat init-containers.yaml
apiVersion: v1
kind: Pod
metadata:
  name: init-loop
spec:
  containers:
  - name: nginx
    image: nginx
    ports:
    - containerPort: 80
    volumeMounts:
    - name: workdir
      mountPath: /usr/share/nginx/html
  initContainers:
  - name: init
    image: centos:centos7
    command:
    - /bin/bash
    - "-c"
    - "while :; do sleep 2; echo hello init container; done"
  volumes:
  - name: workdir
    emptyDir: {}
```

```
$ oc get -f init-containers.yaml
NAME        READY     STATUS     RESTARTS   AGE
nginx       0/1       Init:0/1   0          6m
```
## Demo

## Reference
- [init-containers](https://docs.openshift.com/container-platform/3.6/architecture/core_concepts/containers_and_images.html#init-containers)
- [pods-services-init-containers](https://docs.openshift.com/container-platform/3.6/architecture/core_concepts/pods_and_services.html#pods-services-init-containers)
