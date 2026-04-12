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

variable "ecs_cpu" {
  description = "Fargate task CPU units (256, 512, 1024...)"
  type        = number
  default     = 256
}

variable "ecs_memory" {
  description = "Fargate task memory in MB"
  type        = number
  default     = 512
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
