locals {
  regional_name = format("%s-%s-%s-%s-%s", data.aws_region.current.name, data.aws_subnet.this.availability_zone, var.environment_name, var.product_name, random_string.this.id)
  global_name   = format("%s-%s-%s", var.environment_name, var.product_name, random_string.this.id)
}
