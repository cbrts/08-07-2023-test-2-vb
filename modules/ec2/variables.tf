## EC2
variable "instance_type" {
  description = "Type of instance to deploy"
  type        = string
  default     = "t3.micro"
}

variable "private_subnet_id" {
  description = "Private subnet id to deploy to"
  type        = string
}

## Tags
variable "environment_name" {
  description = "Environemnt name in pascal case"
  type        = string
  default     = "Dev"
}

variable "squad_name" {
  description = "Squad Name in pascal case"
  type        = string
  default     = "SquadA"
}

variable "contact_email" {
  description = "Email of group responsible for infrastructure"
  type        = string
  default    = "SquadA@acme.com"
}

variable "product_name" {
  description = "Product name in pascal case"
  type        = string
  default     = "Nginx"
}

variable "cost_centre" {
  description = "Cost Centre for the application"
  type        = string
  default     = "32413132"
}
