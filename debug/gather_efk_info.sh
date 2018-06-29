
export oc_user=$1
export oc_password=$2
export oc_server_url=$3

function usage(){
cat <<EOF
# Usage :
#        ./gather_efk_info.sh oc_user oc_password oc_api_server_url
#        ./gather_efk_info.sh admin admin https://master.example.com:8443
#        ./gather_efk_info.sh system:admin https://master.example.com:8443

EOF
}


# start script
if [[ $1 == "" ]]
then
  usage
  exit 1  
fi


wget https://raw.githubusercontent.com/openshift/origin-aggregated-logging/master/hack/logging-dump.sh
chmod +x logging-dump.sh
if [[ $oc_user == 'system:admin' ]]
then
   oc_server_url=$2
   oc login -u $oc_user $oc_server_url
else
   oc login -u $oc_user -p $oc_password $oc_server_url
fi

export NAMESPACE=logging
./logging-dump.sh

echo "please upload 'logging-$Date' folder"

# xz --decompress file.xz
