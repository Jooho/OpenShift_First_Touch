# Demo 1 - Base deployment

## Redis
oc apply -f <(istioctl kube-inject -f scripts/applications/deployment/redis.yaml )

## Client
cat scripts/applications/deployment/chat-client-configmap.yaml |sed "s/%CHAT_SERVER_HOSTNAME%/${CHAT_SERVER_HOSTNAME}/g" | oc apply -f -
oc apply -n chat -f <(istioctl kube-inject -f scripts/applications/deployment/chat-client-v1.yaml)
oc apply -n chat -f <(istioctl kube-inject -f scripts/applications/deployment/chat-client-android.yaml)
oc apply -n chat -f <(istioctl kube-inject -f scripts/applications/deployment/chat-client-apple.yaml)

## Server
oc apply -n chat -f <(istioctl kube-inject -f scripts/applications/deployment/chat-server-v1.yaml)

## Auth
oc apply -n chat -f <(istioctl kube-inject -f scripts/applications/deployment/auth.yaml)

## Ingress
oc expose  svc istio-ingressgateway --name chat-client-ingress --hostname="${CHAT_CLIENT_HOSTNAME}" -n istio-system
oc expose  svc istio-ingressgateway --name chat-server-ingress --hostname="${CHAT_SERVER_HOSTNAME}" -n istio-system

## istio 
### Access chat-client v1 only
oc apply -n chat -f scripts/istio/1-0.chat-gateway.yaml 
cat scripts/istio/1-1.vs-client-server.yaml | sed "s/%CHAT_CLIENT_HOSTNAME%/${CHAT_CLIENT_HOSTNAME}/g" | sed "s/%CHAT_SERVER_HOSTNAME%/${CHAT_SERVER_HOSTNAME}/g" | oc apply -f -
cat scripts/istio/1-2.destinationRule-chat-client.yaml | sed "s/%CHAT_CLIENT_HOSTNAME%/${CHAT_CLIENT_HOSTNAME}/g" | oc apply -f -
cat scripts/istio/1-3.destinationRule-chat-server.yaml | sed "s/%CHAT_SERVER_HOSTNAME%/${CHAT_SERVER_HOSTNAME}/g" | oc apply -f -

clear

oc get pod -w 
