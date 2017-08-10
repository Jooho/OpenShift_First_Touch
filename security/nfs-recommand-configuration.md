
NFS Recommended Configuraion
----------------------------

- /etc/exports configuration
```
$ vi /etc/exports

/exports *(root_squash,rw,sync)
```
(note) no_root_squash can allow remote root user to change files that are owned by root user in shared storage server.

- exports folder permission

```
$ chgrp -R 5555 /exports
$ chmod -R 2770 /exports                

$ ls -al /exports 

drwxrwx---  12 root 5555 146 Aug  3 10:42 .
dr-xr-xr-x. 18 root root 259 Jul 25 08:58 ..
drwxrwx---   2 root 5555  29 Jul 25 09:19 lv_001
drwxrwx---   2 root 5555   6 Jul 25 09:21 lv_002

```

- The SupplementalGroups must be set in the DeploymentConfig of the pods.
```
kind:  DeploymentConfig
…
spec:
  template:
    spec:
      securityContext:  
        suplementalGroups: [5555]
…
```

