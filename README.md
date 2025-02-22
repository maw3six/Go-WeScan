# Webshell Finder

Go-WeScan adalah tool berbasis Go yang digunakan untuk mendeteksi keberadaan backdoor (webshell) di dalam website dengan melakukan pemindaian berdasarkan daftar target dan path backdoor yang telah dikonfigurasi sebelumnya.

## 🚀 Fitur Utama
- **Multi-threading**: Memanfaatkan worker pool hingga 50 concurrent requests.
- **Random User-Agent**: Menggunakan User-Agent yang berbeda untuk setiap request.
- **Retry Mechanism**: Otomatis mencoba ulang hingga 3 kali jika terjadi error koneksi.
- **Path Shuffling**: Mengacak daftar path sebelum scanning untuk menghindari pola scanning yang mudah dideteksi.
- **Real-time Logging**: Menyimpan hasil ke `results.txt` secara langsung tanpa delay.

## 📂 Struktur File
```
📁 Webshell-Finder/
├── 📜 main.go             # Kode utama
├── 📜 README.md           # Dokumentasi ini
├── 📜 results.txt         # Hasil scanning (akan dibuat otomatis)
├── 📂 lib/
│   ├── shell.json         # Daftar path backdoor yang dicari
│   ├── useragent.json    # Daftar User-Agent untuk scanning
```

Atau download dari [golang.org](https://golang.org/dl/).

### 5️⃣ Hasil Scanning
Setiap backdoor yang ditemukan akan tersimpan di `results.txt` dengan format:
```
http://example.com/wp-content/uploads/shell.php
http://target.com/adminer.php
```

## ⚠️ Disclaimer
Tool ini hanya digunakan untuk tujuan **pentesting legal** dan **pengujian keamanan**. Penggunaan tanpa izin dapat melanggar hukum yang berlaku. Gunakan dengan tanggung jawab penuh!

## 📜 Lisensi
MIT License. Bebas digunakan dan dimodifikasi dengan tetap menyertakan atribusi.

---
🔥 **Dikembangkan oleh**: @maw3six | Versi: 1.0

