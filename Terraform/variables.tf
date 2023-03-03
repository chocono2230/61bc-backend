variable "env" {
  type = string
}
variable "project" {
  type = string
}
variable "region" {
  type = string
}
variable "accountId" {
  type      = string
  sensitive = true
}
