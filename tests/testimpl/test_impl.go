package testimpl

import (
	"context"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// iamRootARN matches arn:aws:iam::ACCOUNT:root to extract the account ID.
var iamRootARN = regexp.MustCompile(`^arn:aws:iam::(\d{12}):root$`)

// bareAccountID matches a 12-digit AWS account ID.
var bareAccountID = regexp.MustCompile(`^\d{12}$`)

// extractAccountIDsFromPrincipalAWS returns the set of account IDs from a
// Principal.AWS value (string or []string). Handles both "123456789012" and
// "arn:aws:iam::123456789012:root" forms.
func extractAccountIDsFromPrincipalAWS(v interface{}) []string {
	var ids []string
	switch val := v.(type) {
	case string:
		if m := iamRootARN.FindStringSubmatch(val); len(m) > 0 {
			ids = append(ids, m[1])
		} else if bareAccountID.MatchString(val) {
			ids = append(ids, val)
		}
	case []interface{}:
		for _, item := range val {
			ids = append(ids, extractAccountIDsFromPrincipalAWS(item)...)
		}
	}
	return ids
}

// assertPolicyStatementsMatch asserts that the attached policy has the same
// semantic content as the expected policy: Version, Effect, Action, Resource,
// Sid, and Principal.AWS (by account ID equivalence).
func assertPolicyStatementsMatch(t *testing.T, expectedDoc, actualDoc map[string]interface{}) {
	t.Helper()
	assert.Equal(t, expectedDoc["Version"], actualDoc["Version"], "Policy Version must match")

	expectedStmts, ok := expectedDoc["Statement"].([]interface{})
	require.True(t, ok, "Expected policy must have Statement array")
	actualStmts, ok := actualDoc["Statement"].([]interface{})
	require.True(t, ok, "Actual policy must have Statement array")
	require.Len(t, actualStmts, len(expectedStmts), "Statement count must match")

	for i := range expectedStmts {
		exp, ok := expectedStmts[i].(map[string]interface{})
		require.True(t, ok)
		act, ok := actualStmts[i].(map[string]interface{})
		require.True(t, ok)

		assert.Equal(t, exp["Effect"], act["Effect"], "Statement[%d] Effect must match", i)
		assert.Equal(t, exp["Action"], act["Action"], "Statement[%d] Action must match", i)
		assert.Equal(t, exp["Resource"], act["Resource"], "Statement[%d] Resource must match", i)
		assert.Equal(t, exp["Sid"], act["Sid"], "Statement[%d] Sid must match", i)

		expPrincipal, ok := exp["Principal"].(map[string]interface{})
		require.True(t, ok, "Expected policy Statement[%d] must have Principal", i)
		actPrincipal, ok := act["Principal"].(map[string]interface{})
		require.True(t, ok, "Actual policy Statement[%d] must have Principal", i)

		expAWS := expPrincipal["AWS"]
		actAWS := actPrincipal["AWS"]
		expAccountIDs := extractAccountIDsFromPrincipalAWS(expAWS)
		actAccountIDs := extractAccountIDsFromPrincipalAWS(actAWS)
		require.Len(t, expAccountIDs, 1, "Expected policy Statement[%d] Principal.AWS must resolve to one account", i)
		require.Len(t, actAccountIDs, 1, "Actual policy Statement[%d] Principal.AWS must resolve to one account", i)
		assert.Equal(t, expAccountIDs[0], actAccountIDs[0], "Statement[%d] Principal.AWS must refer to same account", i)
	}
}

func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	t.Run("VerifyTerraformOutputs", func(t *testing.T) {
		queueURL := terraform.Output(t, ctx.TerratestTerraformOptions(), "queue_url")
		policyOutput := terraform.Output(t, ctx.TerratestTerraformOptions(), "policy")
		idOutput := terraform.Output(t, ctx.TerratestTerraformOptions(), "id")

		require.NotEmpty(t, queueURL, "queue_url output must not be empty")
		require.NotEmpty(t, policyOutput, "policy output must not be empty")
		assert.Equal(t, queueURL, idOutput, "id and queue_url should be the same")
	})

	t.Run("VerifyPolicyViaAWSAPI", func(t *testing.T) {
		queueURL := terraform.Output(t, ctx.TerratestTerraformOptions(), "queue_url")
		expectedPolicy := terraform.Output(t, ctx.TerratestTerraformOptions(), "policy")

		cfg, err := config.LoadDefaultConfig(context.Background())
		require.NoError(t, err)

		client := sqs.NewFromConfig(cfg)
		attrs, err := client.GetQueueAttributes(context.Background(), &sqs.GetQueueAttributesInput{
			QueueUrl:       aws.String(queueURL),
			AttributeNames: []sqsTypes.QueueAttributeName{sqsTypes.QueueAttributeNamePolicy},
		})
		require.NoError(t, err)

		actualPolicy, ok := attrs.Attributes[string(sqsTypes.QueueAttributeNamePolicy)]
		require.True(t, ok, "Policy attribute must be present on the queue")
		require.NotEmpty(t, actualPolicy, "Policy must not be empty")

		var expectedDoc, actualDoc map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(expectedPolicy), &expectedDoc))
		require.NoError(t, json.Unmarshal([]byte(actualPolicy), &actualDoc))
		assertPolicyStatementsMatch(t, expectedDoc, actualDoc)
	})

	t.Run("SendMessageToQueue", func(t *testing.T) {
		queueURL := terraform.Output(t, ctx.TerratestTerraformOptions(), "queue_url")

		cfg, err := config.LoadDefaultConfig(context.Background())
		require.NoError(t, err)

		client := sqs.NewFromConfig(cfg)
		_, err = client.SendMessage(context.Background(), &sqs.SendMessageInput{
			QueueUrl:    aws.String(queueURL),
			MessageBody: aws.String("test message from functional test"),
		})
		require.NoError(t, err)
	})
}

func TestComposableCompleteReadonly(t *testing.T, ctx types.TestContext) {
	t.Run("VerifyTerraformOutputs", func(t *testing.T) {
		queueURL := terraform.Output(t, ctx.TerratestTerraformOptions(), "queue_url")
		policyOutput := terraform.Output(t, ctx.TerratestTerraformOptions(), "policy")
		idOutput := terraform.Output(t, ctx.TerratestTerraformOptions(), "id")

		require.NotEmpty(t, queueURL, "queue_url output must not be empty")
		require.NotEmpty(t, policyOutput, "policy output must not be empty")
		assert.Equal(t, queueURL, idOutput, "id and queue_url should be the same")
	})

	t.Run("VerifyPolicyViaAWSAPI", func(t *testing.T) {
		queueURL := terraform.Output(t, ctx.TerratestTerraformOptions(), "queue_url")
		expectedPolicy := terraform.Output(t, ctx.TerratestTerraformOptions(), "policy")

		cfg, err := config.LoadDefaultConfig(context.Background())
		require.NoError(t, err)

		client := sqs.NewFromConfig(cfg)
		attrs, err := client.GetQueueAttributes(context.Background(), &sqs.GetQueueAttributesInput{
			QueueUrl:       aws.String(queueURL),
			AttributeNames: []sqsTypes.QueueAttributeName{sqsTypes.QueueAttributeNamePolicy},
		})
		require.NoError(t, err)

		actualPolicy, ok := attrs.Attributes[string(sqsTypes.QueueAttributeNamePolicy)]
		require.True(t, ok, "Policy attribute must be present on the queue")
		require.NotEmpty(t, actualPolicy, "Policy must not be empty")

		var expectedDoc, actualDoc map[string]interface{}
		require.NoError(t, json.Unmarshal([]byte(expectedPolicy), &expectedDoc))
		require.NoError(t, json.Unmarshal([]byte(actualPolicy), &actualDoc))
		assertPolicyStatementsMatch(t, expectedDoc, actualDoc)
	})
}
