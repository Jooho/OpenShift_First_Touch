
resource "libvirt_pool" "default" {
  name = "${var.libvirt_storage_pool_name}"
  type = "dir"
  path = "${var.libvirt_storage_pool_path}"
}


# We fetch the latest centos release image from their mirrors
resource "libvirt_volume" "lb-qcow2" {
  name   = "lb.qcow2"
  pool   = "${libvirt_pool.default.name}"
  source = "${var.lb_vm_volume_source}"
  format = "qcow2"

  provisioner "local-exec" {
     command = <<EOF
       sudo qemu-img resize "${var.libvirt_storage_pool_path}/${libvirt_volume.lb-qcow2.name}" "${var.lb_disk}"G   # TODO remove sudo
       EOF
   }
}

data "template_file" "user_data" {
  template = "${file("${path.module}/user-data")}"
}

data "template_file" "meta_data" {
  template = "${file("${path.module}/meta-data")}"
}


resource "libvirt_cloudinit_disk" "commoninit" {
  name           = "commoninit.iso"
  user_data      = "${data.template_file.user_data.rendered}"
  meta_data      = "${data.template_file.meta_data.rendered}"
  pool           = "${libvirt_pool.default.name}"
}

resource "libvirt_domain" "domain-lb" {
  name   = "lb-${var.cluster_name}"
  memory = "${var.lb_memory}"
  vcpu   = "${var.lb_vcpu}"

  cloudinit = "${libvirt_cloudinit_disk.commoninit.id}"

  network_interface {
    network_id     = "${libvirt_network.ocp_network.id}"
    hostname       = "lb.${var.cluster_name}.${var.network_domain}"
    addresses      = ["${cidrhost(var.network_address,2)}"]
    wait_for_lease = true
 
  }

  boot_device {
    dev = [ "hd", "cdrom"]
  }


  # IMPORTANT: this is a known bug on cloud images, since they expect a console
  # we need to pass it
  # https://bugs.launchpad.net/cloud-images/+bug/1573095
  console {
    type        = "pty"
    target_port = "0"
    target_type = "serial"
  }

  console {
    type        = "pty"
    target_type = "virtio"
    target_port = "1"
  }

  disk {
    volume_id = "${libvirt_volume.lb-qcow2.id}"
  }

  graphics {
    type        = "spice"
    listen_type = "address"
    autoport    = true
  }

  provisioner "local-exec" {
     command = <<EOF
        ansible-playbook -i ./ansible/inventory ./ansible/tasks/haproxy_config.yml -vvv
       EOF
   }

}



#resource "null_resource" "clean_process" {
#
#   provisioner "local-exec" {
#       when = "destroy"
#       command = <<EOF
#          ansible-playbook -i ./ansible/inventory ./ansible/tasks/dns_config.yml -e op=destroy
#        EOF
#   }
#
#}

