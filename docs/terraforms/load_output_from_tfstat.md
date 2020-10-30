

upi
   /prep/tfstate
   /prep/output.tf
   /prep/config.tf

   /ocp4/main.tf


prep/output.tf
~~~
output "ocp_network_id" {
   value = "${libvirt_network.ocp_network.id}"
}
~~~

prep/config.tf
~~~
terraform {
 backend local {
    path = "../prep/terraform.tfstate"
  }
}
~~~

ocp4/main.tf
```
data "terraform_remote_state" "prep" {
  backend = "local"

  config = {
    path = "../prep/terraform.tfstate"
  }
}

resource "null_resource" "load_state" {
      provisioner "local-exec" {
       command = <<EOF
         echo "${data.terraform_remote_state.prep.ocp_network_id}
        EOF
   }

}
```
    