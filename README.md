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
- GitHub Actions CI/CD pipeline
- Automated test execution

## 🛠️ Teknoloji Stack

- **Backend Framework:** Go Fiber v2
- **Database:** PostgreSQL 17
- **ORM:** GORM
- **Authentication:** JWT + BCrypt
- **Documentation:** Swagger/OpenAPI
- **Testing:** Contract Testing with OpenAPI validation
- **CI/CD:** GitHub Actions
- **Environment Management:** Godotenv

## 📋 Gereksinimler

- Go 1.24.6+
- PostgreSQL 17
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

### Tests dosyasındaki manual testler
```bash
go test ./tests -v -timeout=30s 
```

### Contract Testing
```bash
go test ./test/contract -v -timeout=30s
```

### GitHub Actions CI/CD
- Her push ve pull request için otomatik test çalıştırma
- PostgreSQL service container ile database testleri
- Multi-environment test desteği

## 📝 Örnek Kullanım

### 1. Kullanıcı Kaydı
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "hakan",
    "email": "hakan@example.com",
    "password": "1234"
  }'
```

### 2. Giriş Yapma
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "hakan@example.com",
    "password": "1234"
  }'
```

### 3. Görev Oluşturma
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Proje Tamamla",
    "description": "Go Task Management API projesini bitir",
    "status": "pending",
    "priority": "high"
  }'
```

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

## 🚦 Development Workflow

### 1. Feature Development
```bash
# Feature branch oluştur
git checkout -b feature/new-feature




# Commit ve push
git commit -m "feat: new feature description"
git push origin feature/new-feature
```

### 2. Testing
```bash

# Contract testleri
go test ./test/contract -v -timeout=30s
```

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
