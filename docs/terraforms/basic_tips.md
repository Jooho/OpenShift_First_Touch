Terraforms Basic Tips
---------------------

- .gitignore
  ```  
  # Compiled files
  *.tfstate
  *.tfstate.backup
  
  # Module directory
  .terraform/

  .DS_Store

  ./idea
  ```


- Variable

  ```
  vi variable.tf


  variable "stringVariable" {
     type = "string"
     description = "String variable"

  }

  variable "mapVariable" {
     type =  "map"
     description = "Map variable"
     default = {
        test1 = "test1"
        test2 = "test2"
     }
  }
     
  variable "availablility_zone" {
     type =  "map"
     description = "Sample Map variable"
     default = {
        zone1 = "zone31"
        zone2 = "zone2"
     }
  variable "stringVariableFromTerminal" {}
  ```


- Variable file
  ```
  vi variableFile.tfvars

  stringVariable = "test"
  mapVariable = "{test = 'test'}"

  ```

  terraform plan -var-file="variableFile.tfvars"


- How to use variables in tf file
  ```
  resource "aws_subnet" "subnet" {
    availablility_zone = "${var.availability_zones["zone1"]}"
  
  }
  
  ```



- Required backend
  ```
  vi  remote-backend.tf

  terraform {
    required_version = "~> 0.10"
    
    backend "s3" {
      encrypt = true
      bucket = "terraform-backend-dalles"
      key = "ourdatastore/terraform.tfstate"
      region = "us-west-2"
    }
  }
  ```

- Resource
  ```
   resource "aws_vpc" "environment-example-two" {
     cidr_block = "10.0.0.0/16"
     enable_dns_hostnames = true
     enable_dns_support = true
     tags {
       Name = "terraform-aws-vpc-example-two"
     }
  }
  ```

- Provider
  ```
   vi connection.tf
  
   provider "aws" {
     region = "${var.region}"
   }
  ```


- Terraform cli 
  - terraform init
  - terraform plan
  - terraform apply
  - terraform validate
  - terraform fmt


- Compile time pass variables  (normal .tfvars file pass variable at execution time)
  **(Notice)**  terraform validate need some variables at compling time.

  - Copy all variables to terraforms.tfvars file

- start.sh
  ```
  vi start.sh
  
  terraform plan

  echo "yes"| terraform apply -var-file="./test.tfvars"
  
  ```

- TF environment variables
  ```
  TF_VAR_subscription_id="..."
  TF_VAR_client_id="..."
  TF_LOG=TRACE
  ```
