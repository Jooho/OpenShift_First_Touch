export NODE_NAME=$1
export SSH_BASTION_NAMESPACE=openshift-ssh-bastion
export PRIVATE_KEY=~/.ssh/id_rsa

oc project ${SSH_BASTION_NAMESPACE} &> /dev/null
if [[ $? == 1 ]] 
then
  if [ -f "./deploy.sh" ] 
  then
    ./deploy.sh
  else
    curl https://raw.githubusercontent.com/eparis/ssh-bastion/master/deploy/deploy.sh | bash
  fi
fi

oc get service ssh-bastion -n ${SSH_BASTION_NAMESPACE}  -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'

ssh -i ${PRIVATE_KEY} -o StrictHostKeyChecking=no -o ProxyCommand='ssh -i ${PRIVATE_KEY} -A -o StrictHostKeyChecking=no -o ServerAliveInterval=30 -W %h:%p core@$(oc get service -n openshift-ssh-bastion ssh-bastion -o jsonpath="{.status.loadBalancer.ingress[0].hostname}")' core@${NODE_NAME}


