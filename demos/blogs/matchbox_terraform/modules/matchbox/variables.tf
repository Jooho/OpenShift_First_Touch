variable "ssh_authorized_key"{}

variable "matchbox_http_endpoint" {}
variable "matchbox_rpc_endpoint" {}
variable "matchbox_client_cert" {}
variable "matchbox_client_key" {}
variable "matchbox_trusted_ca_cert" {}

variable "rhcos_kernel_path" {}
variable "rhcos_initramfs_path"{}


variable "moduleTest_ign_path" {}
variable "moduleTest_names" {
   type="list"
}
variable "moduleTest_macs"{
   type="list"
}
variable "moduleTest_domains" {
   type="list"
}

