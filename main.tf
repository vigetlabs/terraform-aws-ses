#=====================================
# SES Domain Identity
#=====================================
resource "aws_sesv2_email_identity" "this" {
  count = module.this.enabled ? 1 : 0

  email_identity         = var.domain
  configuration_set_name = aws_sesv2_configuration_set.this[0].configuration_set_name
  tags                   = module.this.tags
}

resource "aws_sesv2_configuration_set" "this" {
  count = module.this.enabled ? 1 : 0

  configuration_set_name = module.this.id
  tags                   = module.this.tags

  dynamic "delivery_options" {
    for_each = var.create_sending_pool ? ["_enable"] : []
    content {
      sending_pool_name = aws_sesv2_dedicated_ip_pool.this[0].pool_name
    }
  }

  dynamic "delivery_options" {
    for_each = !var.create_sending_pool && var.sending_pool_name != "" ? ["_enable"] : []
    content {
      sending_pool_name = var.sending_pool_name
    }
  }
}

resource "aws_sesv2_dedicated_ip_pool" "this" {
  count = module.this.enabled && var.create_sending_pool ? 1 : 0

  pool_name = var.sending_pool_name != "" ? var.sending_pool_name : module.this.id
  tags      = module.this.tags
}

#=====================================
# IAM Group for SES Domain Identity
#=====================================
resource "aws_iam_group" "ses_users" {
  count = module.this.enabled ? 1 : 0

  name = module.this.id
  path = var.group_path
}

## IAM Group Policies for SES Domain Identity
data "aws_iam_policy_document" "ses_group_sending_policy" {
  statement {
    effect = "Allow"

    resources = [
      aws_sesv2_email_identity.this[0].arn,
      aws_sesv2_configuration_set.this[0].arn
    ]

    actions = [
      "ses:SendRawEmail",
      "ses:SendEmail"
    ]

    condition {
      test     = "StringLike"
      variable = "ses:FromAddress"
      values = [
        "*@${var.domain}"
      ]
    }

    condition {
      test     = "StringLike"
      variable = "ses:FeedbackAddress"
      values = [
        "*@${var.domain}"
      ]
    }
  }
}

resource "aws_iam_policy" "ses_sending_policy" {
  count = module.this.enabled ? 1 : 0

  name = module.this.id
  path = var.group_path

  policy = data.aws_iam_policy_document.ses_group_sending_policy.json

  tags = module.this.tags
}

resource "aws_iam_group_policy_attachment" "ses_sending_policy" {
  count = module.this.enabled ? 1 : 0

  group      = aws_iam_group.ses_users[0].name
  policy_arn = aws_iam_policy.ses_sending_policy[0].arn
}
