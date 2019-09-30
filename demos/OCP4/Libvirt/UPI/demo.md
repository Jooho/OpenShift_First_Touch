
## OpenShift UPI Demo script
## Manual way

Config KVM (optional)
```
cd Config_KVM/ 
ansible-playbook config_kvm.yml -vvvv
```

Generate config files
```
ansible-playbook -i config prep/ansible/tasks/generate_config_files.yml -vvv
```

Download essential files
```
ansible-playbook -i config prep/ansible/tasks/download_files.yml -vvv
```

Go to prep folder
```
cd prep
```

Generate Cloud Init files
```
ansible-playbook -i ./ansible/inventory ./ansible/tasks/cloud_init.yml
```

Create DNSMasq config file
```
ansible-playbook -i ./ansible/inventory ./ansible/tasks/dns_config.yml
```

Config Matchbox config & deploy matchbox server
```
ansible-playbook -i ./ansible/inventory ./ansible/tasks/matchbox_config.yml
```

Generate OCP VM module.tf file
```
ansible-playbook -i ./ansible/inventory ./ansible/tasks/
```

Deploy LB, config matchbox and config haproxy
```
terraform init
terraform get 
TF_LOG=trace terraform apply -auto-approve
```

Deploy OCP
```
cd ../ocp4
terraform init
terraform get 
TF_LOG=trace terraform apply -auto-approve
```

Wait for bootstrap
```
cd ../prep
openshift-install --dir ./upi wait-for bootstrap-complete
```

Update registry storage to emptyDir
```
oc --config upi/auth/kubeconfig patch configs.imageregistry.operator.openshift.io cluster --type merge --patch '{"spec":{"storage":{"emptyDir":{}}}}'
```

Remove lb_rm_bootstrap.yml
```
ansible-playbook -i ansible/inventory ansible/tasks/lb_rm_bootstrap.yml -vvv
```

Wait for ocp install complete
```
openshift-install --dir upi/ wait-for install-complete'
```

## jkit commands

```
./jkit.py init
```

**2. Prep**
```
./jkit.py prep -op [apply,dtr]  # dtr = destory
```

**3. OCP**
```
./jkit.py ocp -op [apply,dtr]   # dtr = destory
```
**4. Post**
```
./jkit.py post
```

**oneshot**
```
./jkit.py oneshot
```
