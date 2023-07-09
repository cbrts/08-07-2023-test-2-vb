data "aws_ami" "this" {
  most_recent = true
  owners      = ["amazon"]
  filter {
    name   = "architecture"
    values = ["x86_64"]
  }
  filter {
    name   = "name"
    values = ["al2023-ami-2023*"]
  }
}

data "aws_region" "current" {}

data "aws_subnet" "this" {
  id = var.private_subnet_id
}

resource "random_string" "this" {
  length  = 6
  special = false
}

resource "aws_instance" "this" {
  ami                  = data.aws_ami.this.id
  instance_type        = var.instance_type
  subnet_id            = var.private_subnet_id
  iam_instance_profile = aws_iam_instance_profile.this.id
  security_groups      = [aws_security_group.allow_http_from_alb.id]
  tags = merge(
    {
      "Name" = format("ec2-%s", local.regional_name)
    },
    module.tags.tags
  )
  user_data = file("${path.module}/files/user_data.sh")
}
