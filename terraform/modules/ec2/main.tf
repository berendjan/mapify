resource "aws_security_group" "this" {
  for_each = var.sg_parameters
  vpc_id   = each.value.vpc_id

  dynamic "ingress" {
    for_each = [for rule in each.value.rules : rule if rule.type == "ingress"]
    content {
      from_port   = ingress.value.from_port
      to_port     = ingress.value.to_port
      protocol    = ingress.value.protocol
      cidr_blocks = ingress.value.cidr_blocks
    }
  }

  dynamic "egress" {
    for_each = [for rule in each.value.rules : rule if rule.type == "egress"]
    content {
      from_port   = egress.value.from_port
      to_port     = egress.value.to_port
      protocol    = egress.value.protocol
      cidr_blocks = egress.value.cidr_blocks
    }
  }
  tags = merge(each.value.tags, {
    Name : each.key
  })
}

resource "aws_ebs_volume" "this" {
  for_each          = var.ebs_volume_parameters
  availability_zone = each.value.availability_zone
  size              = each.value.size
  encrypted         = each.value.encrypted
  tags = merge(each.value.tags, {
    Name : each.key
  })
}

# Post-provisioning step, mount EBS: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-using-volumes.html
resource "aws_volume_attachment" "this" {
  for_each    = var.ebs_volume_attachment_parameters
  device_name = each.value.device_name
  volume_id   = aws_ebs_volume.this[each.value.volume_name].id
  instance_id = aws_instance.this[each.value.instance_name].id
}

resource "aws_eip" "this" {
  for_each = var.eip_parameters
  domain   = "vpc"
  instance = aws_instance.this[each.value.instance_name].id
  tags = merge(each.value.tags, {
    Name : each.key
  })
}

resource "aws_instance" "this" {
  for_each               = var.instance_parameters
  ami                    = each.value.ami
  availability_zone      = each.value.availability_zone
  instance_type          = each.value.instance_type
  subnet_id              = each.value.subnet_id
  vpc_security_group_ids = [for name in each.value.vpc_security_group_names : aws_security_group.this[name].id]
  key_name               = aws_key_pair.this[each.value.key_name].key_name
  tags = merge(each.value.tags, {
    Name : each.key
  })
}

resource "aws_key_pair" "this" {
  for_each   = var.keypair_parameters
  key_name   = each.key
  public_key = each.value.public_key
  tags = merge(each.value.tags, {
    Name : each.key
  })
}
