output "vpcs" {
  description = "VPC Outputs"
  value       = { for vpc in aws_vpc.this : vpc.tags.Name => { "cidr_block" : vpc.cidr_block, "id" : vpc.id } }
}

output "subnets" {
  description = "Subnet Outputs"
  value       = { for subnet in aws_subnet.this : subnet.tags.Name => { "id" : subnet.id } }
}
