resource "aws_sesv2_dedicated_ip_pool" "this" {
  pool_name = module.this.id
  tags      = module.this.tags
}

output "pool_name" {
  value = aws_sesv2_dedicated_ip_pool.this.pool_name
}
