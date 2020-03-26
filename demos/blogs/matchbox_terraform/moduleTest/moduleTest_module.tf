module "matchbox" {
  source = "../modules/matchbox"

  matchbox_http_endpoint   = "matchbox.example.com:8080"
  matchbox_rpc_endpoint    = "matchbox.example.com:8081"
  matchbox_client_cert     = "~/.matchbox/client.crt"
  matchbox_client_key      = "~/.matchbox/client.key"
  matchbox_trusted_ca_cert = "~/.matchbox/ca.crt"
  ssh_authorized_key       = "ssh_public_key"

  rhcos_kernel_path = "assets/rhcos_kernel"
  rhcos_initramfs_path = "assets/rhcos_initramfs.img"

  moduleTest_names   = [
      "moduleTest-0",
      "moduleTest-1",
    ]

  moduleTest_macs    = [
          "52:54:00:11:01:20",
          "52:54:00:11:01:21",
      ]

  moduleTest_domains = [
          "moduleTest-0.upi.example.com",
          "moduleTest-1.upi.example.com",

  ]

  moduleTest_ign_path = "./moduleTest.ign"
}
