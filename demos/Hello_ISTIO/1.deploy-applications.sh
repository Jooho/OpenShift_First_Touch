#0 Demo project
oc new-project chat --display-name="Toronto Convergency Event Test Project - Chat Program" 
oc adm policy add-scc-to-user privileged -z default,deployer -n chat
oc delete limitrange --all

# Demo 1 - Base deployment

## Redis
oc apply -f <(istioctl kube-inject -f scripts/applications/deployment/redis.yaml )

## Client
oc apply -n chat -f <(istioctl kube-inject -f scripts/applications/deployment/chat-client-v1.yaml)
oc apply -n chat -f <(istioctl kube-inject -f scripts/applications/deployment/chat-client-android.yaml)
oc apply -n chat -f <(istioctl kube-inject -f scripts/applications/deployment/chat-client-apple.yaml)

## Server
oc apply -n chat -f <(istioctl kube-inject -f scripts/applications/deployment/chat-server-v1.yaml)

## Auth
oc apply -n chat -f <(istioctl kube-inject -f scripts/applications/deployment/auth.yaml)

## Ingress
oc expose  svc istio-ingressgateway --name chat-client-ingress --hostname="chat-client.apps.toronto.openshiftworkshop.com" -n istio-system
oc expose  svc istio-ingressgateway --name chat-server-ingress --hostname="chat-server-example-chat.apps.toronto.openshiftworkshop.com" -n istio-system

## istio 
### Access chat-client v1 only
oc apply -n chat -f scripts/istio/1-0.chat-gateway.yaml 
oc apply -n chat -f scripts/istio/1-1.vs-client-server.yaml  
oc apply -n chat -f scripts/istio/1-2.destinationRule-chat-client.yaml 
oc apply -n chat -f scripts/istio/1-3.destinationRule-chat-server.yaml 


clear

oc get pod -w 
