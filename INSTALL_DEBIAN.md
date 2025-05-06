# Panduan Instalasi Haxorport untuk Debian

Dokumen ini berisi panduan untuk menginstal dan menggunakan Haxorport pada sistem Linux Debian.

## Persyaratan Sistem

- Sistem operasi Linux Debian atau turunannya (Ubuntu, Linux Mint, dll.)
- Akses sudo
- Koneksi internet untuk mengunduh dependensi

## Cara Instalasi

1. Unduh file installer:
   ```
   curl -O https://raw.githubusercontent.com/alwanandri2712/haxorport-go-client/main/install_debian.sh
   ```

2. Berikan izin eksekusi pada file installer:
   ```
   chmod +x install_debian.sh
   ```

3. Jalankan installer:
   ```
   ./install_debian.sh
   ```

4. Ikuti petunjuk yang muncul di layar.

## Apa yang Dilakukan Installer

Installer akan melakukan hal-hal berikut:

1. Menginstal dependensi yang diperlukan (Go, Git, build-essential)
2. Mengkloning repositori haxorport-go-client
3. Mengkompilasi aplikasi
4. Menginstal aplikasi ke `/opt/haxorport`
5. Membuat symlink ke `/usr/local/bin/haxorport` agar dapat diakses dari mana saja

## Cara Penggunaan

Setelah instalasi selesai, Anda dapat menggunakan Haxorport dengan perintah berikut:

- Membuat HTTP tunnel:
  ```
  haxorport http http://localhost:9090
  ```

- Membuat TCP tunnel:
  ```
  haxorport tcp 22
  ```

- Menampilkan bantuan:
  ```
  haxorport --help
  ```

## Cara Menghapus Instalasi

1. Unduh file uninstaller:
   ```
   curl -O https://raw.githubusercontent.com/alwanandri2712/haxorport-go-client/main/uninstall_debian.sh
   ```

2. Berikan izin eksekusi pada file uninstaller:
   ```
   chmod +x uninstall_debian.sh
   ```

3. Jalankan uninstaller:
   ```
   ./uninstall_debian.sh
   ```

4. Ikuti petunjuk yang muncul di layar.

Atau, Anda dapat menghapus instalasi secara manual dengan perintah:
```
sudo rm -rf /opt/haxorport /usr/local/bin/haxorport
```

## Troubleshooting

### Masalah: Perintah `haxorport` tidak ditemukan

Pastikan bahwa symlink telah dibuat dengan benar:
```
ls -l /usr/local/bin/haxorport
```

Jika symlink tidak ada, buat secara manual:
```
sudo ln -sf /opt/haxorport/haxorport /usr/local/bin/haxorport
```

### Masalah: Aplikasi tidak dapat dikompilasi

Pastikan Go terinstal dengan benar:
```
go version
```

Jika Go tidak terinstal, instal secara manual:
```
sudo apt-get update
sudo apt-get install -y golang-go
```

### Masalah: Tidak dapat mengkloning repositori

Pastikan Git terinstal dengan benar:
```
git --version
```

Jika Git tidak terinstal, instal secara manual:
```
sudo apt-get update
sudo apt-get install -y git
```

## Informasi Tambahan

- Direktori instalasi: `/opt/haxorport`
- Binary aplikasi: `/opt/haxorport/bin/haxor`
- Script wrapper: `/opt/haxorport/haxorport`
- Symlink: `/usr/local/bin/haxorport`

## Lisensi

Haxorport dikembangkan oleh Alwan Putra Andriansyah dan terinspirasi oleh ngrok.
