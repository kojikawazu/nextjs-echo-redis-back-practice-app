# ---------------------------------------------
# ElastiCache - Subnet Group
# ---------------------------------------------
resource "aws_elasticache_subnet_group" "elasticache_subnet_group" {
  name       = "${var.project}-${var.environment}-elasticache-subnet-group"
  subnet_ids = [aws_subnet.private_subnet_1a.id, aws_subnet.private_subnet_1c.id]

  tags = {
    Name    = "${var.project}-${var.environment}-elasticache-subnet-group"
    Project = var.project
    Env     = var.environment
  }
}

# ---------------------------------------------
# ElastiCache - Redis Replication Group
# ---------------------------------------------
resource "aws_elasticache_replication_group" "redis" {
  replication_group_id       = "${var.project}-${var.environment}-redis-rg"
  description                = "Redis replication group for ${var.project}-${var.environment}"
  engine                     = "redis"
  engine_version             = "6.x"
  node_type                  = "cache.t3.micro"
  num_cache_clusters         = 1
  automatic_failover_enabled = false
  subnet_group_name          = aws_elasticache_subnet_group.elasticache_subnet_group.name
  security_group_ids         = [aws_security_group.elasticache_sg.id]
  parameter_group_name       = "default.redis6.x"

  tags = {
    Name    = "${var.project}-${var.environment}-redis-replication-group"
    Project = var.project
    Env     = var.environment
  }
}

output "elasticache_primary_endpoint" {
  value = aws_elasticache_replication_group.redis.primary_endpoint_address
}
