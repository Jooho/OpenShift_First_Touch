
resource "matchbox_group" "bootstrap" {

  name = "bootstrap"

  profile = "bootstrap"

  selector = {
    mac = "${var.bootstrap_mac}"
  }
  metadata = {
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}

resource "matchbox_group" "master" {
  count   = "${length(var.master_names)}"
  name    = "${var.cluster_name}-${var.master_names[count.index]}"
  #name    = "${format('%s-%s', var.cluster_name, element(var.master_names, count.index))}" #version 0.12
  profile = "master"

  selector = {
    mac = "${element(var.master_macs, count.index)}"
    #mac = element("${var.master_macs}", "${count.index}") #version 0.12
  }
  metadata = {
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}

resource "matchbox_group" "worker" {
  count   = "${length(var.worker_names)}"
  name    = "${var.cluster_name}-${var.worker_names[count.index]}"
  #name    = "format('%s-%s', var.cluster_name, element(var.master_names, count.index))" version 0.12

  profile = "worker"

  selector = {
    mac = "${element(var.worker_macs, count.index)}"
    #mac = element("${var.worker_macs}", count.index) #version 0.12
  }
  metadata = {
    ssh_authorized_key = "${var.ssh_authorized_key}"
  }
}
