module "nginx_instance" {
  source = "../../modules/ec2"

  instance_type     = "t3.micro"
  private_subnet_id = var.private_subnet_id
  public_subnet_ids = var.public_subnet_ids
  vpc_id            = var.vpc_id
  environment_name  = "Dev"
  squad_name        = "SquadA"
  contact_email     = "SquadA@acme.com"
  product_name      = "Nginx"
  cost_centre       = "32413142"
}
