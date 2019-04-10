
1. Get single string (Gather clusterIP from service object)
```
oc -n default get svc router -o jsonpath='{.spec.clusterIP}{"\n"}'
```
2. Get single string from list
*jsonpath type*
```
oc get pod --output='jsonpath={.items[*].metadata.name}'
```
*go-template type*
```
oc get pod  --output='go-template={{ range .items}}{{.metadata.name}}{{end}}' 
```

*template string*
```
oc get pod --template='{{ range .items}}{{.metadata.name}}{{end}}' 
```

3. Get the word that match your condition( Gather pod names that are not running state in a project)
```
oc get pods -o jsonpath='{.items[?(@.status.phase!="Running")].metadata.name}'
```

4. Get the word that match your condition by go-template (Gather pod names that are not running state in all projects). 
```
oc get pod --all-namespaces --template='{{ range $pod := .items}}{{if ne $pod.status.phase "Running"}} {{$pod.metadata.name}} {{"\n"}}{{end}}{{end}}' 
```

5. Clean all pods that are not running state.
*json-path style:*
```
oc get pods -o jsonpath='{.items[?(@.status.phase!="Running")].metadata.name}' |xargs oc delete pod
```
*go-template style:*
```
oc get pod --all-namespaces --template='{{ range $pod := .items}}{{if ne $pod.status.phase "Running"}} oc delete pod -n {{$pod.metadata.namespace}} {{$pod.metadata.name}} {{"\n"}}{{end}}{{end}}' | bash -
```

6. Gather pod information which are not running state.
```
oc get pods --all-namespaces --template='
 {{- range .items -}} 
   {{- $pod_name:=.metadata.name -}}
   {{- $pod_namespace:=.metadata.namespace -}} 
   {{- if  ne .status.phase "Running" -}}  
       **namespace: {{ $pod_namespace}} **pod: {{ $pod_name }} **Reason: 
         {{- if .status.reason -}} 
             {{- .status.reason -}}
         {{- else if .status.containerStatuses -}}
             {{- range $containerStatus:=.status.containerStatuses -}}
                 {{- if $containerStatus.state.waiting -}}
		            {{- $containerStatus.state.waiting.reason -}}				
   		 {{- else if $containerStatus.state.terminated -}}
                     {{- $containerStatus.state.terminated.reason -}}
                  {{- end -}}
             {{- end -}}
         {{- else -}}
             {{- range $condition:=.status.conditions -}}
                 {{ with $condition.reason -}}
					 {{ if $condition.reason -}}
						 {{- $condition.reason -}}
					 {{- else -}}
						 "NOT SPECIFIED"
					 {{- end -}}
				 {{- end -}}
             {{- end -}}
         {{- end -}}
	{{- else if .status.containerStatuses -}}
	        {{- range $containerStatus:=.status.containerStatuses -}}
                {{- if $containerStatus.state.waiting -}}
			      **namespace: {{ $pod_namespace }} **pod: {{ $pod_name }} **Reason: {{- $containerStatus.state.waiting.reason -}}					 
				{{- end -}}
		    {{- end -}}
   {{ "\n"}}{{- end -}}
{{- end -}}'| tr -s '\n' '\n'
```

7. Gather Elastic Search Data Usage  
```
for es_pod in $(oc get po --selector=component=es --no-headers -o jsonpath='{range .items[?(@.status.phase=="Running")]}{.metadata.name}{"\n"}{end}'); do oc exec $es_pod -- df -h /elasticsearch/persistent; done
```

8. Gather Elastic Search Logging Image version 
```
oc get po -n logging -o 'go-template={{range $pod := .items}}{{if eq $pod.status.phase "Running"}}{{range $container := $pod.spec.containers}}oc exec -c {{$container.name}} {{$pod.metadata.name}} -n logging -- find /root/buildinfo -name Dockerfile-openshift* | grep -o logging.* {{"\n"}}{{end}}{{end}}{{end}}' | bash -
```
9. Gather Metrics Image version
```
oc get po -n openshift-infra -o 'go-template={{range $pod := .items}}{{if eq $pod.status.phase "Running"}}{{range $container := $pod.spec.containers}}oc exec {{$pod.metadata.name}} -n openshift-infra -- find /root/buildinfo -name Dockerfile-openshift* | grep -o metrics.* {{"\n"}}{{end}}{{end}}{{end}}' | bash -
```

10. Gather map value like environment value?
```
oc get dc/router -n default -o jsonpath='{ .spec.template.spec.containers[0].env[?(@.name=="STATS_PASSWORD")].value}'
```

11. Gather string with complex data.
```
oc get dc/router -n default -o jsonpath='{ .spec.template.spec.containers[?(@.name=="router")].volumeMounts[?(@.name=="server-certificate")].mountPath }'
```

12. Delete orphan Docker images 
```
for im in $(docker images|grep '\<none' |awk '{print $3}'); do docker rmi --force $im;done
```
