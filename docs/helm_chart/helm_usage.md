# Helm Usage

## Built-in Object
- **Release**: This object describes the release itself. It has several objects inside of it
  - Release.Name: The release name
  - Release.Time: The time of the release
  - Release.Namespace: The namespace to be released into (if the manifest doesn’t override)
  - Release.Service: The name of the releasing service (always Tiller).
  - Release.Revision: The revision number of this release. It begins at 1 and is incremented for each helm upgrade.
  - Release.IsUpgrade: This is set to true if the current operation is an upgrade or rollback.
  - Release.IsInstall: This is set to true if the current operation is an install.

- **Values**: Values passed into the template from the values.yaml file and from user-supplied files. By default, Values is empty.

- **Chart**: The contents of the Chart.yaml file. Any data in Chart.yaml will be accessible here. For example {{.Chart.Name}}-{{.Chart.Version}} will print out the mychart-0.1.0.
The available fields are listed in the Charts Guide

- **Files**: This provides access to all non-special files in a chart. While you cannot use it to access templates, you can use it to access other files in the chart. See the section Accessing Files for more.
  - Files.Get is a function for getting a file by name (.Files.Get config.ini)
  - Files.GetBytes is a function for getting the contents of a file as an array of bytes instead of as a string. This is useful for things like images.

- **Capabilities**: This provides information about what capabilities the Kubernetes cluster supports.
  - Capabilities.APIVersions is a set of versions.
  - Capabilities.APIVersions.Has $version indicates whether a version (batch/v1) is enabled on the cluster.
  - Capabilities.KubeVersion provides a way to look up the Kubernetes version. It has the following values: Major, Minor, GitVersion, GitCommit, GitTreeState, BuildDate, GoVersion, Compiler, and Platform.
  - Capabilities.TillerVersion provides a way to look up the Tiller version. It has the following values: SemVer, GitCommit, and GitTreeState.

- **Template**: Contains information about the current template that is being executed
  - Name: A namespaced filepath to the current template (e.g. mychart/ templates/mytemplate.yaml)
  - BasePath: The namespaced path to the templates directory of the current chart (e.g. mychart/templates).


## TEMPLATE FUNCTIONS
Helm has over 60 available functions.

**Syntax**
```
{{ functionName arg1 arg2}}
```

*Example*
```
{{ quote .Values.serviceaccount }}
```
*Result*
```
"serviceaccountName"
```

## PIPELINES
One of the powerful features of the template language is its concept of pipelines.

**Syntax**
```
{{ String |functionName }}
```

*Example*
```
{{ .Values.serviceaccount | quote }}
```
*Result*
```
"serviceaccountName"
```

### Useful Functions

- **repeat**
  - `{{ "abc" | repeat 2 }}`
    - `abcabc`
- **upper**
  - `{{ "abc" | repeat 2 | upper }}`
    - `ABCABC`
- **quote**
  - `{{ "abc" | repeat 2 | upper | quote}}`
    - `"ABCABC"`
- **default**
  - `{{ .Values.favorite.drink | default "tea" }}`
  - `{{ .Values.favorite.drink | default (.Chart.Name)}}`
  - `{{ .Values.favorite.drink | default (include "test.fullname")}}`
    - If the drink value is "", default value will be used
- **now**
  - `{{ now | htmlDate }}`
    - `2019-03-22`
- **eq,ne,lt,gt,and,or,not**
  - `{{ if and .Values.fooString (eq .Values.fooString "foo") }}`
    - The variable .Values.fooString exists and is set to "foo" 
  - `{{ if or .Values.anUnsetVariable (not .Values.aSetVariable) }}`
    - unset variables evaluate to false and .Values.setVariable was negated with the not function. 


### Flow Control (Action)
- **if/else**
  ```
  {{ if PIPELINE }}
    # Do something
  {{ else if OTHER PIPELINE }}
    # Do something else
  {{ else }}
    # Default case
  {{ end }}
  ```
  ~~~
    {{ if and .Values.favorite.drink (eq .Values.favorite.drink "coffee") }}{{ indent 2 "mug: true"}}{{ end }}
  ~~~
- **with**
  ```
  {{ with PIPELINE }}
  # restricted scope
  {{ end }}
  ```
  ~~~
  {{- with .Values.favorite }}
  drink: {{ .drink | default "tea" | quote }}
  {{- end }}
  ~~~

- **range**
  ```
  toppings: |-
  {{- range .Values.pizzaToppings }}
  - {{ . | title | quote }}
  {{- end }}
  ```
  ```
  toppings: |-
   {{- range $index, $topping := .Values.pizzaToppings }}
    {{ $index }}: {{ $topping }}
  {{- end }}
  ```
  ```
  {{- range $key, $val := .Values.favorite }}
  {{ $key }}: {{ $val | quote }}
  {{- end }}
  ```


## Named Templates
- `define` variable in _helper.tpl(You can define block of lines.)
  ```
  {{/* Convention for description */}}
  {{- define "nginx.fullname" -}}
  {{- $name := default .Chart.Name .Values.nameOverride -}}
  {{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
  {{- end -}}

  ```
- `template` (You can include `define` block into yaml file. template names are global.)
    ```
    {{- template "nginx.fullname" }}
    ```
- `include`
  - `template` that is substituted in has the text aligned to the right.
  - `template` is action and not a function
  - there is no way to pass the output of a template call to other functions
  - the data is simply inserted inline.
  ```
   {{- include "mychart.app" . | nindent 4 }}
  ```

  - `tuple` function to create a list of files that we loop 
    ```
    config1.toml:
       message = Hello from config 1
    config2.toml:   
       message = This is config 2
     config3.toml:
       message = Goodbye from config 3
    ```

    ```
    {{- range tuple "config1.toml" "config2.toml" "config3.toml" }}
    {{ . }}: |-
     {{ $files.Get . }}
    {{- end }}
    ```
    ```
    config1.toml: |-
      message = Hello from config 1
  
    config2.toml: |-
      message = This is config 2
  
    config3.toml: |-
      message = Goodbye from config 3
    ```

## [Path](https://golang.org/pkg/path)
- base/dir/ext/isAbs/clean
```
{{ "/etc/abc.yaml"|base }}
{{ "/etc/abc.yaml"|dir }}
{{ "/etc/abc.yaml"|ext }}
{{ "/etc/abc.yaml"|isAbs }}   <<= Absolute Path?
{{ "/etc/abc.yaml"|clean }}


abc.yaml
/etc
.yaml
true
/etc/abc.yaml

```

## [Glob](https://godoc.org/github.com/gobwas/glob)
- Test files
  ```
  mkdir foo bar
  cat << EOF > foo/foo.txt
  1txt=abc
  2txt=abc
  3txt=abc
  EOF
  cat << EOF > bar/bar.yaml
  1yaml: abc
  2yaml: abc
  3yaml: abc
  EOF
  ```
- Get yaml files
  ```
  {{ $root := . }}
  {{ range $path, $bytes := .Files.Glob "**.**" }}
  {{ $path }}: |-
  {{ $root.Files.Get $path }}
  {{ end }}
  ```
  ```
  bar/bar.yaml: |-
  1yaml: abc
  2yaml: abc
  3yaml: abc
  ```
- Encode content

  ```
  {{ range $path, $bytes := .Files.Glob "foo/*" }}
  {{ base $path }}: '{{ $root.Files.Get $path | b64enc }}'
  {{ end }}
  ```
  ```
  foo.txt: 'MXR4dD1hYmMKMnR4dD1hYmMKM3R4dD1hYmMK'
  ```
  ```
  echo "MXR4dD1hYmMKMnR4dD1hYmMKM3R4dD1hYmMK"| base64 -d
  1txt=abc
  2txt=abc
  3txt=abc
  ```
- As Config
  ```
  {{- (.Files.Glob "foo/*").AsConfig | nindent 2 }}
  ```
  ```
  foo.txt: |
    1txt=abc
    2txt=abc
    3txt=abc
  ```

- As Secret
  ```
  {{- (.Files.Glob "bar/*").AsSecrets | nindent 2 }}
  ```
  ```
  bar.yaml: MXlhbWw6IGFiYwoyeWFtbDogYWJjCjN5YW1sOiBhYmMK
  ```
  ```
  echo "MXlhbWw6IGFiYwoyeWFtbDogYWJjCjN5YW1sOiBhYmMK" |base64 -d
  1yaml: abc
  2yaml: abc
  3yaml: abc
  ```
- Load each line
  ```
  {{ range .Files.Lines "bar/bar.yaml" }}
     {{ .|upper }}{{ end }}
  ```

  ```
  1YAML: ABC
  2YAML: ABC
  3YAML: ABC
  ```


## Tip
- `{{-`  will remove space
- Why `.` is needed after template ` {{- template "mychart.labels" . }}` to pass [scope](https://helm.sh/docs/chart_template_guide/#setting-the-scope-of-a-template).

- action vs function
  - template is an action so no ways to pass output
  - include is a function so the output can be passed