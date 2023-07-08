# 08-07-2023-test-2-vb

## What is this repo?
A terraform module for a take home interview task.

## What does it do?
* Deploys an AWS EC2 instance with a LoadBalancer (ELB) running nginx on port 80.
* Ensures the instance uses a private IP.
* Restricts Ingress to port 80 only
* Attaches an IAM instance profile for SSM sessions and `s3:List*`
* Asks the user for private_subnet_id and instance_type
* Uses a relevant tagging and naming standard

## Assumptions
* The Terraform code assumes that you already have a suitable VPC and subnets set up to deploy your instance in.
* The Terratest code assumes that the default region will be eu-west-1, usually test resources should be region agnostic but as this has been developed in my personal AWS account, I'm sticking to a single region.
* Terraform state is handeld by the user. For this example it's assuming you have an s3 bucket already prepared.
* SSM sessions have been enabled even though it's not in scope. This was to aid in debugging any issues with nginx.

## Requirements
This has been tested on a windows amd64 machine but should work on any OS as long as the correct binaries are installed.
* go1.20.5 windows/amd64
* Terraform v1.5.1 windows/amd64

## Running the example
If you want Terraform to prompt you for values, remove both the -var flags.
```
cd examples/ec2_nginx

terraform init -backend-config="bucket='YourS3Bucket'" \
-backend-config="key='YourS3BucketKey" \
-backend-config="region='YourS3BucketRegion"

terraform plan -out terraform.plan -var "private_subnet_id='YourSubnetID'" \
-var "instance_type='YourInstanceType'"

terrraform apply terraform.plan
```

## Tests
The tests assert on each requirement listed above. Replace the subnet-id with your current target subnet.
```
go test -v -privateSubnetId 'subnet-0cdaf467e3b2e0ea6'
```

## Issues
You may hit availabiliity issues when running the test suite.
```
None of the given instance types ([t2.micro, t3.micro, t2.small, t3.small]) is available in all the AZs in this region ([eu-west-1a eu-west-1b eu-west-1c]).
```
Hardcoding `testInstanceType` to a value of an instance you know is available in an AZ you want to deploy to will resolve it.
```
testInstanceType := "t3.micro"
```

## What should we improve?
#### Small improvements that would benefit the use case.
* Move user-data to cloud-init
* Set up a DNS record for the load balancer
* Make this multi-AZ
* Source an AMI with depedencies already installed
* Use a container for Nginx
* Use SSL
