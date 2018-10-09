# 0 Demo project
oc new-project chat --display-name="Toronto Convergency Event Test Project - Chat Program" 

## For istio deployment as a sidecar, privileged permission is needed
oc adm policy add-scc-to-user privileged -z default,deployer -n chat

oc delete limitrange --all


