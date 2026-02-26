# 📊 Observability with Prometheus Metrics - Learning Guide

## 🎯 Apa yang sudah kita setup

Kita telah mengimplementasi **HTTP Metrics Pipeline**:

```
┌─────────────────┐
│  HTTP Request   │
└────────┬────────┘
         │
         ▼
┌──────────────────────────────┐
│ CorrelationIDMiddleware      │ ← Add unique ID ke request
└────────┬─────────────────────┘
         │
         ▼
┌──────────────────────────────┐
│ LoggerMiddleware             │ ← Log request details
└────────┬─────────────────────┘
         │
         ▼
┌──────────────────────────────┐
│ MetricsMiddleware            │ ← Record ke Prometheus
│ ├─ Request counter           │
│ ├─ Request duration          │
│ ├─ Response size             │
│ ├─ Active requests           │
│ └─ Errors                    │
└────────┬─────────────────────┘
         │
         ▼
┌──────────────────────────────┐
│ Your Actual Handler          │ ← Proses request
└─────────────────────────────┘
```

## 🚀 Local Testing - Step by Step

### 1️⃣ Start the Application

```bash
cd d:\Yoga\self-projek\GO\go-backend-learn
.\app.exe

# Expected output:
# 🚀 Server running on :8080
```

### 2️⃣ Make Some Test Requests

Open **PowerShell** atau **curl**:

```bash
# Test 1: Health check
curl http://localhost:8080/ping

# Test 2: Create user (akan record metrics)
curl -X POST http://localhost:8080/go-api/auth/register `
  -H "Content-Type: application/json" `
  -d @'
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "role": "user"
}
'@

# Test 3: Get all users
curl http://localhost:8080/go-api/users

# Test 4: Not found (akan record 404 error)
curl http://localhost:8080/go-api/users/invalid-id

# Buat beberapa request untuk populate metrics
for ($i=1; $i -le 10; $i++) {
    curl http://localhost:8080/ping
}
```

### 3️⃣ View Raw Metrics (OpenMetrics Format)

Visit di browser: **http://localhost:8080/metrics**

Output akan keliatan seperti:

```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{endpoint="/ping",method="GET",status="2xx"} 10
http_requests_total{endpoint="/go-api/users",method="GET",status="2xx"} 1
http_requests_total{endpoint="/go-api/users",method="GET",status="4xx"} 1

# HELP http_request_duration_seconds HTTP request latencies in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{endpoint="/ping",method="GET",le="0.005"} 8
http_request_duration_seconds_bucket{endpoint="/ping",method="GET",le="0.01"} 10

# HELP http_response_size_bytes HTTP response size in bytes
# TYPE http_response_size_bytes histogram
http_response_size_bytes_bucket{endpoint="/ping",method="GET",le="1000",status="2xx"} 10
http_response_size_bytes_sum{endpoint="/ping",method="GET",status="2xx"} 150

# HELP http_requests_active Number of currently active HTTP requests
# TYPE http_requests_active gauge
http_requests_active{endpoint="/ping",method="GET"} 0

# HELP http_errors_total Total number of HTTP errors
# TYPE http_errors_total counter
http_errors_total{endpoint="/go-api/users",error_type="client",method="GET"} 1
```

---

## 📈 Setup Prometheus untuk Scrape Metrics

### File: `prometheus.yml`

Buat file `prometheus.yml` di root project:

```yaml
global:
  scrape_interval: 15s # Scrape setiap 15 detik
  evaluation_interval: 15s
  external_labels:
    monitor: "go-app-monitor"

scrape_configs:
  - job_name: "go-app"
    static_configs:
      - targets: ["localhost:8080"]
    metrics_path: "/metrics"
    scrape_interval: 5s # Scrape lebih sering untuk dev
```

### Jalankan Prometheus dengan Docker:

```bash
# Pastikan Docker installed, then:
docker run -d \
  --name prometheus \
  -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus

# Access: http://localhost:9090
```

---

## 🔍 Metrics Explained - Deep Dive

### 1. **HTTP Requests Counter** (`http_requests_total`)

**Apa**: Jumlah total request yang masuk
**Format**: `http_requests_total{endpoint="...",method="GET",status="2xx"}`
**Use case**:

- Track throughput (requests per second)
- Compare endpoints popularity
- Detect traffic anomalies

**Prometheus Query**:

```promql
# Total requests last 1 menit
rate(http_requests_total[1m])

# Requests per endpoint
sum by (endpoint) (rate(http_requests_total[1m]))

# Error rate (4xx dan 5xx)
sum by (endpoint) (rate(http_requests_total{status=~"4xx|5xx"}[1m]))
```

### 2. **Request Duration Histogram** (`http_request_duration_seconds`)

**Apa**: Durasi setiap request dalam detik (histogram dengan bucket)
**Use case**:

- Find slow endpoints
- Calculate percentiles (p50, p95, p99)
- Alert jika duration > threshold

**Prometheus Query**:

```promql
# Average latency per endpoint
avg by (endpoint) (rate(http_request_duration_seconds_sum[5m]))
/
avg by (endpoint) (rate(http_request_duration_seconds_count[5m]))

# 95th percentile latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# P99 latency per endpoint
histogram_quantile(0.99, sum by (endpoint, le) (rate(http_request_duration_seconds_bucket[5m])))
```

### 3. **Response Size** (`http_response_size_bytes`)

**Apa**: Ukuran response body dalam bytes
**Use case**:

- Monitor bandwidth usage
- Find endpoints yang return data besar
- Optimization opportunities

**Prometheus Query**:

```promql
# Total bytes transferred per minute
sum(rate(http_response_size_bytes_sum[1m]))

# Average response size
sum(rate(http_response_size_bytes_sum[5m]))
/
sum(rate(http_response_size_bytes_count[5m]))
```

### 4. **Active Requests** (`http_requests_active`)

**Apa**: Jumlah request yang **sedang running** saat ini (Gauge)
**Use case**:

- Find connection bottlenecks
- Detect slow requests accumulation
- Monitor concurrency

**Prometheus Query**:

```promql
# Current active requests
http_requests_active

# Max active requests observed
max(http_requests_active)

# By endpoint
sum by (endpoint) (http_requests_active)
```

### 5. **HTTP Errors** (`http_errors_total`)

**Apa**: Jumlah error (4xx client, 5xx server)
**Use case**:

- Track application reliability
- Alert on error spike
- Identify problematic endpoints

**Prometheus Query**:

```promql
# Error rate per second
rate(http_errors_total[1m])

# Error breakdown by type
sum by (error_type) (http_errors_total)

# Errors per endpoint
sum by (endpoint) (http_errors_total)
```

---

## 🔗 Correlation ID - Request Tracing

### Apa itu Correlation ID?

**Scenario**:
User request → Service A → Service B → Database

Dengan Correlation ID:

- Setiap service dapat tracing request yang sama
- Logs dari semua service bisa di-correlate
- Debug production issues jadi lebih mudah

### Implementasi di Codebase

**Dalam middleware** (already implemented):

```go
// Every request dapat unique X-Correlation-ID header
X-Correlation-ID: 550e8400-e29b-41d4-a716-446655440000
```

**Dalam logger** (next step):

```go
// Log dapat include correlation ID
logger.Info("Processing user",
  zap.String("correlation_id", GetCorrelationID(c)),
  zap.String("email", userEmail),
)
```

**Use case**:

```
Request in: X-Correlation-ID: abc123
↓
Service A logs: [abc123] Received request
↓
Service A calls Service B dengan X-Correlation-ID: abc123
↓
Service B logs: [abc123] Processing in service B
↓
Service B calls Database dengan correlation ID
↓
Database Connection logs: [abc123] Query executed
↓
Response back dengan: X-Correlation-ID: abc123
```

**Test correlation ID**:

```bash
# Without correlation ID (akan auto-generate)
curl http://localhost:8080/ping

# Dengan custom correlation ID
curl -H "X-Correlation-ID: my-custom-id-123" http://localhost:8080/ping

# Response header akan include correlation ID
# X-Correlation-ID: my-custom-id-123
```

---

## 📚 Key Concepts in Prometheus Metrics

### Metric Types

| Type          | Purpose                               | Example                                 |
| ------------- | ------------------------------------- | --------------------------------------- |
| **Counter**   | Accumulating value (only goes up)     | `http_requests_total` - total requests  |
| **Gauge**     | Point-in-time value (up/down)         | `http_requests_active` - current active |
| **Histogram** | Distribution of values dalam buckets  | `http_request_duration_seconds`         |
| **Summary**   | Like histogram tapi compute quantiles | Rarely used (use histogram instead)     |

### Labels (Dimensions)

```
http_requests_total{endpoint="/ping", method="GET", status="2xx"}
                    ↑                  ↑               ↑
                    Labels untuk query & filter
```

Labels penting:

- `endpoint` - Jenis request
- `method` - HTTP method (GET, POST, etc)
- `status` - Response code (2xx, 4xx, 5xx)

---

## 🛠️ Next: Integration dengan Logging

Saat ini logging dan metrics terpisah. Best practice adalah **combine keduanya**:

```go
// Better logging dengan correlation ID
ctx := c.Request.Context()
correlationID := middleware.GetCorrelationID(c)

logger.Info("User created",
  zap.String("correlation_id", correlationID),
  zap.String("user_id", user.ID.String()),
  zap.String("email", user.Email),
  zap.Int64("duration_ms", duration),
)
```

---

## ✅ Checklist: Metrics Observatory Ready

- [x] Prometheus client installed
- [x] Metrics defined (counter, histogram, gauge)
- [x] MetricsMiddleware implemented
- [x] CorrelationIDMiddleware implemented
- [x] `/metrics` endpoint setup
- [ ] Prometheus server setup (docker run command)
- [ ] Grafana dashboard (next learning phase)
- [ ] Database query metrics (next step)

---

## 📈 Practice Assignments

**Easy**:

1. Make 100 requests ke endpoint berbeda, lihat metrics change
2. Calculate error rate dari metrics
3. Find slowest endpoint via histogram

**Medium**: 4. Setup Prometheus + access dashboard 5. Create custom alert rule (e.g., error_rate > 5%)

**Hard**: 6. Implement database query metrics (similar pattern) 7. Create Grafana dashboard dengan custom queries

---

## 🎓 Learning Resources

- [Prometheus Docs](https://prometheus.io/docs/)
- [Go Prometheus Client](https://github.com/prometheus/client_golang)
- [PromQL Query Language](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Metrics Best Practices](https://prometheus.io/docs/practices/naming/)

---

## 🔧 Troubleshooting

**Q: `/metrics` endpoint returns 404**
A: Check `router.go` - pastikan `r.GET("/metrics", ...)` di setup

**Q: Metrics tidak increment**
A: Make requests ke endpoint terlebih dahulu, metrics tidak auto-populate

**Q: Prometheus query returns nothing**
A: Check di http://localhost:9090/graph, verify scrape_configs correct

**Q: Correlation ID not appearing in logs**
A: Update logger calls to include `GetCorrelationID(c)` as zap field

---

Good luck belajarnya! 🚀
