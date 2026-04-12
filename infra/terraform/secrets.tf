resource "aws_secretsmanager_secret" "db_password" {
  name = "rabi-food/db-password"
}

resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id     = aws_secretsmanager_secret.db_password.id
  secret_string = var.db_password
}

resource "aws_secretsmanager_secret" "auth_secret" {
  name = "rabi-food/auth-secret"
}

resource "aws_secretsmanager_secret_version" "auth_secret" {
  secret_id     = aws_secretsmanager_secret.auth_secret.id
  secret_string = var.auth_secret
}

resource "aws_secretsmanager_secret" "grafana_token" {
  name = "rabi-food/grafana-token"
}

resource "aws_secretsmanager_secret_version" "grafana_token" {
  secret_id     = aws_secretsmanager_secret.grafana_token.id
  secret_string = var.grafana_token
}
