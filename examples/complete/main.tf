// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

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
  name                       = "sqs-queue-policy-example-${random_string.suffix.result}"
  kms_master_key_id          = aws_kms_key.queue.arn
  sqs_managed_sse_enabled    = false
  message_retention_seconds  = var.message_retention_seconds
  visibility_timeout_seconds = var.visibility_timeout_seconds

  tags = var.tags
}

module "queue_policy" {
  source = "../.."

  queue_url = aws_sqs_queue.queue.url
  policy = jsonencode({
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
