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
    The JSON policy document for the SQS queue. Must include a Version identifier
    (e.g., "2012-10-17" or "2008-10-17") as the top-level key per AWS IAM policy
    document requirements. See:
    https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_version.html
  EOT
  type        = string

  validation {
    condition = can(jsondecode(var.policy)) && contains(
      ["2012-10-17", "2008-10-17"],
      try(jsondecode(var.policy)["Version"], "")
    )
    error_message = "The policy document must contain a top-level Version identifier. Valid values are \"2012-10-17\" or \"2008-10-17\" as per AWS IAM policy document requirements."
  }
}
