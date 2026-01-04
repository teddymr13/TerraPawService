# Panduan Migrasi Database Lokal ke Railway via DBeaver

Panduan ini untuk memindahkan seluruh isi database dari Laptop (Lokal) ke Railway (Cloud) menggunakan DBeaver.

## Prasyarat
- **Koneksi Lokal**: Database lama Anda di DBeaver (Source).
- **Koneksi Railway**: Database baru di DBeaver yang terhubung ke Railway Public Network (Target).

---

## Metode 1: Backup & Restore (Disarankan)
Metode ini paling aman karena memindahkan struktur (schema), index, dan data sekaligus.

### Langkah 1: Backup (Export) Database Lokal
1. Di DBeaver, **klik kanan** pada nama database **Lokal** Anda.
2. Pilih menu **Tools** > **Backup**.
3. Centang schema `public`.
4. Klik **Next**.
5. Pada bagian Format, pilih **Plain** (agar menghasilkan file .sql).
6. Tentukan lokasi penyimpanan file di komputer Anda (misal: `backup_terrapaw.sql`).
7. Klik **Start** / **Proceed**.
8. Tunggu hingga proses selesai ("Backup finished").

### Langkah 2: Restore (Import) ke Railway
1. Di DBeaver, **klik kanan** pada nama database **Railway** (koneksi baru).
2. Pilih menu **Tools** > **Restore**.
3. Pada kolom "Backup file", cari dan pilih file `backup_terrapaw.sql` yang tadi dibuat.
4. Klik **Start** / **Proceed**.
5. DBeaver akan mengunggah dan menjalankan script SQL tersebut di server Railway.

### Langkah 3: Verifikasi
1. Klik kanan koneksi Railway > **Refresh**.
2. Buka folder **Tables** dan pastikan semua tabel sudah ada.
3. Klik kanan salah satu tabel > **View Data** untuk memastikan isinya benar.

---

## Metode 2: Copy Tables (Cepat untuk data sedikit)
Jika tabel Anda sedikit, Anda bisa copy-paste langsung antar koneksi.

1. Buka koneksi **Lokal**, blok semua tabel di folder Tables.
2. Klik Kanan > **Copy**.
3. Buka koneksi **Railway**, klik kanan pada folder Tables (di schema public).
4. Klik **Paste**.
5. Akan muncul wizard import data, klik **Next** terus hingga **Finish**.

---

## Setelah Migrasi Selesai
Setelah data berhasil dipindahkan, restart service **Backend** di dashboard Railway (tombol **Redeploy**) untuk memastikan aplikasi membaca data terbaru.
