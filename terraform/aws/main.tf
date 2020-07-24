provider "aws" {
    region = "us-east-1"
}

resource "aws_instance" "approzium" {
    ami = var.ami
    instance_type = var.instance_type
    user_data = templatefile(
        "${path.module}/scripts/install.sh.tpl",
        {
        download-url  = var.download-url
        config        = file("approzium.config.yml")
        extra-install = var.extra-install
        logs-file     = var.logs-file
        }
    )

    tags = {
        Name = var.instance_name
    }

    security_groups = [aws_security_group.approzium.name]
    key_name = var.key-name
}


// Security group for Approzium allows HTTP, and GRPC access
resource "aws_security_group" "approzium" {
    name = "approzium"
    description = "Approzium servers"
}

resource "aws_security_group_rule" "approzium-http-api" {
    security_group_id = aws_security_group.approzium.id
    type = "ingress"
    from_port = 6000
    to_port = 6000
    protocol = "tcp"
}

resource "aws_security_group_rule" "approzium-grpc" {
    security_group_id = aws_security_group.approzium.id
    type = "ingress"
    from_port = 6001
    to_port = 6001
    protocol = "tcp"
}

resource "aws_security_group_rule" "approzium-egress" {
    security_group_id = aws_security_group.approzium.id
    type = "egress"
    from_port = 0
    to_port = 0
    protocol = "-1"
}
