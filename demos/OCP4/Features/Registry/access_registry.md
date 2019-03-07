# Access integrated registry

In OpenShift v4, we can not access to the node directly by default,which means, we can not use service ip.

In order to access internal registry, we normally use 2 hostname: service, route. However, service is not an option by default.

So I tried to use route hostname but I failed with weird error messages. Hence, I will use port-forward here.

## Steps
- Login to Openshift
```
oc login
```

- Enable port-forward 
```
oc get pod 
NAME                                               READY     STATUS    RESTARTS   AGE
cluster-image-registry-operator-6c97bfc564-mr8r8   1/1       Running   1          46h

oc port-forward image-registry-79f9f4749b-7v9w5  5000:5000

```
- Login to registry
```
podman login -u joe -p $(oc whoami -t) localhost:5000 --tls-verify=false
```

- Pull/Tag/Push test image
```
$ podman pull docker.io/busybox
Trying to pull docker.io/busybox...Getting image source signatures
Copying blob sha256:697743189b6d255069caf6c455be10c7f8cae8076c6f94d224ae15cd41420e87
 738.18 KB / 738.18 KB [====================================================] 0s
Copying config sha256:d8233ab899d419c58cf3634c0df54ff5d8acc28f8173f09c21df4a07229e1205
 1.46 KB / 1.46 KB [========================================================] 0s
Writing manifest to image destination
Storing signatures
d8233ab899d419c58cf3634c0df54ff5d8acc28f8173f09c21df4a07229e1205


$ podman tag docker.io/busybox localhost:5000/openshift/busybox


$ podman push localhost:5000/openshift/busybox --tls-verify=false
Getting image source signatures
Copying blob sha256:adab5d09ba79ecf30d3a5af58394b23a447eda7ffffe16c500ddc5ccb4c0222f
 1.35 MB / 1.35 MB [========================================================] 1s
Copying config sha256:d8233ab899d419c58cf3634c0df54ff5d8acc28f8173f09c21df4a07229e1205
 1.46 KB / 1.46 KB [========================================================] 1s
Writing manifest to image destination
Storing signatures

```
