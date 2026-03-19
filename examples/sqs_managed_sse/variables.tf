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

variable "logical_product_family" {
  description = "Logical product family for resource naming."
  type        = string
}

variable "logical_product_service" {
  description = "Logical product service for resource naming."
  type        = string
}

variable "class_env" {
  description = "Class environment for resource naming (e.g., dev, prod)."
  type        = string
}

variable "instance_env" {
  description = "Instance environment number for resource naming (0-999)."
  type        = number
}

variable "instance_resource" {
  description = "Instance resource number for resource naming (0-100)."
  type        = number
}

variable "resource_names_map" {
  description = "Map of resource types to naming config for the resource naming module."
  type = map(object({
    name       = string
    max_length = number
  }))
}

variable "queue_url" {
  description = "The URL of the SQS queue to attach the policy to. If null, a queue is created in the example."
  type        = string
  default     = null
}

variable "policy" {
  description = "The JSON policy document. If null, a policy is built from the created queue."
  type        = string
  default     = null
}

variable "message_retention_seconds" {
  description = "The number of seconds Amazon SQS retains a message."
  type        = number
  default     = 345600 # 4 days
}

variable "visibility_timeout_seconds" {
  description = "The visibility timeout for the queue."
  type        = number
  default     = 30
}

variable "tags" {
  description = "Map of tags to assign to the resources."
  type        = map(string)
  default     = {}
}
