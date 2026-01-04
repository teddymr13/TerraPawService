# Panduan Koneksi DBeaver ke Railway PostgreSQL

Setelah Anda berhasil deploy database PostgreSQL di Railway, ikuti langkah ini untuk menghubungkannya ke DBeaver:

## 1. Dapatkan Kredensial dari Railway
1. Buka dashboard project Anda di [Railway](https://railway.app).
2. Klik pada service **PostgreSQL**.
3. Buka tab **Connect**.
4. Anda akan melihat informasi seperti:
   - **Host**: (contoh: `containers-us-west-123.railway.app`)
   - **Port**: (contoh: `6543`)
   - **User**: `postgres`
   - **Password**: (string acak panjang)
   - **Database**: `railway`

## 2. Setup di DBeaver
1. Buka DBeaver.
2. Klik icon **New Database Connection** (colokan listrik dengan tanda plus) di pojok kiri atas.
3. Pilih **PostgreSQL** dan klik **Next**.
4. Isi form dengan data dari Railway:
   - **Host**: Copy dari field Host Railway.
   - **Port**: Copy dari field Port Railway.
   - **Database**: `railway` (atau nama database Anda).
   - **Username**: `postgres`.
   - **Password**: Copy password dari Railway.
5. Klik **Test Connection**.
   - Jika berhasil, akan muncul popup "Connected".
6. Klik **Finish**.

## 3. Mengelola Data
Sekarang database Railway Anda akan muncul di sidebar DBeaver. Anda bisa:
- Melihat tabel yang dibuat oleh aplikasi Go Anda.
- Menjalankan query SQL manual.
- Mengedit data secara langsung jika diperlukan untuk testing.
