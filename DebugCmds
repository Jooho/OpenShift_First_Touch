[LDAP]

ldapsearch -v -H ldaps://ldap2.example.com:389 -D "cn=read-only-admin,dc=example,dc=com" -w "redhat" -b "dc=example,dc=com" -o ldif-wrap=no  "(&(objectClass=groupOfNames))" -vvvv




[Elastic Search]
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

[Common]
- Image Version Check
```
oc get po -o 'go-template={{range $pod := .items}}{{if eq $pod.status.phase "Running"}}{{range $container := $pod.spec.containers}}oc exec -c {{$container.name}} {{$pod.metadata.name}} -- find /root/buildinfo -name Dockerfile-openshift* | grep -o logging.* {{"\n"}}{{end}}{{end}}{{end}}' | bash -
```
