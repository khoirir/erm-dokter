# Panduan Pengembang & Arsitektur (AGENTS.md)

Dokumen ini adalah referensi utama untuk AI Agent dan pengembang yang melanjutkan pengembangan repositori **`erm-dokter`**.

---

## 1. Ikhtisar Proyek & Arsitektur

* **Bahasa & Runtime**: Go (Golang) 1.22+
* **Routing**: Go Standard Library `net/http` (`http.ServeMux` dengan method routing Go 1.22+).
* **Database**: MySQL / MariaDB (Database SIMRS Khanza).
* **Live Reload**: `air` (sudah mendukung ekstensi `.go`, `.yaml`, `.yml`, `.json`, `.html`, `.env`).
* **Dokumentasi API**: Scalar API Reference terintegrasi di `GET /docs` (Spesifikasi OpenAPI 3.1.0 di `internal/docs/openapi.yaml`).

### Struktur Lapisan (Clean Architecture & DI)
```
internal/
├── auth/           # Login dokter, token JWT, profile
├── master/         # Master data (Penjamin, Depo Farmasi)
├── rawatjalan/     # Antrean pasien dokter, detail & riwayat kunjungan
├── pemeriksaan/    # SOAP & TTV (Ralan & Ranap), riwayat, CRUD, validasi 48 jam
├── obat/           # Pencarian & detail data master obat per depo
├── resep/          # Resep obat dokter (non-racikan, racikan), master aturan pakai & metode racik
├── middleware/     # Auth JWT Middleware (Bearer), Timeout Middleware
├── routes/         # Centralized HTTP route multiplexer
├── di/             # Manual Dependency Injection (providers)
├── pkg/            # Utilitas inti: crypto (AES-GCM), token, database, logger, response
└── shared/         # Model bersama, enum StatusLanjut, validasi waktu & validator
```

---

## 2. Standar Desain & Konvensi Penting

1. **Keamanan & Identifier URL**:
   - Seluruh identifier publik (`id_kunjungan`, `id_pasien`, `id_pemeriksaan`, `id_obat`) **wajib dienkripsi** menggunakan AES-GCM URL-safe (`internal/pkg/crypto`).
   - ID internal database (`no_rawat`, `no_rkm_medis`, `kode_brng`, dll.) tidak boleh diekspos mentah ke client.

2. **Format Response Standar (`internal/pkg/response`)**:
   - Sukses: `response.Success(w, message, data)`
   - Sukses dengan Paginasi: `response.SuccessWithMeta(w, message, data, meta)`
   - Error: Ditangani terpusat melalui `apperror.HandleError(w, err)`.

3. **Pemisahan Validasi (Handler vs Service)**:
   - **Layer Handler**: Validasi format input, sanitasi query param, pengecekan enum `status_lanjut`, decrypt URL token, dan validasi non-database.
   - **Layer Service**: Logika bisnis murni, validasi aturan medis (misal batas 48 jam, kepemilikan dokter), dan pemanggilan repository.

4. **Pencarian Master Data**:
   - Master yang rawan variasi spasi seperti `aturan_pakai` menggunakan pencarian *space-insensitive*:
     `WHERE REPLACE(aturan, ' ', '') LIKE ?` dan `strings.ReplaceAll(keyword, " ", "")`.

---

## 3. Fitur yang Telah Selesai Diimplementasikan

### A. Auth (`/api/v1/auth`)
- `POST /api/v1/auth/login` (Verifikasi password AES Khanza untuk dokter).
- `GET /api/v1/auth/profile` (Profile dokter dari token).

### B. Rawat Jalan (`/api/v1/rawat-jalan`)
- `GET /api/v1/rawat-jalan/antrean` (Daftar antrean dokter dengan filter tanggal, keyword min 3 char, pagination, sorting).
- `GET /api/v1/rawat-jalan/kunjungan/{id_kunjungan}` (Detail kunjungan pasien).
- `GET /api/v1/rawat-jalan/pasien/{id_pasien}/riwayat` (Riwayat kunjungan pasien by RM).

### C. Pemeriksaan Medis / SOAP (`/api/v1/pemeriksaan`)
- `GET /api/v1/pemeriksaan/kesadaran` (Referensi opsi kesadaran).
- `GET /api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}` (Riwayat SOAP kunjungan, filter tanggal, pagination).
- `GET /api/v1/pemeriksaan/pasien/{id_pasien}/{status_lanjut}` (Riwayat SOAP seluruh kunjungan by RM).
- `GET /api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}/{id_pemeriksaan}` (Detail SOAP).
- `POST /api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}` (Simpan SOAP baru).
- `PUT /api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}/{id_pemeriksaan}` (Update SOAP: proteksi 48 jam & dokter pembuat).
- `DELETE /api/v1/pemeriksaan/{id_kunjungan}/{status_lanjut}/{id_pemeriksaan}` (Hapus SOAP: proteksi 48 jam & dokter pembuat).

### D. Master & Obat (`/api/v1/master`, `/api/v1/obat`)
- `GET /api/v1/master/penjamin` (Daftar asuransi/penjamin).
- `GET /api/v1/master/depo` (Daftar depo obat).
- `GET /api/v1/obat` (Pencarian obat per depo, stok > 0, keyword min 3 char, pagination).
- `GET /api/v1/obat/{id_obat}` (Detail obat & harga).

### E. Resep Obat (`/api/v1/resep`)
- `GET /api/v1/resep/aturan-pakai` (Master aturan pakai, space-insensitive search).
- `GET /api/v1/resep/metode-racik` (Master metode racikan: PULV, CAPS, Salep, dll).
- `GET /api/v1/resep/{id_kunjungan}/{status_lanjut}` (Daftar resep per kunjungan, filter tanggal, pagination).
- `GET /api/v1/resep/pasien/{id_pasien}/{status_lanjut}` (Riwayat resep seluruh kunjungan by RM).
- Hierarki data: `Resep` &rarr; `ResepDokter` (non-racikan) & `ResepDokterRacikan` (header: `jumlah_racikan`, `kode_racik`, `metode_racik`) &rarr; `ResepDokterRacikanDetail` (`jumlah` bahan).

---

## 4. Keputusan Bisnis & Catatan Operasional RS (PENTING)

Berdasarkan diskusi mendalam mengenai perilaku operasional riil di Rumah Sakit & database Khanza:

1. **Masalah Billing Parsial Pasien Umum**:
   - Pasien umum yang membayar karcis di loket pendaftaran awal akan langsung memiliki status `reg_periksa.status_bayar = 'Sudah Bayar'`.
   - **Aturan**: `status_bayar == 'Sudah Bayar'` **TIDAK BOLEH** dijadikan indikator bahwa pasien sudah pulang/selesai periksa.

2. **Masalah Status Pelayanan (`reg_periksa.stts`)**:
   - Petugas/perawat di RS umumnya hanya mengubah status menjadi `'Sudah'` (sudah diperiksa dokter) dan jarang sekali mengupdate status saat pasien fisik pulang.
   - **Aturan**: Kolom `stts` tidak dapat dijadikan acuan valid kepulangan fisik pasien.

3. **Masalah Pasien IGD Lintas Hari (Tengah Malam)**:
   - Pasien yang masuk IGD jam 23:30 (tgl 26) dan baru diresepkan obat jam 02:00 dini hari (tgl 27) tidak boleh ditolak hanya karena beda tanggal kalender.
   - **Aturan**: Gunakan **Batas Waktu Relatif (*Time-Window*) 24 Jam** dari waktu registrasi (`tgl_registrasi` + `jam_reg`):
     $$\text{Waktu Sekarang} \le \text{Waktu Registrasi} + 24\text{ Jam}$$

4. **Validasi Status Rawat Inap Aktif vs Rawat Jalan**:
   - Jika pasien memiliki kamar aktif di tabel `kamar_inap` (`stts_pulang = '-'`), maka pasien berstatus **Active Inpatient**.
   - Resep yang dibuat **wajib berstatus `Ranap`** (resep berstatus `Ralan` ditolak agar tidak salah depo dan tidak terjadi klaim ganda BPJS).
   - Jika pasien sudah checkout dari rawat inap (`stts_pulang != '-'` dan ada `tgl_keluar`), kunjungan tersebut sudah selesai &rarr; resep baru ditolak.

5. **Proteksi Resep yang Telah Divalidasi Farmasi**:
   - Di Khanza, jika resep sudah diproses/diserahkan oleh apotek, kolom `resep_obat.tgl_perawatan != '0000-00-00'`.
   - Resep yang sudah divalidasi farmasi **terkunci permanen dari edit dan hapus oleh dokter**.

---

## 5. Roadmap / Tugas Selanjutnya

Tugas prioritas berikutnya yang akan dikerjakan:

1. **Simpan Resep Obat Baru (`POST /api/v1/resep/{id_kunjungan}/{status_lanjut}`)**:
   - Generate nomor resep otomatis (`no_resep` format harian Khanza: `YYYYMMDD0001`).
   - Menyimpan transaksi ke tabel:
     - `resep_obat` (header peresepan).
     - `resep_dokter` (daftar obat non-racikan).
     - `resep_dokter_racikan` (header racikan dengan `no_racik`, `jml_dr`, `kd_racik`, dll).
     - `resep_dokter_racikan_detail` (komposisi obat bahan racikan).
   - Menerapkan Database Transaction (`db.BeginTx`).
   - Menerapkan validasi bisnis:
     - Dokter pembuat diambil dari token JWT (`KodeDokter`).
     - Validasi batas waktu 24 jam dari waktu registrasi untuk Rawat Jalan/IGD.
     - Validasi status rawat inap aktif (tabel `kamar_inap`).
     - Validasi ketersediaan obat & depo.

2. **Hapus / Batalkan Resep Obat (`DELETE /api/v1/resep/{id_kunjungan}/{status_lanjut}/{id_resep}`)**:
   - Pengecekan kepemilikan dokter pembuat.
   - Pengecekan apakah resep sudah divalidasi farmasi (`resep_obat.tgl_perawatan != '0000-00-00'`). Jika sudah diproses apotek &rarr; Tolak penghapusan.

3. **Uji Coba & Dokumentasi OpenAPI**:
   - Melengkapi file `test.http` untuk skenario simpan resep non-racikan dan racikan.
   - Menambahkan skema request body di `internal/docs/openapi.yaml`.
