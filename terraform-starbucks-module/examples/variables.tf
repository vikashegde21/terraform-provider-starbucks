variable "api_key" {
  description = "Starbucks API Key"
  type        = string
  sensitive   = true
}

variable "region" {
  description = "API region"
  type        = string
  default     = "us-west-2"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}
