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
variable "posts_table_gsi_name_all" {
  type = string
}
variable "posts_table_gsi_name_usr" {
  type = string
}
