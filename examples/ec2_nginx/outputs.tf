output "instance_id" {
  value = module.nginx_instance.instance_id
}

output "iam_instance_profile_arn" {
  value = module.nginx_instance.iam_instance_profile_arn
}

output "load_balancer_dns_record" {
  value = module.nginx_instance.load_balancer_dns_record
}

output "alb_sg_group_id" {
  value = module.nginx_instance.alb_sg_group_id
}

output "instance_sg_group_id" {
  value = module.nginx_instance.instance_sg_group_id
}
