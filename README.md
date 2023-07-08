# 08-07-2023-test-2-vb

## What is this repo?
A terraform module for a take home interview task.

## What does it do?
deploys an AWS EC2 instance with a LoadBalancer (ELB) running nginx on port 80

## Assumptions
* The Terraform code assumes that you already have a suitable VPC and subnets set up to deploy your instance in.
* The Terratest code assumes that the default region will be eu-west-1, usually test resources should be region agnostic but as this has been developed in my personal AWS account, I'm sticking to a single region.
* SSM sessions have been enabled even though it's not in scope. This was to aid in debugging any issues with nginx.


## Requirements
This has been tested on a windows amd64 machine but should work on any OS as long as the correct binaries are installed.
* go1.20.5 windows/amd64
* Terraform v1.5.1 windows/amd64
