# 🕷️ Bug Hunting & Security Audit Reports

Dokumen ini mencatat temuan bug, celah keamanan, dan anomali logika yang ditemukan selama sesi simulasi mendalam dan pengujian stres pada sistem FOSSBilling Next-Gen.

---

## 📑 Ringkasan Temuan (Sesi Sept 2024)

| ID | Tingkat Bahaya | Kategori | Judul Bug | Status |
| :--- | :--- | :--- | :--- | :--- |
| **BUG-01** | 🔴 **CRITICAL** | Finansial | Negative Pricing in Cart Checkout | Resolved |
| BUG-02 | 🟢 Info | Validasi | Invalid Registration Payload Handling | Resolved |
| BUG-03 | 🟢 Info | Billing | Promo Discount > Subtotal Overflow | Resolved |
| BUG-04 | 🟡 Medium | Database | Duplicate Email Race Condition | Resolved |
| **BUG-05** | 🔴 **CRITICAL** | Keamanan | IDOR in Invoice Payment via Balance | Resolved |
| **BUG-06** | 🔴 **CRITICAL** | Finansial | Price Injection via ProductID=0 | Resolved |
| **BUG-07** | 🟠 High | Integritas | Promo MaxUses Race Condition | Resolved |
| **BUG-13** | 🔴 **CRITICAL** | Finansial | Price Injection via Invalid ProductID | Resolved |
| **BUG-14** | 🟡 Medium | Logika | Zero/Negative Quantity Acceptance | Resolved |
| **BUG-15** | 🟡 Medium | Keamanan | ClientID Enumeration via Guest API | Resolved |
| **BUG-08** | 🟠 High | Keamanan | Stored XSS in Profile Fields | Resolved |
| **BUG-09** | 🟠 High | Keamanan | Stored XSS in Support Tickets | Resolved |
| **BUG-12** | 🟡 Medium | Keamanan | Weak Password Policy | Resolved |
| **BUG-16** | 🟠 High | Akses | Privilege Escalation (Unsuspend) | Resolved |
| **BUG-17** | 🟠 High | Akses | Missing Auth Check (Activate) | Resolved |
| **BUG-18** | 🟠 High | Integritas | Promo OncePerClient Race Condition | Resolved |
| **BUG-19** | 🔴 **CRITICAL** | Finansial | Deposit Paid but Balance Not Increased | Resolved |
| **BUG-20** | 🟡 Medium | Logika | False Activation on Provisioning Failure | Resolved |

---

## ✅ BUG-20: False Activation on Provisioning Failure (RESOLVED)

### Deskripsi
Ditemukan bahwa jika proses teknis pembuatan akun (provisioning) gagal di server tujuan, sistem tetap membiarkan status pesanan menjadi `active` dan mengirimkan email selamat datang kepada klien. Hal ini menyebabkan ketidaksesuaian data antara billing dan server real.

### Perbaikan
Menambahkan pengecekan status sukses pada `OrderListener`. Jika provisioning gagal, status pesanan otomatis di-*rollback* ke `pending_setup` dan pengiriman email aktivasi dibatalkan. Staf admin dapat melihat kegagalan di log dan mencoba memproses ulang setelah memperbaiki masalah teknis di server.

---

## ✅ BUG-19: Deposit Paid but Balance Not Increased (RESOLVED)

### Deskripsi
Ditemukan bahwa validasi pembatasan penggunaan kupon "Satu kali per klien" (`OncePerClient`) dilakukan di luar mekanisme penguncian database (*lock*). Hal ini memungkinkan satu pengguna untuk mengklaim diskon yang sama berkali-kali jika request dikirimkan secara simultan (konkuren).

### Perbaikan
Memindahkan logika pengecekan riwayat penggunaan kupon (`redemptions`) ke dalam blok *Atomic Lock* pada layer Repository. Sekarang, meskipun request datang bersamaan, hanya satu yang akan diproses, dan sisanya akan ditolak dengan error `promo already used`.

---

## ✅ BUG-16 & BUG-17: Privilege Escalation in Admin Endpoints (RESOLVED)

### Deskripsi
Ditemukan bahwa beberapa endpoint administrasi tingkat tinggi seperti `UnsuspendOrder`, `ActivateOrder`, `SyncOrder`, dan `ChangeOrderPassword` tidak melakukan pengecekan hak akses (*permissions*) staff. Hal ini memungkinkan staf dengan peran terbatas (misal: Support Read-Only) untuk melakukan tindakan modifikasi pada layanan klien.

### Perbaikan
1.  Melakukan audit menyeluruh pada seluruh layer Handler Administrasi (`core/handler/http/admin/`).
2.  Mengintegrasikan `staffService.HasPermission` pada setiap fungsi yang melakukan modifikasi data (Write/Delete).
3.  Memastikan setiap Handler menerima `staffService` melalui konstruktor untuk validasi izin yang konsisten.

---

## ✅ BUG-12: Weak Password Policy (RESOLVED)

### Deskripsi
Sistem sebelumnya hanya memvalidasi panjang minimal password (6 karakter) tanpa mengecek kompleksitas. Hal ini memungkinkan penggunaan password yang sangat lemah seperti `123456` atau `abcdef`.

### Perbaikan
Memperbarui `AuthUsecase.Register` untuk menggunakan `v.CheckPasswordStrength`, yang mewajibkan:
1. Panjang minimal 8 karakter.
2. Mengandung huruf besar (uppercase).
3. Mengandung huruf kecil (lowercase).
4. Mengandung angka (numeric).

---

## ✅ BUG-08: Stored XSS in Profile Fields (RESOLVED)

---

## ✅ BUG-08: Stored XSS in Profile Fields (RESOLVED)

### Deskripsi
Input pada field profil seperti `first_name`, `last_name`, dan `city` tidak disanitasi sebelum disimpan ke database. Hal ini memungkinkan serangan Stored XSS jika data tersebut dirender di area dashboard Admin tanpa *escaping* yang tepat.

### Perbaikan
Mengimplementasikan `security.SanitizeAlphaNumeric` dan `security.SanitizeHTML` pada layer `AuthUsecase.UpdateProfile`.

---

## ✅ BUG-09: Stored XSS in Support Tickets (RESOLVED)

### Deskripsi
Pesan tiket bantuan (`Content`) menerima input HTML mentah termasuk script berbahaya seperti `<script>` atau atribut `onerror`.

### Perbaikan
Menggunakan `security.SanitizeHTML` untuk membersihkan tag dan atribut berbahaya pada seluruh pesan tiket di `SupportService`.

---

## ✅ BUG-13: Price Injection via Invalid ProductID (RESOLVED)

### Deskripsi
Sistem menerima `product_id` yang tidak ada di database dan tetap menggunakan harga yang dikirimkan oleh klien. Hal ini menyebabkan kegagalan dalam proses validasi harga (*Price Override*).

### Perbaikan
Menambahkan pengecekan keberadaan produk di `CartService.CalculateTotals`. Jika `product_id` diberikan tapi tidak ditemukan di database, sistem akan mengembalikan error.

---

## ✅ BUG-15: ClientID Enumeration via Guest API (RESOLVED)

### Deskripsi
Endpoint publik `/api/v1/guest/cart/calculate` sebelumnya menerima field `client_id` secara opsional. Jika diberikan, sistem akan melakukan pencarian profil klien tersebut untuk menghitung pajak. Hal ini memungkinkan penyerang (Guest) untuk mendeteksi apakah suatu `client_id` valid dan menebak lokasi geografis klien berdasarkan pajak yang diterapkan.

### Perbaikan
Memperbarui `guest.CartHandler` (Calculate dan Checkout) untuk selalu memaksa `ClientID = 0`. Request dari Guest kini tidak lagi dapat memicu pencarian profil klien ke database.

---

## ✅ BUG-14: Zero/Negative Quantity Acceptance (RESOLVED)

### Deskripsi
Sistem secara otomatis mengubah quantity nol atau negatif menjadi `1` tanpa memberitahu pengguna. Meskipun terlihat seperti "fitur", hal ini bisa disalahgunakan untuk membuat pesanan yang tidak diinginkan atau membingungkan alur inventaris.

### Perbaikan
Menambahkan validasi ketat pada `CartService`. Sekarang sistem menolak checkout jika quantity item kurang dari atau sama dengan nol.

---

## ✅ BUG-06: Price Injection via ProductID=0 (RESOLVED)

### Deskripsi
Input `product_id` bernilai `0` dianggap valid oleh sistem, melewati pengecekan database, dan menggunakan harga kustom dari payload API.

### Perbaikan
Menambahkan validasi `if it.ProductID <= 0` pada awal siklus perhitungan total keranjang.

---

## ✅ BUG-05: IDOR in Invoice Payment via Balance (RESOLVED)

### Deskripsi
Ditemukan celah keamanan **IDOR (Insecure Direct Object Reference)** pada fitur pembayaran invoice menggunakan saldo. Sebelumnya, sistem tidak memvalidasi apakah klien yang meminta pembayaran adalah pemilik sah dari invoice tersebut.

### Dampak Bisnis
*   **Akses Data Ilegal**: Penyerang bisa menebak ID Invoice klien lain.
*   **Pencurian Saldo**: Jika fungsi tidak divalidasi, penyerang bisa memicu pembayaran invoice mereka sendiri menggunakan saldo korban, atau sebaliknya.

### Rekomendasi Perbaikan
Signature fungsi `PayWithBalance` pada `InvoiceService` diperbarui untuk menerima `clientID` dan memvalidasi kepemilikan sebelum memproses transaksi.

---

## ✅ BUG-01: Negative Pricing in Cart Checkout (RESOLVED)

### Deskripsi
Ditemukan bahwa layer `CartService` tidak melakukan validasi nilai terhadap harga unit (`Price`) pada item yang dikirimkan melalui payload checkout. Hal ini memungkinkan pengguna jahat untuk memanipulasi request API dan menyertakan item dengan harga negatif.

### Dampak Bisnis
*   **Pencurian Layanan**: Pengguna dapat menambahkan item seharga `-Rp 10.000.000` untuk menihilkan total tagihan produk lainnya.
*   **Manipulasi Saldo**: Jika invoice negatif dianggap sah, pembayaran via balance bisa berpotensi menambah saldo pengguna secara ilegal.
*   **Integritas Data**: Laporan finansial dan MRR akan menjadi tidak akurat (terdistorsi oleh angka negatif).

### Bukti (Proof of Concept)
Dijalankan melalui `cmd/demo/bug_hunter.go`:
```go
badCart := &cart.Cart{
    ClientID: 1,
    Items: []cart.CartItem{
        {ProductID: 101, Title: "Hacker Item", Price: decimal.FromFloat(-5000000.00), Quantity: 1},
    },
}
res, _ := cartService.Checkout(ctx, badCart)
fmt.Println(res.Invoice.Total) // Hasil: -5000000.00
```

### Rekomendasi Perbaikan
Tambahkan pengecekan pada `backend-go/core/usecase/cart/cart_service.go` di dalam fungsi `CalculateTotals` atau `Checkout`:
```go
if item.Price < 0 {
    return errors.New("item price cannot be negative")
}
```

---

## ✅ Skenario Teruji (Lulus)

### BUG-02: Invalid Registration Payload
*   **Skenario**: Mendaftarkan user dengan format email salah, password terlalu pendek, dan field wajib kosong.
*   **Hasil**: Sistem berhasil menangkap error melalui layer `validator` dan memberikan respons `400 Bad Request` yang detail.

### BUG-03: Promo Discount Overlap
*   **Skenario**: Menggunakan kupon diskon nominal (Absolute) yang nilainya lebih besar dari total belanja (misal: Diskon Rp 1.000.000 untuk belanja Rp 10.000).
*   **Hasil**: Sistem secara otomatis melakukan *capping* pada angka `0.00`. Invoice tidak menjadi negatif.

### BUG-04: Duplicate Registration
*   **Skenario**: Dua pendaftaran simultan dengan email yang sama.
*   **Hasil**: Layer Repository berhasil melempar `ErrDuplicateEntry` dan mencegah duplikasi data pada basis data.

---

## 🛠️ Metodologi Pengujian
Pengujian dilakukan menggunakan:
1.  **Deep E2E Simulation**: Simulasi alur kerja user dari registrasi hingga terminasi.
2.  **Negative Unit Testing**: Mengirimkan input yang sengaja dirusak ke layer Usecase.
3.  **Boundary Value Analysis**: Menguji nilai batas bawah (angka negatif) pada modul billing.

---
*Terakhir diperbarui: 2024-09-08 oleh AI Assistant (Bug Hunter Mode)*
