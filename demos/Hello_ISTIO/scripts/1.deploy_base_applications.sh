## Client
oc new-app docker.io/ljhiyh/chat-client:v1
oc set probe dc/chat-client --readiness --get-url=http://:4200/ --initial-delay-seconds=10 --timeout-seconds=60
oc scale dc chat-client --replicas=1
oc expose svc chat-client


## Server
oc new-app docker.io/ljhiyh/chat-server:noredis
oc set probe dc/chat-server --readiness --get-url=http://:8080/healthz --initial-delay-seconds=10 --timeout-seconds=60
oc expose svc chat-server --hostname="chat-server-example-chat.apps.toronto.openshiftworkshop.com"


## Auth
oc new-app docker.io/ljhiyh/auth
oc set probe dc/auth --readiness --get-url=http://:8080/hello --initial-delay-seconds=10 --timeout-seconds=60
