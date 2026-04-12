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

resource "aws_secretsmanager_secret" "gotrue_service_key" {
  name = "rabi-food/gotrue-service-key"
}

resource "aws_secretsmanager_secret_version" "gotrue_service_key" {
  secret_id     = aws_secretsmanager_secret.gotrue_service_key.id
  secret_string = var.gotrue_service_key
}

resource "aws_secretsmanager_secret" "gotrue_db_url" {
  name = "rabi-food/gotrue-db-url"
}

resource "aws_secretsmanager_secret_version" "gotrue_db_url" {
  secret_id = aws_secretsmanager_secret.gotrue_db_url.id
  secret_string = "postgres://${var.db_user}:${var.db_password}@${aws_db_instance.postgres.address}:5432/${var.db_name}?search_path=auth&sslmode=require"
}
