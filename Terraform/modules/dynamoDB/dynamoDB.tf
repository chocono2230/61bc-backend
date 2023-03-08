resource "aws_dynamodb_table" "posts_table" {
  name         = "${var.identifier}-ddb-posts"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "gsiSKey"
    type = "S"
  }

  attribute {
    name = "timestamp"
    type = "N"
  }

  attribute {
    name = "userId"
    type = "S"
  }

  global_secondary_index {
    name            = "${var.identifier}-ddb-posts-gsi-alltime"
    hash_key        = "gsiSKey"
    range_key       = "timestamp"
    projection_type = "ALL"
  }

  global_secondary_index {
    name            = "${var.identifier}-ddb-posts-gsi-usrtime"
    hash_key        = "userId"
    range_key       = "timestamp"
    projection_type = "ALL"
  }

}
