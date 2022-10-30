variable "project" {
  type = string
}

variable "region" {
  type = string
}

variable "zone" {
  type = string
}

variable "env_vars" {
  type = list(object({
    value = string
    name  = string
  }))
}
