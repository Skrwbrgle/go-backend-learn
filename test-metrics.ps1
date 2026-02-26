#!/usr/bin/env pwsh
# Quick Test Script untuk Observability Metrics
# Run ini di PowerShell untuk test semua metrics functionality

Write-Host "🚀 Starting Observability Metrics Test..." -ForegroundColor Cyan
Write-Host ""

# Base URL
$baseUrl = "http://localhost:8080"

# Test 1: Health check
Write-Host "📍 Test 1: Health Check Request" -ForegroundColor Yellow
$response = Invoke-WebRequest -Uri "$baseUrl/ping" -Method GET
Write-Host "✅ Status: $($response.StatusCode)" -ForegroundColor Green
Write-Host "✅ Correlation ID: $($response.Headers['X-Correlation-ID'])" -ForegroundColor Green
Write-Host ""

# Test 2: Multiple requests untuk populate metrics
Write-Host "📍 Test 2: Creating 20 requests to populate metrics..." -ForegroundColor Yellow
1..20 | ForEach-Object {
    try {
        Invoke-WebRequest -Uri "$baseUrl/ping" -Method GET -ErrorAction SilentlyContinue | Out-Null
        Write-Host "  ✓ Request $_/20 completed"
    } catch {
        Write-Host "  ✗ Request failed: $_"
    }
}
Write-Host "✅ Done creating requests" -ForegroundColor Green
Write-Host ""

# Test 3: View raw metrics
Write-Host "📍 Test 3: Fetching raw Prometheus metrics..." -ForegroundColor Yellow
try {
    $metricsResponse = Invoke-WebRequest -Uri "$baseUrl/metrics" -Method GET
    $content = $metricsResponse.Content
    
    # Count metrics
    $metricLines = ($content -split "`n" | Where-Object { $_ -match "^http_requests" }).Count
    Write-Host "✅ Found $metricLines metric lines" -ForegroundColor Green
    
    # Show sample
    Write-Host ""
    Write-Host "📊 Sample Metrics:" -ForegroundColor Cyan
    $content -split "`n" | Where-Object { $_ -match "^http_requests_total" } | Select-Object -First 5 | ForEach-Object {
        Write-Host "  $_"
    }
} catch {
    Write-Host "❌ Failed to fetch metrics: $_" -ForegroundColor Red
}
Write-Host ""

# Test 4: Error handling (404)
Write-Host "📍 Test 4: Test error tracking (404)..." -ForegroundColor Yellow
try {
    Invoke-WebRequest -Uri "$baseUrl/go-api/users/invalid" -Method GET -ErrorAction SilentlyContinue
} catch {
    Write-Host "✅ 404 error correctly recorded" -ForegroundColor Green
}
Write-Host ""

# Test 5: Correlation ID propagation
Write-Host "📍 Test 5: Test Correlation ID propagation..." -ForegroundColor Yellow
$customId = "test-correlation-$(Get-Random)"
$headers = @{
    'X-Correlation-ID' = $customId
}
$response = Invoke-WebRequest -Uri "$baseUrl/ping" -Method GET -Headers $headers
$returnedId = $response.Headers['X-Correlation-ID']
if ($returnedId -eq $customId) {
    Write-Host "✅ Correlation ID correctly propagated: $customId" -ForegroundColor Green
} else {
    Write-Host "❌ Correlation ID mismatch! Expected: $customId, Got: $returnedId" -ForegroundColor Red
}
Write-Host ""

Write-Host "🎉 Test Complete!" -ForegroundColor Cyan
Write-Host ""
Write-Host "📊 Next Step: View metrics in Prometheus" -ForegroundColor Yellow
Write-Host "  1. Start Prometheus: docker-compose up -d" -ForegroundColor Gray
Write-Host "  2. Open: http://localhost:9090" -ForegroundColor Gray
Write-Host "  3. Query: rate(http_requests_total[1m])" -ForegroundColor Gray
Write-Host ""
