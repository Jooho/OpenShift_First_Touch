
NFS Custom SCC
---

- Check NFS exports folder uid/gid
```
# ls -al /exports/
total 4
drwxrwxrwx   4 root   root     40 Jul  7 16:41 .
dr-xr-xr-x. 18 root   root    259 Jul  7 16:34 ..
drwxrwxrwx   4 nobody nobody 4096 Jul 10 15:04 pv1
# id nobody
uid=99(nobody) gid=99(nobody) groups=99(nobody)
```

- Create Custom SCC

```
allowHostDirVolumePlugin: false
allowHostIPC: false
allowHostNetwork: false
allowHostPID: false
allowHostPorts: false
allowPrivilegedContainer: false
allowedCapabilities: null
apiVersion: v1
defaultAddCapabilities: null
fsGroup:
  type: MustRunAs
groups:
- system:authenticated
kind: SecurityContextConstraints
metadata:
  annotations:
    kubernetes.io/description: restricted denies access to all host features and requires
      pods to be run with a UID, and SELinux context that are allocated to the namespace.  This
      is the most restrictive SCC.
  name: nfs-scc
priority: 9
readOnlyRootFilesystem: false
requiredDropCapabilities:
- KILL
- MKNOD
- SYS_CHROOT
- SETUID
- SETGID
runAsUser:
  type: MustRunAsRange
seLinuxContext:
  type: MustRunAs
supplementalGroups:       <===== Start
  type: MustRunAs
  ranges:
  -  min: 95              <=== guid is 99(nobody) so i set 95 to 105
     max: 105             <===== End
volumes:
- configMap
- downwardAPI
- emptyDir
- persistentVolumeClaim
- secret
```

