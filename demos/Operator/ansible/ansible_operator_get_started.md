# Ansible Operator Get Started

## Pre-requisites

- Export environment Variable
  ```
  export WORK_DIR=/tmp
  export OPERATOR_VERSION="0.6.0"
  ```

- Download operator-sdk
  ```
  wget https://github.com/operator-framework/operator-sdk/releases/download/v${OPERATOR_VERSION}/operator-sdk-v${OPERATOR_VERSION}-x86_64-linux-gnu
  mv operator-sdk-v0.6.0-x86_64-linux-gnu /usr/bin/operator-sdk
  chmod 777 /usr/bin/.
  ```

- Install Ansible
  ```
  yum install -y ansible 
  ```


  Reference
  - [User Guide](https://github.com/operator-framework/operator-sdk/blob/master/doc/ansible/user-guide.md)

