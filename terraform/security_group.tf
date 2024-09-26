# ---------------------------------------------
# Security Group
# ---------------------------------------------
# App Runner Security Group
resource "aws_security_group" "app_runner_sg" {
  name        = "${var.project}-${var.environment}-app-runner-sg"
  description = "app runner security group"
  vpc_id      = aws_vpc.vpc.id

  tags = {
    Name    = "${var.project}-${var.environment}-app-runner-sg"
    Project = var.project
    Env     = var.environment
  }
}

resource "aws_security_group_rule" "app_runner_in_8080" {
  security_group_id = aws_security_group.app_runner_sg.id
  type              = "ingress"
  protocol          = "tcp"
  from_port         = var.api_port
  to_port           = var.api_port
  cidr_blocks       = [var.igw_address]
}

resource "aws_security_group_rule" "app_runner_out_elasticache" {
  security_group_id = aws_security_group.app_runner_sg.id
  type              = "egress"
  protocol          = "tcp"
  from_port         = var.redis_port
  to_port           = var.redis_port
  cidr_blocks       = [var.vpc_address]
}

resource "aws_security_group_rule" "app_runner_out_internet" {
  security_group_id = aws_security_group.app_runner_sg.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = [var.igw_address]
}

# ElastiCache Security Group
resource "aws_security_group" "elasticache_sg" {
  name        = "${var.project}-${var.environment}-elasticache-sg"
  description = "elasticache security group"
  vpc_id      = aws_vpc.vpc.id

  tags = {
    Name    = "${var.project}-${var.environment}-elasticache-sg"
    Project = var.project
    Env     = var.environment
  }
}

resource "aws_security_group_rule" "elasticache_in_app_runner" {
  security_group_id        = aws_security_group.elasticache_sg.id
  type                     = "ingress"
  protocol                 = "tcp"
  from_port                = var.redis_port
  to_port                  = var.redis_port
  source_security_group_id = aws_security_group.app_runner_sg.id
}

resource "aws_security_group_rule" "elasticache_egress" {
  security_group_id = aws_security_group.elasticache_sg.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = [var.vpc_address]
}
