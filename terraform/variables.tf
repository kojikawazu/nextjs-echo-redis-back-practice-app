# ---------------------------------------------
# Variables
# ---------------------------------------------
variable "project" {
  type = string
}

variable "environment" {
  type = string
}

variable "region" {
  type = string
}

variable "vpc_address" {
  type = string
}

variable "igw_address" {
  type = string
}

variable "public_1a_address" {
  type = string
}

variable "public_1c_address" {
  type = string
}

variable "private_1a_address" {
  type = string
}

variable "private_1c_address" {
  type = string
}

variable "api_port" {
  type = number
}

variable "redis_port" {
  type = number
}

variable "supabase_url" {
  type = string
}

variable "cors_address" {
  type = string
}
