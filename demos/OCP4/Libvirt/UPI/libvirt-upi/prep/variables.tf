variable "cluster_name" {
     default = "ocp4"
}

variable "ocp_version" {
     default = "4.1"
}



variable "network_bridge" {
     default = "virbr-ocp4"
}


variable "network_mode" {
     default = "nat"
  }

variable "network_domain" {
     default = "example.com"
  }

variable "network_address" {
     default = "192.168.222.1/24"
  }


### Storage

variable "libvirt_storage_pool_name" {
     default = "default"
}

variable "libvirt_storage_pool_path" {
     default = "/var/lib/libvirt/images"
}


### LB
variable "lb_memory" {
     description = "MB"
     default = "512"
}

variable "lb_vcpu" {
     default = "1"
}
variable "lb_disk" {
     description = "GB"
     default = "80"
}
variable "lb_vm_volume_source" {
     default = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-1809.qcow2"
}

### Matchbox
variable "rhcos_install_dev" {}
variable "ssh_public_key" {}
variable "matchbox_rpc_endpoint" {
    default="matchbox.example.com:8081"
}

variable "matchbox_http_endpoint" {
    default="matchbox.example.com:8080"
}

variable "matchbox_client_cert" {
    default="~/.matchbox/client.crt"
}

variable "matchbox_client_key" {
    default="~/.matchbox/client.key"
}

variable "matchbox_trusted_ca_cert" {
    default="~/.matchbox/ca.crt"
}


variable "rhcos_kernel_path" {
}

variable "rhcos_initramfs_path" {
}

variable "rhcos_os_image_url"{
}

variable "webserver_doc_path" {
}

variable "webserver_url" {
     default = "http://192.168.222.1"
}
