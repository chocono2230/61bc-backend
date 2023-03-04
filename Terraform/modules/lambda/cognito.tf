resource "aws_cognito_user_pool" "this" {
  name = "${var.identifier}-userpool"
  alias_attributes = [
    "preferred_username",
  ]
  auto_verified_attributes = [
    "email",
  ]
  deletion_protection = "ACTIVE"
  mfa_configuration   = "OFF"

  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  admin_create_user_config {
    allow_admin_create_user_only = true
  }

  email_configuration {
    email_sending_account = "COGNITO_DEFAULT"
  }

  password_policy {
    minimum_length                   = 6
    require_lowercase                = true
    require_numbers                  = true
    require_symbols                  = true
    require_uppercase                = true
    temporary_password_validity_days = 7
  }

  schema {
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true
    name                     = "email"
    required                 = true

    string_attribute_constraints {
      max_length = "2048"
      min_length = "0"
    }
  }

  user_attribute_update_settings {
    attributes_require_verification_before_update = [
      "email",
    ]
  }

  username_configuration {
    case_sensitive = true
  }

  verification_message_template {
    default_email_option = "CONFIRM_WITH_CODE"
  }
}

resource "aws_cognito_user_pool_client" "this" {
  name                    = "${var.identifier}-userpool-client"
  user_pool_id            = aws_cognito_user_pool.this.id
  access_token_validity   = 60
  auth_session_validity   = 3
  enable_token_revocation = true
  explicit_auth_flows = [
    "ALLOW_CUSTOM_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_SRP_AUTH",
  ]
  id_token_validity             = 60
  prevent_user_existence_errors = "ENABLED"
  read_attributes = [
    "email",
    "email_verified",
    "name",
    "nickname",
    "picture",
    "preferred_username",
    "profile",
    "updated_at",
  ]
  refresh_token_validity = 30
  write_attributes = [
    "email",
    "name",
    "nickname",
    "picture",
    "preferred_username",
    "profile",
    "updated_at",
  ]

  token_validity_units {
    access_token  = "minutes"
    id_token      = "minutes"
    refresh_token = "days"
  }
}
