module "tags" {
  source = "../tags"

  squad_name      = var.squad_name
  product         = var.product_name
  environment     = var.environment_name
  contact_email   = var.contact_email
  cost_centre     = var.cost_centre
}
