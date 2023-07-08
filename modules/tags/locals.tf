locals {
  mandatory_tags = {
    SquadName    = title(var.squad_name)
    CostCentre   = var.cost_centre
    ContactEmail = var.contact_email
    Product      = title(var.product)
    Environment  = title(var.environment)
  }
  tags = merge(
    local.mandatory_tags,
    var.additional_tags,
  )
}
