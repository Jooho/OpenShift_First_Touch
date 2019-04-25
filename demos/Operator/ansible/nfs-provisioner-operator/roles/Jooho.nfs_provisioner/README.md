NFS Provisioner
=========

This role help to deploy nfs provisioner on the openshift/kubernetes

Requirements
------------

- python >= 2.7
- openshift >= 0.6
- PyYAML >= 3.11


Role Variables
--------------

| Name                      | Default value                                | Requird | Description |
| ------------------------- | -------------------------------------------- | ------- | ----------- |
| provisioner_name          | nfs-provisioner                              | no      |             |
| serviceaccount            | {{ provisioner_name}}                        | no      |             |
| namespace                 | {{ provisioner_name}}                        | no      |             |
| storageclass.name         | nfs-storageclass                             | no      |             |
| storageclass.provisioner  | example.com/nfs                              | no      |             |
| storageclass.mountOptions | vers=4.1                                     | no      |             |
| deployment.name           | {{ provisioner_name}}                        | no      |             |
| deployment.strategy.type  | RollingUpdate                                | no      |             |
| deployment.mountPath      | /export                                      | no      |             |
| deployment.hostPath       | /exports-nfs                                 | no      |             |
| image.repository          | quay.io/kubernetes_incubator/nfs-provisioner | no      |             |
| image.tag                 | latest                                       | no      |             |
| image.pullPolicy          | IfNotPresent                                 | no      |             |
| service.name              | {{ provisioner_name}}                        | no      |             |
| service.type              | ""                                           | no      | not used    |
| service.port              | 2049                                         | no      |             |
| nodeSelector              | "app: nfs-provisioner"                       | no      |             |
| tolerations               | []                                           | no      |             |
| affinity                  | {}                                           | no      |             |
| replicaCount              | 1                                            | no      |             |
| k8s_force                 | False                                        | no      |             |
| k8s_state                 | pres                                         | no      |             |


Dependencies
------------

None



Example Playbook
----------------
~~~
- name: Example Playbook
  hosts: localhost
  tasks:
   - import_role:
       name: ansible-role-nfs-provisioner
     vars:
       host: 'https://masters-311-0129.example.com:8443'
       api_key: '$TOKEN' 
~~~


MOLECULE
--------
- Converge
  ~~~
  molecule converge
  ~~~



License
-------

BSD/MIT

Author Information
------------------

This role was created in 2019 by [Jooho Lee](http://github.com/jooho).

