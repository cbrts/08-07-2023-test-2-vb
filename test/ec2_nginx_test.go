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

	// grunt_aws "github.com/gruntwork-io/terratest/modules/aws"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// Create an EC2 client
var ec2svc = ec2.New(session.New(&aws.Config{
	Region: aws.String("eu-west-1"),
}))

// Generate a unique ID for the state file
var uniqueId = random.UniqueId()

// Default to my private subnet but allow one to be passed in for other VPCs
var privateSubnetId = flag.String("privateSubnetId", "subnet-0cdaf467e3b2e0ea6", "Private Subnet ID to deploy to")
var backendBucket = flag.String("backendBucket", "cb-infra-states", "Backend S3 bucket to store teststate")
var backendKey = flag.String("backetKey", fmt.Sprintf("%s-terraform.tfstate", uniqueId), "Key for test state location")

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

	// Validate on the instance having only a private IP
	test_structure.RunTestStage(t, "validate_on_private_ip", func() {
		validateInstanceUsesPrivateIP(t, workingDir)
	})

	// Validate on nginx being reachable and returning a 200
	test_structure.RunTestStage(t, "validate_on_nginx_returning_200", func() {
		validateInstanceRunningNginx(t, workingDir)
	})

	// Validate on the instance only allowing port 80
	// test_structure.RunTestStage(t, "validate_on_ingress_80_only", func() {
	// 	validateInstanceIngressRules(t, workingDir)
	// })

	// Validate on the instance having the correct IAM profile
	test_structure.RunTestStage(t, "validate_on_correct_IAM_profile", func() {
		validateInstanceIAMProfile(t, workingDir)
	})
}

func deployUsingTerraform(t *testing.T, testRegion string, workingDir string) {
	// Get recommended instance type based on test region
	// testInstanceType := grunt_aws.GetRecommendedInstanceType(t, testRegion, []string{"t2.micro, t3.micro, t2.small, t3.small"})
	testInstanceType := "t3.micro"

	// Set the Terraform directory to init as well as variables to pass in for the test.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workingDir,

		Vars: map[string]interface{}{
			"instance_type":     testInstanceType,
			"private_subnet_id": *privateSubnetId,
		},
		BackendConfig: map[string]interface{}{
			"bucket": *backendBucket,
			"key":    *backendKey,
			"region": testRegion,
		},
	})

	// Save Terraform options so other test stages can re-use them.
	test_structure.SaveTerraformOptions(t, workingDir, terraformOptions)

	// Run `terraform init` and `terraform apply`
	terraform.InitAndApply(t, terraformOptions)
}

func validateInstanceRunningNginx(t *testing.T, workingDir string) {
	// Load the same options used in the deploy stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Fetch the generated DNS record from the ALB using `terraform output`
	loadbalancerDNSRecord := terraform.Output(t, terraformOptions, "hello_world")

	// Set the URL to verify against using the DNS record above
	instanceUrl := fmt.Sprintf("https://%s", loadbalancerDNSRecord)

	// Set up blank TLS config for use with http_helper
	tlsConfig := tls.Config{}

	// Set up retry configuration as the instance will take time to boot
	retryLimit := 30
	sleepBetweenRetry := 5 * time.Second

	// Verify that we get a 200 when a response is sent
	http_helper.HttpGetWithRetry(t, instanceUrl, &tlsConfig, 200, "", retryLimit, sleepBetweenRetry)
}

func validateInstanceUsesPrivateIP(t *testing.T, workingDir string) {
	// Load the same options used in the deploy stage
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
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}

	// Grab the NetworkInterface data for the instance
	instanceResponseData := *resp.Reservations[0].Instances[0].NetworkInterfaces[0]

	// Assert on a private IP being associated to the instance
	assert.NotEmpty(t, instanceResponseData.PrivateIpAddress)
}

func validateInstanceIAMProfile(t *testing.T, workingDir string) {

}

func undeployUsingTerraform(t *testing.T, workingDir string) {
	// Load the same options used in the deploy stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Destroy the test resources
	terraform.Destroy(t, terraformOptions)
}
