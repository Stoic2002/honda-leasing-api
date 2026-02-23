# Configuration & Dependencies: Honda Leasing API

Dokumen ini menjelaskan pustaka/paket utama yang mendukung jalannya aplikasi, serta tata cara konfigurasi *environment* untuk menjalankan sistem ini baik di tahapan *Development* maupun *Production*.

---

## 1. Library / Package Utama

Aplikasi Honda Leasing API dibangun di atas ekosistem Go yang kuat dengan mengandalkan beberapa modul berlisensi *open-source* populer. Detail modul dapat ditemukan di dalam file `go.mod` proyek:

| Package / Library Utama | Kegunaan | Deskripsi Lebih Detail |
| :--- | :--- | :--- |
| **`github.com/gin-gonic/gin`** | Web Routing & Framework | Digunakan untuk mengatur segala *routing* HTTP. Sangat ringan namun memiliki performa ekstraksi URL parameter dan JSON *binding* yang cepat. |
| **`gorm.io/gorm`** | Object Relational Mapping | Library ORM Go yang populer. Mampu merepresentasikan skema tabel ke dalam sintaks Go murni dan menyederhanakan *query* SQL yang kompleks (*CRUD builder*). |
| **`gorm.io/driver/postgres`** | PostgreSQL Driver | Adapter GORM yang dispesifikasikan agar ORM dikenali dan kompatibel dengan tata cara sintaks dan konektivitas RDBMS PostgreSQL. |
| **`github.com/golang-jwt/jwt/v5`** | User Session Security | Modul untuk men-generate (Sign) serta memvalidasi token otentikasi (JWT) yang aman (*stateless authentication*). |
| **`github.com/spf13/viper`** | Configuration Management | Sangat efektif untuk memuat konfigurasi dari file konfigurasi seperti YAML, JSON, TOML form, atau variabel sistem. Membantu transisi antar *Environments* (Dev/Staging/Prod). |
| **`github.com/google/wire`** | Dependency Injection Tools | Bukanya sebuah modul yang di-*import*, melainkan sebuah alat *Compiler-time injection* yang men-*generate* ulang file dependensi terintegrasi via terminal `wire`. |
| **`github.com/go-playground/validator/v10`** | Request Input Validation | Menambahkan tag komprhensif ke properti *struct* di Go (contoh: `validate:"required,email"`) untuk sanitasi dan filter input dari *client* (mencegah *SQL Injection* logis). |
| **`golang.org/x/crypto/bcrypt`** | Keamanan Data | Implementasi standard utilitas Bcrypt untuk merancang password dengan proses sistem *salt* 1-arah pada algoritma *hash*. |

---

## 2. Environment Setup dan Konfigurasi Penting

Struktur konfigurasi berpusat di direktori `configs/` dan prosesnya dinamis sesuai nilai variabel yang disebut `APP_ENV`.

### A. Lokasi Konfigurasi
Sistem mendeteksi keberadaan file **YAML** menurut nama *environment* pada variabel OS (Secara default `dev` jika tidak disematkan). 
Maka format file konfigurasinya mengacu ke pola `configs/app.[APP_ENV].yaml`.

Contoh nama file yang bisa dibuat (dan dapat dipisahkan):
- `configs/app.dev.yaml`
- `configs/app.staging.yaml`
- `configs/app.prod.yaml`

### B. Struktur Konfigurasi (File YAML)
Berikut adalah konfigurasi vital yang harus diatur dalam sistem:

```yaml
# app.dev.yaml
app:
  port: "8080"           # Port internal HTTP Server
  env: "development"     # Mode eksekusi aplikasi

database:
  host: "localhost"      
  port: "5432"           
  user: "postgres"       
  password: "password123"# Pastikan berbeda untuk Production
  name: "honda_leasing_db"
  sslmode: "disable"     # 'disable' untuk Dev lokal, 'require' untuk AWS RDS / VPS

jwt:
  secret: "KunciRahasia_TandaTanganBackendBbe7102"
  expire_minutes: 60     # Jeda waktu (umur) Expiry Token JWT (60 menit)
  refresh_days: 7        # (Opsional) Apabila ada fitur Refresh Token
```

### C. Makefile Utilities
Daripada menulis skrip bash konfigurasi yang berulang kali di terminal, fungsionalitas otomasi dieksekusi melalui **Makefile** bawaan:

```bash
# Untuk menjalankan server API yang merujuk pada `configs/app.dev.yaml`
make run-dev

# Untuk membangun binary eksekusi file kompilasi akhir (untuk deploy)
make build

# Alat utilitas untuk men-generate ulang injection setelah modifikasi function
make wire

# Seeder: Utilitas yang menyuntikkan data tabel dummy (Role dan sample User default)
make seed-dev
```

### D. Setup Minimum Basis Data (Database)
Sebelum aplikasi dijalankan, sistem tidak akan menggunakan "Automigrate" yang merusak tabel lama, sehingga *prasyarat wajib* adalah pembuatan databasenya di lokal / instance PSQL:
1. Panggil `createdb honda_leasing_db -U postgres` via CLI
2. Jalankan migrasi SQL *(Jika ada di folder migrations, atau melalui trigger manual DB engine).*
3. Akses `make seed-dev` untuk memastikan akun *dummy* Leasing Officer dan Delivery sudah dibentuk.
