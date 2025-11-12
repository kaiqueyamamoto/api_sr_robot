# Prometheus Metrics Documentation

## Overview

This application exposes Prometheus metrics for monitoring performance, authentication, database operations, and more.

## Metrics Endpoint

```
GET http://localhost:8080/metrics
```

---

## Available Metrics

### 1. HTTP Metrics

#### `http_requests_total`

**Type:** Counter  
**Description:** Total number of HTTP requests  
**Labels:**

- `method` - HTTP method (GET, POST, etc.)
- `endpoint` - Request endpoint path
- `status` - HTTP status code

**Example:**

```promql
http_requests_total{method="POST", endpoint="/auth/login", status="200"}
```

#### `http_request_duration_seconds`

**Type:** Histogram  
**Description:** Duration of HTTP requests in seconds  
**Labels:**

- `method` - HTTP method
- `endpoint` - Request endpoint path
- `status` - HTTP status code

**Example:**

```promql
# Average request duration
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# 95th percentile latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

---

### 2. Authentication Metrics

#### `auth_attempts_total`

**Type:** Counter  
**Description:** Total number of authentication attempts  
**Labels:**

- `type` - Authentication type (login, register)
- `status` - Result status (success, failure)

**Example:**

```promql
# Login success rate
rate(auth_attempts_total{type="login", status="success"}[5m])

# Failed login attempts
rate(auth_attempts_total{type="login", status="failure"}[5m])
```

#### `jwt_tokens_issued_total`

**Type:** Counter  
**Description:** Total number of JWT tokens issued

**Example:**

```promql
# Tokens issued per minute
rate(jwt_tokens_issued_total[1m]) * 60
```

#### `jwt_token_validation_failures_total`

**Type:** Counter  
**Description:** Total number of JWT token validation failures  
**Labels:**

- `reason` - Failure reason (expired, invalid, malformed, missing)

**Example:**

```promql
# Expired tokens
jwt_token_validation_failures_total{reason="expired"}

# All validation failures
sum(jwt_token_validation_failures_total)
```

#### `active_users_total`

**Type:** Gauge  
**Description:** Number of currently active users

**Example:**

```promql
active_users_total
```

---

### 3. Database Metrics

#### `database_operations_total`

**Type:** Counter  
**Description:** Total number of database operations  
**Labels:**

- `operation` - Database operation (insert, find, update, delete)
- `collection` - MongoDB collection name
- `status` - Operation status (success, failure)

**Example:**

```promql
# Database operations per second
rate(database_operations_total[1m])

# Failed database operations
database_operations_total{status="failure"}
```

#### `database_operation_duration_seconds`

**Type:** Histogram  
**Description:** Duration of database operations in seconds  
**Labels:**

- `operation` - Database operation
- `collection` - MongoDB collection name

**Example:**

```promql
# Average database query time
rate(database_operation_duration_seconds_sum[5m]) / rate(database_operation_duration_seconds_count[5m])

# Slow queries (95th percentile)
histogram_quantile(0.95, rate(database_operation_duration_seconds_bucket[5m]))
```

---

### 4. Chat Metrics

#### `chat_messages_total`

**Type:** Counter  
**Description:** Total number of chat messages processed

**Example:**

```promql
# Messages per minute
rate(chat_messages_total[1m]) * 60
```

---

### 5. System Metrics

#### `active_connections`

**Type:** Gauge  
**Description:** Number of active connections

**Example:**

```promql
active_connections
```

---

## Useful Queries

### Performance Monitoring

```promql
# Request rate (requests per second)
rate(http_requests_total[5m])

# Error rate
sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))

# Average response time
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# P95 latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# P99 latency
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))
```

### Authentication Monitoring

```promql
# Login success rate
sum(rate(auth_attempts_total{type="login", status="success"}[5m])) /
sum(rate(auth_attempts_total{type="login"}[5m]))

# Failed login attempts (potential attacks)
rate(auth_attempts_total{type="login", status="failure"}[5m])

# Token validation failures by reason
sum by (reason) (rate(jwt_token_validation_failures_total[5m]))
```

### Database Monitoring

```promql
# Database query rate
sum(rate(database_operations_total[5m]))

# Database error rate
sum(rate(database_operations_total{status="failure"}[5m])) /
sum(rate(database_operations_total[5m]))

# Slow database queries
histogram_quantile(0.99, rate(database_operation_duration_seconds_bucket[5m]))
```

### Business Metrics

```promql
# User registration rate
rate(auth_attempts_total{type="register", status="success"}[1h])

# Active users
active_users_total

# Chat activity
rate(chat_messages_total[5m])
```

---

## Setting Up Prometheus

### Option 1: Using Docker Compose

1. Start Prometheus and Grafana:

```bash
docker-compose -f docker-compose.metrics.yml up -d
```

2. Access Prometheus at `http://localhost:9090`
3. Access Grafana at `http://localhost:3000` (admin/admin)

### Option 2: Local Installation

1. Download Prometheus from https://prometheus.io/download/
2. Update `prometheus.yml` with your configuration
3. Run Prometheus:

```bash
./prometheus --config.file=prometheus.yml
```

---

## Grafana Dashboard Setup

### Add Prometheus Data Source

1. Login to Grafana (http://localhost:3000)
2. Go to Configuration → Data Sources
3. Add Prometheus data source
4. URL: `http://prometheus:9090` (if using Docker) or `http://localhost:9090`
5. Click "Save & Test"

### Import Dashboard

Create a new dashboard with panels for:

1. **Request Rate Panel**

   - Query: `rate(http_requests_total[5m])`
   - Visualization: Graph

2. **Response Time Panel**

   - Query: `histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))`
   - Visualization: Graph

3. **Error Rate Panel**

   - Query: `sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))`
   - Visualization: Stat

4. **Active Connections Panel**

   - Query: `active_connections`
   - Visualization: Stat

5. **Auth Success Rate Panel**

   - Query: `sum(rate(auth_attempts_total{status="success"}[5m])) / sum(rate(auth_attempts_total[5m]))`
   - Visualization: Gauge

6. **Database Operations Panel**
   - Query: `sum by (operation) (rate(database_operations_total[5m]))`
   - Visualization: Graph (stacked)

---

## Alerting Examples

### High Error Rate Alert

```yaml
groups:
  - name: api_alerts
    interval: 30s
    rules:
      - alert: HighErrorRate
        expr: |
          sum(rate(http_requests_total{status=~"5.."}[5m])) / 
          sum(rate(http_requests_total[5m])) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value | humanizePercentage }}"
```

### High Latency Alert

```yaml
- alert: HighLatency
  expr: |
    histogram_quantile(0.95, 
      rate(http_request_duration_seconds_bucket[5m])
    ) > 1
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "High response time detected"
    description: "P95 latency is {{ $value }}s"
```

### Failed Login Attempts Alert

```yaml
- alert: SuspiciousLoginAttempts
  expr: |
    rate(auth_attempts_total{type="login", status="failure"}[5m]) > 10
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "High number of failed login attempts"
    description: "{{ $value }} failed logins per second"
```

---

## Best Practices

1. **Monitor continuously** - Set up Grafana dashboards for real-time monitoring
2. **Set up alerts** - Configure Prometheus alerts for critical metrics
3. **Track trends** - Use longer time windows to identify trends
4. **Correlate metrics** - Look at multiple metrics together to diagnose issues
5. **Optimize queries** - Use recording rules for frequently used queries
6. **Set proper retention** - Configure appropriate data retention policies

---

## Troubleshooting

### Metrics not showing up

- Check if the app is running: `http://localhost:8080/health`
- Verify metrics endpoint: `curl http://localhost:8080/metrics`
- Check Prometheus targets: `http://localhost:9090/targets`

### High cardinality issues

- Avoid using user IDs or other high-cardinality values as labels
- Use the existing labels (method, endpoint, status)

### Performance impact

- Metrics collection has minimal overhead
- Histogram buckets are pre-configured for optimal performance

---

## Further Reading

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [PromQL Basics](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/naming/)
