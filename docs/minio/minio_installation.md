# Minio Deployment 

## Manaual Install

- [Blog](https://blog.min.io/object_storage_as_a_service_on_minio/)
 
### Steps
 
 - Download [krew](https://krew.sigs.k8s.io/docs/user-guide/setup/install/)
    ~~~
    (
      set -x; cd "$(mktemp -d)" &&
      OS="$(uname | tr '[:upper:]' '[:lower:]')" &&
      ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" &&
      curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/krew.tar.gz" &&
      tar zxvf krew.tar.gz &&
      KREW=./krew-"${OS}_${ARCH}" &&
      "$KREW" install krew
    )
    
    vi ~/.bashrc
    export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"
    ~~~
  
 - Install minio plugin
    ~~~
    kubectl krew install minio
    ~~~

 - Deploy operator/console
    ~~~
    kubectl minio init
    ~~~
 
 - Add anyuid scc to minio-operator sa
    ~~~   
    oc adm policy add-scc-to-user anyuid -z minio-operator -n minio-operator
    ~~~

 - Deploy minio
    ~~~
    oc new-project minio-tenant-1  
    oc minio tenant create minio-tenant-1       --servers 2    --volumes 10    --capacity 5G    --namespace minio-tenant-1      --storage-class gp2

      Username: admin 
      Password: e5e1b12b-97b5-4af3-aec6-6f14fc694dca 
      Note: Copy the credentials to a secure location. MinIO will not display these again.

    +-------------+------------------------+----------------+--------------+--------------+
    | APPLICATION | SERVICE NAME           | NAMESPACE      | SERVICE TYPE | SERVICE PORT |
    +-------------+------------------------+----------------+--------------+--------------+
    | MinIO       | minio                  | minio-tenant-1 | ClusterIP    | 443          |
    | Console     | minio-tenant-1-console | minio-tenant-1 | ClusterIP    | 9443         |
    +-------------+------------------------+----------------+--------------+--------------+
    wait for 3~4mins
    ~~~
 
 - Access minio operator console
    ~~~
    oc minio proxy
    ~~~
 
 - Access minio console/ create a bucket
   - Get accesskey/secretkey
      ~~~
      oc get secret minio-tenant-1-creds-secret
      ~~~

      ~~~
      accesskey:  4df90cb2-5857-453b-8a59-e80e12fa1067
      oc get secret minio-tenant-1-creds-secret -n minio-tenant-1 -o jsonpath='{ .data.accesskey}'|base64 -d

      secretkey 0a6da0b9-5d7b-4a8c-af1c-eb768a66588a
      oc get secret minio-tenant-1-creds-secret -n minio-tenant-1 -o jsonpath='{ .data.secretkey}'|base64 -d
      ~~~
    - Access minio console
        ~~~
        oc port-forward svc/minio-tenant-1-console 9090:9443 -n minio-tenant-1
        
        google-chrome https://localhost:9090
        ~~~
  
    - Create a bucket
      ~~~
      bucket name:  pachyderm
      ~~~
