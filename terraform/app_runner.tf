# ---------------------------------------------
# App Runner - サービス
# ---------------------------------------------
resource "aws_apprunner_service" "back_echo_service" {
  service_name = "${var.project}-${var.environment}-service"

  source_configuration {
    authentication_configuration {
      access_role_arn = aws_iam_role.apprunner_access_role.arn
    }

    image_repository {
      image_identifier      = "${aws_ecr_repository.echo_back_practice_repo.repository_url}:latest"
      image_repository_type = "ECR"
      image_configuration {
        port = var.api_port

        runtime_environment_variables = {
          SUPABASE_URL = var.supabase_url
          CORS_ADDRESS = var.cors_address
          PORT         = var.api_port
          REDIS_PORT   = var.redis_port
          REDIS_HOST   = aws_elasticache_replication_group.redis.primary_endpoint_address
        }
      }
    }
  }

  instance_configuration {
    cpu               = "1024"
    memory            = "2048"
    instance_role_arn = aws_iam_role.apprunner_access_role.arn
  }

  network_configuration {
    egress_configuration {
      egress_type       = "VPC"
      vpc_connector_arn = aws_apprunner_vpc_connector.vpc_connector.arn
    }
  }

  health_check_configuration {
    protocol            = "HTTP"
    path                = "/"
    interval            = 10
    timeout             = 5
    healthy_threshold   = 1
    unhealthy_threshold = 5
  }

  depends_on = [aws_elasticache_replication_group.redis, aws_iam_role_policy.apprunner_access_policy]

  tags = {
    Name    = "${var.project}-${var.environment}-app-runner"
    Project = var.project
    Env     = var.environment
  }
}

# ---------------------------------------------
# App Runner - VPC Connector
# ---------------------------------------------
resource "aws_apprunner_vpc_connector" "vpc_connector" {
  vpc_connector_name = "${var.project}-${var.environment}-ar-vpc"
  subnets            = [aws_subnet.private_subnet_1a.id, aws_subnet.private_subnet_1c.id]
  security_groups    = [aws_security_group.app_runner_sg.id]

  tags = {
    Name    = "${var.project}-${var.environment}-vpc-connector"
    Project = var.project
    Env     = var.environment
  }
}
