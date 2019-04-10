# Practical Examples


## Gather Pods name
```
oc get pod --output='jsonpath={.items[*].metadata.name}

oc get pod --template='{{ range .items}}{{.metadata.name}}{{end}}' 

oc get pod  --output='go-template={{ range .items}}{{.metadata.name}}{{end}}' 

```

## Gather Not Running Pods name

- Project Level
```
oc get pods -o jsonpath='{.items[?(@.status.phase!="Running")].metadata.name}'

# Result
mysql-1-deploy
```

- Cluster Level
```
 oc get pod --all-namespaces --template='{{ range $pod := .items}}{{if ne $pod.status.phase "Running"}} oc delete pod -n {{$pod.metadata.namespace}} {{$pod.metadata.name}} {{"\n"}}{{end}}{{end}}' 


# Result
 oc delete pod -n chart maudlin-greyhound-drupal-8596dbc8d5-2mn9w 
 oc delete pod -n default mysql-1-deploy
```

**Tip** 
In order to delete not runnin pods, you can add following commands:
```
oc get pods -o jsonpath='{.items[?(@.status.phase!="Running")].metadata.name}' |xargs oc delete pod

oc get pod --all-namespaces --template='{{ range $pod := .items}}{{if ne $pod.status.phase "Running"}} oc delete pod -n {{$pod.metadata.namespace}} {{$pod.metadata.name}} {{"\n"}}{{end}}{{end}}' | bash -
```


## Gather fail reason from not Running Pod.
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
