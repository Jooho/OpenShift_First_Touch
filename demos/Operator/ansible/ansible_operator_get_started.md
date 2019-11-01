# Ansible Operator Get Started

## Pre-requisites

- Export environment Variable
  ```
  export WORK_DIR=/tmp
  export OPERATOR_VERSION="0.12.0"
  ```

- Download operator-sdk
  ```
  wget https://github.com/operator-framework/operator-sdk/releases/download/v${OPERATOR_VERSION}/operator-sdk-v${OPERATOR_VERSION}-x86_64-linux-gnu
  mv operator-sdk-v${OPERATOR_VERSION}-x86_64-linux-gnu /usr/bin/operator-sdk
  chmod 777 /usr/bin/.
  ```

- Install necessary packages for local test
  ```
  sudo yum install -y gcc python-devel ansible
  pip install ansible-runner ansible-runner-http openshift

  # openshift python package for CentOS or RHEL
  # yum install python-openshift
  ```


  Reference
  - [User Guide](https://github.com/operator-framework/operator-sdk/blob/master/doc/ansible/user-guide.md)

