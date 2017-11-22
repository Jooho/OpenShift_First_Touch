ansible-playbook  -vvvv   ./efs.yaml \
-e openshift_provisioners_install_provisioners=True \
-e openshift_provisioners_efs=True \
-e openshift_provisioners_efs_fsid=fs-XXXX\
-e openshift_provisioners_efs_region=us-east-1 \
-e openshift_provisioners_efs_aws_access_key_id=XXXX \
-e openshift_provisioners_efs_aws_secret_access_key='XXXX' \
-e openshift_provisioners_project=openshift-infra \
-e '{openshift_provisioners_efs_nodeselector: {role: "infra"}}' \
-e openshift_provisioners_image_prefix=openshift3/ \
-e openshift_provisioners_image_version=v3.6 \
-e openshift_provisioners_efs_path=/ 

