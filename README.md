# Electronic Medical Record (EMR) Doctor Service 🏥

A high-performance, production-grade backend REST API service for the **Doctor Module (Electronic Medical Record / EMR)**, engineered as a **modern Go rewrite of SIMRS Khanza** (one of the most widely adopted Hospital Information Management Systems in Indonesia).

Built entirely with the **Go Standard Library (`net/http`)**, this service transitions legacy hospital workflows into a fast, stateless, and secure healthcare microservice featuring **Clean Architecture**, **AES-256-GCM URL Obfuscation**, **SOAP & Vital Signs clinical validations**, and interactive **OpenAPI 3.1 & Scalar API Reference**.

---

## 💡 Background & Motivation

**SIMRS Khanza** is a prominent open-source Hospital Management Information System widely deployed across Indonesian hospitals and clinics, historically operating on Java desktop (Swing) and PHP architectures. 

This project modernizes the critical **Doctor & Outpatient Examination** workflow by decoupling it into an independent, high-throughput Go backend:
- **Performance**: Sub-millisecond response times with zero third-party web framework overhead.
- **Security & Patient Privacy**: Compliant with modern healthcare data privacy guidelines by masking and encrypting sensitive patient record numbers in URLs.
- **Developer Experience**: Interactive OpenAPI 3.1 documentation powered by Scalar UI embedded directly into the binary.

---

## 🌟 Key Features

### 🔐 1. Authentication & Medical Authorization
- **Khanza Credential Verification**: Secure doctor authentication against Khanza user databases with encrypted password verification.
- **Stateless JWT Sessions**: HMAC-SHA256 signed bearer tokens with configurable expiration.
- **Ownership & Access Control**: Doctors can only edit, update, or delete examination records created by themselves.

### 🏥 2. Outpatient & Patient Queue Management
- **Doctor Patient Queue**: Filter patient queues by examination status, insurance/payer (`BPJS`, private, cash), visit status (`Ralan`, `Ranap`), and referral type.
- **Search & Pagination**: Full-text keyword search (patient name, medical record number, registration number), date-range filtering, and dynamic multi-column sorting.
- **Patient Detail**: Comprehensive clinical profile and visit metadata.

### 📝 3. Clinical Examination & Medical Records (SOAP & Vital Signs)
- **Standardized SOAP Workflow**: Subjective, Objective, Assessment, and Plan documentation.
- **Vital Signs (TTV) Clinical Boundary Validations**:
  - Blood Pressure (Systolic & Diastolic ranges)
  - Body Temperature, Heart Rate, Respiratory Rate, and Oxygen Saturation (SpO2)
  - Neurological Consciousness Scales: **GCS** (Glasgow Coma Scale) & **AVPU** (*Compos Mentis, Somnolence, Sopor, Coma, Alert, Confusion, Voice, Pain, Unresponsive*)
- **Longitudinal Medical History**: Retrieve a patient's complete cross-visit medical history (`id_pasien`) or single-visit examination records (`id_kunjungan`).
- **48-Hour Medical Audit Lock**: Examination records are legally locked against modification or deletion 48 hours after creation, adhering to medical compliance standards.

### 🛡️ 4. Advanced Security & Resiliency
- **URL Obfuscation (AES-256-GCM)**: Sensitive patient identifiers (`no_rawat`, `no_rkm_medis`, and composite examination keys) are symmetrically encrypted in URL paths, preventing ID enumeration and tampering.
- **Brute-Force Rate Limiting**: Dedicated rate limiter middleware protecting authentication endpoints (10 requests/minute).
- **Context Timeouts**: Automatic 10-second timeout propagation to safeguard against cascading database bottlenecks.
- **CORS Middleware**: Flexible origin handling for modern web and mobile frontend integration.

### 📖 5. Interactive API Documentation
- **Embedded Scalar API Reference**: Modern, responsive UI accessible at `GET /docs`.
- **OpenAPI 3.1.0 Specification**: Raw schema served at `GET /docs/openapi.yaml`.

---

## 🛠️ Architecture & Tech Stack

| Component | Technology | Description |
| :--- | :--- | :--- |
| **Language** | Go 1.22+ | Zero framework bloat; uses standard library `http.ServeMux` |
| **Database** | MySQL 5.7+ / 8.0+ / MariaDB | SIMRS Khanza relational database schema (`sik`) |
| **Driver & Pooling** | `go-sql-driver/mysql` | Native Go MySQL driver with connection pooling |
| **Security / Crypto** | AES-256-GCM, `golang-jwt/jwt/v5` | Token authentication and sensitive ID obfuscation |
| **Documentation** | Scalar API Reference & OpenAPI 3.1 | Embedded interactive API documentation |
| **Benchmarking** | [k6](https://k6.io/) | Load testing and throughput benchmarking scripts |

### Project Layout (Modular Clean Architecture)

```text
erm-dokter/
├── cmd/
│   └── api/
│       └── main.go               # Application entrypoint & dependency injection
├── internal/
│   ├── auth/                     # Doctor authentication & JWT generation
│   ├── config/                   # Environment configuration loader
│   ├── di/                       # Dependency Injection container
│   ├── docs/                     # Scalar UI & OpenAPI 3.1 YAML handlers
│   ├── health/                   # Liveness probe / health check
│   ├── master/                   # Insurance providers, clinics, and medication master data
│   ├── middleware/               # Auth JWT, CORS, Rate Limiter, Context Timeout
│   ├── pasien/                   # Patient demographic domain
│   ├── pemeriksaan/              # Clinical EMR (SOAP & Vital Signs) domain
│   ├── pkg/                      # Shared utility packages:
│   │   ├── crypto/               # AES-256-GCM URL encryption/decryption
│   │   ├── database/             # Database connection pool manager
│   │   ├── logger/               # Structured logging
│   │   ├── response/             # Standardized JSON response envelope
│   │   └── token/                # JWT claim generator & validator
│   ├── rawatjalan/               # Outpatient queue & visit domain
│   ├── resep/                    # Doctor e-prescription domain
│   ├── routes/                   # HTTP Route Registry (Go 1.22+ ServeMux)
│   └── shared/                   # Domain errors, pagination metadata, common types
├── test.http                     # Comprehensive REST Client test suites
├── k6-script.js                  # k6 load testing & stress test script
├── .air.toml                     # Air live-reload configuration
├── .env.example                  # Environment variable template
├── Makefile                      # Build and run automation shortcuts
└── go.mod                        # Go module definition
```

---

## 🚀 Getting Started

### Prerequisites
- [Go](https://go.dev/dl/) 1.22 or higher
- MySQL 5.7+ / 8.0+ or MariaDB 10.3+ with SIMRS Khanza schema (`sik`)
- (Optional) [Air](https://github.com/air-verse/air) for live reloading

### 1. Clone & Configure Environment
```bash
git clone https://github.com/khoirir/erm-dokter.git
cd erm-dokter
cp .env.example .env
```

Configure your database connection and 32-character AES encryption key in `.env`:
```env
APP_ENV=development
APP_PORT=8082
APP_NAME=erm-dokter

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_database_password
DB_NAME=sik

JWT_SECRET=your_super_secret_jwt_key_here
USER_KEY=your_aes_user_key_here
PASSWORD_KEY=your_aes_password_key_here
ENCRYPTION_KEY=your_aes_encryption_key_32_characters_here

CORS_ORIGIN=*
MAX_EDIT_REKAM_MEDIS_JAM=48
```

### 2. Download Dependencies
```bash
go mod tidy
```

### 3. Run the Server
```bash
# Using Go CLI
go run cmd/api/main.go

# Or using Makefile
make run

# Or with live-reload (Air)
air
```

The server will start at `http://localhost:8082`.

---

## 📑 Interactive Documentation (Scalar UI)

Interactive API documentation is embedded directly into the application:
- **Scalar API Reference UI**: [http://localhost:8082/docs](http://localhost:8082/docs)
- **OpenAPI 3.1 YAML Specification**: [http://localhost:8082/docs/openapi.yaml](http://localhost:8082/docs/openapi.yaml)

---

## 🧪 Testing & Benchmarking

### Automated Unit Tests
```bash
go test -v ./...
```

### End-to-End Testing (`test.http`)
The [`test.http`](test.http) file provides executable test scenarios covering Auth, Outpatient Queues, Clinical Examinations, and error responses (`400`, `401`, `403`, `404`).

### Performance & Load Testing (`k6`)
Execute the bundled k6 script to benchmark throughput and latency under load:
```bash
k6 run k6-script.js
```

---

## 📄 License

This project is developed for Electronic Medical Record (EMR) integration and modernization of Hospital Information Systems.
