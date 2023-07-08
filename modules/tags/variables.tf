variable "squad_name" {
  description = "Squad Name in pascal case"
}

variable "cost_centre" {
  description = "Cost Centre for the application"
}

variable "contact_email" {
  description = "Email of group responsible for infrastructure"
}

variable "product" {
  description = "Product name in pascal case"
}

variable "environment" {
  description = "Environemnt name in pascal case"
}

variable "additional_tags" {
  type        = map(string)
  description = "Additional tags to apply"
  default     = {}
}
