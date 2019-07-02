cluster_name="upi"

network_domain="example.com"

ssh_public_key="CHANGE_ME"

webserver_url="http://192.168.222.1:9000"
webserver_doc_path="/var/www/html"

matchbox_rpc_endpoint="matchbox.example.com:8081"
matchbox_client_cert="~/.matchbox/client.crt"
matchbox_client_key="~/.matchbox/client.key"
matchbox_trusted_ca_cert="~/.matchbox/ca.crt"

rhcos_kernel_path="assets/rhcos/4.1/rhcos-kernel"
rhcos_initramfs_path="assets/rhcos/4.1/rhcos-initramfs.img"
rhcos_os_image_url="192.168.222.1:9000/rhcos-bios.raw.gz"
rhcos_install_dev="vda"

