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
if [[ z$(oc get secret htpass-secret -n openshift-config|wc -l) != z ]];
then
  oc get secret htpass-secret -ojsonpath={.data.htpasswd} -n openshift-config | base64 -d > users.htpasswd
  export targetLine=$(cat users.htpasswd  -n |grep $username |awk '{print $1}')
  if [[ z${targetLine} != z ]];
  then
     sed "${targetLine}d" -i users.htpasswd
  fi

  htpasswd -bB users.htpasswd $username $password
  oc create secret generic htpass-secret --from-file=htpasswd=users.htpasswd --dry-run -o yaml -n openshift-config | oc replace -f -

else
  htpasswd -c -B -b htpasswd $username $password
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
fi


oc adm policy add-cluster-role-to-user cluster-admin $username

echo "...waiting for applying htpasswd"
sleep 10

echo "Try oc login --username $username --password $password"
