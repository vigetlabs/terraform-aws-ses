package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// Test module validity
func TestValidate(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		TerraformDir: "../",
	}

	terraform.InitAndValidate(t, terraformOptions)
}

// Test outputs
func TestOutputs(t *testing.T) {
	t.Parallel()

	tc := struct {
		domain string
	}{domain: "example.com"}

	uniqueId := strings.ToLower(random.UniqueId())

	terraformOptions := &terraform.Options{
		TerraformDir: "../",
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
	dkimRecords := terraform.OutputList(t, terraformOptions, "ses_dkim_records")
	assert.Equal(t, 3, len(dkimRecords))
}

// Test sending pool name
func TestSendingPoolName(t *testing.T) {
	t.Parallel()
	uniqueId := strings.ToLower(random.UniqueId())

	// Set up pre-existing sending pool
	sendingPoolTerraformOptions := &terraform.Options{
		TerraformDir: "./fixtures/existing_sending_pool",
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

	// Test different scenarios
	testCases := []struct {
		domain   string
		poolName string
	}{
		{domain: "example.com", poolName: sendingPoolName},
		{domain: "example.com", poolName: ""},
	}

	for _, tc := range testCases {
		terraformOptions := &terraform.Options{
			TerraformDir: "../",
			Vars: map[string]interface{}{
				"namespace":             "vgt",
				"stage":                 "test",
				"component":             "ses",
				"attributes":            []string{uniqueId},
				"domain":                tc.domain,
				"ses_sending_pool_name": tc.poolName,
			},
		}

		defer terraform.Destroy(t, terraformOptions)
		terraform.InitAndApply(t, terraformOptions)

		// Validate sending pool name output
		sendingPool := terraform.Output(t, terraformOptions, "ses_sending_pool_name")
		assert.Equal(t, tc.poolName, sendingPool)
	}
}
