output "instance_ip" {
  description = "The IP address of the instance"
  value = aws_instance.main.public_ip
}