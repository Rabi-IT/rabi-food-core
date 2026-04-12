output "alb_dns" {
  description = "ALB DNS name — use as API endpoint"
  value       = aws_lb.app.dns_name
}

output "ecr_registry_url" {
  description = "ECR registry URL for the app image"
  value       = aws_ecr_repository.app.repository_url
}

output "db_host" {
  description = "RDS endpoint (use as DATABASE_HOST)"
  value       = aws_db_instance.postgres.address
}
