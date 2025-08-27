#!/bin/bash

# Dredd API Test Script
# Bu script otomatik olarak server'ı başlatır ve API testlerini çalıştırır

echo "🚀 Go Task Management API - Otomatik Test Başlatılıyor..."

# Dredd test dizinine geç
cd "$(dirname "$0")"

echo "📋 Test ortamı hazırlanıyor..."

# Dredd testlerini çalıştır (server otomatik başlatılacak)
echo "🔥 API server başlatılıyor ve testler çalıştırılıyor..."
npx dredd --config=dredd-simple.yml

echo "✅ Test tamamlandı!"
