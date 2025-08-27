# Go Task Management API - Dredd Test Suite

Bu proje, Go Task Management API'si için kapsamlı API testlerini içermektedir. Dredd framework kullanılarak OpenAPI 3.0.3 spesifikasyonuna göre otomatik testler yapılmaktadır.

## 🎯 Test Kapsamı

### Authentication Endpoints (3)
- `POST /register` - Kullanıcı kaydı (201, 400)
- `POST /login` - Giriş yapma (200, 401) 
- `POST /logout` - Çıkış yapma (200, 401)

### Task Management Endpoints (6)
- `GET /tasks/public` - Public görevler (200)
- `GET /tasks` - Kullanıcı görevleri (200, 401)
- `POST /tasks` - Görev oluşturma (201, 400, 401)
- `GET /tasks/{id}` - Görev detayı (200, 401, 404)
- `PUT /tasks/{id}` - Görev güncelleme (200, 400, 401, 404)
- `DELETE /tasks/{id}` - Görev silme (200, 401, 404)

**Toplam: 22 Test Senaryosu** ✅

## 🛠️ Teknolojiler

- **Dredd v14.1.0** - API testing framework
- **Node.js** - JavaScript runtime
- **Axios** - HTTP client for hooks
- **OpenAPI 3.0.3** - API spesifikasyonu
- **JWT Authentication** - Token tabanlı kimlik doğrulama

## 📁 Dosya Yapısı

```
dredd_testing/
├── dredd-simple.yml       # Dredd konfigürasyonu
├── hooks_fixed.js         # Test hooks ve authentication logic
├── openapi_fixed.yaml     # Düzeltilmiş OpenAPI spesifikasyonu
├── package.json           # Node.js dependencies
├── run_tests.ps1          # PowerShell test scripti
├── run_tests.sh           # Bash test scripti
├── run_tests.bat          # Windows batch test scripti
└── README.md              # Bu dosya
```

## 🚀 Kurulum

### 1. Dependencies Kurulumu
```bash
npm install
```

### 2. PostgreSQL Database
Veritabanının çalışır durumda olduğundan emin olun:
- Host: localhost:5432
- Database: go_taskmanagement
- User: postgres
- Password: 1234

## ▶️ Testleri Çalıştırma

### Otomatik Test (Önerilen)
```powershell
# PowerShell
.\run_tests.ps1

# Bash
./run_tests.sh

# Windows Batch
run_tests.bat
```

### Manuel Test
```bash
# Server'ı ayrı terminal'de başlat
cd ..
go run main.go

# Testleri çalıştır
npx dredd --config=dredd-simple.yml
```

## 🔧 Test Konfigürasyonu

### dredd-simple.yml
- **Server**: http://localhost:8080 
- **Blueprint**: openapi_fixed.yaml
- **Hooks**: hooks_fixed.js
- **Reporter**: cli (konsol çıktısı)
- **Wait Time**: 3 saniye (server başlatma için)

### hooks_fixed.js Özellikleri
- ✅ **Otomatik User Registration**: Unique test kullanıcıları
- ✅ **JWT Token Management**: Otomatik token alma ve kullanma
- ✅ **Dynamic Task Creation**: Test için gerçek task oluşturma
- ✅ **Smart ID Replacement**: Task ID'lerini dinamik değiştirme
- ✅ **401 Test Scenarios**: Invalid token testleri
- ✅ **Data Cleanup**: Test sonrası temizlik

## 📊 Test Sonuçları

Başarılı test çalıştırması örneği:
```
🚀 Starting Dredd API Tests...
📋 Setting up test environment...
✅ User registered successfully
✅ Authentication successful, token obtained
✅ Test task created with ID: 23
✅ Test environment ready

� Test Results:
✅ POST /register - 201, 400
✅ POST /login - 200, 401  
✅ POST /logout - 200, 401
✅ GET /tasks/public - 200
✅ GET /tasks - 200, 401
✅ POST /tasks - 201, 400, 401
✅ GET /tasks/{id} - 200, 401, 404
✅ PUT /tasks/{id} - 200, 400, 401, 404
✅ DELETE /tasks/{id} - 200, 401, 404

🏁 Dredd API Tests Completed
complete: 22 passing, 0 failing, 0 errors, 0 skipped, 22 total
```

## 🔐 Authentication Test Detayları

### JWT Token Flow
1. **Setup Phase**: Test kullanıcısı kaydedilir
2. **Login**: JWT token alınır
3. **Protected Endpoints**: Token ile API'ye erişim
4. **401 Tests**: Invalid token senaryoları
5. **Cleanup**: Test verileri temizlenir

### Test Scenarios
- **Valid Authentication**: Geçerli credentials ile login
- **Invalid Authentication**: Yanlış credentials ile 401
- **Token Usage**: Protected endpoint'lere token ile erişim
- **Token Validation**: Invalid/expired token ile 401
- **Registration**: Unique kullanıcı kaydı ve duplicate kontrolü

## 🔄 CI/CD Integration

Bu test suite GitHub Actions veya diğer CI/CD pipeline'larına entegre edilebilir:

```yaml
- name: Run API Tests
  run: |
    cd dredd_testing
    npm install
    ./run_tests.sh
```

## 🐛 Troubleshooting

### Common Issues

1. **Server Connection Error**
   ```
   Solution: go run main.go ile server'ı başlatın
   ```

2. **Database Connection Error**
   ```
   Solution: PostgreSQL servisinin çalıştığından emin olun
   ```

3. **Test Timeout**
   ```
   Solution: dredd-simple.yml'de server-wait değerini artırın
   ```

4. **Authentication Failures**
   ```
   Solution: Database'de test kullanıcılarının temizlendiğinden emin olun
   ```

## 📈 Test Coverage

- **HTTP Methods**: GET, POST, PUT, DELETE
- **Status Codes**: 200, 201, 400, 401, 404
- **Authentication**: JWT token validation
- **Data Validation**: Request/response body validation
- **Error Handling**: Error response format validation
- **CRUD Operations**: Complete task lifecycle testing

## 🔮 Gelecek Geliştirmeler

- [ ] Performance testing integration
- [ ] Load testing scenarios
- [ ] API rate limiting tests
- [ ] File upload endpoint tests
- [ ] WebSocket connection tests
- [ ] Microservice integration tests

---

**Geliştirici**: Hakan Çetinkaya  
**Tarih**: Ağustos 2025  
**Status**: ✅ Tüm testler başarılı
