resource "stripe_apple_pay_domain" "app" {
  domain_name = var.domain_name
}

resource "stripe_webhook_endpoint" "subscription" {
  url = "https://${var.domain_name}/subscription/confirm-payment"

  enabled_events = [
    "payment_intent.succeeded",
    "payment_intent.payment_failed",
  ]

  description = "Subscription payment events"
}
