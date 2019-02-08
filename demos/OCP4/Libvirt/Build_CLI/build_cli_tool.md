Build OpenShift Install CLI Tool
--------------------------------

Libvirt is not supported and it is only for Dev purpose. That's why the relase binary `openshift-install` cli does not have the option `libvirt` as a provider.

In order to install OCP4 on KVM, you must build it by yourself.

This doc help you how to build the project by golang on Fedora.


## Manual

** NOTE **
- If you already have golang, do not install `golang-bin` 
- If you don’t have the dep, please install `dep`.

```
# Update
export WORK_DIR=/home/jooho


sudo yum install golang-bin gcc-c++ libvirt-devel
mkdir -p ${WORK_DIR}/dev/git/go/{src,pkg,bin}

echo "export GOBIN=${WORK_DIR}/dev/git/go/bin" >> ~/.bashrc
echo "export GOPATH=${WORK_DIR}/dev/git/go" >> ~/.bashrc
echo "export PATH=${GOBIN}:${PATH}" >> ~/.bashrc

# Dep Install
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

go get github.com/openshift/installer
cd ${WORK_DIR}/dev/git/go/src/openshift/installer
dep ensure
 
# Build
TAG=libvirt hack/build.sh
```

## Ansible

```
ansible-playbook go_build.yml -vvvvv
```
[go_build.yml](./go_build.yml)



## Reference
- https://github.com/openshift/installer/blob/master/docs/dev/dependencies.md
