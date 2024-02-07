output "eips" {
  description = "EIP Outputs"
  value       = { for eip in aws_eip.this : eip.tags.Name => { "public_dns" : eip.public_dns, "public_ip" : eip.public_ip } }
}

output "volumes" {
  description = "EBS Volume Outputs"
  value       = { for ebs in aws_ebs_volume.this : ebs.tags.Name => { "id" : ebs.id } }
}
