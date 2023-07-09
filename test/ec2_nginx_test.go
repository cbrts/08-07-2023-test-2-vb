package test

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// Create an EC2 client
var ec2svc = ec2.New(session.New(&aws.Config{
	Region: aws.String("eu-west-1"),
}))

// Default to my network but allow details to be passed in for other VPC's
var privateSubnetId = flag.String("privateSubnetId", "subnet-0cdaf467e3b2e0ea6", "Private Subnet ID to deploy to")
var publicSubnetIds = flag.String("publicSubnetIds", "[\"subnet-8e94e3c6\", \"subnet-78a70622\"]", "Public Subnet IDs for the load balancer in array format")
var vpcID = flag.String("vpcID", "vpc-d9517bbf", "VPC ID to deploy to")

// Default to my bucket but allow other buckets to be used for test state
var backendBucket = flag.String("backendBucket", "cb-infra-states", "Backend S3 bucket to store teststate")
var backendBucketKey = flag.String("backendBucketKey", "test-states/local-terraform.tfstate", "Key for test state location")
var backendBucketRegion = flag.String("backendBucketRegion", "eu-west-1", "The S3 Bucket region")

func TestNginxInstance(t *testing.T) {
	t.Parallel()

	// Set directory for test fixture and test region.
	workingDir := "../examples/ec2_nginx"
	testRegion := "eu-west-1"

	// At the end of the test, destroy all test resources.
	defer test_structure.RunTestStage(t, "cleanup_terraform", func() {
		undeployUsingTerraform(t, workingDir)
	})

	// Run the test fixture
	test_structure.RunTestStage(t, "deploy_terraform", func() {
		deployUsingTerraform(t, testRegion, workingDir)
	})

	// Get Instance data post deployment
	instanceData := getInstanceData(t, workingDir)

	// Validate on the instance having only a private IP
	test_structure.RunTestStage(t, "validate_on_private_ip", func() {
		validateInstanceUsesPrivateIP(t, workingDir, instanceData)
	})

	// Validate on nginx being reachable and returning a 200
	test_structure.RunTestStage(t, "validate_on_nginx_returning_200", func() {
		validateInstanceRunningNginx(t, workingDir)
	})

	// Validate on the instance having the correct IAM profile
	test_structure.RunTestStage(t, "validate_on_correct_IAM_profile", func() {
		validateInstanceIAMProfile(t, workingDir, instanceData)
	})

	// Validate on the instance only allowing port 80
	test_structure.RunTestStage(t, "validate_on_ingress_80_only", func() {
		validateInstanceIngressRules(t, workingDir)
	})
}

func deployUsingTerraform(t *testing.T, testRegion string, workingDir string) {
	//  Set test instance Type to deploy
	testInstanceType := "t3.micro"

	// Set the Terraform directory to init as well as variables to pass in for the test.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workingDir,

		Vars: map[string]interface{}{
			"instance_type":     testInstanceType,
			"private_subnet_id": *privateSubnetId,
			"public_subnet_ids": *publicSubnetIds,
			"vpc_id":            *vpcID,
		},
		BackendConfig: map[string]interface{}{
			"bucket": *backendBucket,
			"key":    *backendBucketKey,
			"region": *backendBucketRegion,
		},
		Reconfigure: true,
	})

	// Save Terraform options so other test stages can re-use them.
	test_structure.SaveTerraformOptions(t, workingDir, terraformOptions)

	// Run `terraform init` and `terraform apply`
	terraform.InitAndApply(t, terraformOptions)
}

func getInstanceData(t *testing.T, workingDir string) *ec2.DescribeInstancesOutput {
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Fetch the target instanceID using `terraform output`
	instanceId := terraform.Output(t, terraformOptions, "instance_id")

	// Fetch the privateIP to assert on from the SDK
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: []*string{aws.String(instanceId)},
			},
		},
	}

	// Error if the SDK can't return results
	resp, err := ec2svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances", err.Error())
		log.Fatal(err.Error())
	}

	// Return data to use in other tests
	return resp
}

func customValidation(status int, _ string) bool {
	return status == 200
}

func validateInstanceRunningNginx(t *testing.T, workingDir string) {
	// Load the same options used in the deploy stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Fetch the generated DNS record from the ALB using `terraform output`
	loadbalancerDNSRecord := terraform.Output(t, terraformOptions, "load_balancer_dns_record")

	// Set the URL to verify against using the DNS record above
	instanceURL := fmt.Sprintf("http://%s", loadbalancerDNSRecord)

	// Set up blank TLS config for use with http_helper
	tlsConfig := tls.Config{}

	// Set up retry configuration as the instance will take time to boot
	retryLimit := 30
	sleepBetweenRetry := 5 * time.Second

	// Assert on if the request to the load balancer returns a 200 status code
	http_helper.HttpGetWithRetryWithCustomValidation(t, instanceURL, &tlsConfig, retryLimit, sleepBetweenRetry, customValidation)
}

func validateInstanceUsesPrivateIP(t *testing.T, workingDir string, instanceData *ec2.DescribeInstancesOutput) {
	// Fetch the NetworkInterface data for the instance
	networkInterface := instanceData.Reservations[0].Instances[0].NetworkInterfaces[0]

	// Assert on a private IP being associated to the instance
	assert.NotEmpty(t, networkInterface.PrivateIpAddress)
}

func validateInstanceIAMProfile(t *testing.T, workingDir string, instanceData *ec2.DescribeInstancesOutput) {
	// Fetch the target IAM Instance profile data for the instance
	actualIAMInstanceProfile := instanceData.Reservations[0].Instances[0].IamInstanceProfile.Arn

	// Load the same options used in the deploy stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Fetch the target IAM instance profile using `terraform output`
	expectedIAMInstanceProfile := terraform.Output(t, terraformOptions, "iam_instance_profile_arn")

	// Assert on the value being the same as the output
	assert.Equal(t, expectedIAMInstanceProfile, *actualIAMInstanceProfile)
}

func validateInstanceIngressRules(t *testing.T, workingDir string) {
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Fetch the target security Group IDs using `terraform output`
	albGroupID := terraform.Output(t, terraformOptions, "alb_sg_group_id")
	instanceGroupID := terraform.Output(t, terraformOptions, "instance_sg_group_id")

	// Fetch the ALB security group from the SDK
	albGroupParams := &ec2.DescribeSecurityGroupsInput{
		GroupIds: aws.StringSlice([]string{albGroupID}),
	}

	// Fetch the instance security group from the SDK
	instanceGroupParams := &ec2.DescribeSecurityGroupsInput{
		GroupIds: aws.StringSlice([]string{instanceGroupID}),
	}

	// Error if the SDK can't return results
	albResp, err := ec2svc.DescribeSecurityGroups(albGroupParams)
	if err != nil {
		fmt.Println("there was an error listing security groups", err.Error())
		log.Fatal(err.Error())
	}
	instanceResp, err := ec2svc.DescribeSecurityGroups(instanceGroupParams)
	if err != nil {
		fmt.Println("there was an error listing security groups", err.Error())
		log.Fatal(err.Error())
	}

	// Assert that the alb accepts traffic on port 80 from the internet
	assert.Equal(t, int64(80), *albResp.SecurityGroups[0].IpPermissions[0].FromPort)
	assert.Equal(t, int64(80), *albResp.SecurityGroups[0].IpPermissions[0].ToPort)
	assert.Equal(t, "0.0.0.0/0", *albResp.SecurityGroups[0].IpPermissions[0].IpRanges[0].CidrIp)

	// Assert that the instance accepts trafficon port 80 from ALB with a SG whitelist
	assert.Equal(t, int64(80), *instanceResp.SecurityGroups[0].IpPermissions[0].FromPort)
	assert.Equal(t, int64(80), *instanceResp.SecurityGroups[0].IpPermissions[0].ToPort)
	assert.NotEmpty(t, *instanceResp.SecurityGroups[0].IpPermissions[0].UserIdGroupPairs[0].GroupId)
}

func undeployUsingTerraform(t *testing.T, workingDir string) {
	// Load the same options used in the deploy stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Destroy the test resources
	terraform.Destroy(t, terraformOptions)
}
