
variable "domain_name" {}
variable "cluster_name" {}
variable "ssh_authorized_key"{}
variable "webserver_doc_path"{}

variable "matchbox_http_endpoint" {}
variable "matchbox_rpc_endpoint" {}
variable "matchbox_client_cert" {}
variable "matchbox_client_key" {}
variable "matchbox_trusted_ca_cert" {}

variable "rhcos_install_dev" {} 
variable "rhcos_os_image_url" {}
variable "rhcos_kernel_path" {}
variable "rhcos_initramfs_path"{}

variable "bootstrap_ign_path" {}
variable "bootstrap_names" {}
variable "bootstrap_mac" {}
variable "bootstrap_domains" {}

variable "master_ign_path" {}
variable "master_names" {
   type="list"
}
variable "master_macs" {
   type="list"
}
variable "master_domains" {
   type="list"
}

variable "worker_ign_path" {}
variable "worker_names" {
   type="list"
}
variable "worker_macs"{
   type="list"
}
variable "worker_domains" {
   type="list"
}
