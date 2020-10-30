# Build image
- enable Cloud Build Api
- enable Container Registry api

- Command
  ~~~
  gcloud builds submit --tag gcr.io/${GOOGLE_CLOUD_PROJECT}/quickstart-image .
  ~~~

Build command using yaml file

cloudbuild.yaml
~~~
steps:
- name: 'gcr.io/cloud-builders/docker'
  args: [ 'build', '-t', 'gcr.io/$PROJECT_ID/quickstart-image', '.' ]
images:
- 'gcr.io/$PROJECT_ID/quickstart-image'
~~~

Command:
```
gcloud builds submit --config cloudbuild.yaml .
```


# Kubectl

```
kubectl config view
```

```
gcloud container clusters get-credentials [Cluster_Name] --zone [zone_name]
```


#GKE Cluster by Cloud Shell

## GKE cluster up
```
export my_zone=us-central1-a
export my_cluster=standard-cluster-1

gcloud container clusters create $my_cluster --num-nodes 3 --zone $my_zone --enable-ip-alias
```

## GKE cluster modify
```
gcloud container clusters resize $my_cluster --zone $my_zone --size=4
```

## Create kubeconfig to connect to GKE
```
gcloud container clusters get-credentials $my_cluster --zone $my_zone
```

## Check GKE Cluster info
```
kubectl config view
kubectl cluster-info
kubectl config current-context
kubectl config get-contexts
```

## Change GKE Cluster
```
kubectl config use-context gke_${GOOGLE_CLOUD_PROJECT}_us-central1-a_standard-cluster-1
```

## GKE Node resource 
```
kubectl top nodes
```

## bash completion
```
source <(kubectl completion bash)
```

## file copy to container
```
kubectl cp ~/test.html $my_nginx_pod:/usr/share/nginx/html/test.html
```

## Expose pod
```
kubectl expose pod $my_nginx_pod --port 80 --type LoadBalancer
```

## Connect to a pod
```
kubectl exec -it new-nginx /bin/bash
```



## Check auth list
```
gcloud auth list
gcloud config list project
```


## Auto scale applications
```
kubectl autoscale deployment web --max 4 --min 1 --cpu-percent 1

kubectl get hpa

apiVersion: autoscaling/v1
kind: HorizontalPodAutoscaler
metadata:
  annotations:
    autoscaling.alpha.kubernetes.io/conditions:  [...]   
    autoscaling.alpha.kubernetes.io/current-metrics: [...]
  creationTimestamp: 2018-11-14T02:59:28Z
  name: web
  namespace: default
  resourceVersion: "14588"
  selfLink: /apis/autoscaling/v1/namespaces/[...]
spec:
  maxReplicas: 4
  minReplicas: 1
  scaleTargetRef:
    apiVersion: extensions/v1beta1
    kind: Deployment
    name: web
  targetCPUUtilizationPercentage: 1
status:
  currentCPUUtilizationPercentage: 0
  currentReplicas: 1
  desiredReplicas: 1

```

## Auto scale node
```
gcloud container node-pools create "temp-pool-1" \
--cluster=$my_cluster --zone=$my_zone \
--num-nodes "2" --node-labels=temp=true --preemptible
```

### only specific label should run on a node
```
kubectl taint node -l temp=true nodetype=preemptible:NoExecute

tolerations:
- key: "nodetype"
  operator: Equal
  value: "preemptible"
```


## Enable Network Policy
```
gcloud container clusters create $my_cluster --num-nodes 2 --enable-ip-alias --zone $my_zone --enable-network-policy
```
## Network Policy Example (ingress)
```
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: hello-allow-from-foo
spec:
  policyTypes:
  - Ingress
  podSelector:
    matchLabels:
      app: hello
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: foo
```
## Network Policy Example (egress)
```
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: foo-allow-to-hello
spec:
  policyTypes:
  - Egress
  podSelector:
    matchLabels:
      app: foo
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: hello
  - to:
    ports:
    - protocol: UDP
      port: 53
```

## Deploy Test pod
```
kubectl run test-1 --labels app=foo --image=alpine --restart=Never --rm --stdin --tty
```

## VPC native Kubernetes cluster
```
gcloud container clusters create $my_cluster \
   --num-nodes 3 --enable-ip-alias --zone $my_zone
```

# Enable Pod Security policy
```
gcloud beta container clusters update $my_cluster --zone $my_zone --enable-pod-security-policy
```

# Credential rotate (Add new IP)
```
gcloud container clusters update $my_cluster --zone $my_zone --start-credential-rotation
```

# Finish Credential rotate (remove previous IP)
```
gcloud container clusters update $my_cluster --zone $my_zone --complete-credential-rotation
```

# Role
## Create clusterrolebinding
```
kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user [USERNAME_1_EMAIL]
```

## create Rolebinding
```
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: username2-editor
  namespace: production
subjects:
- kind: User
  name: [USERNAME_2_EMAIL]
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```


# Create cloud sql instance
```
gcloud sql instances create sql-instance --tier=db-n1-standard-2 --region=us-central1
```

# Kubernetes explain
```
kubectl explain deployment
kubectl explain deployment --recursive
kubectl explain deployment.metadata.name
```

## Kubernetes Service IP 
```
curl -ks https://`kubectl get svc frontend -o=jsonpath="{.status.loadBalancer.ingress[0].ip}"`
```

## pod name and image version
```
kubectl get pods -o jsonpath --template='{range .items[*]}{.metadata.name}{"\t"}{"\t"}{.spec.containers[0].image}{"\n"}{end}'
```

# Gather project_id
```
export PROJECT_ID=$(gcloud config list --format 'value(core.project)')
```

# Push image
```
gcloud docker -- push gcr.io/${PROJECT_ID}/slack-codelab:v1
```

# comput instance list
```
gcloud compute instances list
```

##Regional Kubernetes Engine Cluster 
```
CLUSTER_VERSION=$(gcloud container get-server-config --region us-west1 --format='value(validMasterVersions[0])')

export CLOUDSDK_CONTAINER_USE_V1_API_CLIENT=false

gcloud container clusters create repd \
  --cluster-version=${CLUSTER_VERSION} \
  --machine-type=n1-standard-4 \
  --region=us-west1 \
  --num-nodes=1 \
  --node-locations=us-west1-a,us-west1-b,us-west1-c
  ```


## Regional StorageClass
```
kubectl apply -f - <<EOF
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: repd-west1-a-b-c
provisioner: kubernetes.io/gce-pd
parameters:
  type: pd-standard
  replication-type: regional-pd
  zones: us-west1-a, us-west1-b, us-west1-c
EOF
```

## Waiting for external IP assigned
```
while [[ -z $SERVICE_IP ]]; do SERVICE_IP=$(kubectl get svc wp-repd-wordpress -o jsonpath='{.status.loadBalancer.ingress[].ip}'); echo "Waiting for service external IP..."; sleep 2; done; echo http://$SERVICE_IP/admin
```

## Delete compute instance 
```
gcloud compute instance-groups managed delete ${IG} --zone ${ZONE}

// IG=gke-repd-default-pool-59739f96-grp
// ZONE=us-west1-a
```

# Create Tiller service account
```
kubectl create serviceaccount --namespace kube-system tiller
```

## Create clusterrolebinding
```
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
```

## Patch tiller to use service account
```
kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'  
```

## install helm tiller
```
helm init --service-account tiller --upgrade
```


# Allocate a Static IP
```
gcloud compute addresses create endpoints-ip --region us-central1
```

## Gather allocated static IPs
```
gcloud compute addresses list
```

# Deploy the Cloud Endpoints
```
gcloud endpoints services deploy openapi.yaml
```

# Cert (SSL with Let's Encrypt)
```
helm install --name cert-manager --version v0.3.2 \
    --namespace kube-system stable/cert-manager

export EMAIL=ahmet@example.com

cat letsencrypt-issuer.yaml | sed -e "s/email: ''/email: $EMAIL/g" | kubectl apply -f-
```
letsencrypt-issuer.yaml
```
# Copyright 2018 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
apiVersion: certmanager.k8s.io/v1alpha1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
spec:
  acme:
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    email: ''
    privateKeySecretRef:
      name: letsencrypt-staging
    http01: {}
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: ''
    privateKeySecretRef:
      name: letsencrypt-prod
    http01: {}
```

## ingress-tls.yaml
```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: esp-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    kubernetes.io/tls-acme: “true”
    certmanager.k8s.io/cluster-issuer: letsencrypt-prod
  labels:
    type: endpoints-app
spec:
  tls:
  - hosts:
    - api.endpoints.[MY-PROJECT].cloud.goog
    secretName: esp-tls
  rules:
  - host: api.endpoints.[MY-PROJECT].cloud.goog
    http:
      paths:
      - backend:
          serviceName: esp-srv
          servicePort: 80
```


# kubeflow 
## set up
```
export KUBEFLOW_TAG=0.4.0-rc.2

git clone https://github.com/kubeflow/examples.git
cd examples
git checkout v${KUBEFLOW_TAG}

pip install --user pyyaml

```
### Insstall ksonnet
```
export KS_VER=0.13.1
export KS_BIN=ks_${KS_VER}_linux_amd64
wget -O /tmp/${KS_BIN}.tar.gz https://github.com/ksonnet/ksonnet/releases/download/v${KS_VER}/${KS_BIN}.tar.gz
mkdir -p ${HOME}/bin
tar -xvf /tmp/${KS_BIN}.tar.gz -C ${HOME}/bin
export PATH=$PATH:${HOME}/bin/${KS_BIN}
```

## Install kfctl
```
wget -P /tmp https://github.com/kubeflow/kubeflow/archive/v${KUBEFLOW_TAG}.tar.gz
mkdir -p ${HOME}/src
tar -xvf /tmp/v${KUBEFLOW_TAG}.tar.gz -C ${HOME}/src
cd ${HOME}/src/kubeflow-${KUBEFLOW_TAG}/scripts
ln -s kfctl.sh kfctl
export PATH=$PATH:${HOME}/src/kubeflow-${KUBEFLOW_TAG}/scripts
cd ${HOME}
```


## Set GCP project ID
```
export PROJECT_ID=<gcp_project_id>

export ZONE=us-central1-a
gcloud config set project ${PROJECT_ID}
gcloud config set compute/zone ${ZONE}
gcloud config set container/new_scopes_behavior true
```

## Docker login
```
gcloud auth configure-docker
```

## Create GKE service account (for accessing to storage buckets)
```
export SERVICE_ACCOUNT=user-gcp-sa
export SERVICE_ACCOUNT_EMAIL=${SERVICE_ACCOUNT}@${PROJECT_ID}.iam.gserviceaccount.com
gcloud iam service-accounts create ${SERVICE_ACCOUNT} \
  --display-name "GCP Service Account for use with kubeflow examples"

gcloud projects add-iam-policy-binding ${PROJECT_ID} --member \
serviceAccount:${SERVICE_ACCOUNT_EMAIL} \
--role=roles/storage.admin
```

Generate a credential file for upload to the cluster
```
export KEY_FILE=${HOME}/secrets/${SERVICE_ACCOUNT_EMAIL}.json
gcloud iam service-accounts keys create ${KEY_FILE} \
  --iam-account ${SERVICE_ACCOUNT_EMAIL}
```

## Create storage bucket
```
export BUCKET_NAME=kubeflow-${PROJECT_ID}
gsutil mb -c regional -l us-central1 gs://${BUCKET_NAME}
```

## Create a cluster
```
kfctl init kubeflow-qwiklab --platform gcp --project ${PROJECT_ID}
```

```
cd kubeflow-qwiklab
kfctl generate platform
sed -i 's/n1-standard-8/n1-standard-4/g' gcp_config/cluster.jinja
kfctl apply platform
```

## Upload service account credentials
```
kubectl create secret generic user-gcp-sa \
  --from-file=user-gcp-sa.json="${KEY_FILE}"
```

Install kubeflow
```
kfctl generate k8s
kfctl apply k8s
```

#Delete cluster
```
gcloud container clusters delete kubeflow-qwiklab --zone=us-central1-a
```

#Create VM instance
```
gcloud compute instances create "my-vm-2" \
--machine-type "n1-standard-1" \
--image-project "debian-cloud" \
--image "debian-9-stretch-v20190213" \
--subnet "default"
```