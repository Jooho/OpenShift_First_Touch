
locals {
  kernel_args = [
    "console=tty0",
    "coreos.inst.ignition_url=${var.matchbox_http_endpoint}/ignition?mac=$${mac:hexhyp}",
  ]

  pxe_kernel = "${var.rhcos_kernel_path}"
  pxe_initrd = "${var.rhcos_initramfs_path}"

}


resource "matchbox_profile" "moduleTest" {
  name  = "moduleTest"

  kernel = "${local.pxe_kernel}"

  initrd = ["${local.pxe_initrd}"]

  args = "${local.kernel_args}"
 
  raw_ignition = "${file(var.moduleTest_ign_path)}"

}

