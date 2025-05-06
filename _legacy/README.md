# Haxorport Client

Haxorport Client adalah aplikasi command-line untuk membuat tunnel HTTP dan TCP yang memungkinkan Anda mengekspos layanan lokal ke internet melalui server Haxorport.

## Fitur

- Tunneling HTTP dengan subdomain kustom
- Tunneling TCP dengan port remote
- Autentikasi tunnel (basic auth dan header)
- Konfigurasi melalui file YAML atau flag command line
- Logging dengan berbagai level (debug, info, warn, error)
- Reconnect otomatis saat koneksi terputus
- Manajemen tunnel melalui CLI

## Instalasi

### Dari Source

1. Clone repository:
   ```bash
   git clone https://github.com/haxorport/client.git
   cd client
   ```

2. Build aplikasi:
   ```bash
   go build -o hxp ./cmd/hxp
   ```

3. (Opsional) Pindahkan executable ke direktori dalam PATH:
   ```bash
   sudo mv hxp /usr/local/bin/
   ```

## Penggunaan

### Konfigurasi

Haxorport Client dapat dikonfigurasi melalui file konfigurasi YAML atau flag command line. Secara default, aplikasi akan mencari file konfigurasi di lokasi berikut:

1. Path yang ditentukan dengan flag `--config`
2. `~/.haxorport/config.yaml`
3. `./config.yaml`

Contoh file konfigurasi:

```yaml
# Alamat server haxorport
server_address: "localhost"

# Port untuk control plane
control_port: 8080

# Token untuk autentikasi ke server
auth_token: "rahasia"

# Level logging (debug, info, warn, error)
log_level: "info"

# Path ke file log (kosong untuk stdout)
log_file: ""

# Daftar tunnel yang akan dibuat saat startup
tunnels:
  # Contoh tunnel HTTP
  - name: "web-app"
    type: "http"
    local_port: 3000
    subdomain: "webapp"
    auth:
      type: "basic"
      username: "user"
      password: "pass"

  # Contoh tunnel TCP
  - name: "database"
    type: "tcp"
    local_port: 5432
    remote_port: 12345
```

### Command Line Interface

Haxorport Client menyediakan antarmuka command-line yang user-friendly dengan berbagai subcommand:

#### Membuat Tunnel HTTP

```bash
hxp http 8080 --subdomain=myapp
```

Opsi tambahan:
- `--subdomain`: Subdomain untuk tunnel HTTP (opsional)
- `--auth-user` dan `--auth-pass`: Kredensial untuk autentikasi basic
- `--auth-header-name` dan `--auth-header-value`: Header untuk autentikasi header

#### Membuat Tunnel TCP

```bash
hxp tcp 5432 --remote-port=12345
```

Opsi tambahan:
- `--remote-port`: Port remote untuk tunnel TCP (opsional)

#### Melihat Status Tunnel

```bash
hxp status
```

#### Mengelola Konfigurasi

```bash
# Melihat konfigurasi
hxp config show

# Mengatur nilai konfigurasi
hxp config set server_address example.com
hxp config set control_port 8080
hxp config set auth_token my-secret-token

# Mendapatkan nilai konfigurasi
hxp config get server_address
```

### Opsi Global

Opsi berikut dapat digunakan dengan semua subcommand:

- `--config`: Path ke file konfigurasi
- `--server`: Alamat server haxorport
- `--port`: Port control plane server
- `--token`: Token autentikasi
- `--log-level`: Level logging (debug, info, warn, error)
- `--log-file`: Path ke file log (kosong untuk stdout)

## Contoh Penggunaan

### Mengekspos Aplikasi Web Lokal

```bash
# Aplikasi web berjalan di localhost:3000
hxp http 3000 --subdomain=myapp
```

Setelah perintah di atas dijalankan, aplikasi web Anda akan tersedia di `http://myapp.haxorport.local` (atau domain yang dikonfigurasi di server).

### Mengekspos Server Database

```bash
# Server PostgreSQL berjalan di localhost:5432
hxp tcp 5432 --remote-port=12345
```

Setelah perintah di atas dijalankan, server database Anda akan tersedia di `haxorport.local:12345` (atau domain yang dikonfigurasi di server).

## Pengembangan

### Struktur Proyek

```
client/
├── cmd/
│   └── hxp/
│       └── main.go           # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── config.go         # Manajemen konfigurasi
│   ├── logger/
│   │   └── logger.go         # Logging
│   ├── proto/
│   │   └── message.go        # Protokol komunikasi
│   └── tunnel/
│       ├── client.go         # Klien WebSocket
│       ├── http.go           # Implementasi tunnel HTTP
│       ├── tcp.go            # Implementasi tunnel TCP
│       └── tunnel.go         # Implementasi tunnel dasar
├── config.example.yaml       # Contoh file konfigurasi
├── go.mod                    # Dependensi Go
└── README.md                 # Dokumentasi
```

### Kompilasi

```bash
# Kompilasi untuk platform saat ini
go build -o hxp ./cmd/hxp

# Kompilasi untuk Linux
GOOS=linux GOARCH=amd64 go build -o hxp-linux-amd64 ./cmd/hxp

# Kompilasi untuk Windows
GOOS=windows GOARCH=amd64 go build -o hxp-windows-amd64.exe ./cmd/hxp

# Kompilasi untuk macOS
GOOS=darwin GOARCH=amd64 go build -o hxp-macos-amd64 ./cmd/hxp
```

## Lisensi

MIT
