provider "matchbox" {
  endpoint    = "${var.matchbox_rpc_endpoint}"
  client_cert = "${file(var.matchbox_client_cert)}"
  client_key  = "${file(var.matchbox_client_key)}"
  ca          = "${file(var.matchbox_trusted_ca_cert)}"
}

