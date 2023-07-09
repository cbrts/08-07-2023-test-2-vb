# The ALB and target are using local.global_name due to naming length constraints

resource "aws_lb" "this" {
  name               = format("alb-%s", local.global_name)
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.allow_http_from_internet.id]
  subnets            = var.public_subnet_ids
  tags = module.tags.tags
}

resource "aws_lb_listener" "this" {
  load_balancer_arn = aws_lb.this.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
}

resource "aws_lb_target_group" "this" {
  name     = format("tg-%s", local.global_name)
  port     = 80
  protocol = "HTTP"
  vpc_id   = var.vpc_id
  tags = module.tags.tags
}

resource "aws_lb_target_group_attachment" "this" {
  target_group_arn = aws_lb_target_group.this.arn
  target_id        = aws_instance.this.id
  port             = 80
}
