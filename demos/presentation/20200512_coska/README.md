# COSKA Presentation

Topic: Jay's Container History ... and OpenShift 4

Presentation Doc: 

Where: *Online*


## Demos

### NFS Provisioner Operator 

*Create a base directory*
```
cd /tmp
mkdir COSKA
cd COSKA
```

*Clone Git repository*
```
git clone https://github.com/Jooho/jhouse_openshift.git
```

*Update target node lavel*
```
oc label node $NODE_NAME app=nfs-provisioner 

(ex) oc label node ip-10-0-129-143.ec2.internal app=nfs-provisioner 
```

*Deploy NFS provisioner operator pod*
```
cd jhouse_openshift/demos/Operator/helm/nfs-provisioner-operator/

oc new-project nfs-provisioner-operator
oc create -f deploy/crds/jhouse_v1alpha1_nfs_crd.yaml 
oc create -f deploy/service_account.yaml
oc create -f deploy/role.yaml

# This NFS Provisioner Operator is from a tutorial of another demos, so you need to update the namespace and image
export OPERATOR_NAMESPACE=$(oc config view --minify -o jsonpath='{.contexts[0].context.namespace}')
sed -i "s|REPLACE_NAMESPACE|$OPERATOR_NAMESPACE|g" deploy/role_binding.yaml

oc create -f deploy/role_binding.yaml

sed -i 's|REPLACE_IMAGE|quay.io/jooholee/nfs-provisioner-operator:v0.0.1|g'  deploy/operator.yaml
oc create -f deploy/operator.yaml
```

*Deploy NFS Provisioner*
```
oc create -f deploy/crds/jhouse_v1alpha1_nfs_cr.yaml

oc logs $NFS_Provisioner_Operator_Pod

oc project nfs-provisioner
oc logs $NFS_Provisioner_Pod
```

*Make the NFS provisioenr as a default Storageclass*
```
oc patch storageclass nfs-storageclass  -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
oc patch storageclass gp2  -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}'
```

*Check PVC*

Test PVC is alreay there because the helm operator has test.

(Just in case, this is sample script to create a pvc using the nfs sc)
```
echo "# Source: nfs-provisioner/templates/tests/test-pvc.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
 name: test-pvc
 namespace: default
spec:
 accessModes:
  - ReadWriteMany
 resources:
   requests:
     storage: 1Gi
 storageClassName: nfs-storageclass"|oc create -f -
```

*Delete NFS Service Object*
```
oc delete service nfs-provisioner -n nfs-provisioner

oc get svc -w
```


**Clean Up**
```
cd jhouse_openshift/demos/Operator/helm/nfs-provisioner-operator/
oc project nfs-provisioner-operator
oc delete -f deploy/crds/jhouse_v1alpha1_nfs_cr.yaml 
oc delete -f deploy/operator.yaml
oc delete -f deploy/role_binding.yaml
oc delete -f deploy/role.yaml
oc delete -f deploy/service_account.yaml
oc delete -f deploy/crds/jhouse_v1alpha1_nfs_crd.yaml 
oc delete project nfs-provisioner-operator nfs-provisioner
```


### Service Mesh


*Deploy sample application `bookInfo`*
```
oc new-project bookinfo

oc apply -n bookinfo -f https://raw.githubusercontent.com/Maistra/istio/maistra-1.1/samples/bookinfo/platform/kube/bookinfo.yaml
oc get pod
oc get svc
oc get route
```

*Create istio objects*
```
oc apply -n bookinfo -f https://raw.githubusercontent.com/Maistra/istio/maistra-1.1/samples/bookinfo/networking/bookinfo-gateway.yaml
oc apply -n bookinfo -f https://raw.githubusercontent.com/Maistra/istio/maistra-1.1/samples/bookinfo/networking/destination-rule-all.yaml

```

*Add bookinfo project to ServiceMesh Member*
```
oc -n istio-system patch --type='json' smmr default -p '[{"op": "add", "path": "/spec/members", "value":["'"bookinfo"'"]}]'    # oc edit ServiceMeshMemberRoll default -n istio-system
```

*Check productPage*
```
export GATEWAY_URL=$(oc -n istio-system get route istio-ingressgateway -o jsonpath='{.spec.host}')
curl -o /dev/null -s -w "%{http_code}\n" http://$GATEWAY_URL/productpage
```

*Enable auto route creation feature*
```
 oc -n istio-system patch smcp --type='json'  basic-install  -p '[{"op": "replace", "path": "/spec/istio/gateways/istio-ingressgateway/ior_enabled", "value": true}]'
```

*Recreate GateWay of bookinfo*
```
oc delete -n bookinfo -f https://raw.githubusercontent.com/Maistra/istio/maistra-1.1/samples/bookinfo/networking/bookinfo-gateway.yaml

echo "apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: bookinfo-gateway
spec:
  selector:
    istio: ingressgateway # use istio default controller
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - 'coska.apps.coska.jlee.rhcee.support'
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: bookinfo
spec:
  hosts:
  - '*'
  gateways:
  - bookinfo-gateway
  http:
  - match:
    - uri:
        exact: /productpage
    - uri:
        prefix: /static
    - uri:
        exact: /login
    - uri:
        exact: /logout
    - uri:
        prefix: /api/v1/products
    route:
    - destination:
        host: productpage
        port:
          number: 9080"|oc create -f -

oc get route -n istio-system
```

*Kiali console*
- Request (traffic animation)


*Jason Review2*
```
echo "apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: reviews
spec:
  hosts:
  - reviews
  http:
  - match:
    - headers:
        end-user:
          exact: jason
    route:
    - destination:
        host: reviews
        subset: v2
  - route:
    - destination:
        host: reviews
        subset: v1" |oc create -f -

```

**Clean Up**

```
oc delete project bookinfo
oc -n istio-system patch --type='json' smmr default -p '[{"op": "remove", "path": "/spec/members", "value":["'"bookinfo"'"]}]'
```


