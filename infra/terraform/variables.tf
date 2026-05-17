variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "api_port" {
  description = "Port the API listens on"
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

variable "grafana_stack_slug" {
  description = "Grafana Cloud stack slug — subdomain of yourstack.grafana.net"
  type        = string
}

variable "grafana_cloud_api_key" {
  description = "Grafana Cloud organization API key (grafana.com → My Account → API Keys)"
  type        = string
  sensitive   = true
}

variable "domain_name" {
  description = "API Domain (ex: api.rabi.food)"
  type        = string
}

variable "gotrue_service_key" {
  description = "Shared secret used by the app to authenticate against the GoTrue Admin API (GOTRUE_OPERATOR_TOKEN)"
  type        = string
  sensitive   = true
}

variable "grafana_url" {
  description = "Grafana Cloud instance URL (ex: https://yourstack.grafana.net)"
  type        = string
}

variable "grafana_admin_token" {
  description = "Grafana Cloud Service Account token with Editor or Admin role (for alerting provisioning)"
  type        = string
  sensitive   = true
}

variable "grafana_loki_datasource_name" {
  description = "Name of the Loki datasource in Grafana Cloud"
  type        = string
  default     = "grafana-loki"
}

variable "alert_email_to" {
  description = "Email address to receive alerts"
  type        = string
}

variable "stripe_secret_key" {
  description = "Stripe secret key (sk_live_... or sk_test_...)"
  type        = string
  sensitive   = true
}

variable "stripe_publishable_key" {
  description = "Stripe publishable key (pk_live_... or pk_test_...)"
  type        = string
}

