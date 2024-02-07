variable "sg_parameters" {
  description = "SG parameters"
  type = map(object({
    vpc_id = string
    rules = list(object({
      from_port   = number
      to_port     = number
      protocol    = string
      cidr_blocks = list(string)
      type        = string // "ingress" or "egress"
    }))
    tags = optional(map(string), {})
  }))
  default = {}
}

variable "ebs_volume_parameters" {
  description = "EBS Volume Parameters"
  type = map(object({
    availability_zone = string
    size              = number
    encrypted         = optional(bool, false)
    tags              = optional(map(string), {})
  }))
  default = {}
}

variable "ebs_volume_attachment_parameters" {
  description = "EBS Volume Attachment Parameters"
  type = map(object({
    device_name   = string // /dev/sd[f-p]
    volume_name   = string
    instance_name = string
  }))
  default = {}
}

variable "eip_parameters" {
  description = "EIP Parameters"
  type = map(object({
    instance_name = string
    tags          = optional(map(string), {})
  }))
  default = {}
}

variable "instance_parameters" {
  description = "Instance Parameters"
  type = map(object({
    ami                      = string
    availability_zone        = string
    instance_type            = string
    subnet_id                = string
    vpc_security_group_names = list(string)
    key_name                 = string
    tags                     = optional(map(string), {})
  }))
  default = {}
}

variable "keypair_parameters" {
  description = "SSH Public Key Parameters"
  type = map(object({
    public_key = string
    tags       = optional(map(string), {})
  }))
  default = {}
}
