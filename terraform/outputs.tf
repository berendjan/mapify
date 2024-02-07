
output "aws_ami_name" {
  description = "Amazon Machine Image (AMI) Name"
  value       = data.aws_ami.debian.name
}


output "aws_ami_id" {
  description = "Amazon Machine Image (AMI) Id"
  value       = data.aws_ami.debian.id
}

output "vpcs" {
  description = "VPC Outputs"
  value       = module.vpc.vpcs
}

output "subnets" {
  description = "Subnet Outputs"
  value       = module.vpc.subnets
}

# output "eips" {
#   description = "EIP Outputs"
#   value       = module.ec2.eips
# }

# output "volumes" {
#   description = "EBS Volume Outputs"
#   value       = module.ec2.volumes
# }

# output "instances" {
#   description = "EC2 Instance Outputs"
#   value       = module.ec2.instances
# }
