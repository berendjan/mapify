provider "aws" {
  region = var.region
}

# Uncomment after initial provisioning of resources
terraform {
  backend "s3" {
    bucket         = "haic-studios-tfstate"
    key            = "global/s3/terraform.tfstate"
    region         = "eu-west-2"
    dynamodb_table = "haic-studios-tfstate-lock"
    encrypt        = true
  }
}

# retrieve instance type
data "aws_ami" "debian" {
  most_recent = true
  name_regex  = "^debian-12-arm64-[0-9]{8}-[0-9]{4}$"

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "architecture"
    values = ["arm64"]
  }
}

module "backend" {
  source = "./modules/backend"
  name   = "haic-studios"
}

module "vpc" {
  source = "./modules/vpc"
  vpc_parameters = {
    vpc1 = {
      cidr_block = "10.0.0.0/16"
    }
  }
  subnet_parameters = {
    subnet1 = {
      cidr_block        = "10.0.1.0/24"
      vpc_name          = "vpc1"
      availability_zone = "${var.region}a"
    }
  }
  igw_parameters = {
    igw1 = {
      vpc_name = "vpc1"
    }
  }
  rt_parameters = {
    rt1 = {
      vpc_name = "vpc1"
      routes = [
        {
          cidr_block = "0.0.0.0/0"
          gateway_id = "igw1"
        }
      ]
    }
  }
  rt_association_parameters = {
    rta1 = {
      subnet_name = "subnet1"
      rt_name     = "rt1"
    }
  }
}

module "ec2" {
  source = "./modules/ec2"
  sg_parameters = {
    sg1 = {
      vpc_id = module.vpc.vpcs["vpc1"]["id"]
      rules = [
        { from_port = 22, to_port = 22, protocol = "tcp", cidr_blocks = ["0.0.0.0/0"], type = "ingress" },
        { from_port = 80, to_port = 80, protocol = "tcp", cidr_blocks = ["0.0.0.0/0"], type = "ingress" },
        { from_port = 443, to_port = 443, protocol = "tcp", cidr_blocks = ["0.0.0.0/0"], type = "ingress" },
        { from_port = 0, to_port = 0, protocol = "all", cidr_blocks = ["0.0.0.0/0"], type = "egress" },
      ]
    }
  }
  ebs_volume_parameters = {
    ebs1 = {
      availability_zone = "${var.region}a"
      size              = var.volume_size
    }
  }
  ebs_volume_attachment_parameters = {
    ebsa1 = {
      device_name   = "/dev/sdf"
      volume_name   = "ebs1"
      instance_name = "instance1"
    }
  }
  eip_parameters = {
    eip1 = {
      instance_name = "instance1"
    }
  }
  instance_parameters = {
    instance1 = {
      ami                      = "ami-0d2ef3b94aac558f9"
      availability_zone        = "${var.region}a"
      instance_type            = var.instance_type
      subnet_id                = module.vpc.subnets["subnet1"]["id"]
      vpc_security_group_names = ["sg1"]
      key_name                 = "keypair1"
    }
  }
  keypair_parameters = {
    keypair1 = {
      public_key = file("~/.ssh/id_ed25519.pub")
    }
  }
}

