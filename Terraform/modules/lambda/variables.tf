variable "identifier" {
  type = string
}
variable "env" {
  type = string
}
variable "region" {
  type = string
}
variable "accountId" {
  type      = string
  sensitive = true
}
variable "posts_table_name" {
  type = string
}
