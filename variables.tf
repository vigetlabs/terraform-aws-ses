variable "domain" {
  description = "The domain for the SES identity."
  type        = string
}

variable "ses_sending_pool_name" {
  type        = string
  description = "The name of the SES sending pool to associate the domain with."
  default     = ""
}

variable "create_sending_pool" {
  type        = bool
  description = "Whether to create a sending pool for the domain."
  default     = false
}

variable "ses_group_path" {
  type        = string
  description = "The IAM Path of the group and policy to create"
  default     = "/"
}
