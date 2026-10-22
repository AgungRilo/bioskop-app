# Bioskop App

Aplikasi web sederhana untuk mengelola data bioskop, dibuat dengan **Go**, framework **Gin**, dan database **PostgreSQL**.

## Struktur Project

```
bioskop-app/
├── main.go                            # entry point: koneksi database & menjalankan server
├── go.mod / go.sum                    # dependency (gin, lib/pq)
├── config/
│   └── database.go                    # koneksi PostgreSQL & pembuatan tabel otomatis
├── models/
│   └── bioskop.go                     # struct Bioskop
├── controllers/
│   └── bioskop_controller.go          # handler endpoint bioskop
└── routes/
    └── routes.go                      # daftar route API
```

## Model Data

| Field  | Tipe    | Keterangan                    |
|--------|---------|-------------------------------|
| id     | integer | primary key, auto increment   |
| nama   | string  | wajib diisi                   |
| lokasi | string  | wajib diisi                   |
| rating | float   | opsional                      |

Tabel `bioskop` dibuat otomatis saat server pertama kali dijalankan.

## Cara Menjalankan

1. Pastikan PostgreSQL sudah berjalan, lalu buat database:
   ```
   createdb bioskop_db
   ```
2. Jalankan server:
   ```
   go run main.go
   ```
3. Server berjalan di `http://localhost:8080`.

### Konfigurasi (environment variable)

| Variable      | Default      |
|---------------|--------------|
| `DB_HOST`     | `localhost`  |
| `DB_PORT`     | `5432`       |
| `DB_USER`     | `agungrilo`  |
| `DB_PASSWORD` | *(kosong)*   |
| `DB_NAME`     | `bioskop_db` |
| `PORT`        | `8080`       |

Contoh: `DB_USER=postgres DB_PASSWORD=secret PORT=9000 go run main.go`

## Daftar API

| Method | Endpoint   | Keterangan                 |
|--------|------------|----------------------------|
| POST   | `/bioskop` | Menambahkan bioskop baru   |
| GET    | `/bioskop` | Menampilkan semua bioskop  |
| GET    | `/bioskop/:id` | Menampilkan detail bioskop |
| PUT    | `/bioskop/:id` | Memperbarui data bioskop   |
| DELETE | `/bioskop/:id` | Menghapus bioskop          |

### POST /bioskop

Menambahkan bioskop baru. `nama` dan `lokasi` tidak boleh kosong, `rating` harus 0–5.

**Request body (JSON)**
```json
{
  "nama": "CGV Grand Indonesia",
  "lokasi": "Jakarta",
  "rating": 4.5
}
```

**Response `201 Created`**
```json
{
  "message": "bioskop berhasil ditambahkan",
  "data": {
    "id": 1,
    "nama": "CGV Grand Indonesia",
    "lokasi": "Jakarta",
    "rating": 4.5
  }
}
```

**Response `400 Bad Request`**, jika nama/lokasi kosong
```json
{ "error": "Nama dan Lokasi tidak boleh kosong" }
```

**Response `400 Bad Request`**, jika format JSON tidak valid
```json
{ "error": "data JSON tidak valid: ..." }
```

**Contoh curl**
```
curl -X POST http://localhost:8080/bioskop \
  -H "Content-Type: application/json" \
  -d '{"nama":"CGV Grand Indonesia","lokasi":"Jakarta","rating":4.5}'
```

### GET /bioskop

Menampilkan semua data bioskop, diurutkan berdasarkan `id`.

**Response `200 OK`**
```json
{
  "data": [
    {
      "id": 1,
      "nama": "CGV Grand Indonesia",
      "lokasi": "Jakarta",
      "rating": 4.5
    }
  ]
}
```

**Contoh curl**
```
curl http://localhost:8080/bioskop
```

### GET /bioskop/:id

Menampilkan detail satu bioskop.

- `200 OK` → `{ "data": { "id": 1, "nama": "...", "lokasi": "...", "rating": 4.5 } }`
- `400 Bad Request` → `{ "error": "ID harus berupa angka positif" }`
- `404 Not Found` → `{ "error": "bioskop dengan ID 1 tidak ditemukan" }`

```
curl http://localhost:8080/bioskop/1
```

### PUT /bioskop/:id

Memperbarui data bioskop. Body dan validasi sama dengan `POST /bioskop`.

- `200 OK` → `{ "message": "bioskop berhasil diperbarui", "data": { ... } }`
- `400 Bad Request` → ID tidak valid, JSON tidak valid, nama/lokasi kosong, atau rating di luar 0–5
- `404 Not Found` → bioskop tidak ditemukan

```
curl -X PUT http://localhost:8080/bioskop/1 \
  -H "Content-Type: application/json" \
  -d '{"nama":"CGV Grand Indonesia","lokasi":"Jakarta Pusat","rating":4.7}'
```

### DELETE /bioskop/:id

Menghapus bioskop.

- `200 OK` → `{ "message": "bioskop berhasil dihapus" }`
- `400 Bad Request` → ID tidak valid
- `404 Not Found` → bioskop tidak ditemukan

```
curl -X DELETE http://localhost:8080/bioskop/1
```

### Endpoint tidak dikenal

Semua route yang tidak terdaftar mengembalikan `404` dengan `{ "error": "endpoint tidak ditemukan" }`.
