data "aws_ssm_parameter" "ecs_ami" {
  name = "/aws/service/ecs/optimized-ami/amazon-linux-2/recommended/image_id"
}

resource "aws_ecs_cluster" "main" {
  name = "rabi-food"

  tags = { Name = "rabi-food-cluster" }
}

resource "aws_cloudwatch_log_group" "app" {
  name              = "/ecs/rabi-food-app"
  retention_in_days = 30
}

resource "aws_cloudwatch_log_group" "alloy" {
  name              = "/ecs/rabi-food-alloy"
  retention_in_days = 7
}

resource "aws_cloudwatch_log_group" "gotrue" {
  name              = "/ecs/rabi-food-gotrue"
  retention_in_days = 30
}

locals {
  alloy_config = <<-ALLOY
    prometheus.scrape "app" {
      targets         = [{"__address__" = "localhost:${var.app_port}"}]
      forward_to      = [prometheus.remote_write.grafana_cloud.receiver]
      scrape_interval = "15s"
    }

    prometheus.remote_write "grafana_cloud" {
      endpoint {
        url = env("GRAFANA_PROMETHEUS_URL")
        basic_auth {
          username = env("GRAFANA_PROMETHEUS_INSTANCE_ID")
          password = env("GRAFANA_TOKEN")
        }
      }
    }
  ALLOY
}

resource "aws_ecs_task_definition" "app" {
  family                   = "rabi-food-app"
  network_mode             = "awsvpc"
  requires_compatibilities = ["EC2"]
  cpu                      = var.ecs_cpu
  memory                   = var.ecs_memory
  execution_role_arn       = aws_iam_role.ecs_execution.arn

  container_definitions = jsonencode([
    {
      name  = "app"
      image = "${aws_ecr_repository.app.repository_url}:latest"

      portMappings = [{
        containerPort = var.app_port
        protocol      = "tcp"
      }]

      environment = [
        { name = "ENV",           value = "production" },
        { name = "APP_PORT",      value = tostring(var.app_port) },
        { name = "DATABASE_HOST", value = aws_db_instance.postgres.address },
        { name = "DATABASE_NAME", value = var.db_name },
        { name = "DATABASE_USER", value = var.db_user },
        { name = "DATABASE_PORT", value = "5432" },
        # GoTrue runs as a sidecar — reachable only on localhost within the task
        { name = "GOTRUE_URL",    value = "http://localhost:9999" },
      ]

      secrets = [
        { name = "DATABASE_PASSWORD",  valueFrom = aws_secretsmanager_secret.db_password.arn },
        { name = "AUTH_SECRET",        valueFrom = aws_secretsmanager_secret.auth_secret.arn },
        { name = "GOTRUE_SERVICE_KEY", valueFrom = aws_secretsmanager_secret.gotrue_service_key.arn },
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.app.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    },
    {
      name      = "gotrue"
      image     = "supabase/gotrue:v2.189.0-rc.13"
      essential = true

      # No portMappings — GoTrue is not publicly reachable, only via localhost within the task
      environment = [
        { name = "GOTRUE_API_HOST",           value = "0.0.0.0" },
        { name = "GOTRUE_API_PORT",           value = "9999" },
        { name = "API_EXTERNAL_URL",          value = "http://localhost:9999" },
        { name = "GOTRUE_DB_DRIVER",          value = "postgres" },
        { name = "GOTRUE_JWT_EXP",            value = "3600" },
        { name = "GOTRUE_DISABLE_SIGNUP",     value = "true" },
        { name = "GOTRUE_MAILER_AUTOCONFIRM", value = "true" },
        { name = "GOTRUE_SITE_URL",           value = "http://localhost:3000" },
        { name = "GOTRUE_LOG_LEVEL",          value = "info" },
      ]

      secrets = [
        { name = "GOTRUE_JWT_SECRET",      valueFrom = aws_secretsmanager_secret.auth_secret.arn },
        { name = "GOTRUE_DB_DATABASE_URL", valueFrom = aws_secretsmanager_secret.gotrue_db_url.arn },
        { name = "GOTRUE_OPERATOR_TOKEN",  valueFrom = aws_secretsmanager_secret.gotrue_service_key.arn },
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.gotrue.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    },
    {
      name      = "alloy"
      image     = "grafana/alloy:v1.14.0"
      essential = false

      entryPoint = ["/bin/sh", "-c"]
      command    = ["printf '%s' \"$ALLOY_CONFIG\" > /tmp/config.alloy && exec /bin/alloy run /tmp/config.alloy"]

      environment = [
        { name = "ALLOY_CONFIG",                   value = local.alloy_config },
        { name = "GRAFANA_PROMETHEUS_URL",          value = var.grafana_prometheus_url },
        { name = "GRAFANA_PROMETHEUS_INSTANCE_ID",  value = var.grafana_prometheus_instance_id },
      ]

      secrets = [
        { name = "GRAFANA_TOKEN", valueFrom = aws_secretsmanager_secret.grafana_token.arn },
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.alloy.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "ecs"
        }
      }
    }
  ])

  lifecycle {
    ignore_changes = [container_definitions]
  }
}

resource "aws_ecs_service" "app" {
  name            = "rabi-food-app"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.app.arn
  desired_count = 1

  capacity_provider_strategy {
    capacity_provider = aws_ecs_capacity_provider.ec2.name
    weight            = 1
  }

  network_configuration {
    subnets         = [aws_subnet.public_a.id, aws_subnet.public_b.id]
    security_groups = [aws_security_group.ecs_tasks.id]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.app.arn
    container_name   = "app"
    container_port   = var.app_port
  }

  depends_on = [aws_lb_listener.http_redirect, aws_ecs_capacity_provider.ec2]

  lifecycle {
    ignore_changes = [task_definition, desired_count]
  }
}

resource "aws_launch_template" "ecs_instance" {
  name_prefix   = "rabi-food-ecs-"
  image_id      = data.aws_ssm_parameter.ecs_ami.value
  instance_type = var.ec2_instance_type

  iam_instance_profile {
    name = aws_iam_instance_profile.ecs_instance.name
  }

  network_interfaces {
    associate_public_ip_address = true
    security_groups             = [aws_security_group.ecs_tasks.id]
  }

  user_data = base64encode("#!/bin/bash\necho ECS_CLUSTER=${aws_ecs_cluster.main.name} >> /etc/ecs/ecs.config\n")

  tag_specifications {
    resource_type = "instance"
    tags          = { Name = "rabi-food-ecs-instance" }
  }
}

resource "aws_autoscaling_group" "ecs" {
  name                = "rabi-food-ecs-asg"
  min_size            = 1
  max_size            = 1
  desired_capacity    = 1
  vpc_zone_identifier = [aws_subnet.public_a.id, aws_subnet.public_b.id]

  launch_template {
    id      = aws_launch_template.ecs_instance.id
    version = "$Latest"
  }

  tag {
    key                 = "AmazonECSManaged"
    value               = "true"
    propagate_at_launch = true
  }
}

resource "aws_ecs_capacity_provider" "ec2" {
  name = "rabi-food-ec2"

  auto_scaling_group_provider {
    auto_scaling_group_arn = aws_autoscaling_group.ecs.arn

    managed_scaling {
      status          = "ENABLED"
      target_capacity = 100
    }
  }
}

resource "aws_ecs_cluster_capacity_providers" "main" {
  cluster_name       = aws_ecs_cluster.main.name
  capacity_providers = [aws_ecs_capacity_provider.ec2.name]

  default_capacity_provider_strategy {
    capacity_provider = aws_ecs_capacity_provider.ec2.name
    weight            = 1
  }
}
