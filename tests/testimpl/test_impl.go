package testimpl

import (
	"context"
	"encoding/json"
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
		assert.Equal(t, expectedDoc["Version"], actualDoc["Version"], "Policy Version must match")
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
		assert.Equal(t, expectedDoc["Version"], actualDoc["Version"], "Policy Version must match")
	})
}
