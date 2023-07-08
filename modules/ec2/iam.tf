data "aws_iam_policy_document" "this" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

data "aws_iam_policy" "ssm" {
  name = "AmazonSSMManagedInstanceCore"
}

resource "aws_iam_policy" "s3" {
  name = format("iam-policy-%s-%s-%s", var.environment_name, var.product_name, random_id.this.id)
  policy = jsonencode({
    Version : "2012-10-17"
    Statement : [
      {
        Action : [
          "s3:List*"
        ]
        Effect : "Allow"
        Resource : "arn:aws:s3:::test-bucket"
      },
    ]
  })
  tags = module.tags.tags
}

resource "aws_iam_role" "this" {
  name                = format("iam-role-%s-%s-%s", var.environment_name, var.product_name, random_id.this.id)
  assume_role_policy  = data.aws_iam_policy_document.this.json
  managed_policy_arns = [data.aws_iam_policy.ssm.arn, aws_iam_policy.s3.arn]
  tags                = module.tags.tags
}

resource "aws_iam_instance_profile" "this" {
  name = format("iam-profile-%s-%s-%s", var.environment_name, var.product_name, random_id.this.id)
  role = aws_iam_role.this.name
  tags = module.tags.tags
}
