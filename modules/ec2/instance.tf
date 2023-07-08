locals {
  formatted_instance_name = format("ec2-%s-%s-%s-%s-%s", data.aws_region.current.name, data.aws_subnet.this.availability_zone, var.environment_name, var.product_name, random_id.this.id)
}

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

resource "random_id" "this" {
  byte_length = 4
}

resource "aws_instance" "this" {
  ami           = data.aws_ami.this.id
  instance_type = var.instance_type
  subnet_id     = var.private_subnet_id
  tags = merge(
    {
      "Name" = local.formatted_instance_name
    },
      module.tags.tags
  )
}
