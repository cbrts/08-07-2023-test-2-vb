resource "aws_security_group" "allow_http_from_alb" {
  name        = format("sg1-%s,", local.regional_name)
  description = "Allow http inbound traffic from the associated ALB"
  vpc_id      = var.vpc_id

  ingress {
    description     = "http from the ALB"
    from_port       = 80
    to_port         = 80
    protocol        = "tcp"
    security_groups = [aws_security_group.allow_http_from_internet.id]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
  tags = module.tags.tags
}

resource "aws_security_group" "allow_http_from_internet" {
  name        = format("sg2-%s", local.regional_name)
  description = "Allow http from the internet"
  vpc_id      = var.vpc_id

  ingress {
    description      = "http form the internet"
    from_port        = 80
    to_port          = 80
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
  tags = module.tags.tags
}
