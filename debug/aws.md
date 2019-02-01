## [TAG]

- Get Tag with VPC id
```
# Replace $VPC_ID
aws --region us-east-1 ec2 describe-vpcs --vpc-ids $VPC_ID
#ex) aws --region us-east-1 ec2 describe-vpcs --vpc-ids vpc-04aa480d1bf9ef648

```

- Get Tag with role name
```
# Replace $ROLE_NAME

aws --region us-east-1 iam list-role-tags --role-name $ROLE_NAME
#ex) aws --region us-east-1 iam list-role-tags --role-name ocp4-bootstrap-role
#False
#TAGS	openshiftClusterID	b78c01d7-xxxx
#TAGS	kubernetes.io/cluster/ocp4	owned
```


## [Resources]

- All resources with Tag
```
# KEY is from metadata.json

{"clusterName":"ocp4","clusterID":"b78c01d7-xxx","aws":{"region":"us-east-1","identifier":[{"openshiftClusterID":"b78c01d7-xxxx"},{"kubernetes.io/cluster/ocp4":"owned"}]}}    # kubernetes.io/cluster/ocp4 

# Replace $KEY
 aws --region us-east-1 resourcegroupstaggingapi get-resources --query "ResourceTagMappingList[?Tags[? Key == '$KEY' && Value == 'owned']].ResourceARN" --output text
 ex) aws --region us-east-1 resourcegroupstaggingapi get-resources --query "ResourceTagMappingList[?Tags[? Key == 'kubernetes.io/cluster/ocp4' && Value == 'owned']].ResourceARN" --output text
 ```
 
 ## [CLUSTER_UUID]
 
 - Get CLUSTER UUID by cluster name
 ```
 # Replace $CLUSTER_NAME
  $ aws ec2 --region us-east-1 describe-vpcs --output json | jq '.Vpcs[] | {"name": ([.Tags[] | select(.Key == "Name") | .Value][0]), "openshiftClusterID": ([.Tags[] | select(.Key == "openshiftClusterID") | .Value][0])} | select(.name | contains("$CLUSTER_NAME"))' 
  ex) $ aws ec2 --region us-east-1 describe-vpcs --output json | jq '.Vpcs[] | {"name": ([.Tags[] | select(.Key == "Name") | .Value][0]), "openshiftClusterID": ([.Tags[] | select(.Key == "openshiftClusterID") | .Value][0])} | select(.name | contains("ocp4"))'
  ```
 
 ## [PROFILE]
 
 - Get profiles by name starting with cluster name
 ```
 # Replace $CLUSTER_NAME
 aws --region us-east-1 iam list-instance-profiles --output json | jq '[.InstanceProfiles[] | select(.InstanceProfileName | startswith("${CLUSTER_NAME}-"))]'
 ex) aws --region us-east-1 iam list-instance-profiles --output json | jq '[.InstanceProfiles[] | select(.InstanceProfileName | startswith("ocp4-"))]'
 ```
 
 - Delete profile
 ```
 # Replace $PROFILE_NAME
 aws --region us-east-1 iam delete-instance-profile --instance-profile-name $PROFILE_NAME
 ex) aws --region us-east-1 iam delete-instance-profile --instance-profile-name ocp4-bootstrap-profile
 ```
 
 - Delete roles from instance profile
 **NOTE: If some roles attached to a profile, it can not be deleted unless the roles are removed**
  ```
  aws --region us-east-1 iam remove-role-from-instance-profile --instance-profile-name $PROFILE_NAME --role-name $ROLE_NAME
  ex)aws --region us-east-1 iam remove-role-from-instance-profile --instance-profile-name ocp4-on-aws-master-profile --role-name ocp4-on-aws-master-role
  ```
 
