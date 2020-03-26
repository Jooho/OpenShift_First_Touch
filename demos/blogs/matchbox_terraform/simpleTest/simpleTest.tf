# Initialize Matchbox Terraform Provider.
provider "matchbox" {
  endpoint    = "matchbox.example.com:8081"
  client_cert = "${file("~/.matchbox/client.crt")}"
  client_key  = "${file("~/.matchbox/client.key")}"
  ca          = "${file("~/.matchbox/ca.crt")}"
}

# Create a group
resource "matchbox_group" "simpleTest" {

  name = "simpleTest"

  profile = "simpleTest"

  selector = {
    mac = "10:10:10:10:10:10"             # find VM by mac
  }
  metadata = {
    ssh_authorized_key = "ssh_authorized_key"
  }
}

# Create a profile
resource "matchbox_profile" "simpleTest" {
  name  = "simpleTest"

  kernel = "assets/rhcos_kernel"

  initrd = ["assets/rhcos_initramfs.img"]

  args = ["simpleTest"]
}
