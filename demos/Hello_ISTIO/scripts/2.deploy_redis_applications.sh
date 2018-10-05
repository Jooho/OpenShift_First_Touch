## Redis
oc new-app --template=redis-ephemeral -p REDIS_PASSWORD=redis

sleep 5

## Server
oc new-app docker.io/ljhiyh/chat-server:redis
oc set probe dc/chat-server --readiness --get-url=http://:8080/healthz --initial-delay-seconds=10 --timeout-seconds=60
#oc expose svc chat-server --hostname="chat-server-example-chat.apps.toronto.openshiftworkshop.com"


