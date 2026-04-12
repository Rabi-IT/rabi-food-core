variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "app_port" {
  description = "Port the app listens on"
  type        = number
  default     = 3000
}

variable "ec2_instance_type" {
  description = "EC2 instance type for ECS cluster nodes"
  type        = string
  default     = "t3.small"
}

variable "ecs_cpu" {
  description = "ECS task CPU units (256, 512, 1024...)"
  type        = number
  default     = 512
}

variable "ecs_memory" {
  description = "ECS task memory in MB"
  type        = number
  default     = 1024
}

variable "db_name" {
  description = "Database name"
  type        = string
  default     = "rabi_food"
}

variable "db_user" {
  description = "Database username"
  type        = string
  default     = "postgres"
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "auth_secret" {
  description = "JWT auth secret"
  type        = string
  sensitive   = true
}

variable "grafana_token" {
  description = "Grafana Cloud API token (MetricsPublisher + LogsPublisher scopes)"
  type        = string
  sensitive   = true
}

variable "grafana_prometheus_url" {
  description = "Grafana Cloud Prometheus remote_write URL"
  type        = string
}

variable "grafana_prometheus_instance_id" {
  description = "Grafana Cloud Prometheus numeric instance ID (used as basic_auth username)"
  type        = string
}

variable "domain_name" {
  description = "Domínio da API (ex: api.rabi.food)"
  type        = string
}
