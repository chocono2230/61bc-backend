locals {
  golang_codedir_local_path          = "${path.module}/src"
  golang_binary_local_path           = "${path.module}/bin/lambda"
  golang_zip_local_path              = "${path.module}/archive/lambda.zip"
  golang_zip_base64sha256_local_path = "${local.golang_zip_local_path}.base64sha256"
  golang_zip_s3_key                  = "archive/lambda.zip"
  golang_zip_base64sha256_s3_key     = "${local.golang_zip_s3_key}.base64sha256.txt"
}
