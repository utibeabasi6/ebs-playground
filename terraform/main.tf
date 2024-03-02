locals {
  device_letter_list = ["f", "g", "h"]
}

data "aws_vpc" "main" {
    default = true
    }

resource "aws_key_pair" "main" {
  key_name   = var.instance_name
  public_key = var.instance_public_key
}

resource "aws_instance" "main" {
  ami           = var.instance_ami
  instance_type = var.instance_type
  availability_zone = var.ebs_availability_zone
  key_name = aws_key_pair.main.key_name
  security_groups = [ aws_security_group.main.name ]
  user_data_base64 = filebase64("../scripts/bootstrap.sh")

  tags = {
    Name = var.instance_name
  }
}

resource "aws_ebs_volume" "main" {
  count = var.ebs_volume_count
  availability_zone = var.ebs_availability_zone
  size              = var.ebs_volume_size

  tags = {
    Name = "${var.instance_name}-volume-${count.index}"
  }
}

resource "aws_volume_attachment" "ebs_attatchment" {
  count = var.ebs_volume_count
  device_name = "/dev/sd${local.device_letter_list[count.index]}"
  volume_id   = aws_ebs_volume.main[count.index].id
  instance_id = aws_instance.main.id
}


resource "aws_security_group" "main" {
  name        = "main"
  description = "Allow traffic"
  vpc_id      = data.aws_vpc.main.id

  tags = {
    Name = "main"
  }
}

resource "aws_vpc_security_group_ingress_rule" "allow_ssh" {
  security_group_id = aws_security_group.main.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 22
  ip_protocol       = "tcp"
  to_port           = 22
}

resource "aws_vpc_security_group_ingress_rule" "allow_port_3000" {
  security_group_id = aws_security_group.main.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 3000
  ip_protocol       = "tcp"
  to_port           = 3000
}

resource "aws_vpc_security_group_egress_rule" "allow_all" {
  security_group_id = aws_security_group.main.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = -1
}