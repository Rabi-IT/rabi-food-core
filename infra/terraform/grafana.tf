data "grafana_cloud_stack" "main" {
  slug = var.grafana_stack_slug
}

data "grafana_data_source" "loki" {
  name = var.grafana_loki_datasource_name
}

resource "grafana_folder" "alerts" {
  title = "Alerts"
}

resource "grafana_contact_point" "email" {
  name = "email"

  email {
    addresses               = [var.alert_email_to]
    disable_resolve_message = false
  }
}

resource "grafana_notification_policy" "main" {
  contact_point = grafana_contact_point.email.name
  group_by      = ["alertname", "app"]

  group_wait      = "30s"
  group_interval  = "5m"
  repeat_interval = "3h"
}

resource "grafana_rule_group" "log_errors" {
  name             = "log-errors"
  folder_uid       = grafana_folder.alerts.uid
  interval_seconds = 60

  rule {
    name      = "LogErrorDetected"
    condition = "threshold"
    for       = "30s"

    labels = {
      severity = "warning"
    }

    annotations = {
      summary     = "Error detected in logs"
      description = "Errors were detected in the last 5 minutes."
    }

    data {
      ref_id         = "query"
      datasource_uid = data.grafana_data_source.loki.uid
      query_type     = "range"

      relative_time_range {
        from = 300
        to   = 0
      }

      model = jsonencode({
        expr      = "count_over_time({job=\"docker\"} |= \"error\" [5m])"
        queryType = "range"
        refId     = "query"
      })
    }

    data {
      ref_id         = "threshold"
      datasource_uid = "__expr__"

      relative_time_range {
        from = 300
        to   = 0
      }

      model = jsonencode({
        type  = "threshold"
        refId = "threshold"
        conditions = [{
          evaluator = { params = [0], type = "gt" }
          operator  = { type = "and" }
          query     = { params = ["query"] }
          reducer   = { type = "last" }
          type      = "query"
        }]
      })
    }

    notification_settings {
      contact_point = grafana_contact_point.email.name
    }
  }
}
