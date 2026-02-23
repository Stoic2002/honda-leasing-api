# Honda Leasing API

Honda Leasing API adalah layanan backend (RESTful API) yang dirancang untuk menangani seluruh *lifecycle* kredit sepeda motor Honda. Mulai dari melihat katalog, pengajuan pinjaman oleh *customer*, proses verifikasi oleh pihak *leasing*, hingga sistem checklist pengiriman unit oleh divisi *delivery*.

Proyek ini mengadopsi struktur ruang lingkup **Modular Monolith** dengan perpaduan arsitektur berlapis (Handler, Service, Repository) serta telah dibekali dengan **Google Wire** untuk Dependency Injection otomatis yang bersih.

## 🛠 Tech Stack
- **Bahasa Pemrograman**: Go 1.24+
- **Framework Web**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **Database & ORM**: PostgreSQL & GORM
- **Authentication**: JWT (JSON Web Token)
- **Dependency Injection**: [Google Wire](https://github.com/google/wire)
- **Dokumentasi API**: Swagger OpenAPI
- **Lain-lain**: Go Playground Validator, Viper (untuk file `.yaml`), Bcrypt.

---

## 🚀 Pre-requisites (Prasyarat Instalasi)
Pastikan Anda telah memasang beberapa hal berikut pada sistem Anda:
1. **[Go](https://go.dev/)** (minimal versi 1.24)
2. **PostgreSQL** lokal untuk Database, jalankan dan buat sebuah database kosong (misal: `honda_leasing_db`).
3. (Opsional namun disarankan) **`golangci-lint`** untuk linter.

---

## ⚙️ Persiapan Basis Data (Setup Configuration)

Meskipun contoh file config (`app.example.yaml` / `app.dev.yaml`) sudah ada, pastikan Anda menyesuaikan kredensial koneksinya:

1. Pergi ke folder `configs/`.
2. Duplikat atau periksa struktur `app.dev.yaml`:

```yaml
app:
  port: "8080"
  env: "development"

database:
  host: "localhost"
  port: "5432"
  user: "postgres"        # Sesuaikan user postgres Anda
  password: "password123" # Sesuaikan password
  name: "honda_leasing_db"# Sesuaikan nama Database Anda
  sslmode: "disable"

jwt:
  secret: "YOUR_SUPER_SECRET_KEY"
  expire_minutes: 60
  refresh_days: 7
```

---

## 🏃 Cara Menjalankan Project (How to Run)

Karena project ini telah difasilitasi dengan `Makefile`, proses pengelolaan, build, dan *running* menjadi sangat mudah.

### 1. Menjalankan Mode Development
Perintah ini akan menyalakan server di port `8080` dan akan menggunakan parameter yang ada pada `configs/app.dev.yaml`.

```bash
make run-dev
```

### 2. Membangkitkan Data Dummy (Seeder)
Jika Anda membutuhkan data pengguna (Admin/Officer, Delivery, dan Customer dummy) agar mudah diloginkan, jalankan *seeder db*. Pastikan struktur tabel Anda masih sinkron.

```bash
make seed-dev
```

### 3. Build & Run Binary (Untuk Production/Staging)
Jika ingin melakukan kompilasi program menjadi binary executable di folder `bin/honda-leasing-api`:

```bash
# Melakukan build kompilasi
make build

# Menjalankan spesifik environment
make run-staging
# atau
make run-prod
```

### 4. Code Generation (Google Wire)
Jika Anda memodifikasi parameter atau *dependencies* baru pada fitur (contoh menambahkan `NewService()` pada `internal/...`), Anda harus memperbarui berkas *injection*:

```bash
make wire
```
Jika Anda belum menginstall *Wire* tool secara global, jalankan dulu `go install github.com/google/wire/cmd/wire@latest`.

---

## 📚 Dokumentasi Endpoint (Swagger UI)

Setelah server API Anda berjalan (*running*), dokumentasi API yang interaktif (Swagger UI) dapat segera diakses via browser pada rute:
👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

---
