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

logical_product_family  = "lpf"
logical_product_service = "lps"
class_env               = "dev"
instance_env            = 1
instance_resource       = 1

resource_names_map = {
  sqsqueue1 = {
    name       = "sqsqueue1"
    max_length = 80
  }
}

tags = {
  Environment = "test"
  Module      = "sqs-queue-policy"
  Example     = "sqs_managed_sse"
}
