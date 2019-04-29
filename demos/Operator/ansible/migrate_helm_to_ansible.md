# Migrate Helm Chart to Ansible Role

Helm chart templates can be resued to ansible templates but you need to change the format. This tutorial will explain the simple way to change the format.

## Test Environment
- Fedora 28
- ansible 2.7.7
- [NFS Provisioner Helm Chart Repository](https://github.com/Jooho/jhouse_openshift/tree/master/demos/Operator/helm/helm-charts/nfs-provisioner)
- [NFS Provisioner Ansible Role Repository](https://github.com/Jooho/ansible-role-nfs-provisioner) - This repository is new one that has nothing.

## Migrate Steps

### Create Structure of ansible role
- Create ansible role 
  - [molecule installation](../../Ansible_Molecule/molecule/install.md)
  ```
  cd /tmp
  ansible-galaxy init ansible-role-nfs-provisioner
  # or
  # molecule init role -r ansible-role-nfs-provisioner
  
  cd ansible-role-nfs-provisioner/
  mkdir templates
  ```

- Initialize git and add remote url
  ~~~
  git init
  git remote add origin https://github.com/Jooho/  ansible-role-nfs-provisioner.git
  ~~~

- Clone the Helm Chart git repository 
  ~~~
  git clone git@github.com:Jooho/jhouse_openshift.git /tmp/jhouse_openshift
  ~~~
- Copy it to ansible role folder 
  ~~~
  cp -R /tmp/jhouse_openshift/demos/Operator/helm/helm-charts/nfs-provisioner/templates/* ./templates/.
  cd ./templates
  ~~~

### Principle Change Rule for Operator
- Do not use `default` anymore
  - ex) `default (.Release.Namespace)`
  - Use default values in `defaults/main.yml`
- Avoid special variables(Magic variables)
  - ex) `role_name`
  - Use default values in `defaults/main.yml`
- Change file type of `yaml` to `yaml.j2`
  - `role.yaml` to `role.yaml.j2`
- Remove unnecessary files
  - _helpers.tpl
  - NOTES.txt


### Migrate chart template one by one
- `include "nfs-provisioner.name`
- `{{ .Chart.name }}`
  - add `provisioner_name` to ./defaults/main.yaml
  - replace it to `{{ provisioner_name }}`
  - it is used as global value for another variables.
  
- `.Values.XXX`
  - add XXX to ./defaults/main.yaml
  - replace it to `{{ XXX }}`
  
- `{{ .Values.serviceaccount | default (include "nfs-provisioner.name" .) }}`
  - add `serviceaccount` to ./defaults/main.yaml
  - replace it to `{{ serviceaccount }}`

- `{{ .Values.namespace | default (.Release.Namespace) }}`
  - add `namespace` to ./defaults/main.yaml
  - replace it to `{{ namespace }}`


- `{{ .Values.deployment.name| default (include "nfs-provisioner.name" .) }}`
  - add `deployment.name` to ./defaults/main.yaml
  - replace it to `{{ deployment.name }}`


- `{{ .Values.replicaCount | default 1 }}` 
  - add `replicaCount` to ./defaults/main.yaml
  - replace it to `{{ replicaCount }}`

- `{{ .Values.deployment.strategy.type  }}`
  - add `deployment.strategy.type` to ./defaults/main.yaml
  - replace it to `{{ deployment.strategy.type }}`

- `{{ .Values.image.repository }}:{{ .Values.image.tag }}`
  - add `image.repository` to ./defaults/main.yaml
  - add `image.tag` to ./defaults/main.yaml
  - replace it to `{{ image.repository }}:{{ image.tag }}`

- `{{ with }}` 
  - replace it to jinja2 template
    ```
     {% if nodeSelector is defined and nodeSelector != '' -%}
      nodeSelector:
        {{nodeSelector}}
      {%- endif -%}
    ```


#### Steps
- Remove unnecessary files
  ```
  rm -rf _helpers.tpl NOTES.txt ./tests/test-connection.yaml ../handlers 
  mkdir ../tests
  mv ./tests/test-*  ../tests/.
  rm -rf ./tests
  ```

- Change file type to j2
  ```
  for file in $(ls); do mv $file ${file}.j2; done
  ```

- Change files one by one
- Lint files
  ```
  cd ..;molecule lint
  ```
- Update tasks/main.yml
  ```
  ---
  # - name: debugging purpose
  #   template: 
  #     src: deployment.yaml.j2
  #     dest: /tmp/deployment.yaml
  - name: Create objects
    k8s:
      state: "{{k8s_state}}"
      definition: "{{ lookup('template', item) }}"
      api_key: "{{api_key}}"
      host: "{{host}}"
      verify_ssl: no
      force: "{{k8s_force}}"
    with_items:
      - 'namespace.yaml.j2'
      - 'serviceaccount.yaml.j2'
      - 'scc.yaml.j2'
      - 'storageclass.yaml.j2'
      - 'service.yaml.j2'
      - 'clusterrole.yaml.j2'
      - 'clusterrolebinding.yaml.j2'
      - 'role.yaml.j2'
      - 'rolebinding.yaml.j2'
      - 'deployment.yaml.j2'
  ```

- Create templates/namespace.yaml
  ```
  apiVersion: v1
  kind: Namespace
  metadata:
    annotations:
      app.kubernetes.io/name: {{ provisioner_name }}
    name: {{ namespace }}
  spec:
    finalizers:
    - kubernetes
  ```
- Update modeclue/default/playbook.yml 
  ```
  ---
  - name: Converge
    hosts: localhost
    tasks:
     - import_role:
         name: ansible-role-nfs-provisioner
       vars:
         host: 'https://masters-311-0129.example.com:8443'      #Update
         api_key: '_gQPBF8-Dkg_fREUDk1Y1oGEo2FFnG58lqpSEblil0A' #Update
  
  ```

- Update defatults/main.yml
  ```
  ---
  provisioner_name: nfs-provisioner
  serviceaccount: "{{ provisioner_name }}"
  namespace: "{{ provisioner_name }}"
  
  k8s_force: False
  k8s_state: present
  
  storageclass:
    name: nfs-storageclass
    provisioner: example.com/nfs
    mountOptions: "vers=4.1"
  
  deployment:
    name: "{{ provisioner_name }}"
    strategy:
      type: "RollingUpdate"
    mountPath: "/export"
    hostPath: "/exports-nfs"
  
  replicaCount: 1
  
  image:
    repository: "quay.io/kubernetes_incubator/nfs-provisioner"
    tag: "latest"
    pullPolicy: "IfNotPresent"
  
  service:
    name: "{{ provisioner_name }}"
    type: ""
    port: "2049"
  
  nameOverride: ""
  fullnameOverride: ""
  nodeSelector: "app: nfs-provisioner"
  
  tolerations: []
  affinity: {} 
  ```

- Deploy ansible 
  ```
  molecule converge
  ```
  
- Error will happen (hostpath related)
  ```
  hostpath: {{ hostpath }}

  ==>
  hostPath: {{ hostPath }}
  ``` 
  


### Useful commands in vi to replace string
```
%s/.Values.//g

```

### Test Ansible Script

```
oc project nfs-provisioner
oc create -f ./tests/test-pvc.yaml
oc get pvc
```

### Tip
- debug log
  ```
  provisioner:
    log: True
    debug: True
    options:
      vvv: True
  ```



### Clean Up

```
oc delete pvc test-pvc
oc delete project nfs-provisioner
```