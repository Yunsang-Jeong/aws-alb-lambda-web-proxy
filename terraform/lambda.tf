################################################################################
# lambda

resource "aws_lambda_function" "this" {
  # General
  function_name = join(local.delimiter, [local.name_tag_prefix, "lambda"])

  # Runtime
  role    = aws_iam_role.lambda_service_role.arn
  runtime = "go1.x"
  handler = "main"
  timeout = 10

  # Source code
  filename         = "../bin/main.zip"
  package_type     = "Zip"
  publish          = true
  source_code_hash = data.archive_file.zip.output_base64sha256

  vpc_config {
    subnet_ids         = [ 
      for identifier, subnet_id in module.network.subnet_ids: 
      subnet_id 
      if can(regex("^private\\-", identifier))
    ]
    security_group_ids = [aws_security_group.lambda.id]
  }

}
################################################################################


################################################################################
# Service role for the lambda
resource "aws_iam_role" "lambda_service_role" {
  name               = join(local.delimiter, [local.name_tag_prefix, "lambda", "service", "role"])
  assume_role_policy = data.aws_iam_policy_document.lambda_service_role.json
}

resource "aws_iam_role_policy" "lambda_cloudwatch" {
  name   = join(local.delimiter, [local.name_tag_prefix, "lambda", "cloudwatch", "role"])
  role   = aws_iam_role.lambda_service_role.id
  policy = data.aws_iam_policy_document.cloudwatch_policy.json
}

resource "aws_iam_role_policy_attachment" "AWSLambdaVPCAccessExecutionRole" {
  role       = aws_iam_role.lambda_service_role.id
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
}
################################################################################


################################################################################
# Invoke

resource "aws_lambda_permission" "this" {
  function_name = aws_lambda_function.this.arn
  
  statement_id  = "AllowExecutionFromalb"
  principal     = "elasticloadbalancing.amazonaws.com"
  action        = "lambda:InvokeFunction"
  source_arn    = aws_lb_target_group.this.arn
}
################################################################################


################################################################################
# Build 

resource "null_resource" "makefile" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command = "make -f ../Makefile clean"
  }

  provisioner "local-exec" {
    command = "make -f ../Makefile build"
  }
}

data "archive_file" "zip" {
  type        = "zip"
  source_dir  = "../bin/"
  output_path = "../bin/main.zip"

  depends_on = [
    null_resource.makefile
  ]
}
################################################################################