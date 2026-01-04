# Panduan Koneksi DBeaver ke Railway PostgreSQL

Karena Railway memisahkan akses jaringan, Anda harus menggunakan **Public Network** untuk menghubungkan DBeaver (dari laptop Anda) ke database.

## 1. Dapatkan Kredensial (Public Network)
1. Buka dashboard project Anda di [Railway](https://railway.app).
2. Klik pada service **PostgreSQL**.
3. Buka tab **Connect**.
4. **PENTING:** Klik tab **Public Network** (bukan Private Network).
5. Salin **Connection URL** yang terlihat seperti:
   `postgresql://postgres:PASSWORD@host.railway.net:PORT/railway`

   Atau ambil detail manualnya:
   - **Host**: (misal: `ballast.proxy.rlwy.net`)
   - **Port**: (misal: `24293`)
   - **User**: `postgres`
   - **Password**: (Klik tombol mata/copy untuk melihat)
   - **Database**: `railway`

## 2. Setup di DBeaver
1. Buka DBeaver.
2. Klik icon **New Database Connection** (colokan listrik dengan tanda plus) di pojok kiri atas.
3. Pilih **PostgreSQL** dan klik **Next**.
4. Isi form dengan data dari tab Public Network tadi:
   - **Host**: `ballast.proxy.rlwy.net` (sesuaikan dengan dashboard Anda)
   - **Port**: `24293` (sesuaikan dengan dashboard Anda)
   - **Database**: `railway`
   - **Username**: `postgres`
   - **Password**: Paste password yang Anda salin.
5. Klik **Test Connection**.
   - Jika berhasil, akan muncul popup "Connected".
6. Klik **Finish**.

## 3. Catatan untuk Backend
Untuk aplikasi Backend Go yang berjalan **di dalam Railway**, jangan gunakan URL publik ini. Gunakan variabel internal agar koneksi lebih cepat dan aman:
- Di menu **Variables** backend service Anda, gunakan `${{ Postgres.DATABASE_URL }}`.
