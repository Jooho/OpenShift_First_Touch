Build OpenShift Install CLI Tool
--------------------------------

Libvirt is not supported and it is only for Dev purpose. That's why the relase binary `openshift-install` cli does not have the option `libvirt` as a provider.

In order to install OCP4 on KVM, you must build it by yourself.

This doc help you how to build the project by golang on Fedora.

## Video
[![Build Openshift New Generation Installer](http://img.youtube.com/vi/ZD5SnNyOuy4/0.jpg)](https://www.youtube.com/embed/ZD5SnNyOuy4)


## Manual

**NOTE**
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
TAGS=libvirt hack/build.sh
```

## Ansible

### Params

|name| default|required|description|
|-----|-----|-----|-----|
|work_dir|/tmp|no| The location where to clone openshift installer repository|
|git_repo_branch |release-4.1| no| The tag name which openshift installer branch is |

```
ansible-playbook go_build.yml -e git_repo_branch=release-4.1 -vvvvv

```
[go_build.yml](./go_build.yml)

**Default openshift-install binary location:** */tmp/go/src/github.com/openshift/install/bin/openshift-install*

### TroubleShooting

- CPU issue

  Error Msg
  ```
    "/usr/lib/golang/pkg/tool/linux_amd64/link: signal: killed"
  ```
  Solution: Need more cpu resources.

## Reference
- https://github.com/openshift/installer/blob/master/docs/dev/dependencies.md
