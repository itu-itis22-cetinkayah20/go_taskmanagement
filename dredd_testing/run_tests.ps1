# Dredd API Test Script for PowerShell
# Bu script otomatik olarak server'ı başlatır ve API testlerini çalıştırır

Write-Host "🚀 Starting Dredd API Tests with Automatic Server Start..." -ForegroundColor Green
Write-Host "📋 Dredd will automatically start the Go server" -ForegroundColor Yellow
Write-Host "⏳ Please wait while server starts and tests execute..." -ForegroundColor Cyan

# Run Dredd with automatic server startup
Write-Host "`n🧪 Running Dredd tests..." -ForegroundColor Green

# Dredd test dizinine geç
Set-Location $PSScriptRoot

Write-Host "📋 Test ortamı hazırlanıyor..." -ForegroundColor Yellow

# Dredd testlerini çalıştır (server otomatik başlatılacak)
Write-Host "🔥 API server başlatılıyor ve testler çalıştırılıyor..." -ForegroundColor Cyan
npx dredd --config=dredd-simple.yml

Write-Host "✅ Test tamamlandı!" -ForegroundColor Green
Read-Host -Prompt "Devam etmek için Enter'a basın"
