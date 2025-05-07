# Haxor Client

Haxor Client adalah aplikasi klien untuk layanan Haxorport, yang memungkinkan Anda membuat tunnel HTTP dan TCP untuk mengekspos layanan lokal ke internet.

## Fitur

- Tunnel HTTP/HTTPS: Ekspos layanan web lokal dengan subdomain kustom, mendukung protokol HTTP dan HTTPS
- Tunnel TCP: Ekspos layanan TCP lokal dengan port remote
- Autentikasi: Lindungi tunnel dengan autentikasi basic atau header
- Konfigurasi: Kelola konfigurasi dengan mudah melalui CLI
- Reconnect otomatis: Koneksi akan otomatis terhubung kembali jika terputus

## Arsitektur

Haxor Client dibangun dengan arsitektur heksagonal (ports and adapters) yang memisahkan domain bisnis dari infrastruktur teknis. Arsitektur ini memungkinkan:

1. **Pemisahan Concern**: Domain bisnis terpisah dari detail teknis
2. **Testability**: Komponen dapat diuji secara terpisah
3. **Fleksibilitas**: Implementasi infrastruktur dapat diganti tanpa mengubah domain bisnis

Struktur proyek:

```
haxor-client/
├── cmd/                    # Command-line interface
├── internal/               # Kode internal
│   ├── domain/             # Domain layer
│   │   ├── model/          # Model domain
│   │   └── port/           # Port (interface)
│   ├── application/        # Application layer
│   │   └── service/        # Service
│   ├── infrastructure/     # Infrastructure layer
│   │   ├── config/         # Implementasi konfigurasi
│   │   ├── transport/      # Implementasi komunikasi
│   │   └── logger/         # Implementasi logger
│   └── di/                 # Dependency injection
├── scripts/                # Script untuk build dan run
└── main.go                 # Entry point
```

## Instalasi

### Dari Source

1. Clone repositori:
   ```
   git clone https://github.com/haxorport/haxor-client.git
   cd haxor-client
   ```

2. Build aplikasi:
   ```
   ./scripts/build.sh
   ```

3. (Opsional) Pindahkan binary ke direktori dalam PATH:
   ```
   sudo cp bin/haxor /usr/local/bin/
   ```

### Dari Binary

1. Download binary terbaru dari [releases](https://github.com/haxorport/haxor-client/releases)
2. Ekstrak dan pindahkan ke direktori dalam PATH

## Penggunaan

### Konfigurasi

Sebelum menggunakan Haxor Client, Anda perlu mengatur konfigurasi:

```
haxor config set server_address haxorport.com
haxor config set control_port 8080
haxor config set auth_token your-auth-token
```

Untuk melihat konfigurasi saat ini:

```
haxor config show
```

### Tunnel HTTP

Membuat tunnel HTTP untuk layanan web lokal:

```
haxor http --port 8080 --subdomain myapp
```

Dengan autentikasi basic:

```
haxor http --port 8080 --subdomain myapp --auth basic --username user --password pass
```

Dengan autentikasi header:

```
haxor http --port 8080 --subdomain myapp --auth header --header "X-API-Key" --value "secret-key"
```

### Tunnel HTTPS

Haxorport sekarang mendukung tunnel HTTPS secara otomatis. Ketika klien terhubung ke server, server akan mendeteksi apakah permintaan datang melalui HTTP atau HTTPS dan meneruskan informasi ini ke klien. Klien kemudian akan menggunakan skema yang sesuai saat membuat permintaan ke layanan lokal.

Untuk menggunakan tunnel HTTPS:

1. Pastikan server haxorport dikonfigurasi dengan benar untuk mendukung HTTPS
2. Jalankan klien seperti biasa:
   ```
   haxorport http http://localhost:9090 -c config.yaml
   ```
3. Akses layanan Anda melalui HTTPS:
   ```
   https://your-subdomain.haxorport.online
   ```

### Solusi Sementara dengan SSH Tunnel

Jika Anda mengalami masalah koneksi dengan tunnel HTTPS, Anda dapat menggunakan SSH tunnel sebagai solusi sementara:

```bash
ssh -R <remote_port>:localhost:<local_port> -i <path_to_key> -N -f user@server
```

Contoh:
```bash
ssh -R 9090:localhost:9090 -i /path/to/key.pem -N -f root@example.com
```

Ini akan meneruskan permintaan yang diterima di port 9090 di server ke port 9090 di komputer lokal Anda.

### Tunnel TCP

Membuat tunnel TCP untuk layanan TCP lokal:

```
haxor tcp --port 22 --remote-port 2222
```

Jika `--remote-port` tidak ditentukan, server akan menetapkan port remote secara otomatis.

### Menambahkan Tunnel ke Konfigurasi

Anda dapat menambahkan tunnel ke konfigurasi untuk digunakan nanti:

```
haxor config add-tunnel --name web --type http --port 8080 --subdomain myapp
haxor config add-tunnel --name ssh --type tcp --port 22 --remote-port 2222
```

## Pengembangan

### Prasyarat

- Go 1.21 atau lebih baru
- Git

### Setup Pengembangan

1. Clone repositori:
   ```
   git clone https://github.com/haxorport/haxor-client.git
   cd haxor-client
   ```

2. Install dependensi:
   ```
   go mod download
   ```

3. Jalankan aplikasi dalam mode pengembangan:
   ```
   ./scripts/run.sh
   ```

### Struktur Kode

- **Domain Layer**: Berisi model domain dan port (interface)
- **Application Layer**: Berisi service yang mengimplementasikan use case
- **Infrastructure Layer**: Berisi implementasi konkret dari port
- **CLI Layer**: Berisi command-line interface menggunakan Cobra
- **DI Layer**: Berisi container untuk dependency injection

## Lisensi

MIT License
