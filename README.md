# ERM Dokter (Electronic Medical Record / Rekam Medis Elektronik)

Backend REST API service untuk modul **Dokter (Rekam Medis Elektronik / ERM)** pada Sistem Informasi Manajemen Rumah Sakit (SIMRS Khanza). Dibangun menggunakan **Go (Golang)** dengan prinsip Clean Architecture, keamanan tingkat lanjut (*URL Obfuscation with AES-GCM*), autentikasi JWT, validasi klinis medis, dan dokumentasi interaktif **OpenAPI 3.1 & Scalar API Reference**.

---

## 🌟 Fitur Utama

- **🔐 Autentikasi & Otorisasi**:
  - Login dokter menggunakan ID/Username Khanza dengan password terenkripsi.
  - JWT Bearer Authentication dengan masa berlaku sesi yang aman.
  - Proteksi hak akses rekam medis: Dokter hanya dapat mengubah/menghapus rekam medis miliknya sendiri.
- **🏥 Rawat Jalan & Antrean Pasien**:
  - Referensi opsi filter antrean (Status Pemeriksaan, Penjamin/Asuransi, Status Lanjut, Jenis Antrean).
  - Daftar antrean pasien per dokter dengan dukungan paginasi (`page`, `limit`), pencarian bebas (`keyword`), rentang `tanggal`, serta pengurutan (`order_by`, `sort_order`).
  - Detail data kunjungan pasien.
- **📝 Rekam Medis & Pemeriksaan Klinis (SOAP & TTV)**:
  - Referensi tingkat kesadaran valid (GCS & AVPU: *Compos Mentis, Somnolence, Sopor, Coma, Alert, Confusion, Voice, Pain, Unresponsive*).
  - Riwayat pemeriksaan berdasarkan nomor kunjungan/registrasi (`id_kunjungan`).
  - Riwayat rekam medis lengkap seorang pasien lintas seluruh kunjungan (`id_pasien`).
  - Detail satu catatan pemeriksaan (`id_pemeriksaan`).
  - Simpan catatan pemeriksaan baru (Rawat Jalan & Rawat Inap) dengan validasi klinis (suhu, tensi sistolik/diastolik, nadi, respirasi, SpO2, GCS) dan validasi waktu sebelum registrasi.
  - Update & Hapus catatan pemeriksaan medis (dengan batas waktu proteksi maksimal 48 jam sejak pemeriksaan dibuat).
- **🛡️ Keamanan & Performa**:
  - **URL ID Obfuscation**: Mengenkripsi `no_rawat`, `no_rkm_medis`, dan *composite key* pemeriksaan pada URL path menggunakan algoritma **AES-256-GCM**.
  - **Rate Limiting Middleware**: Perlindungan terhadap brute-force login (maksimal 10 request/menit).
  - **Context Timeout Middleware**: Timeout otomatis 10 detik untuk menjaga ketahanan server.
  - **CORS Middleware**: Konfigurasi origin dan header yang fleksibel untuk integrasi frontend.
- **📖 Dokumentasi API Interaktif**:
  - Embedded **Scalar API Reference UI** bertema modern di endpoint `GET /docs`.
  - Spesifikasi OpenAPI 3.1.0 di endpoint `GET /docs/openapi.yaml`.

---

## 📁 Struktur Folder Project

```text
erm-dokter/
├── cmd/
│   └── api/
│       └── main.go               # Application entry point
├── internal/
│   ├── auth/                     # Domain autentikasi dokter & JWT
│   ├── config/                   # Konfigurasi aplikasi & load .env
│   ├── di/                       # Dependency Injection provider container
│   ├── docs/                     # Handler Scalar UI & OpenAPI 3.1 YAML spec
│   ├── health/                   # Health check endpoint
│   ├── master/                   # Master data (penjamin, depo, obat) & opsi referensi
│   ├── middleware/               # Auth JWT, CORS, Rate Limiter, Timeout
│   ├── pasien/                   # Domain & model data pasien
│   ├── pemeriksaan/              # Domain Rekam Medis Klinis (SOAP & TTV)
│   ├── pkg/                      # Shared utility packages:
│   │   ├── crypto/               # AES-GCM encryption/decryption helper
│   │   ├── database/             # Database connection & pooling
│   │   ├── logger/               # Structured logging
│   │   ├── response/             # Standardized JSON response helper
│   │   └── token/                # JWT claim generator & validator
│   ├── rawatjalan/               # Domain antrean & kunjungan rawat jalan
│   ├── resep/                    # Domain resep obat dokter (e-resep)
│   ├── routes/                   # HTTP Route Registry (Go 1.22+ ServeMux)
│   └── shared/                   # Shared types, PaginationMeta, AppError
├── test.http                     # File pengujian REST Client (VS Code)
├── .air.toml                     # Konfigurasi live-reload (Air)
├── .env.example                  # Template environment variables
├── Makefile                      # Shortcut command automasi
└── go.mod                        # Go module dependencies
```

---

## ⚙️ Persyaratan Sistem

- **Go (Golang)**: Versi 1.22 atau lebih baru
- **Database**: MySQL 5.7+ / 8.0+ atau MariaDB 10.3+ (Skema Database SIMRS Khanza)
- **Tooling Opsional**: [Air](https://github.com/air-verse/air) untuk live reload saat pengembangan.

---

## 🚀 Panduan Instalasi & Menjalankan

### 1. Clone Repository & Setup Environment
```bash
git clone https://github.com/khoirir/erm-dokter.git
cd erm-dokter
cp .env.example .env
```

Sesuaikan konfigurasi koneksi database dan kunci enkripsi pada file `.env`:
```env
APP_ENV=development
APP_PORT=8082
APP_NAME=erm-dokter

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=sik

JWT_SECRET=your_super_secret_jwt_key_here
USER_KEY=your_aes_user_key_here
PASSWORD_KEY=your_aes_password_key_here
ENCRYPTION_KEY=your_aes_encryption_key_32_characters_here

CORS_ORIGIN=*
MAX_EDIT_REKAM_MEDIS_JAM=48
```

### 2. Download Dependensi
```bash
go mod tidy
```

### 3. Menjalankan Aplikasi

- **Menggunakan Live Reload (Air)**:
  ```bash
  air
  ```
- **Menggunakan Go CLI Standar**:
  ```bash
  go run cmd/api/main.go
  ```
- **Menggunakan Makefile**:
  ```bash
  make run
  ```

Server akan aktif secara default di `http://localhost:8082`.

---

## 🧪 Pengujian (Testing)

### Menjalankan Automated Unit Tests:
```bash
go test ./... -v
```

### Menggunakan VS Code REST Client (`test.http`):
File [`test.http`](file:///d:/PROGRAMMING/DEVELOP_GO/erm-dokter/test.http) berisi skenario pengujian lengkap dari Auth, Rawat Jalan, Master Penjamin, hingga Rekam Medis Pemeriksaan (Simpan, Update, Hapus, dan Validasi Error 400/403/404).

---

## 📑 Dokumentasi API (Scalar UI)

Dokumentasi API interaktif sudah terpasang secara bawaan. Buka browser dan akses:

- **Scalar API Reference UI**: [http://localhost:8082/docs](http://localhost:8082/docs)
- **Raw OpenAPI 3.1 YAML**: [http://localhost:8082/docs/openapi.yaml](http://localhost:8082/docs/openapi.yaml)

### Konvensi REST API:
- **Parameter Kontrol (English)**:
  - `page`: Nomor halaman (contoh: `?page=1`)
  - `limit`: Jumlah data per halaman (contoh: `?limit=20`)
  - `keyword`: Pencarian nama pasien / no RM / no rawat
  - `order_by`: Kolom pengurutan (`waktu_registrasi`, `nama_pasien`)
  - `sort_order`: Arah pengurutan (`ASC`, `DESC`)
- **Filter Domain Database (Bahasa Indonesia)**:
  - `tanggal`: Rentang tanggal registrasi/perawatan (`YYYY-MM-DD,YYYY-MM-DD`)
  - `kode_penjamin`: Kode asuransi/penjamin (misal: `BPJ`)
  - `status_pemeriksaan`: `Belum`, `Sudah`, `Batal`, `Berkas Diterima`, `Dirujuk`, `Meninggal`, `Dirawat`, `Pulang Paksa`
  - `status_bayar`: `Sudah Bayar`, `Belum Bayar`
  - `status_lanjut`: `Ralan`, `Ranap`, `Semua`
  - `jenis_antrean`: `Rujukan`, `Bukan Rujukan`
- **Format Response Paginasi**:
  ```json
  {
    "success": true,
    "message": "Berhasil mengambil daftar antrean dokter",
    "data": [...],
    "meta": {
      "total_records": 45,
      "total_pages": 3,
      "current_page": 1,
      "per_page": 20
    }
  }
  ```

---

## 📄 Lisensi

Proyek ini dikembangkan untuk kebutuhan integrasi modul Dokter Rekam Medis Elektronik pada SIMRS.
