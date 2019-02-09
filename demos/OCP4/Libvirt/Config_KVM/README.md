How to Config KVM for OpenShift 4
---------------------------------

KVM is not supported platform for OpenShift 4 and it is only for development purpose.

However, it does not mean that you can not install OCP 4 on KVM. This doc will help you how to configure KVM and host environment.

**Manual way**
~~~
# kernel parameters
sysctl net.ipv4.ip_forward=1

# Update the parameters permanently
$ echo "net.ipv4.ip_forward = 1" | sudo tee /etc/sysctl.d/99-ipforward.conf
$ sudo sysctl -p /etc/sysctl.d/99-ipforward.conf

# Firewalld
$ sudo firewall-cmd --get-active-zones
$ sudo firewall-cmd --zone=FedoraWorkstation --add-source=192.168.126.0/24
$ sudo firewall-cmd --zone=FedoraWorkstation --add-port=16509/tcp

# NetworkManager
$ vi /etc/NetworkManager/NetworkManager.conf 
...
[main]
dns=dnsmasq             # ⇐ Add
..

# DNS
$ echo server=/tt.testing/192.168.126.1 | sudo tee /etc/NetworkManager/dnsmasq.d/openshift.conf

# Download terraform libvirt plugins (terraforms do not support libvirt yet)
$ GOBIN=~/.terraform.d/plugins go get -u github.com/dmacvicar/terraform-provider-libvirt

# terraform cache setting
$ cat <<EOF > $HOME/.terraformrc
 plugin_cache_dir = "$HOME/.terraform.d/plugin-cache"
 EOF

~~~


**Ansible Playbook**
Originally, this playbook is from [this PR](https://github.com/openshift/installer/blob/ffb427c07a24c30a17a2b13b4eb5096cb2f32609/hack/ocp_libvirt_setup.yaml). This is not mergered yet so I can not send PR after fix issues. I copied it and updated. If you wnat to see original file, please see the PR.

```
ansible-playbook config_kvm.yml
```

[config_kvm.yml](./config_kvm.yml)
