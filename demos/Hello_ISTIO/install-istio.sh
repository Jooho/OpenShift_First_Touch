##ISTIO
curl -L https://github.com/istio/istio/releases/download/1.0.0/istio-1.0.0-linux.tar.gz | tar xz

cd istio-1.0.0
export ISTIO_HOME=`pwd`
export PATH=$ISTIO_HOME/bin:$PATH

oc apply -f install/kubernetes/helm/istio/templates/crds.yaml
sleep 1
oc apply -f install/kubernetes/istio-demo.yaml
oc project istio-system
oc adm policy add-scc-to-group anyuid system:authenticated

oc expose svc istio-ingressgateway
oc expose svc servicegraph
oc expose svc grafana
oc expose svc prometheus
oc expose svc tracing
