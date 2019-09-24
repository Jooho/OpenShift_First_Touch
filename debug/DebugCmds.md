## [Glusterfs-CNS]

- Mapping PVC and glusterfs ID
  ```
  oc get pv -o template --template='{{ range .items }}{{ $type := index .metadata.annotations "gluster.org/type" }}{{ if $type }}PV Name: {{ .metadata.name }}  Volume Type: gluster-{{ $type }}  {{ if eq "file" $type }}Heketi Volume ID: {{ index .metadata.annotations "gluster.kubernetes.io/heketi-volume-id" }}  Gluster Volume Name: {{ .spec.glusterfs.path }}{{ println }}{{ end }}{{ if eq "block" $type }}Heketi BlockVolume ID: {{ index .metadata.annotations "gluster.org/volume-id" }}{{ println }}{{ end }}{{ end }}{{ end }}'

  e.g.
  PV Name: pvc-1a23d25a-5d6c-11e8-bbe9-005056917fec  Volume Type: gluster-file  Heketi Volume ID: 045037c8debbdf69064ce4bba587f2bf  Gluster Volume Name: vol_045037c8debbdf69064ce4bba587f2bf                   
  ```
- Get gluster volume info
  ```
  $ oc rsh $gluster_pod

  $> gluster volume info all
  ```


## [LDAP]

- Check LDAP with ldapsearch
  ```
  ldapsearch -v -H ldaps://ldap2.example.com:389 -D "cn=read-only-admin,dc=example,dc=com" -w "redhat" -b "dc=example,dc=com" -o ldif-wrap=no  "(&(objectClass=groupOfNames))" -vvvv
  ```


## [Elastic Search]
- Thread pool
  ```
  curl -s -k --cert /etc/elasticsearch/secret/admin-cert --key /etc/elasticsearch/secret/admin-key "https://localhost:9200/_cat/thread_pool?v"

  curl -s -k --cert /etc/elasticsearch/secret/admin-cert --key /etc/elasticsearch/secret/admin-key "https://localhost:9200/_cat/thread_pool?v&h=bq,br,ba,sq,sr,sa,gq,gr,ga,bc,sc,gc"
  ```

- Health Check
  ```
  curl -s --key /etc/elasticsearch/secret/admin-key --cert /etc/elasticsearch/secret/admin-cert --cacert /etc/elasticsearch/secret/admin-ca "https://localhost:9200/_cat/health?v"
  curl -s --key /etc/elasticsearch/secret/admin-key --cert /etc/elasticsearch/secret/admin-cert --cacert /etc/elasticsearch/secret/admin-ca "https://localhost:9200/_cat/nodes?v"
  ```
- Indice Dump
  ```
  curl -s -k --cert /etc/elasticsearch/secret/admin-cert --key /etc/elasticsearch/secret/admin-key "https://logging-es:9200/_cat/indices?v"
  ```
  - Shard Dump
  ```
  curl -s --key /etc/elasticsearch/secret/admin-key --cert /etc/elasticsearch/secret/admin-cert --cacert /etc/elasticsearch/secret/admin-ca "https://localhost:9200/_cat/shards?v"
  ```

- Pending Tasks
  ```
  curl -s --key /etc/elasticsearch/secret/admin-key --cert /etc/elasticsearch/secret/admin-cert --cacert /etc/elasticsearch/secret/admin-ca "https://localhost:9200/_cluster/pending_tasks"
  ```

## [Storage]

- Disk Usage
  ```
   du -d 1 -m -x /run
  ```

- Find deleted files but still remains
  ```
  lsof -nP | grep '(deleted)'
  ```

## [Memory]
- Shared Memory
  ```
  ipcs -u
  ipcs -m
  ipcs -s -t
  ```

## [EAP]

- Generate JDR
  ```
  export POD=%POD_NAME

  oc rsync $POD:$(oc exec $POD /opt/eap/bin/jdr.sh | grep "JDR location" | awk '{print $3}') .

  or 

  xdiag.sh -p $POD -j -o ./
  ```

- Heap Dump
  ```
  export POD=%POD_NAME

  oc exec $POD -- /usr/lib/jvm/java-1.8.0-openjdk/bin/jmap -J-d64 -dump:format=b,file='/opt/eap/standalone/tmp/heap.hprof' $(oc exec $POD ps aux | grep java | awk '{print $2}'); oc rsync $POD:/opt/eap/standalone/tmp/heap.hprof .

  or

  xdiag.sh -p $POD -m -o ./
  ```

- Thread Dump
  ```
  export POD=%POD_NAME

  PID=$(oc exec $POD ps aux | grep java | awk '{print $2}'); oc exec $POD -- bash -c "for x in {1..10}; do jstack -l $PID >> /opt/eap/standalone/tmp/jstack.out; sleep 2; done"; oc rsync $POD:/opt/eap/standalone/tmp/jstack.out .
  ```

## [ETCD]
- Backup
  - Snapshot
    ```
    export ETCDCTL_API=3
    source /etc/etcd/etcd.conf
    etcdctl --cert=$ETCD_PEER_CERT_FILE --key=$ETCD_PEER_KEY_FILE --cacert=$ETCD_TRUSTED_CA_FILE --endpoints=$ETCD_LISTEN_CLIENT_URLS snapshot save  /backup/etcd-config-$(date +%Y%m%d)/backup.db  
    ```
  - Check Snapshot
    ```
    etcdctl --write-out=table snapshot status  /backup/etcd-config-$(date +%Y%m%d)/backup.db
    ```
  - [Restore](https://docs.openshift.com/container-platform/3.11/admin_guide/assembly_restore-etcd-quorum.html#backing-up-etcd_restore-etcd-quorum)

  - [Space Quota](https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/maintenance.md#space-quota)

- Unit is masked
  ```
  ex)
  /etc/systemd/system/firewalld.service -> /dev/null
  
  rm -rf /etc/systemd/system/firewalld.service
  systemctl daemon-reload
  ```
  
- etcd wal file is not loaded due to `permission denied`
  ```
  chown -R etcd.etcd /var/lib/etcd/

  restorecon -Rv /var/lib/etcd

  netstat -tunlp|grep 2380
  tcp        0      0 10.10.181.87:2380       0.0.0.0:*               LISTEN      40090/etcd          
  tcp        0      0 10.10.181.87:2380       0.0.0.0:*               LISTEN      40090/etcd          

  kill -9 40090

  systemctl start etcd
  ```
- etcd member list & add
  ```
  # etcd server (https:// :2379,)
  export etcd_members=https://10.10.182.77:2379,https://10.10.178.126:2379,https://10.10.178.125:2379

  # ETCD v2:
  etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers $etcd_members member list
  etcdctl --cert-file=/etc/etcd/peer.crt   --key-file=/etc/etcd/peer.key  --ca-file=/etc/etcd/ca.crt  --peers $etcd_members member add  $(hostname) https://10.10.182.77:2380 

  # ETCD v3:
  etcdctl3 --endpoints $etcd_members member list
  etcdctl3 member add  dhcp182-77.gsslab.rdu2.redhat.com  --peer-urls="https://10.10.182.77:2380"
  ```
## [Certs]
- Verify certificate by ca cert
  ```
  openssl verify -CAfile /etc/origin/master/ca.crt /etc/origin/master/master.server.crt

  openssl verify -CAfile /etc/origin/master/ca-bundle.crt /etc/origin/master/master.server.crt
  ```
- Compare certificate-authority-data
  ```
  grep certificate-authority-data /etc/origin/master/admin.kubeconfig | awk '{ print $2 }' | base64 -d | md5sum
  grep certificate-authority-data /root/.kube/config | awk '{ print $2 }' | base64 -d | md5sum
  md5sum /etc/origin/master/ca-bundle.crt
  ```
- Login with openssl
  ```
  openssl s_client -CAfile /etc/origin/master/ca-bundle.crt -connect oc-master.domain.com:8443
  ```
  
- Execute oc cmd with admin.kubeconfig
  ```
  oc status --config=/etc/origin/master/admin.kubeconfig
  ```

- Check issuer on all masters
  ```
  openssl x509 -noout -issuer -in /etc/origin/master/ca.crt
  md5sum /etc/origin/master/ca.crt 
  ```

- Reference
  - https://github.com/openshift/openshift-ansible/issues/3784

## [Network]
- Long life connection
```
yum install of bcc bcc-tools  kernel-devel-`uname -r`
sudo /usr/share/bcc/tools/tcplife
```

## [Common]
- Image Version Check
  - Logging
   ```
   oc get po -n logging -o 'go-template={{range $pod := .items}}{{if eq $pod.status.phase "Running"}}{{range $container := $pod.spec.containers}}oc exec -c {{$container.name}} {{$pod.metadata.name}} -n logging -- find /root/buildinfo -name Dockerfile-openshift* | grep -o logging.* {{"\n"}}{{end}}{{end}}{{end}}' | bash -
   ```
  - Metrics
  ```
  oc get po -n openshift-infra -o 'go-template={{range $pod := .items}}{{if eq $pod.status.phase "Running"}}{{range $container := $pod.spec.containers}}oc exec {{$pod.metadata.name}} -n openshift-infra -- find /root/buildinfo -name Dockerfile-openshift* | grep -o metrics.* {{"\n"}}{{end}}{{end}}{{end}}' | bash -
  ```


## [OpenShift/Kubernetes]
- Get event regarding node
```
oc get event --all-namespaces --field-selector=involvedObject.kind=Node -o yaml
```
