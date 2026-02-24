# API Reference

This document provides a comprehensive list of all available REST API endpoints in the Honda Leasing API, categorized by domain module. It includes the required request payloads and the expected response structures based on the internal Data Transfer Objects (DTOs).

## Table of Contents
1. [Authentication Module](#1-authentication-module)
2. [Master Data Module](#2-master-data-module)
3. [Catalog Module](#3-catalog-module)
4. [Leasing (Customer) Module](#4-leasing-customer-module)
5. [Officer Module](#5-officer-module)
6. [Finance Module](#6-finance-module)

---

## 1. Authentication Module
Base Path: `/api/v1`

### 1.1. Login
- **Endpoint**: `POST /auth/login`
- **Access**: Public
- **Description**: Authenticate user and receive JWT tokens.
- **Request Body**:
  ```json
  {
    "email": "string (Required)",
    "password": "string (Required)"
  }
  ```
- **Response** (Success):
  ```json
  {
    "access_token": "string",
    "refresh_token": "string",
    "role": "string"
  }
  ```

### 1.2. Refresh Token
- **Endpoint**: `POST /auth/refresh`
- **Access**: Public (Requires valid Refresh Token)
- **Request Body**:
  ```json
  {
    "refresh_token": "string (Required)"
  }
  ```
- **Response** (Success): _Returns new access and refresh tokens identical to Login response._

### 1.3. Get Profile
- **Endpoint**: `GET /user/me`
- **Access**: Private (Requires standard Bearer `access_token`)
- **Response** (Success):
  ```json
  {
    "user_id": 1,
    "email": "user@example.com",
    "full_name": "string",
    "role": "CUSTOMER"
  }
  ```

---

## 2. Master Data Module
Base Path: `/api/v1/master`  
**Access:** Public (Used for Frontend Dropdowns)

### 2.1. Get Provinces
- **Endpoint**: `GET /provinces`
- **Response**: Array of:
  ```json
  {
    "prov_id": 1,
    "prov_name": "Jawa Barat"
  }
  ```

### 2.2. Get Kabupatens
- **Endpoint**: `GET /kabupatens?prov_id=X`
- **Response**: Array of:
  ```json
  {
    "kab_id": 1,
    "kab_name": "Kota Bandung"
  }
  ```

### 2.3. Get Kecamatans
- **Endpoint**: `GET /kecamatans?kab_id=X`
- **Response**: Array of:
  ```json
  {
    "kec_id": 1,
    "kec_name": "Andir"
  }
  ```

### 2.4. Get Kelurahans
- **Endpoint**: `GET /kelurahans?kec_id=X`
- **Response**: Array of:
  ```json
  {
    "kel_id": 1,
    "kel_name": "Campaka"
  }
  ```

---

## 3. Catalog Module
Base Path: `/api/v1/catalog`  
**Access:** Private (Authenticated Users)

### 3.1. Get All Motors
- **Endpoint**: `GET /motors`
- **Response**: Array of `MotorResponse`.

### 3.2. Get Motor By ID
- **Endpoint**: `GET /motors/:id`
- **Response**:
  ```json
  {
    "motor_id": 1,
    "merk": "Honda",
    "tahun": 2026,
    "warna": "Hitam",
    "nomor_rangka": "MH1...",
    "nomor_mesin": "JF1...",
    "cc_mesin": "156.9 cc",
    "nomor_polisi": "D 1234 ABC",
    "status_unit": "ready",
    "harga_otr": 33400000,
    "motor_type": {
      "moty_id": 1,
      "moty_name": "Maxi"
    },
    "assets": [
      {
        "moas_id": 1,
        "file_name": "image.jpg",
        "file_size": 1024,
        "file_type": "jpg",
        "file_url": "https://..."
      }
    ],
    "created_at": "2026-01-01T00:00:00Z"
  }
  ```

### 3.3. Get Leasing Products
- **Endpoint**: `GET /leasing-products`
- **Response**: Array of:
  ```json
  {
    "product_id": 1,
    "kode_produk": "DP-RINGAN",
    "nama_produk": "DP Ringan 24 Bulan",
    "tenor_bulan": 24,
    "dp_persen_min": 10.0,
    "dp_persen_max": 20.0,
    "bunga_flat": 1.2,
    "admin_fee": 350000,
    "asuransi": true
  }
  ```

---

## 4. Leasing (Customer) Module
Base Path: `/api/v1/customer`  
**Access:** Private (`CUSTOMER` Role Only)

### 4.1. Submit Order
- **Endpoint**: `POST /orders`
- **Request Body**:
  ```json
  {
    "motor_id": 1,
    "product_id": 1,
    "nilai_kendaraan": 33400000,
    "dp_dibayar": 5000000,
    "tenor_bulan": 24
  }
  ```
- **Response**:
  ```json
  {
    "contract_id": 1,
    "contract_number": "KTR-2026-001",
    "request_date": "2026-01-01T00:00:00Z",
    "tenor_bulan": 24,
    "nilai_kendaraan": 33400000,
    "dp_dibayar": 5000000,
    "pokok_pinjaman": 28400000,
    "total_pinjaman": 36579200,
    "cicilan_per_bulan": 1524133,
    "status": "draft",
    "customer_id": 1,
    "motor_id": 1,
    "product_id": 1,
    "created_at": "2026-01-01T00:00:00Z"
  }
  ```

### 4.2. Get My Orders
- **Endpoint**: `GET /orders`
- **Response**: Array of My Orders Brief.

### 4.3. Get Contract Progress (Tasks tracking)
- **Endpoint**: `GET /orders/:id/progress`
- **Response**: Array of:
  ```json
  {
    "task_id": 1,
    "task_name": "Input Pengajuan & Unggah Dokumen",
    "status": "completed",
    "sequence_no": 1,
    "actual_startdate": "2026-01-01T00:00:00Z",
    "actual_enddate": "2026-01-02T00:00:00Z",
    "created_at": "2026-01-01T00:00:00Z"
  }
  ```

---

## 5. Officer Module
Base Path: `/api/v1/officer`  
**Access:** Private (Roles: `ADMIN_CABANG`, `SALES`, `SURVEYOR`, `FINANCE`, `COLLECTION`)

### 5.1. Get Incoming Orders (Leasing Requests)
- **Endpoint**: `GET /orders?page=1&limit=10`
- **Response**: Array of `IncomingOrderResponse` (paginated).

### 5.2. Get My Tasks (Assigned pipeline steps)
- **Endpoint**: `GET /tasks?page=1&limit=10`
- **Response**: Array of:
  ```json
  {
    "task_id": 1,
    "task_name": "Survei Lapangan",
    "status": "pending",
    "sequence_no": 3,
    "contract_id": 1,
    "created_at": "2026-01-01T00:00:00Z"
  }
  ```

### 5.3. Process Task (Move task to next stage)
- **Endpoint**: `POST /tasks/:taskId/process`
- **Request Format**: `multipart/form-data`
- **Fields**:
  - `notes` (string): Text notes outlining findings or processing results.
  - `attributes[<Name>]`: Dynamic properties required by the stage. Can be text (e.g. `attributes[Penghasilan Bulanan]="10000000"`) or `*File Uploads*` (e.g. `attributes[Foto KTP]=<file>`).
- **Response**: `{ message: "Task successfully processed..." }`

---

## 6. Finance Module
Base Path: `/api/v1/finance`

### 6.1. Get Payment Schedules
- **Endpoint**: `GET /schedules?contract_id=1`
- **Access**: Private (Role: `FINANCE` or `CUSTOMER` who owns the contract)
- **Description**: Returns the scheduled amortizations. Dynamic Late Fee is generated on-the-fly depending on the current date relative to `jatuh_tempo` (due date).
- **Response**: Array of:
  ```json
  {
    "schedule_id": 1,
    "angsuran_ke": 1,
    "jatuh_tempo": "2026-02-01T00:00:00Z",
    "pokok": 664062,
    "margin": 115104,
    "late_fee": 0,
    "total_tagihan": 779166,
    "status_pembayaran": "unpaid",
    "tanggal_bayar": null
  }
  ```

### 6.2. Webhook: Process Payment
- **Endpoint**: `POST /payments/webhook`
- **Access**: Configured API Key / Integration Layer (Publicly route-able for external Gateway)
- **Request Body**:
  ```json
  {
    "nomor_bukti": "INV-12345",
    "jumlah_bayar": 779166,
    "metode_pembayaran": "Bank Transfer",
    "provider": "BCA VA",
    "contract_id": 1,
    "schedule_id": 1
  }
  ```
- **Response**: `{ message: "Payment webhook processed successfully" }`
