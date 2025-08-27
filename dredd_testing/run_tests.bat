@echo off
rem Dredd API Test Script for Windows
rem Bu script otomatik olarak server'ı başlatır ve API testlerini çalıştırır

echo 🚀 Go Task Management API - Otomatik Test Başlatılıyor...

rem Dredd test dizinine geç
cd /d "%~dp0"

echo 📋 Test ortamı hazırlanıyor...

rem Dredd testlerini çalıştır (server otomatik başlatılacak)
echo 🔥 API server başlatılıyor ve testler çalıştırılıyor...
npx dredd --config=dredd-simple.yml

echo ✅ Test tamamlandı!
pause
