
resource "matchbox_group" "moduleTest" {

  count   = "${length(var.moduleTest_macs)}"
  name    = "${var.moduleTest_names[count.index]}"
  #name    = "format('%s-%s', var.cluster_name, element(var.master_names, count.index))" version 0.12

  profile = "moduleTest"

  selector = {
    mac = "${element(var.moduleTest_macs, count.index)}"
    #mac = element("${var.moduleTest_macs}", count.index) #version 0.12
  }
  metadata = {
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }

}

