resource "aws_ecs_cluster" "main" {
  name = "rabi-food"

  tags = { Name = "rabi-food-cluster" }
}

resource "aws_cloudwatch_log_group" "app" {
  name              = "/ecs/rabi-food-app"
  retention_in_days = 30
}

resource "aws_ecs_task_definition" "app" {
  family                   = "rabi-food-app"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.ecs_cpu
  memory                   = var.ecs_memory
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([{
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
    ]

    secrets = [
      { name = "DATABASE_PASSWORD", valueFrom = aws_secretsmanager_secret.db_password.arn },
      { name = "AUTH_SECRET",       valueFrom = aws_secretsmanager_secret.auth_secret.arn },
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.app.name
        "awslogs-region"        = var.aws_region
        "awslogs-stream-prefix" = "ecs"
      }
    }
  }])

  lifecycle {
    ignore_changes = [container_definitions]
  }
}

resource "aws_ecs_service" "app" {
  name            = "rabi-food-app"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.app.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [aws_subnet.public_a.id, aws_subnet.public_b.id]
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.app.arn
    container_name   = "app"
    container_port   = var.app_port
  }

  depends_on = [aws_lb_listener.app]

  lifecycle {
    ignore_changes = [task_definition, desired_count]
  }
}
