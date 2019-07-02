# Config Matchbox by terraform
This tutorial show how to config matchboxy by terrraform matchbox provider.

[Sample](./sample)
[Real OCP4 sample](https://github.com/Jooho/jhouse_openshift/tree/master/demos/OCP4/Libvirt/UPI/libvirt-upi)

## Relationship between components

- assets
  - contains kernel/initrmfs file (or you can specify URL to download)
- group
  - find VM to match profile by 
    - mac
    - uuid
    - hostname
    - label
    - etc
- profile
  - Specify 
    - kernel parameters
    -  kernel file
    -  initrmfs file
    -   ignition
- ignition
  - ignition files

```
group(profile: bootstrap)  
    |
    +----------->   profile (ignition: bootstrap.ign) 
                                |
                                +------------> inigition (file=> bootstrap.ign)
                                (kernel/initrmfs: kernel/initrmfs.img)
                                    |
                                    +------------> assets (file => kernel/initrmfs.img)
                                            
```                                    
### provider
```
provider "matchbox" {
  endpoint    = "${var.matchbox_rpc_endpoint}"
  client_cert = "${file(var.matchbox_client_cert)}"
  client_key  = "${file(var.matchbox_client_key)}"
  ca          = "${file(var.matchbox_trusted_ca_cert)}"
}
```

### group

```
resource "matchbox_group" "bootstrap" {

  name = "bootstrap"

  profile = "bootstrap"

  selector = {
    mac = "${var.bootstrap_mac}"             # find VM by mac
  }
  metadata = {
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}
```

### profile
```
locals {
  kernel_args = [
    "console=tty0",
    "console=ttyS1,115200n8",
    "rd.neednet=1",

    # "rd.break=initqueue"
    "coreos.inst=yes",

    "coreos.inst.image_url=${var.rhcos_os_image_url}",
    "coreos.inst.install_dev=${var.rhcos_install_dev}",
    "coreos.inst.skip_media_check",
    "coreos.inst.ignition_url=${var.matchbox_http_endpoint}/ignition?mac=$${mac:hexhyp}",
  ]

  pxe_kernel = "${var.rhcos_kernel_path}"
  pxe_initrd = "${var.rhcos_initramfs_path}"
}


resource "matchbox_profile" "bootstrap" {
  name  = "bootstrap"

  kernel = "assets/rhcos_kernel"

  initrd = ["assets/rhcos_initramfs.img}"]

  args = "${local.kernel_args}"

  raw_ignition = "${file(var.bootstrap_ign_path)}"

}

```

