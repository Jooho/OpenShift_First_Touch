# htpasswdConfig.sh #Directory_NAME"

if [[ $# != 1 ]];
then 
   echo "# Usage"
   echo "    ./htpasswdConfig.sh #Directory_NAME"
   echo "# Example"
   echo "    ./htpasswdConfig.sh rhv"
   exit  1
fi

export KUBECONFIG="$(pwd)/$1/auth/kubeconfig"

API_SERVER=$(oc status|grep api|awk -F'server' '{print $2}'|tr -d ' ')

echo -e "API SERVER: $API_SERVER"

if [[ z$API_SERVER == z ]];
then
  "Install Directory is not right"
fi

echo -n "User Name(joe):"
read username

echo -n  "Password(redhat):"
read password

if [[ z$username == z ]];
then 
   username=joe
fi

if [[ z$password == z ]];
then 
   password=redhat
fi


echo $username $password
htpasswd -c -B -b htpasswd joe redhat

oc create secret generic htpass-secret --from-file=htpasswd=./htpasswd -n openshift-config

echo "apiVersion: config.openshift.io/v1
kind: OAuth
metadata:
  name: cluster
spec:
  identityProviders:
  - name: my_htpasswd_provider
    mappingMethod: claim
    type: HTPasswd
    htpasswd:
      fileData:
        name: htpass-secret"| oc apply -f -

oc adm policy add-cluster-role-to-user cluster-admin joe

echo "...waiting for applying htpasswd"
sleep 10

echo "oc login --user $username --password $password"
oc login --user $username --password $password
