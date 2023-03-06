resource "aws_dynamodb_table" "basic-dynamodb-table" {
  name         = "${var.identifier}-ddb-post"
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
    name            = "${var.identifier}-ddb-post-gsi-alltime"
    hash_key        = "gsiSKey"
    range_key       = "timestamp"
    projection_type = "ALL"
  }

  global_secondary_index {
    name            = "${var.identifier}-ddb-post-gsi-usrtime"
    hash_key        = "userId"
    range_key       = "timestamp"
    projection_type = "ALL"
  }

}
