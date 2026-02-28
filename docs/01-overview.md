# Project Overview: Honda Leasing API

## 1. Tujuan dan Fungsi Utama Project

**Honda Leasing API** adalah sebuah layanan backend berbasis RESTful API yang dirancang secara khusus untuk mengelola seluruh siklus hidup (lifecycle) kredit atau leasing sepeda motor Honda. Sistem ini memfasilitasi berbagai proses bisnis mulai dari pendaftaran dan login, melihat katalog motor, hingga proses kompleks seperti pengajuan kredit dan pengiriman unit.

Aplikasi ini melayani dua jenis atau **Role** pengguna utama:
1. **Customer**: Pengguna yang dapat melihat katalog motor dan melakukan pengajuan kredit (contract).
2. **Leasing Officer (Admin)**: Pengguna internal yang bertugas melakukan verifikasi, memproses approval (persetujuan), hingga mengelola data pengiriman unit motor dari customer secara dinamis berdasar _Sequence Task_.

Fungsi utama dari sistem meliputi:
- Otentikasi dan Otorisasi menggunakan JWT (JSON Web Token).
- Manajemen data *customer* dan proses contract/pengajuan kredit.
- Sistem tracking status urutan (Sequence) *leasing tasks* untuk segala operasi internal (Admin, Finance, Surveyor, Delivery).

---

## 2. Tech Stack yang Digunakan

Proyek ini dibangun menggunakan teknologi modern yang berfokus pada performa, skalabilitas, dan kemudahan *maintenance*:

| Kategori | Teknologi/Library Utama | Kegunaan |
| :--- | :--- | :--- |
| **Bahasa Pemrograman** | [Go (Golang) 1.24+](https://go.dev/) | Bahasa utama dengan performa tinggi untuk membangun sistem backend. |
| **Web Framework** | [Gin Gonic (v1.11.0)](https://github.com/gin-gonic/gin) | Web framework yang ringan dan sangat cepat untuk Go. |
| **Database** | PostgreSQL | Relational Database Management System (RDBMS) yang digunakan sebagai basis data relasional. |
| **ORM** | [GORM (v1.31.1)](https://gorm.io/) | Object Relational Mapping untuk menjembatani kode Go dengan instruksi query SQL. |
| **Dependency Injection** | [Google Wire (v0.7.0)](https://github.com/google/wire) | Compiler tools dari Google untuk *Dependency Injection* (DI) secara otomatis saat *compile-time*. |
| **Authentication** | `golang-jwt/jwt/v5` | Implementasi standar JSON Web Token untuk session dan sekuritas endpoint. |
| **Validation** | `go-playground/validator/v10` | Library untuk memvalidasi *request payload* dari client secara otomatis. |
| **Configuration** | [Viper (v1.21.0)](https://github.com/spf13/viper) | Manajer konfigurasi untuk membaca file `.yaml` dan *environment variables*. |
| **Password Hashing** | `golang.org/x/crypto/bcrypt` | Fungsi hash satu arah yang sangat kuat untuk mengamankan data password pengguna. |

---

## 3. Struktur Folder dan Arsitektur High-Level

Proyek ini dipisahkan menjadi beberapa *layer* tanggung jawab yang dikenal dengan gaya **Modular Monolith** atau pendekatan mirip **Clean Architecture**. Pemisahan antar domain dan layer fungsional sangat kentara di folder `internal/`.

```text
honda-leasing-api/
├── Makefile                # Kumpulan perintah singkat (build, run, migrate, wire, seed)
├── README.md               # Dokumentasi awal / startup guide proyek
├── go.mod & go.sum         # Definisi module dan library pihak ketiga untuk Go
├── configs/                # Folder penyimpanan konfigurasi YAML (app.example.yaml, app.dev.yaml)
├── cmd/
│   ├── api/                # Titik berat (Entry point) utama. Berisi main.go untuk menjalankan server API
│   │   ├── main.go         # Bootstrapping (baca config, DB koneksi, mounting routing)
│   │   ├── wire.go         # Definisi awal provider untuk digenerate oleh Google Wire
│   │   └── wire_gen.go     # Hasil auto-generate dari Google Wire untuk Dependency Injection
│   └── seed/               # Script mandiri untuk memasukkan data awal dummy (seeder) ke Database
├── internal/               # Folder inti kode aplikasi (business logic), disegmentasi berdasarkan Domain/Modul
│   ├── auth/               # Modul terkait User Management dan Authentication (Login/Register)
│   ├── catalog/            # Modul terkait katalog produk sepeda motor (Data motor)
│   ├── finance/            # Modul terkait jadwal angsuran (Schedules) dan denda keterlambatan (Late Fees)
│   ├── leasing/            # Modul utama untuk proses pengajuan dan tracking dokumen leasing (Customer)
│   ├── master/             # Modul untuk hierarki master data kewilayahan (Provinsi -> Kelurahan)
│   ├── officer/            # Modul tunggal dinamis untuk urutan task admin leasing (Surveyor, Delivery, Approval)
│   ├── domain/             # (Opsional) Inti dari entities (struct definition, model DB) dari berbagai modul
│   ├── infrastructure/     # Kode yang bertugas berinteraksi dengan dunia luar: (Koneksi GORM Postgres, Redis, dll)
│   └── middleware/         # Gin middleware global (JWT Authorization check, Error Handler, Logger)
├── docs/                   # (Baru) Dokumentasi menyeluruh mengenai analisis teknikal sistem
├── pkg/                    # Package utility untuk helper logic yang sifatnya general (Hash, Response Formatter)
├── migrations/             # Berkas SQL atau file migrasi database (schema setup awal)
└── scripts/                # Utility script tambahan (Bash/Batch)
```

**Konsep Lapisan Aplikasi (Layered Design):**
Setiap modul di dalam `internal/` pada umumnya mengikuti pembagian 3 lapisan utama:
1. **Handler (Controller/Delivery Layer):** Lapisan terluar yang menangani input JSON/Form dari HTTP `request`, memvalidasinya dengan validator, memanggil fungsi `Service`, dan mengembalikan JSON respon.
2. **Service (Business Layer):** Tempat keberadaan logika bisnis yang sesungguhnya. Memastikan segala aturan (rules/constraint) untuk leasing atau aplikasi dapat terpenuhi.
3. **Repository (Data Access Layer):** Lapisan terbawah yang bertugas mengakses langsung dan memanipulasi entitas di dalam Database menggunakan `GORM`.

Semua komunikasi di dalam codebase tersebut dirajut (wiring) secara otomatis menggunakan *Google Wire*, sehingga mengurangi jumlah kode repetitif untuk melakukan *instantiation* objek satu per satu.
