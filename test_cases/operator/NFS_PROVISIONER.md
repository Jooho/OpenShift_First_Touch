# NFS Provisioner Go Operator


## NFS Provisioner Type

1. [NFS client provisioner](https://github.com/Jooho/openshift-first-touch/blob/master/docs/storage/nfs/nfs-client-provisioner.md)
  - When you have NFS server, you can use the NFS server for openshift.
  - Using this NFS client provisioner, you can create PV with your own NFS server.

2. [NFS provisioner](https://github.com/Jooho/openshift-first-touch/blob/master/docs/storage/nfs/nfs-provisioner.md)
  - This NFS provisioner will deploy NFS server and provide StorageClass with it.
  - **This is a main feature to be implemented from the demo**


## Test Environment
- CRC
- Basic Installation (3 master x 3 workers)

## Deployment Methods

- [Template](./template/README.md)
- Ansible 
- Helm
- Helm Operator
- Ansible Operator
- [Go Operator](./go-operator/README.md)
- OLM 
- OPM
- OperatorHub




## Useful tools

- [Ansible Role for deploying NFS server](https://github.com/Jooho/ansible-role-nfs-server)

