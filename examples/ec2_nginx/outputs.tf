output "hello_world" {
  value = "Hello, World!"
}

output "instance_id" {
  value = module.nginx_instance.instance_id
}
