# Pachyderm Install

## Pachyderm deployment manually 
Based on the [doc](https://docs.pachyderm.com/latest/deploy-manage/deploy/openshift/)

### Steps

- Create a new project
   ~~~
   oc new-project pachyderm-doc-based
   ~~~

- Download pachctl from [here](https://github.com/pachyderm/pachyderm/releases)


- Create manifest.json for packyderm deployment 
   *Syntax format*
    ~~~
    ./pachctl deploy custom --persistent-disk aws --object-store s3 \
        gp2 10 \
        'pachyderm' 'ACCESSKEY' 'SECRETKEY' 'minio-tenant-1-hl.minio-tenant-1.svc.cluster.local:9000' \
        --static-etcd-volume=gp2 --local-roles --dry-run > manifest.json
    ~~~
    *Real cmd*
    ~~~
    ./pachctl deploy custom --persistent-disk aws --object-store s3 \
        any-string 10 \
        'pachyderm' "$(oc get secret minio-tenant-1-creds-secret -n minio-tenant-1 -o jsonpath='{ .data.accesskey}'|base64 -d)" "$(oc get secret minio-tenant-1-creds-secret -n minio-tenant-1 -o jsonpath='{ .data.secretkey}'|base64 -d)" 'minio.minio-tenant-1.svc.cluster.local' \
        --static-etcd-volume=gp2 --local-roles --dry-run --no-verify-ssl  --require-critical-servers-only  > manifest.json
    ~~~ 
    
    *Real Data*
    ~~~
    pachctl deploy custom --persistent-disk aws --object-store s3 \
    any-string 10 \
    'pachyderm' 'b236287d-a4ee-49f3-a305-bfa1e701b79f' '9bb676e3-76b0-454e-bf12-962b58c23c86' 'minio.minio-tenant-1.svc.cluster.local' \
    --static-etcd-volume=gp2 ---local-roles --dry-run --no-verify-ssl  --require-critical-servers-only  > manifest.json
    ~~~


- Edit manifest.json
  - Delete PersistentVolume
  - Delete PersistentVolumeClaim volumeName
    ~~~
            "spec": {
            "accessModes": [
            "ReadWriteOnce"
            ],
            "resources": {
            "requests": {
                "storage": "10Gi"
            }
        }
        },
    ~~~
  - Change port number
    - Add 1 infront of all port under 1000
    - Add following params to env of pachd deployment
        ~~~
        "env": [
            {
            "name": "WORKER_USES_ROOT",
            "value": "false"
            },
            {
            "name": "PORT",
            "value": "1650"
            },
            {
            "name": "HTTP_PORT",
            "value": "1652"
            },
            {
            "name": "PEER_PORT",
            "value": "1653"
            },
            {
            "name": "PPS_WORKER_GRPC_PORT",
            "value": "1680"
            },
        ...
        ~~~ 

- Add anyuid scc to default sa
   ~~~
   oc adm policy add-scc-to-user anyuid -z default
   ~~~
- Create all objects
  ~~~
  oc create -f manifest.json
  ~~~
  
   


**Issues**

- Port is not changed properly

    pachd-564d596988-rmlnf
    ~~~
    2021-05-31T14:07:16Z ERROR error setting up and/or running Prometheus Server: listen tcp :656: bind: permission denied
    2021-05-31T14:07:16Z WARNING pfs-over-HTTP - TLS disabled: could not stat public cert at /pachd-tls-cert/tls.crt: stat /pachd-tls-cert/tls.crt: no such file or directory
    2021-05-31T14:07:16Z WARNING s3gateway TLS disabled: could not stat public cert at /pachd-tls-cert/tls.crt: stat /pachd-tls-cert/tls.crt: no such file or directory
    2021-05-31T14:07:16Z ERROR error setting up and/or running S3 Server: listen tcp :600: bind: permission denied
    2021-05-31T14:07:16Z ERROR error setting up and/or running Githook Server: listen tcp :655: bind: permission denie
    ~~~

- adding anyuid scc is not a good idea 
