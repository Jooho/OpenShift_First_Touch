# Jhouse Openshift

## Docs

Volume
------
- [EmptyDir](./docs/volume/emptyDir.adoc)

Memory
------
- [Shared Memory](./docs/memory/shared_memory.adoc)

ETCD
----
- [etcd_recovery_unhealthy_member](./docs/etcd/etcd_recovery_unhealthy_member.md)


Router
------
- [sharding](./docs/router/sharding.md)


Security
--------
- [nfs-custom-scc](./docs/security/nfs-custom-scc.md)
- [nfs-recommand-configuration](./docs/security/nfs-recommand-configuration.md)
- [selinux](./docs/security/selinux.md)

EFK
---
- [debugging-gathering-efk-object](./docs/efk/debugging-gathering-efk-object.md)


Useful
-------
- [dockerfile-from-image](./useful/dockerfile-from-image.md)

Sample Script
------------
- cicd
  - [bluegreen-pipeline.yaml](./sample_scripts/cicd/jenkins-build-script/bluegreen-pipeline.yaml)
- [network](./sample_scripts/network)
  - egress-policy
  - egress-router-service
  - egress-router
  - ingress
  - mysql-external
  - mysql-ingress
- pv
  - [create-pv-yaml.sh](./sample_scripts/pv/create-pv-yaml.sh)
  - [create-pvc-yaml.sh](./sample_scripts/pv/create-pvc-yaml.sh)
- quota
  - [defaultProjectTemplate.yaml](./sample_scripts/quota/defaultProjectTemplate.yaml)
  - [pod-with-resources.yaml](./sample_scripts/quota/pod-with-resources.yaml)
  - [pod-without-resources.yaml](./sample_scripts/quota/pod-without-resources.yaml)
  - [quota.yaml](./sample_scripts/quota/quota.yaml)


APM
----
- [Scouter-APM-Tool](https://github.com/Jooho/scouter-docker)
