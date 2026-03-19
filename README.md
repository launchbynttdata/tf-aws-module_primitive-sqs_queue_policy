# Terraform AWS SQS Queue Policy Module

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)

## Overview

This Terraform module wraps the [aws_sqs_queue_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue_policy) resource to attach an IAM policy document to an SQS queue. The module validates that the policy document contains a Version identifier (`2012-10-17` or `2008-10-17`) as required by AWS IAM policy documents.

## Pre-Commit hooks

[.pre-commit-config.yaml](.pre-commit-config.yaml) defines pre-commit hooks for Terraform, Go, and common linting tasks.

`commitlint` enforces conventional commit message format. See [commitlint-config-conventional](https://github.com/conventional-changelog/commitlint/tree/master/@commitlint/config-conventional#type-enum) for the type enum.

`detect-secrets-hook` prevents new secrets from being introduced into the baseline.

Install hooks:

```
pre-commit install
pre-commit install --hook-type commit-msg
```

## Usage

```hcl
module "queue_policy" {
  source = "terraform.registry.launch.nttdata.com/module_primitive/sqs_queue_policy/aws"
  version = "~> 1.0"

  queue_url = aws_sqs_queue.example.url
  policy    = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowSendMessage"
        Effect = "Allow"
        Principal = {
          AWS = "arn:aws:iam::123456789012:root"
        }
        Action   = "sqs:SendMessage"
        Resource = aws_sqs_queue.example.arn
      }
    ]
  })
}
```

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.5 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.14 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | ~> 5.14 |

## Resources

| Name | Type |
|------|------|
| [aws_sqs_queue_policy.queue_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue_policy) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_queue_url"></a> [queue_url](#input\_queue\_url) | The URL of the SQS Queue to which to attach the policy. | `string` | n/a | yes |
| <a name="input_policy"></a> [policy](#input\_policy) | The JSON policy document for the SQS queue. Must include a Version identifier (e.g., "2012-10-17" or "2008-10-17") as the top-level key per AWS IAM policy document requirements. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_id"></a> [id](#output\_id) | The ID of the resource (same as the queue URL). |
| <a name="output_queue_url"></a> [queue_url](#output\_queue\_url) | The URL of the SQS queue to which the policy is attached. |
| <a name="output_policy"></a> [policy](#output\_policy) | The policy document attached to the queue. |

## Policy Document Validation

The module validates that the policy document:

1. Is valid JSON
2. Contains a top-level `Version` key with value `2012-10-17` or `2008-10-17`

AWS may hang indefinitely when creating or updating an SQS queue policy without an explicit Version identifier. See the [Terraform AWS provider documentation](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue_policy) for details.

## Testing

```bash
make configure
make check
```

Ensure AWS credentials are configured (e.g., `AWS_PROFILE`, `AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY`, or default credential chain).

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.5 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.14 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 5.100.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_sqs_queue_policy.queue_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue_policy) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_queue_url"></a> [queue\_url](#input\_queue\_url) | The URL of the SQS Queue to which to attach the policy. | `string` | n/a | yes |
| <a name="input_policy"></a> [policy](#input\_policy) | The JSON policy document for the SQS queue. Must include a Version identifier<br/>(e.g., "2012-10-17" or "2008-10-17") as the top-level key per AWS IAM policy<br/>document requirements. See:<br/>https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_version.html | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_id"></a> [id](#output\_id) | The ID of the resource (same as the queue URL). |
| <a name="output_queue_url"></a> [queue\_url](#output\_queue\_url) | The URL of the SQS queue to which the policy is attached. |
| <a name="output_policy"></a> [policy](#output\_policy) | The policy document attached to the queue. |
<!-- END_TF_DOCS -->
