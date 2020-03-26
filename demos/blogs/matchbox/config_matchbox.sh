

export MATCHBOX_IP=$(ip a | grep $(route |grep default|awk '{print $8}')|grep inet|awk '{print $2}'|cut -d/ -f1)
export MATCHBOX_HOSTNAME=matchbox.example.com

echo "$MATCHBOX_IP $MATCHBOX_HOSTNAME" >> /etc/hosts 


mkdir /etc/matchbox
mkdir ~/.matchbox
mkdir /var/lib/matchbox/assets -p

cd /tmp
rm -rf matchbox
git clone https://github.com/coreos/matchbox.git 

export SAN=DNS.1:${MATCHBOX_HOSTNAME},IP.1:${MATCHBOX_IP}
cd ./matchbox/scripts/tls ; ./cert-gen

cp ca.crt ~/.matchbox/ca.crt
cp ca.crt /etc/matchbox/ca.crt
cp client.* ~/.matchbox/. -R
cp server.* /etc/matchbox/. -R
mkdir -p /var/lib/matchbox/assets

mkdir -p /var/lib/matchbox/assets

podman run -d --name matchbox_server --net host --rm -v /var/lib/matchbox:/var/lib/matchbox:Z -v /etc/matchbox:/etc/matchbox:Z,ro quay.io/coreos/matchbox:latest -address=0.0.0.0:8080 -rpc-address=0.0.0.0:8081 -log-level=debug
