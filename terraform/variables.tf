variable "instance_type" {
  type        = string
  description = "AWS EC2 instance type"
  default     = "t4g.small"
}

variable "volume_size" {
  type        = number
  description = "AWS EBS volume size"
  default     = 8
}

variable "region" {
  type        = string
  description = "AWS region"
  default     = "eu-west-2"
}

# variable "public_key" {
#   type        = string
#   description = "Public key"
# }

