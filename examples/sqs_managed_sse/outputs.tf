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

output "id" {
  description = "The ID of the queue policy (same as the queue URL)."
  value       = module.queue_policy.id
}

output "queue_url" {
  description = "The URL of the SQS queue."
  value       = module.queue_policy.queue_url
}

output "policy" {
  description = "The policy document attached to the queue."
  value       = module.queue_policy.policy
}
