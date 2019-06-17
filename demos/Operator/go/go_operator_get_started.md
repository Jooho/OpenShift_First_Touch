# Go Operator get started


## Pre-requisites

- Export environment Variable
  ```
  export WORK_DIR=/tmp
  export OPERATOR_VERSION="0.6.0"
  ```

- Download operator-sdk
  ```
  wget https://github.com/operator-framework/operator-sdk/releases/download/v${OPERATOR_VERSION}/operator-sdk-v${OPERATOR_VERSION}-x86_64-linux-gnu
  mv operator-sdk-v${OPERATOR_VERSION}-x86_64-linux-gnu /usr/bin/operator-sdk
  chmod 777 /usr/bin/.
  ```

- Install GoLang
  ```
  sudo yum install -y golang-bin gcc-c++ libvirt-devel
  mkdir -p ${WORK_DIR}/go/{src,pkg,bin}
  
  echo "export GOBIN=${WORK_DIR}/dev/git/go/bin" >> ~/.bashrc
  echo "export GOPATH=${WORK_DIR}/dev/git/go" >> ~/.bashrc
  echo "export PATH=${GOBIN}:${PATH}" >> ~/.bashrc
  source ~/.bashrc
  ```

- Install Dep 
  ```
  # Dep Install
  curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
  mv ${WORK_DIR}/dev/git/go/dep /usr/bin/
  ```


## Clean up
   ```
   rm -rf /usr/bin/{dep,operator-sdk}   
   ```

## Reference
- [User Guide](https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md)