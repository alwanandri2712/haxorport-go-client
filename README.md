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

### Instalasi Mudah (Semua OS)

Gunakan script installer otomatis yang mendukung Linux, macOS, dan Windows (via WSL):

```bash
# Download dan jalankan installer
curl -sSL https://raw.githubusercontent.com/alwanandri2712/haxorport-go-client/main/install.sh | bash
```

Script installer akan:
- Mendeteksi OS Anda secara otomatis
- Menginstal dependensi yang diperlukan
- Mengkompilasi dan menginstal haxorport
- Membuat konfigurasi default

Setelah instalasi, Anda dapat langsung menggunakan perintah `haxorport`.

### Instalasi Manual

#### Dari Source

1. Clone repositori:
   ```bash
   git clone https://github.com/alwanandri2712/haxorport-go-client.git
   cd haxorport-go-client
   ```

2. Build aplikasi:
   ```bash
   # Pastikan Go sudah terinstal
   go build -o bin/haxor main.go
   ```

3. (Opsional) Pindahkan binary ke direktori dalam PATH:
   ```bash
   # Linux/macOS
   sudo cp bin/haxor /usr/local/bin/
   
   # Windows (PowerShell Admin)
   Copy-Item .\bin\haxor.exe -Destination "$env:ProgramFiles\haxorport\"
   ```

#### Dari Binary

1. Download binary terbaru dari [releases](https://github.com/alwanandri2712/haxorport-go-client/releases)
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

Haxorport sekarang mendukung tunnel HTTPS secara otomatis dengan arsitektur reverse connection. Ketika klien terhubung ke server, server akan mendeteksi apakah permintaan datang melalui HTTP atau HTTPS dan meneruskan permintaan tersebut ke klien melalui koneksi WebSocket. Klien kemudian akan membuat permintaan ke layanan lokal dan mengirim respons kembali ke server.

Keunggulan arsitektur reverse connection:

1. **Tidak memerlukan SSH tunnel**: Anda tidak perlu mengatur SSH tunnel untuk mengakses layanan lokal
2. **Penggantian URL otomatis**: URL lokal dalam respons HTML akan otomatis diganti dengan URL tunnel
3. **Dukungan HTTPS**: Akses layanan lokal melalui HTTPS tanpa perlu mengonfigurasi TLS di layanan lokal
4. **Subdomain kustom**: Gunakan subdomain yang mudah diingat untuk mengakses layanan lokal

Untuk menggunakan tunnel HTTPS:

1. Pastikan server haxorport dikonfigurasi dengan benar untuk mendukung HTTPS
2. Jalankan klien dengan menentukan port lokal dan subdomain:
   ```
   haxor-client http --port 8080 --subdomain myapp
   ```
3. Akses layanan Anda melalui HTTPS:
   ```
   https://myapp.haxorport.online
   ```

Semua link dan referensi dalam halaman web Anda akan otomatis diubah untuk menggunakan URL tunnel, sehingga navigasi di situs web berfungsi dengan benar.

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
   git clone https://github.com/alwanandri2712/haxorport-go-client.git
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
