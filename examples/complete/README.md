# SQS Queue Policy Example

This example creates an SQS queue with KMS encryption and attaches a queue policy using the primitive module.

## Usage

```hcl
module "queue_policy" {
  source = "../.."

  queue_url = aws_sqs_queue.queue.url
  policy    = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowCurrentAccountSendMessage"
        Effect = "Allow"
        Principal = {
          AWS = data.aws_caller_identity.current.account_id
        }
        Action   = "sqs:SendMessage"
        Resource = aws_sqs_queue.queue.arn
      }
    ]
  })
}
```

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_message_retention_seconds"></a> [message_retention_seconds](#input_message_retention_seconds) | The number of seconds Amazon SQS retains a message. | `number` | `345600` | no |
| <a name="input_visibility_timeout_seconds"></a> [visibility_timeout_seconds](#input_visibility_timeout_seconds) | The visibility timeout for the queue. | `number` | `30` | no |
| <a name="input_tags"></a> [tags](#input_tags) | Map of tags to assign to the resources. | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_id"></a> [id](#output_id) | The ID of the queue policy (same as the queue URL). |
| <a name="output_queue_url"></a> [queue_url](#output_queue_url) | The URL of the SQS queue. |
| <a name="output_policy"></a> [policy](#output_policy) | The policy document attached to the queue. |

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.5 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.14 |
| <a name="requirement_random"></a> [random](#requirement\_random) | ~> 3.6 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_random"></a> [random](#provider\_random) | 3.8.1 |
| <a name="provider_aws"></a> [aws](#provider\_aws) | 5.100.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_queue_policy"></a> [queue\_policy](#module\_queue\_policy) | ../.. | n/a |

## Resources

| Name | Type |
|------|------|
| [aws_kms_alias.queue](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kms_alias) | resource |
| [aws_kms_key.queue](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/kms_key) | resource |
| [aws_sqs_queue.queue](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue) | resource |
| [random_string.suffix](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/string) | resource |
| [aws_caller_identity.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/caller_identity) | data source |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_message_retention_seconds"></a> [message\_retention\_seconds](#input\_message\_retention\_seconds) | The number of seconds Amazon SQS retains a message. | `number` | `345600` | no |
| <a name="input_visibility_timeout_seconds"></a> [visibility\_timeout\_seconds](#input\_visibility\_timeout\_seconds) | The visibility timeout for the queue. | `number` | `30` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Map of tags to assign to the resources. | `map(string)` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_id"></a> [id](#output\_id) | The ID of the queue policy (same as the queue URL). |
| <a name="output_queue_url"></a> [queue\_url](#output\_queue\_url) | The URL of the SQS queue. |
| <a name="output_policy"></a> [policy](#output\_policy) | The policy document attached to the queue. |
<!-- END_TF_DOCS -->
