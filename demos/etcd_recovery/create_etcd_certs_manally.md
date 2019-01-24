Create ETCD certificates manually
--------------------------------

When you have certificate fail issue with ETCD, you can execute ansible script.

However, you can also generate the certs manually.

This doc show you how to generate ETCD cert with specific scenario.

## TEST Environment ##

### The ETCD node that will generate certificate for target ETCD member ###
- vm125.gsslab.rdu2.redhat.com

### Target ETCD member ###
- dhcp182-77.gsslab.rdu2.redhat.com(10.10.182.77)

### Export target ETCD member hostname ###
```
export target_etcd=dhcp182-77.gsslab.rdu2.redhat.com
```

## Check CA folder ##
**`/etc/etcd/ca` folder must have certs (these files should be on the first master node)**

```
$ ls /etc/etcd/ca
 ca.crt  ca.key  certs  crl  fragments  generated_certs  index.txt  index.txt.attr  index.txt.attr.old  index.txt.old  openssl.cnf  serial  serial.old
```

*The `/etc/etcd/ca/ca.crt` must be the same as `/etc/etcd/ca.crt` file.*
```
diff /etc/etcd/ca/ca.crt /etc/etcd/ca.crt
```

**NOTE: You must execute this command on a ETCD node where ETCD is running well in the cluster  or the node where contain ca files!!**
In this case, `vm125.gsslab.rdu2.redhat.com` is the ETCD member where this command will be execueted.

```
cd /etc/etcd
export NEW_ETCD=dhcp182-77.gsslab.rdu2.redhat.com      # <==== Update
export CN=$NEW_ETCD
export SAN="IP:10.10.182.77,DNS.1:$NEW_ETCD"           # <==== Update
export PREFIX="./generated_certs/etcd-$CN/"
mkdir -p $PREFIX
openssl req -new -keyout ${PREFIX}server.key   -config ca/openssl.cnf   -out ${PREFIX}server.csr   -reqexts etcd_v3_req -batch -nodes   -subj /CN=$CN
openssl ca -name etcd_ca -config ca/openssl.cnf   -out ${PREFIX}server.crt   -in ${PREFIX}server.csr   -extensions etcd_v3_ca_server -batch
openssl req -new -keyout ${PREFIX}peer.key   -config ca/openssl.cnf   -out ${PREFIX}peer.csr   -reqexts etcd_v3_req -batch -nodes   -subj /CN=$CN
openssl ca -name etcd_ca -config ca/openssl.cnf   -out ${PREFIX}peer.crt   -in ${PREFIX}peer.csr   -extensions etcd_v3_ca_peer -batch
#cp etcd.conf ${PREFIX}  # <==== if the etcd node is fresh new node
cp ca.crt ${PREFIX}
tar -czvf ${PREFIX}${CN}.tgz -C ${PREFIX} .
scp ${PREFIX}${CN}.tgz  $CN:/etc/etcd/
````

## If the etcd member is not joined to the cluster yet, please execute this.
```
export ETCD_CA_HOST="vm125.gsslab.rdu2.redhat.com"
export NEW_ETCD="dhcp182-77.gsslab.rdu2.redhat.com"
export NEW_ETCD_IP="10.10.182.77"

etcdctl -C https://${ETCD_CA_HOST}:2379 \
  --ca-file=/etc/etcd/ca.crt     \
  --cert-file=/etc/etcd/peer.crt     \
  --key-file=/etc/etcd/peer.key member add ${NEW_ETCD} https://${NEW_ETCD_IP}:2380
```


**NOTE: You must execute this command on a new or recovered ETCD node!!**
Go to `dhcp182-77.gsslab.rdu2.redhat.com`

```
systemctl stop etcd
tar -xf /etc/etcd/${NEW_ETCD}.tgz /etc/etcd/. --overwrite
chown -R etcd:etcd /etc/etcd/*

systemctl start restart
```
