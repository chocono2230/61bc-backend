resource "aws_api_gateway_rest_api" "this" {
  name = "${var.identifier}-apigateway"
  body = data.template_file.openapi.rendered

  endpoint_configuration {
    types = ["EDGE"]
  }
}

data "template_file" "openapi" {
  template = file("${path.module}/openapi.yaml")

  vars = {
    name                = "${var.identifier}-apigateway"
    auth                = "${var.identifier}-apigateway-auth"
    auth_provider_arn   = aws_cognito_user_pool.this.arn
    integration_uri     = aws_lambda_function.this.invoke_arn
    credential_role_arn = aws_iam_role.api2lambda.arn
  }
}

# IAM Role for API Gateway Execution Lambda
resource "aws_iam_role" "api2lambda" {
  name                = "${var.identifier}-apigateway-lambda-role"
  managed_policy_arns = [aws_iam_policy.api2lambda.arn]

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "apigateway.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "api2lambda" {
  name = "${var.identifier}-apigateway-lambda-policy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "lambda:InvokeFunction"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

# deploy
resource "aws_api_gateway_deployment" "this" {
  depends_on = [
    aws_api_gateway_rest_api.this
  ]

  rest_api_id = aws_api_gateway_rest_api.this.id

  triggers = {
    redeployment = sha256(file("${path.module}/openapi.yaml"))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "this" {
  deployment_id = aws_api_gateway_deployment.this.id
  rest_api_id   = aws_api_gateway_rest_api.this.id
  stage_name    = var.env

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.apigateway_accesslog.arn
    format          = "$context.identity.sourceIp $context.identity.caller $context.identity.user [$context.requestTime] \"$context.httpMethod $context.resourcePath $context.protocol\" $context.status $context.responseLength $context.requestId"
  }
}

resource "aws_api_gateway_method_settings" "all" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  stage_name  = aws_api_gateway_stage.this.stage_name
  method_path = "*/*"

  settings {
    metrics_enabled = true
    logging_level   = "INFO"
  }

  depends_on = [
    aws_api_gateway_account.this
  ]
}

resource "aws_api_gateway_account" "this" {
  cloudwatch_role_arn = aws_iam_role.apigateway_putlog.arn
}

resource "aws_cloudwatch_log_group" "apigateway_accesslog" {
  name = "${var.identifier}-apigateway-log"
}

# Lambda
resource "aws_lambda_permission" "this" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "apigateway.amazonaws.com"

  # More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
  source_arn = "arn:aws:execute-api:${var.region}:${var.accountId}:${aws_api_gateway_rest_api.this.id}/${var.env}/*"
}

# IAM Role for API Gateway Put CloudWatch Logs
resource "aws_iam_role" "apigateway_putlog" {
  name = "${var.identifier}-apigateway-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "apigateway.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
  ]
}
