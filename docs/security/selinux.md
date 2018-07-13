# Selinux CheatSheet


NFS
---
Selinux Util Packages
```
yum install -y policycoreutils libselinux-utils policycoreutils-python 
```

virt_sandbox_use_nfs looks deprecated : https://github.com/openshift/openshift-docs/issues/4746
```
setsebool -P virt_use_nfs 1
```

Change Selinux for hostPath
---------------------------
```
semanage fcontext -a -t svirt_sandbox_file_t "/path(/.*)?"`

chcon -R /path should do it
```
