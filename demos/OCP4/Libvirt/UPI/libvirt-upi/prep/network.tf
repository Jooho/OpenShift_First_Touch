resource "libvirt_network" "ocp_network" {
   name = "${var.cluster_name}"
   mode = "${var.network_mode}"
   domain = "${var.cluster_name}.${var.network_domain}"
 
   addresses = ["${var.network_address}"]
 
   bridge = "${var.network_bridge}"
 
   dns = {
     enabled = true
     local_only = false
   }
#   provisioner "local-exec" {
#       command = <<EOF
#          ansible-playbook -i ./ansible/inventory ./ansible/tasks/matchbox_config.yml -e @ansible/defaults/main.yml
#        EOF
#   }

   provisioner "local-exec" {
       command = <<EOF
          ansible-playbook -i ./ansible/inventory ./ansible/tasks/ocp_module.yml -e @ansible/defaults/main.yml
        EOF
   }
   xml {
       xslt = "${file("bootp.xsl")}"
   }


  depends_on = [ "module.matchbox" ]
}

#  dns = {
#    enabled = true
#    local_only = true
#     forwarders = [
#       {
#         address = "${cidrhost(var.network_address,1)}"
#         domain = "${var.cluster_name}.${var.network_domain}"
#       }
#     ]
#  }

