# System Architecture: Honda Leasing API

## 1. Pola Arsitektur yang Diterapkan

Proyek Honda Leasing API menggunakan pola **Layered Architecture** dengan sedikit sentuhan **Clean Architecture**. Sistem difokuskan untuk memecah secara jelas tanggung jawab dari masing-masing lapisan, sehingga memudahkan *testing*, *maintenance*, dan meminimalisir saling silang dependensi antar komponen (Low Coupling, High Cohesion). Pemisahan fungsional juga ditata dengan gaya **Modular Monolith**, di mana fitur dipisah per *domain/bussiness unit*, seperti `auth`, `catalog`, `leasing`, `delivery`, dan `officer`.

Setiap domain utamanya terbagi atas tiga lapisan inti:
1. **Handler Layer (Presentation):**
   - Bertugas menerima request HTTP dari pengguna (berupa JSON, Path Params, atau Query Params).
   - Memvalidasi bentuk data input menggunakan library `go-playground/validator`.
   - Mengubah struktur data input ke struktur yang dikenali lapisan Service, kemudian memanggil Service.
   - Mengembalikan response HTTP (berupa JSON atau error handling seragam) ke *client*.
2. **Service Layer (Business Logic):**
   - Jantung dari setiap modul. Di sini seluruh perhitungan, pengecekan, atau aturan bisnis didefinisikan secara khusus (e.g. validasi apakah customer dilarang mengajukan leasing dua kali, dsb.).
   - Service tidak peduli apakah data datang dari HTTP, gRPC, maupun *event broker*. 
   - Memanggil fungsionalitas di lapisan bawahnya yaitu Repository.
3. **Repository Layer (Data Access):**
   - Menangani perihal koneksi dan eksekusi query ke dalam bentuk persisten (Database PostgreSQL).
   - Seluruh detail teknis terkait ORM GORM berada di lapisan ini. Service tidak perlu tahu query SQL apa yang dieksekusi, ia cukup menerima atau mengirim data utuh (murni struct Go).

## 2. Diagram Alur (Data Flow)

Berikut adalah diagram alur secara keseluruhan yang merepresentasikan pola Request-Response sederhana dalam sistem aplikasi ini:

```mermaid
sequenceDiagram
    participant Client
    participant Router (Gin)
    participant Middleware (JWT/RBAC)
    participant Handler
    participant Service
    participant Repository
    participant PostgreSQL
    
    Client->>Router (Gin): HTTP POST /api/v1/leasing (Token)
    Router (Gin)->>Middleware (JWT/RBAC): Intercept & Validasi Token
    
    alt Invalid Token / Unauthorized
        Middleware (JWT/RBAC)-->>Client: 401/403 Error JSON
    else Token Valid & Allowed
        Middleware (JWT/RBAC)->>Handler: Teruskan Request
    end
    
    Handler->>Handler: Parse JSON & Validate Payload
    Handler->>Service: Call Create() via Interface
    
    Service->>Service: Terapkan Aturan Bisnis 
    Service->>Repository: Call Save() via Interface
    
    Repository->>PostgreSQL: Execute INSERT Query (GORM)
    PostgreSQL-->>Repository: Return Row / ID
    
    Repository-->>Service: Return Entity Struct
    Service-->>Handler: Return Result
    Handler-->>Client: HTTP 201 Created (JSON Response)
```

## 3. Komponen-Komponen Sistem (Component Relations)

Terdapat berbagai macam komponen dalam sistem ini yang dirakit menjadi satu kesatuan fungsional menggunakan **Google Wire**:

```mermaid
graph TD
    A["<b>HTTP Server / Router</b><br/>(Goroutine - Gin Framework)"] --> B["<b>Middleware Config</b>"]
    A --> C["<b>Handlers</b>"]
    
    B --> B1[Logging]
    B --> B2[Recovery]
    B --> B3["Auth & RBAC Validator"]

    C --> |Depends on via Interface| D["<b>Services</b>"]
    
    D --> E1[Auth Service]
    D --> E2[Catalog Service]
    D --> E3[Leasing Service]
    D --> E4[Officer Service]
    D --> E5[Delivery Service]
    
    E1 --> |Depends on via Interface| F["<b>Repositories</b><br/>(PostgreSQL Implementations)"]
    E2 --> F
    E3 --> F
    E4 --> F
    E5 --> F
    
    F --> G[("<b>PostgreSQL Database</b>")]
    
    style A fill:#4dabf7,stroke:#1971c2,stroke-width:2px,color:white
    style D fill:#69db7c,stroke:#2b8a3e,stroke-width:2px,color:white
    style F fill:#ffa94d,stroke:#e8590c,stroke-width:2px,color:white
    style G fill:#ced4da,stroke:#495057,stroke-width:2px
```

**Penjelasan Hubungan:**
1. **API Router (`gin-gonic/gin`)** akan menugaskan setiap jalur endpoint kepada `Handler` yang berkaitan. Jalur (Routes) diiringi dengan filter khusus dari `Middleware`, seperti blokade akses jika *role* tidak sesuai dengan peruntukan rute (RBAC).
2. **Handlers** tidak membuat sendiri instansi `Service`, hal ini diinjeksi (*Dependency Injection*) saat fase kompilasi oleh **Google Wire**. Handler hanya berinteraksi dengan Service menggunakan kontrak *Interface*.
3. **Services** melakukan perhitungan inti atau menyusun data lintas entitas, dan memanggil fungsi interface dari `Repository`.
4. Secara keseluruhan, ketergantungan mengarah "ke dalam", lapisan terdalam (Service dan Repository Interface) terisolasi dari detail framework (seperti representasi context dari Gin).
