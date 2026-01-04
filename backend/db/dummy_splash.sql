-- Hapus data lama agar bersih
TRUNCATE TABLE splash_events RESTART IDENTITY;

-- 1. Natal & Tahun Baru 2025 (Aktif Sekarang: 27 Des 2025)
-- Icon: Christmas Cat (Transparent PNG)
INSERT INTO splash_events (event_name, image_url, start_date, end_date, is_active) 
VALUES ('Merry Christmas & Happy New Year', 'https://cdn-icons-png.flaticon.com/512/4514/4514868.png', '2024-12-20 00:00:00', '2026-01-05 23:59:59', TRUE);

-- 2. Imlek 2025 (Tahun Ular)
-- Icon: Snake/Chinese Dragon
INSERT INTO splash_events (event_name, image_url, start_date, end_date, is_active) 
VALUES ('Gong Xi Fa Cai 2025', 'https://cdn-icons-png.flaticon.com/512/10603/10603779.png', '2025-01-20 00:00:00', '2025-02-05 23:59:59', TRUE);

-- 3. Ramadhan 2025
-- Icon: Lantern/Ketupat
INSERT INTO splash_events (event_name, image_url, start_date, end_date, is_active) 
VALUES ('Marhaban Ya Ramadhan', 'https://cdn-icons-png.flaticon.com/512/4305/4305432.png', '2025-02-28 00:00:00', '2025-03-30 23:59:59', TRUE);

-- 4. Idul Fitri 2025
INSERT INTO splash_events (event_name, image_url, start_date, end_date, is_active) 
VALUES ('Selamat Hari Raya Idul Fitri', 'https://cdn-icons-png.flaticon.com/512/5024/5024252.png', '2025-03-31 00:00:00', '2025-04-07 23:59:59', TRUE);

-- 5. HUT RI ke-80
INSERT INTO splash_events (event_name, image_url, start_date, end_date, is_active) 
VALUES ('Dirgahayu Republik Indonesia ke-80', 'https://cdn-icons-png.flaticon.com/512/323/323315.png', '2025-08-10 00:00:00', '2025-08-20 23:59:59', TRUE);

-- 6. Halloween 2025
INSERT INTO splash_events (event_name, image_url, start_date, end_date, is_active) 
VALUES ('Spooky Season', 'https://cdn-icons-png.flaticon.com/512/2301/2301148.png', '2025-10-25 00:00:00', '2025-11-01 23:59:59', TRUE);
