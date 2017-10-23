# Init Container

## Description
You run init containers in the same pod as your application container to create the environment your application requires or to satisfy any preconditions the application might have. You can run utilities that you would otherwise need to place into your application image. You can run them in different file system namespaces (view of the same file system) and offer them different secrets than your application container.

Init containers run to completion and each container must finish before the next one starts. The init containers will honor the restart policy. Leverage initContainers in the podspec.

## Demo Scenario(from kubernetes)
Sometimes, another hostname should be resolvable before starting container. In this case, we want to validate it then if fails, the container haven't to be deployed.

With init container, it is much easier because init container will start before main container start.


*1. Create pod has 2 init contaiers that wait until SVC hostname is resolved.*
```
echo "
apiVersion: v1
kind: Pod
metadata:
  name: myapp-pod
  labels:
    app: myapp
spec:
  containers:
  - name: myapp-container
    image: busybox
    command: ['sh', '-c', 'echo The app is running! && sleep 3600']
  initContainers:
  - name: init-myservice
    image: busybox
    command: ['sh', '-c', 'until nslookup myservice; do echo waiting for myservice; sleep 2; done;']
  - name: init-mydb
    image: busybox
    command: ['sh', '-c', 'until nslookup mydb; do echo waiting for mydb; sleep 2; done;']" |oc create -f -
```

*2. Check init container*
As command in init container say, it tries to resolve myservice.

```
oc describe pod myapp-pod

...
Init Containers:
  init-myservice:
    Container ID:	docker://f7d759ffd757d517e62ca904d21e4577d260a8f392773558ef57008fc9606638
    ...
  init-mydb:
    Container ID:	
....

oc logs -f myapp-pod -c init-myservice

waiting for myservice
nslookup: can't resolve 'myservice'
Server:    10.10.181.196
Address 1: 10.10.181.196 dhcp181-196.gsslab.rdu2.redhat.com
nslookup: can't resolve 'myservice'
```

*3. Create SVC object*
```
echo "
kind: Service
apiVersion: v1
metadata:
  name: myservice
spec:
  ports:
  - protocol: TCP
    port: 80
    targetPort: 9376
---
kind: Service
apiVersion: v1
metadata:
  name: mydb
spec:
  ports:
  - protocol: TCP
    port: 80
    targetPort: 9377" |oc create -f -

```

*4. Check init container again*
```
oc logs -f myapp-pod -c init-myservice

...
waiting for myservice
Server:    10.10.181.196
Address 1: 10.10.181.196 dhcp181-196.gsslab.rdu2.redhat.com
Name:      myservice
Address 1: 172.30.28.98 myservice.init-container.svc.cluster.local


```


*5. After the pre-condition meet, the main container will start.*
```
oc get pod

NAME        READY     STATUS    RESTARTS   AGE
myapp-pod   1/1       Running   0          14m
[root@dhcp181-97 init_container]# 

```

## Reference
- [init-containers](https://docs.openshift.com/container-platform/3.6/architecture/core_concepts/containers_and_images.html#init-containers)
- [pods-services-init-containers](https://docs.openshift.com/container-platform/3.6/architecture/core_concepts/pods_and_services.html#pods-services-init-containers)
