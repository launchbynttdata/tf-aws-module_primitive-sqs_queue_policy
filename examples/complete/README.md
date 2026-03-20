# SQS Queue Policy Example

This example creates an SQS queue with KMS encryption and attaches a queue policy using the primitive module. It uses the resource naming module for consistent naming.

## Usage

```hcl
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

module "resource_names" {
  source   = "terraform.registry.launch.nttdata.com/module_library/resource_name/launch"
  version  = "~> 2.0"

  for_each = var.resource_names_map

  logical_product_family  = var.logical_product_family
  logical_product_service = var.logical_product_service
  class_env               = var.class_env
  instance_env            = var.instance_env
  instance_resource       = var.instance_resource
  cloud_resource_type     = each.value.name
  maximum_length          = each.value.max_length
  region                  = join("", split("-", data.aws_region.current.name))
}

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "aws_kms_key" "queue" {
  description             = "KMS key for SQS queue encryption"
  deletion_window_in_days = 7
}

resource "aws_kms_alias" "queue" {
  name          = "alias/sqs-queue-policy-example-${random_string.suffix.result}"
  target_key_id = aws_kms_key.queue.key_id
}

resource "aws_sqs_queue" "queue" {
  name                       = module.resource_names["sqsqueue1"].minimal_random_suffix
  kms_master_key_id          = aws_kms_key.queue.arn
  message_retention_seconds  = var.message_retention_seconds
  visibility_timeout_seconds = var.visibility_timeout_seconds
  tags                       = var.tags
}

locals {
  built_policy = jsonencode({
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

module "queue_policy" {
  source = "../.."

  queue_url = coalesce(var.queue_url, aws_sqs_queue.queue.url)
  policy    = coalesce(var.policy, local.built_policy)
}
```

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.9 |
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
| <a name="module_resource_names"></a> [resource\_names](#module\_resource\_names) | terraform.registry.launch.nttdata.com/module_library/resource_name/launch | ~> 2.0 |
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
| <a name="input_logical_product_family"></a> [logical\_product\_family](#input\_logical\_product\_family) | Logical product family for resource naming. | `string` | n/a | yes |
| <a name="input_logical_product_service"></a> [logical\_product\_service](#input\_logical\_product\_service) | Logical product service for resource naming. | `string` | n/a | yes |
| <a name="input_class_env"></a> [class\_env](#input\_class\_env) | Class environment for resource naming (e.g., dev, prod). | `string` | n/a | yes |
| <a name="input_instance_env"></a> [instance\_env](#input\_instance\_env) | Instance environment number for resource naming (0-999). | `number` | n/a | yes |
| <a name="input_instance_resource"></a> [instance\_resource](#input\_instance\_resource) | Instance resource number for resource naming (0-100). | `number` | n/a | yes |
| <a name="input_resource_names_map"></a> [resource\_names\_map](#input\_resource\_names\_map) | Map of resource types to naming config for the resource naming module. | <pre>map(object({<br/>    name       = string<br/>    max_length = number<br/>  }))</pre> | n/a | yes |
| <a name="input_queue_url"></a> [queue\_url](#input\_queue\_url) | The URL of the SQS queue to attach the policy to. If null, a queue is created in the example. | `string` | `null` | no |
| <a name="input_policy"></a> [policy](#input\_policy) | The JSON policy document. If null, a policy is built from the created queue. | `string` | `null` | no |
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
