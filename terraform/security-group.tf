data "http" "my_publid_ip" {
  url = "http://ipv4.icanhazip.com"
}

resource "aws_security_group" "alb" {
  name_prefix = join(local.delimiter, [local.name_tag_prefix, "alb", "sg"])
  description = "poc"
  vpc_id      = module.network.vpc_id

  ingress {
    description      = "TLS from Public"
    from_port        = 443
    to_port          = 443
    protocol         = "tcp"
    cidr_blocks      = ["${chomp(data.http.my_publid_ip.body)}/32"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "lambda" {
  name_prefix = join(local.delimiter, [local.name_tag_prefix, "lambda", "sg"])
  description = "poc"
  vpc_id      = module.network.vpc_id

  ingress {
    description      = "TLS from ALB"
    from_port        = 443
    to_port          = 443
    protocol         = "tcp"
    security_groups  = [aws_security_group.alb.id]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
  }
}