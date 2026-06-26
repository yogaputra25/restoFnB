# DESIGN.md — restoFnB Design System

Semua aturan desain berikut WAJIB diikuti oleh agent AI saat membuat atau mengubah komponen/halaman apapun.

> ⚠️ **Aturan ini binding.** Jangan membuat file CSS, komponen Vue, gambar, atau struktur halaman yang melanggar aturan di bawah ini.

---

## 1. Warna (Color Palette)

Hanya gunakan warna dari palette ini. Dilarang menggunakan warna purple, biru neon, atau cyber color manapun.

| Token | Hex | Penggunaan |
|-------|-----|-----------|
| `--color-primary` | `#C65D3B` | Tombol utama, harga, link aktif, badge |
| `--color-secondary` | `#D97706` | Aksen sekunder |
| `--color-accent` | `#F4B400` | Rating bintang, label promo |
| `--color-success` | `#2E7D32` | Status sukses, order berhasil |
| `--color-danger` | `#C0392B` | Error, hapus, warning |
| `--color-background` | `#FAF8F5` | Latar halaman (warm off-white) |
| `--color-surface` | `#FFFFFF` | Kartu, modal, sheet |
| `--color-surface-secondary` | `#F5F3EF` | Hover state, skeleton bg |
| `--color-border` | `#E8E3DD` | Garis pemisah, outline |
| `--color-text-primary` | `#2F2F2F` | Teks utama |
| `--color-text-secondary` | `#6B6B6B` | Teks bantuan, label |

---

## 2. Tipografi (Typography)

| Level | Font | Ukuran |
|-------|------|--------|
| Hero heading | `Outfit` | `56px` |
| Page title | `Outfit` | `36px` |
| Section heading | `Outfit` | `28px` |
| Card title | `Outfit` | `22px` |
| Body | `Inter` | `16px` |
| Small | `Inter` | `14px` |
| Caption | `Inter` | `12px` |

Gunakan CSS custom properties: `--font-heading`, `--font-body`, `--font-size-*`.

---

## 3. Spasi (Spacing)

Gunakan sistem 8px multiplier:

```
8px  (--space-1)
16px (--space-2)
24px (--space-3)
32px (--space-4)
40px (--space-5)
48px (--space-6)
56px (--space-7)
64px (--space-8)
```

Jangan gunakan nilai spacing di luar kelipatan 8px (kecuali 4px untuk kasus tepi).

---

## 4. Border Radius

| Elemen | Radius |
|--------|--------|
| Card | `20px` |
| Button | `14px` |
| Input | `14px` |
| Image | `18px` |

Gunakan CSS custom properties: `--radius-card`, `--radius-button`, `--radius-input`, `--radius-image`.

---

## 5. Bayangan (Shadows)

Hanya gunakan soft natural shadow. Dilarang menggunakan efek glow/neon.

```
--shadow-sm: 0 1px 3px 0 rgba(47, 47, 47, 0.06)
--shadow-md: 0 4px 12px 0 rgba(47, 47, 47, 0.08)
--shadow-lg: 0 8px 24px 0 rgba(47, 47, 47, 0.10)
```

---

## 6. Font

- **Heading**: `'Outfit', sans-serif` — Google Fonts
- **Body**: `'Inter', sans-serif` — Google Fonts

Dilarang menggunakan: SF Pro Display, Circular, Gotham, Avenir, Proxima Nova, Helvetica Neue.

---

## 7. Icon

Hanya gunakan **Lucide Icons** untuk semua icon. Dilarang menggunakan emoji sebagai icon (emoji hanya untuk decorative/empty state/illustration).

Auto-import via `lucide-vue-next`.

Pengecualian: SVG inline untuk icon yang tidak ada di Lucide (misal: logo custom).

---

## 8. Efek Hover & Transisi

| Elemen | Efek | Durasi |
|--------|------|--------|
| Button hover | `scale(1.03)` | `250ms ease` |
| Card hover | `translateY(-4px)` + image `scale(1.05)` + shadow naik | `250ms ease` |
| Link hover | Color transition ke `--color-primary` | `150ms ease` |
| Focus input | Border + ring `--color-primary` | `150ms ease` |

Gunakan CSS custom property `--transition-base` untuk 250ms, `--transition-fast` untuk 150ms.

---

## 9. Animasi

Gunakan **AOS** (Animate on Scroll) untuk entrance animation saat scroll.

- Library: `aos` (MIT, free)
- Init di `plugins/aos.client.ts` (client-only)
- Config: `duration: 800`, `once: true`, `offset: 100`
- Gunakan attribute: `data-aos="fade-up"`, `data-aos="fade-left"`, `data-aos="fade-right"`, `data-aos="zoom-in"`

Micro-interactions CSS (tanpa library tambahan):

| Interaksi | Animasi |
|-----------|---------|
| Favorite heart | `@keyframes heart-bounce` (scale 1.3x → 1x, 300ms) |
| Toast notification | Slide down from top (300ms enter, 200ms exit) |
| Page transition | Fade `opacity` via Vue `<Transition mode="out-in">` |
| Cart badge | Bounce scale on count change |

Dilarang menggunakan: GSAP, Motion.dev, Framer Motion, atau library animasi berbayar.

---

## 10. Loading State

Gunakan **Skeleton loading** (bukan spinner) untuk semua konten.

- Komponen: `<AppSkeleton>`
- Variants: `text`, `card`, `image`, `circle`
- Animasi: gradient shimmer via CSS `@keyframes shimmer`

Spinner hanya digunakan di dalam tombol (`loading` state) dan untuk proses transien singkat.

---

## 11. Empty State

Setiap daftar/data kosong WAJIB menampilkan:

1. Ilustrasi (dari unDraw/ManyPixels/Open Doodles gratis)
2. Judul ramah
3. Deskripsi
4. CTA button (jika relevan)

Gunakan komponen: `<AppEmptyState>`

---

## 12. Error Pages

- **404**: Tema "Dish not found" — ilustrasi + "Back to Menu" button
- **500**: Tema "Kitchen is on fire" — ilustrasi + "Try Again" button

Gunakan komponen: `<AppErrorPage>`

---

## 13. API Response Format

Semua endpoint backend WAJIB mengembalikan JSON dengan format berikut:

### Success (single item)
```json
{
  "success": true,
  "code": 200,
  "message": "Success",
  "timestamp": "2026-06-26T12:00:00Z",
  "data": {}
}
```

### Success (list dengan pagination)
```json
{
  "success": true,
  "code": 200,
  "message": "Success",
  "timestamp": "2026-06-26T12:00:00Z",
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "limit": 10,
      "totalItems": 120,
      "totalPages": 12
    }
  }
}
```

### Validation Error (400)
```json
{
  "success": false,
  "code": 400,
  "message": "Validation Error",
  "errors": [
    { "field": "email", "message": "Email is required" }
  ]
}
```

### Error lainnya
- **401**: `{ "success": false, "code": 401, "message": "Unauthorized" }`
- **403**: `{ "success": false, "code": 403, "message": "Forbidden" }`
- **404**: `{ "success": false, "code": 404, "message": "Resource not found" }`
- **500**: `{ "success": false, "code": 500, "message": "Internal Server Error" }`

Gunakan helper dari `backend/pkg/response/response.go`:
- `response.Success(w, data)`
- `response.Created(w, data)`
- `response.List(w, items, page, limit, total)`
- `response.Error(w, status, msg)`
- `response.ValidationError(w, errors)`
- `response.Unauthorized(w)`
- `response.Forbidden(w)`
- `response.NotFound(w)`
- `response.ServerError(w)`

---

## 14. Komponen UI (Wajib Pakai)

Jangan buat ulang pattern yang sudah ada. Gunakan komponen berikut:

| Komponen | Lokasi | Fungsi |
|----------|--------|--------|
| `<AppButton>` | `components/ui/AppButton.vue` | Tombol dengan variant primary/secondary/danger/success |
| `<AppInput>` | `components/ui/AppInput.vue` | Input dengan floating label + validasi |
| `<AppCard>` | `components/ui/AppCard.vue` | Kartu menu dengan image, rating, favorite, add-to-cart |
| `<AppSkeleton>` | `components/ui/AppSkeleton.vue` | Skeleton loading (text/card/image/circle) |
| `<AppEmptyState>` | `components/ui/AppEmptyState.vue` | Empty state dengan ilustrasi + CTA |
| `<AppErrorPage>` | `components/ui/AppErrorPage.vue` | Halaman error 404/500 |
| `<AppNavbar>` | `components/ui/AppNavbar.vue` | Navigasi sticky dengan cart + search |

Semua komponen auto-import tanpa prefix (konfigurasi di `nuxt.config.ts`).

---

## 15. Layout Grid

| Viewport | Kolom |
|----------|-------|
| ≥1024px (desktop) | 4 |
| ≥640px (tablet) | 2 |
| <640px (mobile) | 1 |

Gunakan Tailwind: `grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4`

---

## 16. Gambar

- Semua gambar WAJIB pakai `loading="lazy"`
- Hero image: ukuran besar, `object-cover`
- Food card image: aspect ratio `4:3`, `object-cover`
- Sumber gambar hanya dari: **Unsplash, Pexels, Pixabay** (royalty-free)
- Format: WebP prefered, fallback JPG/PNG

---

## 17. Ilustrasi

Gunakan hanya dari sumber gratis:

- unDraw (https://undraw.co/illustrations)
- Storyset (https://storyset.com/) — dengan atribusi
- ManyPixels Free Gallery
- Open Doodles (https://www.opendoodles.com/)

---

## 18. Kode & Performa

- **Lazy loading** — semua gambar, komponen heavy, AOS library (client-only)
- **CSS custom properties** — token desain via `tokens.css`, bukan hardcode
- **Responsive-first** — mobile-first breakpoints
- **Semantic HTML** — heading levels, landmark elements, alt text
- **WCAG AA** — kontras warna cukup untuk aksesibilitas
- **Tree shaking** — hindari import library besar yang tidak dipakai
- **Code splitting** — halaman terpisah per route (Nuxt otomatis)

---

## 19. Kepatuhan Lisensi

**DILARANG KERAS menggunakan aset berbayar tanpa lisensi:**

- ❌ Font: SF Pro Display, Circular, Gotham, Avenir, Proxima Nova, Helvetica Neue
- ❌ Icon: FontAwesome Pro, Material Design Icons (beberapa varian)
- ❌ Library animasi: GSAP, Motion.dev (untuk React — tidak relevan), Framer Motion
- ❌ Ilustrasi: Humaaans (beberapa varian), Icons8 (tanpa atribusi)

✅ SEMUA aset yang digunakan WAJIB gratis / open-source / MIT license.
