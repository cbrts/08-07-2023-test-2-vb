output "instance_id" {
  value = aws_instance.this.id
}

output "iam_instance_profile_arn" {
  value = aws_iam_instance_profile.this.arn
}

output "load_balancer_dns_record" {
  value = aws_lb.this.dns_name
}

output "alb_sg_group_id" {
  value = aws_security_group.allow_http_from_internet.id
}

output "instance_sg_group_id" {
  value = aws_security_group.allow_http_from_alb.id
}
