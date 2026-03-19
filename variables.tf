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

variable "queue_url" {
  description = "The URL of the SQS Queue to which to attach the policy."
  type        = string
}

variable "policy" {
  description = <<-EOT
    The JSON policy document for the SQS queue. Must include Version = "2012-10-17"
    as the top-level key. AWS may hang indefinitely without an explicit version;
    "2012-10-17" is required per the Terraform AWS provider resource documentation.
    See: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue_policy
  EOT
  type        = string

  validation {
    condition     = can(jsondecode(var.policy)) && try(jsondecode(var.policy)["Version"], "") == "2012-10-17"
    error_message = "The policy document must contain a top-level Version = \"2012-10-17\" identifier. AWS may hang without it; see the aws_sqs_queue_policy resource documentation."
  }
}
