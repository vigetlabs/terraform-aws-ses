package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// Test module validity
func TestValidate(t *testing.T) {
	t.Parallel()

	// Copy the terraform folder to a temp folder
	tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, "..", ".")
	defer os.RemoveAll(tempTestFolder)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,
	}

	terraform.InitAndValidate(t, terraformOptions)
}

// Test outputs
func TestOutputs(t *testing.T) {
	t.Parallel()
	uniqueId := strings.ToLower(random.UniqueId())
	domain := fmt.Sprintf("example%s.com", uniqueId)

	tc := struct {
		domain string
	}{domain: domain}

	// Copy the terraform folder to a temp folder
	tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, "..", ".")
	defer os.RemoveAll(tempTestFolder)

	terraformOptions := &terraform.Options{
		TerraformDir: tempTestFolder,
		Vars: map[string]interface{}{
			"namespace":  "vgt",
			"stage":      "test",
			"component":  "ses",
			"attributes": []string{uniqueId},
			"domain":     tc.domain,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Validate email identity output
	emailIdentity := terraform.Output(t, terraformOptions, "email_identity")
	assert.Equal(t, tc.domain, emailIdentity)

	// Validate iam sending group output
	sendingGroup := terraform.Output(t, terraformOptions, "iam_sending_group_name")
	assert.Equal(t, fmt.Sprintf("vgt-test-ses-%s", uniqueId), sendingGroup)

	// Validate ses dkim records output
	dkimRecords := terraform.OutputList(t, terraformOptions, "dkim_records")
	assert.Equal(t, 3, len(dkimRecords))
}

// Test sending pool with pre-existing sending pool
func TestExistingSendingPoolName(t *testing.T) {
	t.Parallel()
	uniqueId := strings.ToLower(random.UniqueId())

	// Copy the sending pool folder to a temp folder
	tempSendingPoolFolder := test_structure.CopyTerraformFolderToTemp(t, ".", "fixtures/existing_sending_pool")
	defer os.RemoveAll(tempSendingPoolFolder)

	// Set up pre-existing sending pool
	sendingPoolTerraformOptions := &terraform.Options{
		TerraformDir: tempSendingPoolFolder,
		Vars: map[string]interface{}{
			"namespace":  "vgt",
			"stage":      "test",
			"component":  "ses",
			"attributes": []string{uniqueId},
		},
	}

	defer terraform.Destroy(t, sendingPoolTerraformOptions)
	terraform.InitAndApply(t, sendingPoolTerraformOptions)

	sendingPoolName := terraform.Output(t, sendingPoolTerraformOptions, "pool_name")
	domain := fmt.Sprintf("example%s.com", uniqueId)

	// Test different scenarios
	testCases := []struct {
		domain   string
		poolName string
	}{
		{domain: domain, poolName: sendingPoolName},
		{domain: domain, poolName: ""},
	}

	// Copy the terraform folder to a temp folder
	tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, "..", ".")
	defer os.RemoveAll(tempTestFolder)

	for _, tc := range testCases {
		terraformOptions := &terraform.Options{
			TerraformDir: tempTestFolder,
			Vars: map[string]interface{}{
				"namespace":         "vgt",
				"stage":             "test",
				"component":         "ses",
				"attributes":        []string{uniqueId},
				"domain":            tc.domain,
				"sending_pool_name": tc.poolName,
			},
		}

		defer terraform.Destroy(t, terraformOptions)
		terraform.InitAndApply(t, terraformOptions)

		// Validate sending pool name output
		sendingPool := terraform.Output(t, terraformOptions, "sending_pool_name")
		assert.Equal(t, tc.poolName, sendingPool)
	}
}

// Test sending pool
func TestSendingPool(t *testing.T) {
	t.Parallel()
	uniqueId := strings.ToLower(random.UniqueId())
	domain := fmt.Sprintf("example%s.com", uniqueId)

	// Test different scenarios
	testCases := []struct {
		domain      string
		poolName    string
		poolEnabled bool
	}{
		// Test sending pool disabled
		{domain: domain, poolName: "", poolEnabled: false},
		// Test sending pool enabled
		{domain: domain, poolName: fmt.Sprintf("vgt-test-ses-%s", uniqueId), poolEnabled: true},
		// Test sending pool enabled with name
		{domain: domain, poolName: fmt.Sprintf("test-%s", uniqueId), poolEnabled: true},
	}

	// Copy the terraform folder to a temp folder
	tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, "..", ".")
	defer os.RemoveAll(tempTestFolder)

	for _, tc := range testCases {
		terraformOptions := &terraform.Options{
			TerraformDir: tempTestFolder,
			Vars: map[string]interface{}{
				"namespace":           "vgt",
				"stage":               "test",
				"component":           "ses",
				"attributes":          []string{uniqueId},
				"domain":              tc.domain,
				"sending_pool_name":   tc.poolName,
				"create_sending_pool": tc.poolEnabled,
			},
		}

		defer terraform.Destroy(t, terraformOptions)
		terraform.InitAndApply(t, terraformOptions)

		// Validate sending pool name output
		sendingPool := terraform.Output(t, terraformOptions, "sending_pool_name")
		assert.Equal(t, tc.poolName, sendingPool)
	}
}
