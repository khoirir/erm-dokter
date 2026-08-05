# ERM Dokter (Electronic Medical Record / Rekam Medis Elektronik)

Backend REST API service untuk modul **Dokter (Rekam Medis Elektronik)** dalam Sistem Informasi Manajemen Rumah Sakit (SIMRS).

## 📁 Struktur Folder Project (*Standard Go Project Layout*)

```text
erm-dokter/
├── cmd/                # Entrypoint aplikasi (main.go)
│   └── api/
│       └── main.go
├── internal/           # Private application code (Clean Architecture)
│   ├── config/         # Environment & App configuration
│   ├── domain/         # Domain Entities / Interfaces
│   ├── handler/        # HTTP Handlers / Controllers
│   ├── middleware/     # Custom HTTP Middlewares (Auth, CORS, Logger)
│   ├── repository/     # Data Access Layer (DB queries / ORM)
│   └── usecase/        # Core Business Logic
├── pkg/                # Public/Shared packages & utilities
│   ├── logger/         # Structured logger helper
│   └── response/       # Standardized JSON response helper
├── api/                # API Spec / OpenAPI / Swagger files
├── configs/            # Config templates & schema
├── scripts/            # Database migrations & build scripts
├── .env.example        # Environment variable template
├── .gitignore
├── Makefile            # Project command shortcuts
└── go.mod
```

## 🚀 Cara Menjalankan Project

1. Copy `.env.example` menjadi `.env`
   ```bash
   cp .env.example .env
   ```

2. Jalankan aplikasi
   ```bash
   go run cmd/api/main.go
   # atau menggunakan Makefile:
   make run
   ```

3. Uji endpoint Health Check
   ```bash
   curl http://localhost:8080/health
   ```
