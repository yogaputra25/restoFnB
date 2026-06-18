# Daftar Login

## Seed Data (Auto)

| Role | Chain ID | Email | Password |
|------|----------|-------|----------|
| Admin | `bcfe86d3-0d20-4136-a981-8b64a41174fc` | `admin@seed.com` | `belum bisa` |
| Cashier | `bcfe86d3-0d20-4136-a981-8b64a41174fc` | `cashier@seed.com` | `seed123` |

> Chain ID berubah tiap reset volume. Cek ulang dengan:
> `docker logs restofnb-backend-1 2>&1 | Select-String "created chain"`

Cara pakai:
1. Buka `http://localhost:3000/login`
2. Isi Chain ID (UUID dari logs)
3. Isi Email & Password sesuai tabel
4. Admin → redirect ke `/admin`, Cashier → redirect ke `/cashier`

## Manual (Sebelum Seed)

| Role | Chain ID | Email | Password |
|------|----------|-------|----------|
| Admin | `54704e14-fae0-4b81-8a81-d987ffda544d` | `admin@test.com` | `admin123` |

> Catatan: Akun manual ini sudah terhapus kalau volume Docker di-reset.
