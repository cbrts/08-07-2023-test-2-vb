package test

import (
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

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

	// Assert on the test fixture
	test_structure.RunTestStage(t, "validate_terraform", func() {
		validateInstanceRunningNginx(t, workingDir)
	})
}

func deployUsingTerraform(t *testing.T, testRegion string, workingDir string) {
	// Generate a unique id to we don't have naming clashes
	uniqueId := random.UniqueId()

	// Set the instance name based on the tagging standard and unique ID
	testInstanceName := fmt.Sprintf("terratest-ec2-ew1-1a-dev-nginx-%s", uniqueId)

	// Get recommended instance type based on test region
	testInstanceType := aws.GetRecommendedInstanceType(t, testRegion, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})

	// Set the Terraform directory to init as well as variables to pass in for the test.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: workingDir,

		Vars: map[string]interface{}{
			"aws_region":    testRegion,
			"instance_name": testInstanceName,
			"instance_type": testInstanceType,
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
	loadbalancerDNSRecord := terraform.Output(t, terraformOptions, "loadbalancer_dns_record")

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

func undeployUsingTerraform(t *testing.T, workingDir string) {
	// Load the same options used in the deploy stage
	terraformOptions := test_structure.LoadTerraformOptions(t, workingDir)

	// Destroy the test resources
	terraform.Destroy(t, terraformOptions)
}
