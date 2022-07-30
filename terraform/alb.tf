resource "aws_lb" "this" {
  name               = join(local.delimiter, [local.name_tag_prefix, "alb"])
  internal           = false
  load_balancer_type = "application"

  security_groups    = [
    aws_security_group.alb.id
  ]

  subnets            = [ 
    for identifier, subnet_id in module.network.subnet_ids: 
    subnet_id 
    if can(regex("^public\\-", identifier))
  ]
}

resource "aws_lb_target_group" "this" {
  name        = join(local.delimiter, [local.name_tag_prefix, "alb", "tg"])
  target_type = "lambda"
}

resource "aws_lb_target_group_attachment" "this" {
  target_group_arn = aws_lb_target_group.this.arn
  target_id        = aws_lambda_function.this.arn
  depends_on       = [aws_lambda_permission.this]
}

resource "aws_lb_listener" "this" {
  load_balancer_arn = aws_lb.this.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this.arn
  }
}
