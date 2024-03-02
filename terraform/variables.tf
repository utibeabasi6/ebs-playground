variable "instance_name" {
    type = string
    description = "The name of the AWS instance"
    default = "ebs-playground"
}

variable "instance_ami" {
    type = string
    description = "The AMI id for the AWS instance"
    default = "ami-07d9b9ddc6cd8dd30"
}

variable "instance_public_key" {
    type = string
    description = "The public key for SSH"
}

variable "instance_type" {
    type = string
    description = "The AWS instance type"
    default = "t2.micro"
}

variable "ebs_volume_count" {
    type = number
    description = "Number of ebs volumes to create"
    default = 3
}

variable "ebs_volume_size" {
    type = number
    description = "Size (in GB) of each EBS volumes"
    default = 10
}

variable "ebs_availability_zone" {
    type = string
    description = "Availability zone to deploy the EBS volumes in"
    default = "us-east-1a"
}