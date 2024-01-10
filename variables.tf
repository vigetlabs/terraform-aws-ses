variable "domain" {
  description = "The domain for the SES identity."
  type        = string
}

variable "sending_pool_name" {
  type        = string
  description = <<EOT
    Override the default sending pool name. If not provided, the sending pool name will use the context module id.
    Note: If you are using an existing sending pool, create_sending_pool must be set to false.
  EOT
  default     = ""
}

variable "create_sending_pool" {
  type        = bool
  description = "Whether to create a sending pool for the domain."
  default     = false
}

variable "sending_pool_scaling_mode" {
  type        = string
  description = "(Optional) IP pool scaling mode. Valid values: STANDARD, MANAGED."
  default     = "MANAGED"
}

variable "group_path" {
  type        = string
  description = "The IAM Path of the group and policy to create"
  default     = "/"
}

variable "allowed_sending_addresses" {
  type        = list(string)
  description = <<EOT
    A list of email addresses that are allowed to send email from the domain.
    If this list is empty, the domain will be configured to allow any email address to send email from the domain.
  EOT
  default     = []
}
