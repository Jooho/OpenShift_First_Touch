module "matchbox_ocp" {
  source = "../modules/matchbox_ocp"
  
  cluster_name             = "OCP4"
  network_domain             = "ocp4"
  matchbox_http_endpoint   = "matchbox.example.com:8080"
  matchbox_rpc_endpoint    = "matchbox.example.com:8081"
  matchbox_client_cert     = "~/.matchbox/client.crt"
  matchbox_client_key      = "~/.matchbox/client.key"
  matchbox_trusted_ca_cert = "~/.matchbox/ca.crt"
  ssh_authorized_key       = "ssh_public_key"


  # configuations
  rhcos_kernel_path    = "assets/rhcos_kernel"
  rhcos_initramfs_path = "assets/rhcos_initramfs.img"
  rhcos_os_image_url   = "192.168.222.1:9000/rhcos-bios.raw.gz"
  rhcos_install_dev    = "vda"



  # machines
  bootstrap_names = ["bootstrap"]
  bootstrap_mac = "52:54:00:11:00:20"
  bootstrap_domains = ["bootstrap.${cluster_name}.${network_domain}"]
  bootstrap_ign_path = "./bootstrap.ign"   # This should be http server document path in real world (/var/www/html/bootstrap.ign)

  master_names   = [
      "master-0",
    ]

  master_macs    = [
          "52:54:00:11:01:20",
      ]

  master_domains = [
          "master-0.upi.example.com",

  ]

  master_ign_path = "./master.ign"       # This should be http server document path in real world (/var/www/html/worker.ign)

  worker_names   = [
      "worker-0",
    ]

  worker_macs    = [
          "52:54:00:11:02:20",
      ]

  worker_domains = [
          "worker-0.upi.example.com",

  ]

  worker_ign_path = "./worker.ign"     # This should be http server document path in real world (/var/www/html/worker.ign)
}
