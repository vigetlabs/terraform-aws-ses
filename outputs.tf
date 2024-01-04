#=====================================
# SES information
#=====================================
output "email_identity" {
  description = "The email identity."
  value       = try(aws_sesv2_email_identity.this[0].email_identity, "")
}

output "ses_sending_pool_name" {
  description = "The name of the SES sending pool to associate the domain with."
  value       = try(aws_sesv2_dedicated_ip_pool.this[0].pool_name, aws_sesv2_configuration_set.this[0].delivery_options[0].sending_pool_name, "")
}

output "iam_sending_group_name" {
  description = "The IAM group name."
  value       = try(aws_iam_group.ses_users[0].name, "")
}

#=====================================
# DNS Record Data
#=====================================
output "ses_dkim_records" {
  description = "The DNS records required for Amazon SES validation and DKIM setup."
  value = try(flatten([
    for attribute in aws_sesv2_email_identity.this[0].dkim_signing_attributes : [
      for token in attribute.tokens : {
        name    = "${token}._domainkey.${var.domain}"
        type    = "CNAME"
        ttl     = "600"
        records = ["${token}.dkim.amazonses.com"]
      }
    ]
  ]), [])
}
