# Technical Details & Code: Honda Leasing API

Dokumen ini menguraikan detail teknis pada level kode sumber (*source code*), pola desain (*design patterns*) yang diterapkan, serta representasi alur data.

---

## 1. Modul dan Komponen Utama

Secara garis besar, aplikasi ini memiliki 6 modul fungsional di dalam folder `internal/` yang semuanya saling independen, kecuali ketika dideklarasikan keanggotaannya melalui *Interface*.

| Nama Modul | Tanggung Jawab |
| :--- | :--- |
| `auth` | Mengurus registrasi *customer*, *login* seluruh peran pengguna, dan pembuatan token JWT. |
| `catalog` | Mengelola data Master Product (informasi unit motor, harga, stok). |
| `leasing` | Menerima pengajuan contract kredit motor oleh customer, menghitung tenor, memeriksa riwayat pesanan (My Contracts). |
| `master` | Mengelola referensi hierarki Data Wilayah (Provinsi hingga Kelurahan). |
| `finance` | Memproses kalkulasi jadwal angsuran / Late Fees. |
| `officer` | Memeriksa pesanan (Incoming Contracts), memanipulasi progress dan transisi *task status* sesuai *sequence level*. |

### 1.1 Contoh Implementasi Layer: Modul `Leasing`

Pola kode secara seragam dieksekusi melalui **Handler \-\> Service \-\> Repository**. 
Sebagai contoh, inilah potongan kode pada layer Handler untuk fungsi pembuatan pesanan kredit (Contract):

```go
// File: internal/leasing/handler/http.go

// LeasingHandler merupakan representasi level delivery API
type LeasingHandler struct {
	service leasing.Service // Injeksikan interface service ke struct ini
}

func NewLeasingHandler(service leasing.Service) *LeasingHandler {
	return &LeasingHandler{service: service}
}

func (h *LeasingHandler) SubmitContract(c *gin.Context) {
	// 1. Ambil userID dari JWT yang sudah di-set oleh middleware Auth
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(http.StatusUnauthorized, "Unauthorized"))
		return
	}
    
	userID := userIDVal.(int64)

	// 2. Bind Body JSON ke Data Transfer Object (DTO)
	var req SubmitContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(http.StatusBadRequest, err.Error()))
		return
	}

	// 3. Mapping DTO menuju domain model untuk dikirim ke Layer Service
	input := leasing.SubmitContractInput{
		UserID:         userID,
		MotorID:        req.MotorID,
		ProductID:      req.ProductID,
		NilaiKendaraan: req.NilaiKendaraan,
		DpDibayar:      req.DpDibayar,
		TenorBulan:     req.TenorBulan,
	}

	// 4. Eksekusi Business Logic di Service
	cont, err := h.service.SubmitContract(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(http.StatusInternalServerError, err.Error()))
		return
	}

	// 5. Kembalikan Response seragam ke Client
	contractResp := toContractResponse(*cont)
	c.JSON(http.StatusCreated, response.Success(http.StatusCreated, "Contract submitted successfully", contractResp))
}
```

---

## 2. Design Pattern yang Diterapkan

Sistem API pada Honda Leasing ini menerapkan banyak pola arsitektur perangkat lunak untuk memastikan kualitas dan skalabilitas:

1. **Repository Pattern**
   - **Tujuan**: Memisahkan logika kueri basis data (seperti `GORM` operations) dari aturan bisnis (Business Logic). 
   - **Cara Kerja**: Service tidak mengakses DB *context* secara langsung. Ia akan mendefinisikan *interface* kumpulan method yang dibutuhkannya (Misal: `FindByID`, `InsertContract`). *Repository* bertugas menyiapkan implementasi nyata menuju PostgreSQL.

2. **Dependency Injection (Compile-Time DI dengan Google Wire)**
   - **Tujuan**: Melonggarkan perangkapan antar komponen (Decoupling) dan memudahkan dalam *Unit Testing* menggunakan **Mocking**.
   - **Cara Kerja**: Alih-alih melakukan `service := NewService(NewRepository(NewDB()))` di `main.go`, semua struktur saling keterkaitan dikerjakan secara reaktif oleh fungsi generator `wire_gen.go`.

3. **Data Transfer Object (DTO) & Mapping**
   - **Tujuan**: Mencegah kebocoran data (*Over-posting* data) antara request body yang masuk, entitas database murni, dan respon JSON yang keluar.
   - **Cara Kerja**: Aplikasi memiliki *struct-struct* spesifik seperti `SubmitContractRequest` (input HTTP) dan `ContractResponse` (output HTTP). *Struct* ini terpisah dari Model Database asli (seperti `contract.Contract` entity).

4. **Factory Pattern / Construction Pattern**
   - Dapat dengan jelas dilihat dari kehadiran fungsi pembangkit awal secara seragam di seluruh repository, misalnya `func NewLeasingHandler(...) *LeasingHandler`, `NewService(...)`, dan seterusnya.

---

## 3. Detail Alur Data (Data Flow Lifecycle)

### Alur Input Parameter (Menerima Request HTTP)
1. Permintaan mendarat di Gin Framework HTTP (`Router`).
2. Proses melewati **Middleware**: Pengecekan standar JWT Header `Authorization: Bearer <token>`. Jika menggunakan sistem Role Based Access Control (RBAC), Middleware segera menolak bila rute yang memanggil bukan peran yang sesuai (Misal: rute `/officer/tasks` dipanggil dengan token milik *Customer*).
3. Payload dikonversi (`BindJSON`) ke struct spesifik (DTO) di **Handler**. Error tipe data atau validasi *mandatory fields* terjadi pada langkah ini.

### Alur Proses (Business Logic execution)
4. Handler mentransfer data dan *User ID* (terekstraksi dari token) turun ke layer **Service**.
5. Service memvalidasi logis terkait bisnis aturan Honda Leasing. Contoh: 
   - Apakah harga rasio cicilan masuk akal? 
   - Apakah pengguna sudah memiliki kontrak yang sedang aktif/berjalan?
6. Jika perlu manipulasi data atau mendapatkan data unit, Service meminta data persisten kepada **Repository** (melalui Interface kontrak).

### Alur Respon (Database ke JSON)
7. Repository menerjemahkan parameter menjadi query SQL, lalu **GORM DB Driver** mengambil *row* database. Mapping manual terjadi untuk mengubah row dari DB ke Entity Domain murni `struct`.
8. Struct dikembalikan *(bubble-up)* hingga Handler.
9. Di lapisan Handler, *Domain Entity* direstrukturisasi ulang memakai Response DTO (`toContractResponse()`) yang membungkus output agar tidak mengekspos field sensitif (seperti _Deleted At_ atau _Updated At_).
10. Handler membungkus respons menggunakan standarisasi struktur generik utilitas dari paket `pkg/response`. Format JSON output konsisten:
```json
{
  "code": 201,
  "message": "Contract submitted successfully",
  "data": {
    "contract_id": 14002,
    "status": "pending_approval"
  }
}
```
