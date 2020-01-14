#!/usr/bin/env python

import click
import os
import sys
import re
from io import StringIO

@click.command()
@click.argument('cmd',
              default='init',
              type=click.Choice(['init', 'prep', 'ocp', 'update', 'clean', 'oneshot' , 'post']))
@click.option('-t','--type',
              default='inventory',
              type=click.Choice(['inventory', 'ocp', 'all', 'conf']))
@click.option('-op','--operate',
              default='apply',
              type=click.Choice(['apply', 'dtr']))
@click.help_option('--help', '-h')
@click.option('-v', '--verbose', count=True)
def launch(cmd=None,
           operate=None,
           status=None,
           type=None,
           clusterName=None,
           verbose=0):


    pat = re.compile("cluster_name")
    with open ('config', 'rt') as configfile:  
       for line in configfile:      
           if(pat.search(line) != None):
               clusterName=str(line.split("=")[1]).replace("\"", "").replace("\n","")
           

    if verbose > 0:
        verbosity = '-' + 'v' * verbose
    else:
        verbosity = ''

# Construct ansible command
    if cmd == 'init':
       status = os.system(
         'ansible-playbook %s  -i config prep/ansible/tasks/generate_config_files.yml --flush-cache; \
         cd prep;  \
         ansible-playbook %s  -i ansible/inventory ansible/tasks/cloud_init.yml  --flush-cache' \

                % (verbosity, verbosity)
       )

    if cmd == 'prep':
       if operate == 'apply':
          status = os.system(
           'cd prep;  \
            ansible-playbook %s  -i ansible/inventory ansible/playbooks/prep.yml  -e @ansible/defaults/main.yml --flush-cache; \
            terraform init ; \
            terraform get ; \
            terraform apply -auto-approve' 

                % (verbosity)
          )
       if operate == 'dtr':
          status = os.system(
           'cd prep;  \
            terraform destroy -auto-approve' 
          )

    if cmd == 'ocp':
       if operate == 'apply':
          status = os.system(
           'cd ocp4;  \
            terraform init ; \
            terraform get ; \
            terraform apply -auto-approve' 
          )
       if operate == 'dtr':
          status = os.system(
           'cd ocp4;  \
            terraform destroy -auto-approve' 
          )



    if cmd == 'oneshot':
       status = os.system(
        'ansible-playbook %s  -i config prep/ansible/tasks/generate_config_files.yml --flush-cache; \
         cd prep;  \
         ansible-playbook %s  -i ansible/inventory ansible/tasks/cloud_init.yml  --flush-cache; \
         ansible-playbook %s  -i ansible/inventory ansible/playbooks/prep.yml  -e @ansible/defaults/main.yml --flush-cache; \
         terraform init ; terraform get ; terraform apply -auto-approve; \
         cd ../ocp4;  terraform init ; terraform get ; terraform apply -auto-approve; \
         cd ../prep; sudo openshift-install --dir %s wait-for bootstrap-complete'

                % (verbosity, verbosity, verbosity, clusterName)
       )

    if cmd == 'post':
        status = os.system(
        'cd ./prep; \
         oc --config %s/auth/kubeconfig patch configs.imageregistry.operator.openshift.io cluster --type merge --patch \'{"spec":{"storage":{"emptyDir":{}}}}\'; \
         ansible-playbook %s -i ansible/inventory ansible/tasks/lb_rm_bootstrap.yml ; \
         sudo openshift-install --dir %s/ wait-for install-complete'

                % (clusterName, verbosity, clusterName)
       )


    if cmd == 'update':
       if type == 'inventory':
          status = os.system(
           'ansible-playbook %s -i config prep/ansible/tasks/generate_config_files.yml --flush-cache'
   
                   % (verbosity )
           )
       if type == 'ocp':
          status = os.system(
           'ansible-playbook %s  -i config prep/ansible/tasks/generate_config_files.yml --flush-cache; \
            cd prep;  \
            ansible-playbook -i ansible/inventory ansible/tasks/ocp_vm_config.yml  %s -e @ansible/defaults/main.yml'
                   % (verbosity, verbosity )
           )
       if type == 'ocp_module':
          status = os.system(
           'cd prep;  \
            ansible-playbook -i ansible/inventory ansible/tasks/ocp_module.yml  %s -e @ansible/defaults/main.yml'
                   % (verbosity, verbosity )
           )

    if cmd == 'clean':
       status = os.system(
        'cd ocp4 ; terraform destroy -auto-approve ; \
         cd ../prep ; terraform destroy -auto-approve ; \
         ansible-playbook %s -i ../config ansible/tasks/clean.yml --flush-cache'
                % (verbosity)
       )




    # Exit appropriately
    if os.WIFEXITED(status) and os.WEXITSTATUS(status) != 0:
        sys.exit(os.WEXITSTATUS(status))


if __name__ == '__main__':
    launch(auto_envvar_prefix='OCP4_UPI_KVM')

