
# OpenShift 4 deployment demo using installer 0.14.0 

~~~

# Setup
git clone https://github.com/Jooho/jhouse_openshift.git 

# Build Installer
cd jhouse_openshift/demos/OCP4/Libvirt/Build_CLI
ansible-playbook ./go_build.yml -e git_repo_branch=v0.14.0 -vvvv

# Config KVM on Fedora
cd ../Config_KVM
ansible-playbook ./config_kvm.yml

# Create OCP4 manifests
cd /tmp/go/src/github.com/openshift/installer/bin
./openshift-install create manifests --dir=./ocp4 --log-level=debug

vi ./ocp4/manifests/cluster-ingress-02-config.yml
 domain: apps.tt.testing  # <=== update

# Deploy OCP4 
./openshift-install create cluster --dir=./ocp4 --log-level=debug


# Destory OCP4
./openshift-install destroy cluster --dir=./ocp4 --log-level=debug
~~~



Tip. 
- Use oc command using kubeconfig on master node
  
~~~
# Copy kubeconfig content
cat ./ocp4/auth/kubeconfig

ssh core@192.168.126.11
vi kubeconfig

# Paste the config

export KUBECONFIG=$(pwd)/kubeconfig
~~~

- Check pod information
  
~~~
crictl ps -v

crictl inspectp ${podid}
~~~


*SSHTunnel*
~~~
ssh -L 192.168.126.51:443:10.129.0.2:443 root@192.168.126.51
~~~
