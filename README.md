# 08-07-2023-test-2-vb

## What is this repo?
A terraform module for a take home interview task.

## What does it do?
* Deploys an AWS EC2 instance with a LoadBalancer (ALB) running nginx on port 80.
* Ensures the instance uses a private IP.
* Restricts Ingress to port 80 only
* Attaches an IAM instance profile for SSM sessions and `s3:List*`
* Asks the user for relevant deployment paramters
* Uses a relevant tagging and naming standard

## Assumptions
* The Terraform code assumes that you already have a suitable VPC and subnets set up to deploy your instance in.
* The Terratest code assumes that the default region will be eu-west-1, usually test resources should be region agnostic but as this has been developed in my personal AWS account, I'm sticking to a single region.
* Terraform state is handeld by the user. For this example it's assuming you have an s3 bucket already prepared.
* SSM sessions have been enabled even though it's not in scope. This was to aid in debugging any issues with nginx.
* It's assumed that you are authenticated to your own AWS account.

## Requirements
This has been tested on a windows amd64 machine but should work on any OS as long as the correct binaries are installed.
* go1.20.5 windows/amd64
* Terraform v1.5.1 windows/amd64

## Running the example manually
If you want Terraform to prompt you for values, remove both the -var flags. Replace the values of each backend-config and var with your environment.
```
cd examples/ec2_nginx

terraform init -backend-config="bucket=cb-infra-states" \
-backend-config="key=test-states/local.state" \
-backend-config="region='eu-west-1"

terraform plan -out terraform.plan -var "private_subnet_id=subnet-0cdaf467e3b2e0ea6" \
-var "instance_type=t3.micro" \
-var "vpc_id=vpc-d9517bbf" \
-var "public_subnet_ids=[\"subnet-8e94e3c6\", \"subnet-78a70622\"]"

terrraform apply terraform.plan

terraform destroy -var "private_subnet_id=subnet-0cdaf467e3b2e0ea6" \
-var "instance_type=t3.micro" \
-var "vpc_id=vpc-d9517bbf" \
-var "public_subnet_ids=[\"subnet-8e94e3c6\", \"subnet-78a70622\"]"
```

## Tests
The tests assert on each requirement listed above. Replace the values of each backend-config and var with your environment.
```
go test -v -publicSubnetIds "[\"subnet-8e94e3c6\", \"subnet-78a70622\"]" \
-vpcID "vpc-d9517bbf" \
-privateSubnetId "subnet-0cdaf467e3b2e0ea6" \
-backendBucket "cb-infra-states" \
-backendKey "test-states/local-terraform.tfstate" \
-backendBucketRegion "eu-west-1"
```
To only run only the tests and ignore the Terraform lifecycle you can set:
```
export SKIP_deploy_terraform=true \
export SKIP_cleanup_terraform=true
```

## What should we improve?
#### Small improvements that would benefit the use case.
* Move user-data to cloud-init
* Set up a DNS record for the load balancer
* Make this multi-AZ
* Source an AMI with depedencies already installed
* Use a container for Nginx
* Use SSL
