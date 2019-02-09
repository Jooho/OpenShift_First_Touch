How to Config KVM for OpenShift 4
---------------------------------

KVM is not supported platform for OpenShift 4 and it is only for development purpose.

However, it does not mean that you can not install OCP 4 on KVM. This doc will help you how to configure KVM and host environment.

**Prerequisite**
In order to install OCP4 on KVM, you need to build the NextGen installer on you end, which means you have to install golang package.
The ansible playbook does not install golang package but it checks the version of go.

- Install golang. (Doc - [How to build NextGen OpenShift installer](../Build_CLI/README.md)) 


**Manual way**
[The official doc](https://github.com/openshift/installer/blob/master/docs/dev/libvirt-howto.md) show you how to set up KVM enviroment for openshift 4. Hence, I don't repeat the same steps here. If you have to do it manually, please refer the doc. 


**Ansible Playbook**

Originally, this playbook is from [this PR](https://github.com/openshift/installer/blob/ffb427c07a24c30a17a2b13b4eb5096cb2f32609/hack/ocp_libvirt_setup.yaml). This is not mergered yet so I can not send PR after some issues fixed. I copied the origial playbook and updated. If you wnat to see the original file, please see the PR.

```
ansible-playbook config_kvm.yml
```

This is updated file([config_kvm.yml](./config_kvm.yml)).


