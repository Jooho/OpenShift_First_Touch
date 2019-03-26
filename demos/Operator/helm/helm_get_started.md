# Helm Get Started

- Components
- Installation
  - HELM client
  - Tiller


## Components
- HELM client
- Tiller

Tiller is api server that understand helm chart and do actions. This component is deployed inside cluster normally but it can be running indepentantly.


## Installation

### Install HELM

There are many ways to install helm client. please check [here](https://helm.sh/docs/using_helm/#installing-helm)

Here, I will use script way
```
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get > get_helm.sh
chmod 700 get_helm.sh
./get_helm.sh
```

### Install Tiller

[offical doc](https://helm.sh/docs/using_helm/#role-based-access-control)

- Create RBAC 
```
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tiller
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: tiller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: tiller
    namespace: kube-system
```

- Deploy cluster-wide tiller
~~~
helm init
~~~



## Commands

helm create mychart

helm install ./mychart
helm install --debug --dry-run ../mychart
helm install -f myvals.yaml ./mychart
helm install --dry-run --debug --set favoriteDrink=slurm ./mychart
helm install stable/drupal --set image=my-registry/drupal:0.1.0 --set livenessProbe.exec.command=[cat,docroot/CHANGELOG.txt] --set livenessProbe.httpGet=null --dry-run
helm list

helm get manifest $NAME


### Built-in Object
```
Release: This object describes the release itself. It has several objects inside of it:
  Release.Name: The release name
  Release.Time: The time of the release
  Release.Namespace: The namespace to be released into (if the manifest doesn’t override)
  Release.Service: The name of the releasing service (always Tiller).
  Release.Revision: The revision number of this release. It begins at 1 and is incremented for each helm upgrade.
  Release.IsUpgrade: This is set to true if the current operation is an upgrade or rollback.
  Release.IsInstall: This is set to true if the current operation is an install.

Values: Values passed into the template from the values.yaml file and from user-supplied files. By default, Values is empty.

Chart: The contents of the Chart.yaml file. Any data in Chart.yaml will be accessible here. For example {{.Chart.Name}}-{{.Chart.Version}} will print out the mychart-0.1.0.
The available fields are listed in the Charts Guide

Files: This provides access to all non-special files in a chart. While you cannot use it to access templates, you can use it to access other files in the chart. See the section Accessing Files for more.
  Files.Get is a function for getting a file by name (.Files.Get config.ini)
  Files.GetBytes is a function for getting the contents of a file as an array of bytes instead of as a string. This is useful for things like images.

Capabilities: This provides information about what capabilities the Kubernetes cluster supports.
  Capabilities.APIVersions is a set of versions.
  Capabilities.APIVersions.Has $version indicates whether a version (batch/v1) is enabled on the cluster.
  Capabilities.KubeVersion provides a way to look up the Kubernetes version. It has the following values: Major, Minor, GitVersion, GitCommit, GitTreeState, BuildDate, GoVersion, Compiler, and Platform.
  Capabilities.TillerVersion provides a way to look up the Tiller version. It has the following values: SemVer, GitCommit, and GitTreeState.

Template: Contains information about the current template that is being executed
  Name: A namespaced filepath to the current template (e.g. mychart/templates/mytemplate.yaml)
  BasePath: The namespaced path to the templates directory of the current chart (e.g. mychart/templates).
```

### Functions
- repeat 5
- upper
- quote
- default "tea"
  - drink: {{ .Values.favorite.drink | default "tea" | quote }}
- now
  - `{{ now | htmlDate }}` ==> date: 2019-03-22




## Tip
- `{{-`  will remove space
- `{{indent 2 "mug:true"}}` indent configuration
- `with` Change scope
  ```
  {{- with .Values.favorite }}
  drink: {{ .drink | default "tea" | quote }}
  food: {{ .food | upper | quote }}
  {{- end }}
  release: {{ .Release.Name }}
  ```
- `range` for statement
  ```
    {{- range .Values.pizzaToppings }}
    - {{ . | title | quote }}
    {{- end }}
  ```
  ```
  {{- range $key, $val := .Values.favorite }}
  {{ $key }}: {{ $val | quote }}
  {{- end }}
  ```

- variable - it can be used in another scope
  - `{{- $relname := .Release.Name -}}`
  - `release: {{ $relname }}`
- `define` variable in _helper.tpl
  ```
  {{- define "nginx.fullname" -}}
  {{- $name := default .Chart.Name .Values.nameOverride -}}
  {{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
  {{- end -}}

  ```
- `template` 
- `include`
  ```
   {{- include "mychart.app" . | nindent 4 }}
  ```

  - `tuple` function to create a list of files that we loop throug
    ```
    {{- range tuple "config1.toml" "config2.toml" "config3.toml" }}
    {{ . }}: |-
     {{ $files.Get . }}
    {{- end }}
    ```


- Why `.` is needed after template ` {{- template "mychart.labels" . }}` to pass [scope](https://helm.sh/docs/chart_template_guide/#setting-the-scope-of-a-template).

- action vs function
  - template is action so no ways to pass output
  - include is function so the output 

## Tip
- [Trick](https://github.com/helm/helm/blob/master/docs/charts_tips_and_tricks.md)
- [Best Practice](https://helm.sh/docs/chart_best_practices/#the-chart-best-practices-guide)

