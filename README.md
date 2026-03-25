# Gin Boilerplate API

Boilerplate REST API menggunakan **Go + Gin + PostgreSQL + JWT**.

## 📁 Struktur Project

```
gin-boilerplate/
├── cmd/
│   └── main.go                    # Entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Konfigurasi app & env
│   ├── database/
│   │   └── postgres.go            # Koneksi PostgreSQL
│   ├── handler/
│   │   ├── auth_handler.go        # Handler register & login
│   │   ├── dermaga_handler.go     # Handler manajemen dermaga
│   │   └── user_handler.go        # Handler profile user
│   ├── middleware/
│   │   ├── auth.go                # JWT middleware
│   │   └── middleware.go          # Logger & CORS
│   ├── model/
│   │   ├── dto/                   # Data Transfer Objects (Request/Response)
│   │   │   ├── dermaga_dto.go
│   │   │   └── user_dto.go
│   │   └── entity/                # GORM Entities
│   │       ├── cabang.go
│   │       ├── master.go          # Entity Dermaga
│   │       └── user.go
│   ├── repository/
│   │   ├── cabang_repository.go
│   │   ├── dermaga_repository.go
│   │   └── user_repository.go
│   └── service/
│       ├── auth_service.go
│       ├── dermaga_service.go
│       └── user_service.go
├── pkg/
│   └── utils/
│       ├── jwt.go                 # JWT generate & validate
│       └── response.go            # Standard API response
├── .env.example                   # Template environment
├── go.mod
├── Makefile
└── README.md
```

## 🚀 Cara Menjalankan

### 1. Clone dan setup environment
```bash
cp .env.example .env
# Edit .env sesuai konfigurasi Anda
```

### 2. Install dependencies
```bash
go mod tidy
```

### 3. Jalankan aplikasi
```bash
make run
# atau
go run cmd/main.go
```

## 📡 API Endpoints

### Public Routes

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/health` | Health check |
| POST | `/api/v1/auth/register` | Daftar akun baru |
| POST | `/api/v1/auth/login` | Login & dapatkan token |

### Protected Routes (Butuh JWT Token)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/api/v1/users/profile` | Lihat profil sendiri |
| GET | `/api/v1/dermaga` | List semua dermaga (Paginated) |
| GET | `/api/v1/dermaga/:id` | Lihat detail dermaga |
| POST | `/api/v1/dermaga` | Tambah dermaga baru |
| PUT | `/api/v1/dermaga/:id` | Update data dermaga |
| DELETE | `/api/v1/dermaga/:id` | Hapus dermaga |

## 📝 Contoh Request

### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "username": "johndoe",
    "password": "password123",
    "kd_cabang": "100",
    "kd_terminal": "01"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get Profile
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer <your_jwt_token>"
```

### List Dermaga (Pagination)
```bash
curl -X GET "http://localhost:8080/api/v1/dermaga?page=1&size=10" \
  -H "Authorization: Bearer <your_jwt_token>"
```

### Get Dermaga by ID
```bash
curl -X GET http://localhost:8080/api/v1/dermaga/1 \
  -H "Authorization: Bearer <your_jwt_token>"
```

### Create Dermaga
```bash
curl -X POST http://localhost:8080/api/v1/dermaga \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "nm_dermaga": "Dermaga Tanjung Priok",
    "kd_dermaga": "TPK01",
    "posisi_awal": 0,
    "posisi_akhir": 100,
    "keterangan": "Dermaga utama",
    "status": "active"
  }'
```

### Update Dermaga
```bash
curl -X PUT http://localhost:8080/api/v1/dermaga/1 \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "nm_dermaga": "Dermaga Tanjung Priok Updated",
    "kd_dermaga": "TPK01",
    "posisi_awal": 0,
    "posisi_akhir": 200,
    "keterangan": "Dermaga utama updated",
    "status": "active"
  }'
```

### Delete Dermaga
```bash
curl -X DELETE http://localhost:8080/api/v1/dermaga/1 \
  -H "Authorization: Bearer <your_jwt_token>"
```

## 📦 Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL
- **Auth**: JWT ([golang-jwt](https://github.com/golang-jwt/jwt))
- **Password**: bcrypt
- **Config**: godotenv
