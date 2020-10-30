# Create Infra node

### Create a test MC for infra
~~~
oc create -f ./infra.mc.yaml
~~~

### Create a infra MCP
~~~
oc create -f ./infra.mcp.yaml
~~~

### Add infra label and remove worker label from a node
~~~
oc label node $NODE_NAME  node-role.kubernetes.io/infra= 
oc label node $NODE_NAME  node-role.kubernetes.io/worker- 
~~~

