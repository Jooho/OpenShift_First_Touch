# Operator Tutorial Overview

This is practical example to make you be familiar with operator.

I will create 3 operators by each type for [NFS provisioner](https://github.com/Jooho/openshift-first-touch/blob/master/docs/storage/nfs/nfs-client-provisioner.md) that was already developed.

This used template way to create several objects 
- openshift-rbac.yaml
- openshift-scc.yaml
- sa.yaml
- deployment.yaml
- storageclass.yaml

## Helm Chart
From this tutorial, I will recreate it using helm chart first instead of template.

- Config Map
  - nfs-provisioner-cm.yml
    - It contains NFS server & NFS provisioner information
      - namespace_name
      - nfs_path
      - nfs_server
- Service Account
  - nfs-provisioner-sa.yml
    - This sa will have power to gather pv/pvc/storageClass/events/svc/pod security policies
    - nfs-provisioner
- Cluster Role
  - nfs-provisioner-runner
- Cluster Role Binding
  - run-nfs-provisioner
- Role
  - leader-locking-nfs-provisioner
- Role Binding
  - leader-locking-nfs-provisioner
- Service 
  - nfs-provisioner
- Deployment
  - nfs-provisioner


## Operator Control Logic

- Helm Operator
  - 
