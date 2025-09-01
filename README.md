# Go Task Management API

Bu proje, modern web teknolojileri kullanılarak geliştirilmiş, JWT tabanlı kimlik doğrulama sistemi ile kullanıcıların görev yönetimini sağlayan kapsamlı bir REST API'dir.

## 🚀 Özellikler

### 🔐 Kimlik Doğrulama & Güvenlik
- JWT (JSON Web Token) tabanlı kimlik doğrulama
- BCrypt ile şifre hashleme
- Bearer token ile API endpoint koruması
- Middleware tabanlı authorization

### 📊 Veritabanı Yönetimi
- PostgreSQL 17 veritabanı desteği
- GORM ORM ile gelişmiş veritabanı yönetimi
- Otomatik database migration
- Soft delete desteği
- Test ve production ortamları için ayrı database konfigürasyonu

### 🎯 Görev Yönetimi
- Kullanıcı bazlı görev oluşturma, güncelleme ve silme
- Görev durumu takibi (pending, in_progress, completed)
- Öncelik seviyeleri (low, medium, high)
- Public görevler desteği
- Detaylı görev filtreleme

### 📚 API Dokümantasyonu
- Swagger/OpenAPI 3.0.3 desteği
- Interaktif API dokümantasyonu (`/swagger/` endpoint)
- Kapsamlı endpoint açıklamaları ve örnekler

### 🧪 Test & Kalite Güvencesi
- Unit testler ve integration testler
- Contract testing (OpenAPI spec validation)
- **Dredd API testing framework** ile 22 test senaryosu (100% başarı)
- **Otomatik server başlatma** ve test execution
- **JWT authentication flow** testleri
- **Cross-platform test scripts** (PowerShell, Bash, Batch)
- GitHub Actions CI/CD pipeline
- Automated test execution

## 🛠️ Teknoloji Stack

- **Backend Framework:** Go Fiber v2
- **Database:** PostgreSQL 17
- **ORM:** GORM
- **Authentication:** JWT + BCrypt
- **Documentation:** Swagger/OpenAPI
- **Testing:** Contract Testing + Dredd API Testing
- **CI/CD:** GitHub Actions
- **Environment Management:** Godotenv

## 📋 Gereksinimler

- Go 1.24.6+
- PostgreSQL 17
- Node.js 16+ (Dredd testleri için)
- Git

## ⚡ Hızlı Başlangıç

### 1. Projeyi Klonlayın
```bash
git clone https://github.com/itu-itis22-cetinkayah20/go_taskmanagement.git
cd go_taskmanagement
```

### 2. Bağımlılıkları Yükleyin
```bash
go mod tidy
```

### 3. Environment Değişkenlerini Ayarlayın
`.env` dosyası oluşturun:
```env
# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=1234
DB_NAME=go_taskmanagement
DB_PORT=5432
DB_SSLMODE=disable

# Test Database Configuration
TEST_DB_HOST=localhost
TEST_DB_USER=postgres
TEST_DB_PASSWORD=1234
TEST_DB_NAME=go_taskmanagement_test
TEST_DB_PORT=5432
TEST_DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key
```

### 4. PostgreSQL Veritabanını Hazırlayın
```bash
# PostgreSQL bağlantısı (psql)
createdb go_taskmanagement
createdb go_taskmanagement_test

# Veya PowerShell script ile otomatik kurulum
.\scripts\setup_db.ps1
```

### 5. Sunucuyu Başlatın
```bash
go run main.go
```

Sunucu `http://localhost:8080` adresinde çalışacaktır.

## 📖 API Dokümantasyonu

Swagger UI: `http://localhost:8080/swagger/`

### 🔓 Public Endpoints
- `POST /register` — Kullanıcı kaydı
- `POST /login` — Giriş ve JWT token alma
- `GET /tasks/public` — Herkesin görebileceği örnek görevler

### 🔐 Protected Endpoints (JWT Required)
- `GET /tasks` — Kullanıcının kendi görevleri
- `POST /tasks` — Yeni görev ekleme
- `GET /tasks/{id}` — Görev detayları
- `PUT /tasks/{id}` — Görev güncelleme
- `DELETE /tasks/{id}` — Görev silme
- `POST /logout` — Çıkış

## 🧪 Test Senaryoları

### Manual API Tests
```bash
go test ./tests -v -timeout=30s 
```

### Contract Testing
```bash
go test ./test/contract -v -timeout=30s
```

### 🎯 Dredd API Testing (22 Test Senaryosu - 100% Başarı)
```bash
cd dredd_testing
.\run_tests.ps1    # PowerShell (Önerilen)
# veya
./run_tests.sh     # Bash
# veya
run_tests.bat      # Windows Batch
```

#### Dredd Test Kapsamı:
- ✅ **Authentication Tests**: Register, Login, Logout (3 endpoint)
- ✅ **Task Management Tests**: CRUD operations (6 endpoint)
- ✅ **Error Scenarios**: 400, 401, 404 status testleri
- ✅ **JWT Token Flow**: Otomatik authentication ve token management
- ✅ **Dynamic Testing**: Real-time task creation ve ID replacement

**Toplam: 22 Test Senaryosu - Tümü Başarılı** ✅

### GitHub Actions CI/CD
- Her push ve pull request için otomatik test çalıştırma
- PostgreSQL service container ile database testleri
- Multi-environment test desteği
- **Dredd API testing integration** ile otomatik API validasyonu
- **22 test senaryosu** ile kapsamlı endpoint testing
- **Paralel job execution** ile hızlı test pipeline


## 🏗️ Proje Yapısı

```
go_taskmanagement/
├── .github/
│   └── workflows/          # GitHub Actions CI/CD
├── cmd/
│   └── test_db_connection/ # Database connection test utility
├── database/
│   └── database.go         # Database connection and migrations
├── docs/                   # Swagger documentation (auto-generated)
├── dredd_testing/          # 🎯 Dredd API testing (22 tests - 100% pass)
│   ├── dredd-simple.yml    # Dredd configuration
│   ├── hooks_fixed.js      # Test hooks & authentication logic
│   ├── openapi_fixed.yaml  # OpenAPI spec aligned with API
│   ├── run_tests.ps1       # PowerShell automation script
│   ├── run_tests.sh        # Bash automation script
│   ├── run_tests.bat       # Windows batch script
│   └── README.md           # Detailed testing documentation
├── handlers/
│   ├── user_handlers.go    # Authentication endpoints
│   └── task_handlers.go    # Task management endpoints
├── middleware/
│   └── auth.go            # JWT authentication middleware
├── models/
│   ├── user.go            # User model
│   └── task.go            # Task model
├── scripts/               # Database setup scripts
├── test/
│   ├── contract/          # Contract testing
│   └── testdata/          # Test data and OpenAPI specs
├── tests/                 # Manual API tests
├── main.go               # Application entry point
├── go.mod               # Go module dependencies
└── README.md           # Bu dosya
```

## 🔧 Konfigürasyon

### Database Ayarları
- **Development:** `go_taskmanagement` database
- **Testing:** `go_taskmanagement_test` database
- Otomatik migration ve test data seeding
- Connection pooling ve logging

### JWT Ayarları
- Token süresi: Konfigürasyona göre
- Secret key: Environment variable veya fallback
- Secure header validation

### CORS Ayarları
- Tüm origin'lere izin (development için)
- Production için kısıtlama önerilir



### 1. Testing
```bash

# Contract testleri
go test ./test/contract -v -timeout=30s
### 3. Documentation
```bash
# Swagger documentation güncelle
swag init
```

## 🌟 Öne Çıkan Özellikler

### Güvenlik
- ✅ JWT token tabanlı authentication
- ✅ BCrypt password hashing
- ✅ SQL injection koruması (GORM ORM)
- ✅ CORS middleware

### Performance
- ✅ Connection pooling
- ✅ Efficient database queries
- ✅ Preloading optimization
- ✅ Soft delete for data integrity

### Maintainability
- ✅ Clean architecture pattern
- ✅ Separation of concerns
- ✅ Comprehensive error handling
- ✅ Structured logging

### Testing
- ✅ Unit tests
- ✅ Integration tests
- ✅ Contract testing with OpenAPI validation
- ✅ **Dredd API testing** framework ile 22 test senaryosu
- ✅ **Otomatik authentication flow** testleri
- ✅ **Cross-platform automation** (PowerShell/Bash/Batch)
- ✅ **100% test başarı oranı** achieved
- ✅ CI/CD automation

## 🤝 Katkıda Bulunma

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 Lisans

Bu proje MIT lisansı altında yayınlanmıştır. Detaylar için `LICENSE` dosyasına bakınız.

## 📞 İletişim

- **Geliştirici:** Hakan Çetinkaya
- **GitHub:** [@itu-itis22-cetinkayah20](https://github.com/itu-itis22-cetinkayah20)
- **Proje Repo:** [go_taskmanagement](https://github.com/itu-itis22-cetinkayah20/go_taskmanagement)

## 🚀 Sonraki Adımlar

- [ ] Redis cache entegrasyonu
- [ ] Rate limiting middleware
- [ ] Email notification system
- [ ] Task assignment ve team features
- [ ] Advanced filtering ve pagination
- [ ] Mobile API optimization
- [ ] Docker containerization
- [ ] Kubernetes deployment
