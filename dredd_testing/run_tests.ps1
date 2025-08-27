# Dredd API Test Script for PowerShell
# Bu script otomatik olarak server'ı başlatır ve API testlerini çalıştırır

Write-Host "🚀 Go Task Management API - Otomatik Test Başlatılıyor..." -ForegroundColor Green

# Dredd test dizinine geç
Set-Location $PSScriptRoot

Write-Host "📋 Test ortamı hazırlanıyor..." -ForegroundColor Yellow

# Dredd testlerini çalıştır (server otomatik başlatılacak)
Write-Host "🔥 API server başlatılıyor ve testler çalıştırılıyor..." -ForegroundColor Cyan
npx dredd --config=dredd-simple.yml

Write-Host "✅ Test tamamlandı!" -ForegroundColor Green
Read-Host -Prompt "Devam etmek için Enter'a basın"
