resource "aws_lambda_function" "this" {
  function_name    = "${var.identifier}-lambda"
  s3_bucket        = aws_s3_bucket.this.bucket
  s3_key           = data.aws_s3_object.golang_zip.key
  role             = aws_iam_role.lambda.arn
  handler          = "lambda"
  source_code_hash = data.aws_s3_object.golang_zip_hash.body
  runtime          = "go1.x"
  timeout          = "10"
}

resource "aws_iam_role" "lambda" {
  name = "${var.identifier}-lambda-role"

  assume_role_policy = <<-EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "lambda.amazonaws.com"
        },
        "Effect": "Allow",
        "Sid": ""
      }
    ]
  }
  EOF
  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
  ]
}

resource "null_resource" "this" {
  depends_on = [aws_s3_bucket.this]

  triggers = {
    code_diff = sha256(join("", [
      for file in fileset(local.golang_codedir_local_path, "./**/*.go")
      : filebase64("${local.golang_codedir_local_path}/${file}")
    ]))
  }

  provisioner "local-exec" {
    command     = "GOARCH=amd64 GOOS=linux go build -o ../bin/lambda"
    working_dir = local.golang_codedir_local_path
  }
  provisioner "local-exec" {
    command = "zip -j ${local.golang_zip_local_path} ${local.golang_binary_local_path}"
  }
  provisioner "local-exec" {
    command = "aws s3 cp ${local.golang_zip_local_path} s3://${aws_s3_bucket.this.bucket}/${local.golang_zip_s3_key}"
  }
  provisioner "local-exec" {
    command = "openssl dgst -sha256 -binary ${local.golang_zip_local_path} | openssl enc -base64 | tr -d \"\n\" > ${local.golang_zip_base64sha256_local_path}"
  }
  provisioner "local-exec" {
    command = "aws s3 cp ${local.golang_zip_base64sha256_local_path} s3://${aws_s3_bucket.this.bucket}/${local.golang_zip_base64sha256_s3_key} --content-type \"text/plain\""
  }
}

resource "aws_s3_bucket" "this" {
  bucket        = "${var.identifier}-lambda-s3"
  force_destroy = var.env == "tst" ? true : false
}
resource "aws_s3_bucket_acl" "this" {
  bucket = aws_s3_bucket.this.bucket
  acl    = "private"
}
resource "aws_s3_bucket_server_side_encryption_configuration" "this" {
  bucket = aws_s3_bucket.this.bucket

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}
resource "aws_s3_bucket_public_access_block" "this" {
  bucket = aws_s3_bucket.this.bucket

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

data "aws_s3_object" "golang_zip" {
  depends_on = [null_resource.this]

  bucket = aws_s3_bucket.this.bucket
  key    = local.golang_zip_s3_key
}
data "aws_s3_object" "golang_zip_hash" {
  depends_on = [null_resource.this]

  bucket = aws_s3_bucket.this.bucket
  key    = local.golang_zip_base64sha256_s3_key
}
