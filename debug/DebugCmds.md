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
~~~
ldapsearch -v -H ldaps://ldap2.example.com:389 -D "cn=read-only-admin,dc=example,dc=com" -w "redhat" -b "dc=example,dc=com" -o ldif-wrap=no  "(&(objectClass=groupOfNames))" -vvvv
~~~



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


## [Common]
- Image Version Check
```
oc get po -o 'go-template={{range $pod := .items}}{{if eq $pod.status.phase "Running"}}{{range $container := $pod.spec.containers}}oc exec -c {{$container.name}} {{$pod.metadata.name}} -- find /root/buildinfo -name Dockerfile-openshift* | grep -o logging.* {{"\n"}}{{end}}{{end}}{{end}}' | bash -
```
