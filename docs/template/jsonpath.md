# jsonpath Cheat Sheet


### Syntax

- Single string
  - Example
    *Gather clusterIP from service object*
    ```
    oc -n default get svc router -o jsonpath='{.spec.clusterIP}{"\n"}'
    ```
    - Result
      ```
      172.30.47.131
      ```
 
- List
  - Example 1
    *Gather all pod names*
    ```
    oc get pod -o jsonpath='{ .items[*].metadata.name }'
    ``` 
    - Result 
      ```
      master-api-master1.example.com master-controllers-master1.example.com master-etcd-master1.example.com tiller-deploy-76cc4d8dd7-q4xkq
      ```

  - Example 2 (using `range`)
    *Gather all pod names*
    ```
    oc get pod -o jsonpath='{ range .items[*]} {.metadata.name}{end}'
    ```
    - Result 
      ```
      master-api-master1.example.com master-controllers-master1.example.com master-etcd-master1.example.com tiller-deploy-76cc4d8dd7-q4xkq
      ```

- Conditional
  - Example
    *Gather single string(`STATS_PASSWORD`) from list(`env`)*
    ```
    oc get dc/router -o jsonpath='{ .spec.template.spec.containers[0].env[?(@.name=="STATS_PASSWORD")].value}'
    ```
    - Result
      ```
      svnApw1zfV
      ```

- Conditional + List
  - Example 1
    *Gather pod name that are running only*
    ```
     oc -n default get pod -o jsonpath='{range .items[?(@.status.phase=="Running")]}{.metadata.name}{"\n"}{end}' 
    ```
    - Result  
      ```
      docker-registry-1-mvz4s
      dockergc-dm5jv
      ...
  - Example 2 
    *Gather specific container's mountPath*
    ```
    oc get dc/router -o jsonpath='{ .spec.template.spec.containers[?(@.name=="router")].volumeMounts[?(@.name=="server-certificate")].mountPath }'
    ```
    - Result
      ```
      /etc/pki/tls/private
      ```


## Tip
When you want to use Conditional with list, you should not iterate the list in range.

```
# Bad Example
 oc get dc/router -o jsonpath='{ range .spec.template.spec.containers[0].env[*] } {[?(@.name == "STATS_PASSWORD")] @.value}{end}'

# Good Example
 oc get dc/router -o jsonpath='{ range .spec.template.spec.containers[0].env } {[?(@.name == "STATS_PASSWORD")] @.value}{end}'
 ## Same result
 oc get dc/router -o jsonpath='{ range .spec.template.spec.containers[0].env } {[?(@.name == "STATS_PASSWORD")].value}{end}'
```

The reason why the first script have issue is because it returns `map[...]map[...]`. However, the second script returns `[map[...]map[...]map[...]]`

In order to fix the issue, the condition statement should be included in the first range.

```
oc get dc/router -o jsonpath='{ range .spec.template.spec.containers[0].env[?(@.name == "STATS_PASSWORD")]} { @.value} {end}'
```

However, if you want to gather all names, the first script is a good option

```
oc get dc/router -o jsonpath='{ range .spec.template.spec.containers[0].env[*] }{@.value} {end}'

# Result
/etc/pki/tls/private /etc/pki/tls/private/tls.crt     false   ...
```
