# domjudge-helper

Alat command-line untuk batch upload test case ke DomJudge menggunakan headless browser automation.

A command-line tool for batch uploading test cases to DomJudge using headless browser automation.

---

## Daftar Isi / Table of Contents

* [Bahasa Indonesia](#bahasa-indonesia)
* [English](#english)

---

# Bahasa Indonesia

## Features

* Auto-discovery problem dari folder test case
* Upload paralel antar problem
* Upload sequential per test case dalam satu problem
* Live progress rendering ala Docker Compose
* Binary portable tanpa dependensi tambahan
* Tidak memerlukan Node.js
* Tidak perlu install browser manual
* Chromium di-download otomatis saat pertama dijalankan
* Build multi-platform

---

## Persyaratan

* Go 1.21+ (hanya jika build dari source)

Jika menggunakan binary prebuilt:

* Tidak perlu install Go
* Tidak perlu install Node.js
* Tidak perlu install browser

Chromium akan otomatis di-download dan disimpan ke local cache saat pertama dijalankan.

---

## Instalasi

### Dari Binary (Direkomendasikan)

Ambil binary yang sesuai platform dari folder `bin/`, lalu letakkan bersama:

* `config.json`
* folder `test-cases/`

Kemudian jalankan langsung.

### Dari Source

```bash
git clone https://github.com/helscape/domjudge-helper-go.git
cd domjudge-helper-go
go mod tidy
task build
```

---

## Konfigurasi

Salin config contoh:

```bash
cp config.example.json config.json
```

Isi `config.json`:

```json
{
  "server_url": "https://domjudge.example.com",
  "admin_username": "admin",
  "admin_password": "yourpassword",
  "testcases_dir": "./test-cases",
  "concurrency": 4,
  "retry_count": 3,
  "retry_delay_seconds": 2,
  "timeout_seconds": 30
}
```

| Field                 | Keterangan                          |
| --------------------- | ----------------------------------- |
| `server_url`          | Base URL DomJudge                   |
| `testcases_dir`       | Direktori folder test case          |
| `concurrency`         | Jumlah upload paralel antar problem |
| `retry_count`         | Jumlah retry per test case          |
| `retry_delay_seconds` | Delay antar retry                   |
| `timeout_seconds`     | Timeout browser action              |

---

## Cara Kerja

* Setiap subfolder di dalam `test-cases/` dianggap sebagai ID problem DomJudge
* Problem di-upload secara paralel
* Test case di-upload sequential dalam satu problem
* Upload dilakukan menggunakan headless Chromium automation

---

## Struktur Folder Test Case

Nama folder harus sesuai dengan ID problem di DomJudge.

```text
test-cases/
  14/
    01
    01.a
    02
    02.a

  15/
    01
    01.a
```

* File input: `01`, `02`, `03`, ...
* File output: `01.a`, `02.a`, `03.a`, ...

---

## Penggunaan

```bash
# Menggunakan config.json default
./domjudge-helper

# Menggunakan config custom
./domjudge-helper -config /path/to/config.json
```

---

## Build

Memerlukan [Task](https://taskfile.dev):

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

```bash
# Build platform saat ini
task build

# Build semua platform
task all-platforms

# Build platform tertentu
task build-windows
task build-linux
task build-mac-arm
task build-mac-amd
```

Output binary tersedia di folder `bin/`.

---

## Contoh Output

```text
╔═══════════════════════════════════════════════════════════════╗
║                 DOMjudge Test Case Uploader                   ║
╚═══════════════════════════════════════════════════════════════╝
  Server  : https://domjudge.example.com
  Mode    : Browser (Headless Chromium)
  Source  : ./test-cases

Scanning for test cases...
  Found 2 problem(s).

Resolving problem names...
  + [7] -> Asa Dengan Koleksi Angkanya
  + [8] -> FaSHioN

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  [Asa Dengan Koleksi Angkanya][7]
  Uploading test case 4/11

  [FaSHioN][8]
  Uploading test case 12/25

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## Troubleshooting

### Nama problem hanya muncul sebagai ID

DomJudge tidak dapat mengambil nama problem. Upload tetap berjalan normal menggunakan ID problem sebagai label.

### Upload gagal atau timeout

Cek:

* koneksi ke server DomJudge
* kredensial admin/jury
* nilai `timeout_seconds`

Naikkan `retry_count` jika koneksi tidak stabil.

### Browser gagal dijalankan

Chromium mungkin gagal di-download saat pertama dijalankan.

Coba:

* cek koneksi internet
* hapus folder cache `rod-browser`
* jalankan ulang aplikasi

---

## Contributing

Pull request dan issue terbuka.

Untuk perubahan besar:

1. buka issue terlebih dahulu
2. diskusikan pendekatan implementasi
3. buat pull request setelah disepakati

---

# English

## Features

* Automatic problem discovery from testcase folders
* Parallel uploads between problems
* Sequential testcase uploads within a problem
* Docker Compose-style live progress rendering
* Portable standalone binary
* No Node.js required
* No manual browser installation required
* Chromium downloaded automatically on first run
* Cross-platform builds

---

## Requirements

* Go 1.21+ (only required for building from source)

If using a prebuilt binary:

* No Go installation required
* No Node.js required
* No browser installation required

Chromium is automatically downloaded and cached on first run.

---

## Installation

### From Binary (Recommended)

Take the appropriate binary from the `bin/` directory and place it alongside:

* `config.json`
* `test-cases/`

Then run directly.

### From Source

```bash
git clone https://github.com/helscape/domjudge-helper-go.git
cd domjudge-helper-go
go mod tidy
task build
```

---

## Configuration

Copy the example config:

```bash
cp config.example.json config.json
```

Example `config.json`:

```json
{
  "server_url": "https://domjudge.example.com",
  "admin_username": "admin",
  "admin_password": "yourpassword",
  "testcases_dir": "./test-cases",
  "concurrency": 4,
  "retry_count": 3,
  "retry_delay_seconds": 2,
  "timeout_seconds": 30
}
```

| Field                 | Description                                 |
| --------------------- | ------------------------------------------- |
| `server_url`          | Base URL of the DomJudge instance           |
| `testcases_dir`       | Directory containing testcase folders       |
| `concurrency`         | Number of parallel uploads between problems |
| `retry_count`         | Retry attempts per testcase                 |
| `retry_delay_seconds` | Delay between retries                       |
| `timeout_seconds`     | Browser action timeout                      |

---

## How It Works

* Each subfolder inside `test-cases/` is treated as a DomJudge problem ID
* Problems are uploaded concurrently
* Testcases within a problem are uploaded sequentially
* Uploads are performed using headless Chromium automation

---

## Test Case Folder Structure

Folder names must match the DomJudge problem ID.

```text
test-cases/
  14/
    01
    01.a
    02
    02.a

  15/
    01
    01.a
```

* Input files: `01`, `02`, `03`, ...
* Output files: `01.a`, `02.a`, `03.a`, ...

---

## Usage

```bash
# Using default config.json
./domjudge-helper

# Using a custom config path
./domjudge-helper -config /path/to/config.json
```

---

## Build

Requires [Task](https://taskfile.dev):

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

```bash
# Build current platform
task build

# Build all platforms
task all-platforms

# Build specific platforms
task build-windows
task build-linux
task build-mac-arm
task build-mac-amd
```

Generated binaries are available in the `bin/` directory.

---

## Sample Output

```text
╔═══════════════════════════════════════════════════════════════╗
║                 DOMjudge Test Case Uploader                   ║
╚═══════════════════════════════════════════════════════════════╝
  Server  : https://domjudge.example.com
  Mode    : Browser (Headless Chromium)
  Source  : ./test-cases

Scanning for test cases...
  Found 2 problem(s).

Resolving problem names...
  + [7] -> Asa Dengan Koleksi Angkanya
  + [8] -> FaSHioN

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  [Asa Dengan Koleksi Angkanya][7]
  Uploading test case 4/11

  [FaSHioN][8]
  Uploading test case 12/25

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## Troubleshooting

### Problem names only appear as IDs

DomJudge could not resolve problem names. Uploads will still continue normally using the problem ID as labels.

### Uploads fail or timeout

Check:

* server connectivity
* admin/jury credentials
* `timeout_seconds` configuration

Increase `retry_count` for unstable connections.

### Browser failed to start

Chromium may have failed to download during first run.

Try:

* checking internet connectivity
* deleting the `rod-browser` cache folder
* restarting the application

---

## Contributing

Pull requests and issues are welcome.

For major changes:

1. open an issue first
2. discuss the implementation approach
3. submit a pull request afterward

---

Made With 🩷 by Fabian A. (Helscape)
