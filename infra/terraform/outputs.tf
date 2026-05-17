output "alb_dns" {
  description = "ALB DNS name — aponte um CNAME para este valor no Netlify"
  value       = aws_lb.api.dns_name
}

output "acm_cname_name" {
  description = "CNAME name para validação do certificado no Netlify"
  value       = tolist(aws_acm_certificate.api.domain_validation_options)[0].resource_record_name
}

output "acm_cname_value" {
  description = "CNAME value para validação do certificado no Netlify"
  value       = tolist(aws_acm_certificate.api.domain_validation_options)[0].resource_record_value
}

output "ecr_registry_url" {
  description = "ECR registry URL for the api image"
  value       = aws_ecr_repository.api.repository_url
}

output "db_host" {
  description = "RDS endpoint (use as DATABASE_HOST)"
  value       = aws_db_instance.postgres.address
}

output "stripe_webhook_secret" {
  description = "Stripe webhook signing secret — use as STRIPE_WEBHOOK_SECRET on API"
  value       = stripe_webhook_endpoint.subscription.secret
  sensitive   = true
}
