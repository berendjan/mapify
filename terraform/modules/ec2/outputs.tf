output "eips" {
  description = "EIP Outputs"
  value       = { for eip in aws_eip.this : eip.tags.Name => { "public_dns" : eip.public_dns, "public_ip" : eip.public_ip } }
}

output "volumes" {
  description = "EBS Volume Outputs"
  value       = { for ebs in aws_ebs_volume.this : ebs.tags.Name => { "id" : ebs.id } }
}

output "instances" {
  description = "EC2 Instance Outputs"
  value       = { for ec2 in aws_aws_instance.this : ec2.tags.Name => { "public_dns" : ec2.public_dns, "public_ip" : ec2.public_ip, "private_dns" : ec2.private_dns, "private_ip" : ec2.private_ip } }
}
