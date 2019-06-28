resource "libvirt_volume" "vm-qcow2" {
  name   = "${var.vm_name}.qcow2"
  pool   = "${var.vm_pool_name}"
  format = "qcow2"
  size   = "${var.vm_size}"
}

// set boot order hd, network
resource "libvirt_domain" "domain-vm-qcow2" {
  name   = "${var.vm_name}"
  running = true
  memory = "${var.vm_memory}"
  vcpu   = "${var.vm_vcpu}"

  network_interface {
    network_id     = "${var.vm_network_id}"
    hostname       = "${var.vm_hostname}"
    addresses      = ["${var.vm_ip}"]
    mac            = "${var.vm_mac}"
    wait_for_lease = true
  }

  boot_device {
    dev = [ "hd", "network"]
  }

  disk {
    volume_id = "${libvirt_volume.vm-qcow2.id}"
  }
 
  graphics {
    type        = "vnc"
    listen_type = "address"
  }

   
}

terraform {
  required_version = ">= 0.11"
}

