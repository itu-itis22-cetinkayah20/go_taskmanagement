# Go Task Management API

Bu proje, JWT tabanlı kimlik doğrulama ile kullanıcıların görev ekleyip yönetebileceği basit bir Go REST API örneğidir.

## Başlatma

1. Gerekli bağımlılıkları yükleyin:
   ```
   go mod tidy
   ```
2. Sunucuyu başlatın:
   ```
   go run main.go
   ```
   Sunucu `localhost:8080` adresinde çalışacaktır.

## API Endpointleri

- `POST /register` — Kullanıcı kaydı
- `POST /login` — Giriş ve JWT token alma
- `GET /tasks/public` — Herkesin görebileceği örnek görevler
- `GET /tasks` — Kullanıcının kendi görevleri (JWT ile)
- `POST /tasks` — Yeni görev ekleme (JWT ile)
- `GET /tasks/{id}` — Görev detayları (JWT ile)
- `PUT /tasks/{id}` — Görev güncelleme (JWT ile)
- `DELETE /tasks/{id}` — Görev silme (JWT ile)
- `POST /logout` — Çıkış (JWT ile)

## Testleri Çalıştırma

1. Sunucunun açık olduğundan emin olun (`go run main.go`).
2. Ayrı bir terminalde testleri başlatın:
   ```
   go test -v ./tests
   ```
   Her bir test başarılı olduğunda terminalde bilgilendirici mesajlar göreceksiniz.

## Örnek Test Akışı
- Kayıt
- Giriş
- Public görevleri listeleme
- Görev ekleme
- Görevleri listeleme
- Görev detay görüntüleme
- Görev güncelleme
- Görev silme
- Çıkış

## Notlar
- Tüm veriler RAM'de tutulur, sunucu yeniden başlatılırsa veriler silinir.
- JWT gerektiren endpointlere erişmek için önce `/login` ile token alınmalı ve isteklerde `Authorization: Bearer <token>` header'ı kullanılmalıdır.
