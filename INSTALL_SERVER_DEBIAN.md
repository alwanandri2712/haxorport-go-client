# Panduan Instalasi Haxorport Server untuk Debian

Dokumen ini berisi panduan untuk menginstal dan menggunakan Haxorport Server pada sistem Linux Debian.

## Persyaratan Sistem

- Sistem operasi Linux Debian atau turunannya (Ubuntu, Linux Mint, dll.)
- Akses sudo
- Koneksi internet untuk mengunduh dependensi

## Cara Instalasi

1. Unduh file installer:
   ```
   curl -O https://raw.githubusercontent.com/alwanandri2712/haxorport-go-server/main/install_server_debian.sh
   ```

2. Berikan izin eksekusi pada file installer:
   ```
   chmod +x install_server_debian.sh
   ```

3. Jalankan installer:
   ```
   ./install_server_debian.sh
   ```

4. Ikuti petunjuk yang muncul di layar.

## Apa yang Dilakukan Installer

Installer akan melakukan hal-hal berikut:

1. Menginstal dependensi yang diperlukan (Go, Git, build-essential)
2. Mengkloning repositori haxorport-go-server
3. Mengkompilasi aplikasi server
4. Menginstal aplikasi server ke `/opt/haxorport-server`
5. Membuat service systemd untuk menjalankan server sebagai background service
6. Mengaktifkan dan menjalankan service

## Cara Penggunaan

Setelah instalasi selesai, Haxorport Server akan berjalan sebagai service di latar belakang. Anda dapat mengelola service dengan perintah berikut:

- Melihat status service:
  ```
  sudo systemctl status haxorport-server.service
  ```

- Menjalankan service:
  ```
  sudo systemctl start haxorport-server.service
  ```

- Menghentikan service:
  ```
  sudo systemctl stop haxorport-server.service
  ```

- Memulai ulang service:
  ```
  sudo systemctl restart haxorport-server.service
  ```

- Melihat log service:
  ```
  sudo journalctl -u haxorport-server.service
  ```

## Cara Menghapus Instalasi

1. Unduh file uninstaller:
   ```
   curl -O https://raw.githubusercontent.com/alwanandri2712/haxorport-go-server/main/uninstall_server_debian.sh
   ```

2. Berikan izin eksekusi pada file uninstaller:
   ```
   chmod +x uninstall_server_debian.sh
   ```

3. Jalankan uninstaller:
   ```
   ./uninstall_server_debian.sh
   ```

4. Ikuti petunjuk yang muncul di layar.

Atau, Anda dapat menghapus instalasi secara manual dengan perintah:
```
sudo systemctl stop haxorport-server.service
sudo systemctl disable haxorport-server.service
sudo rm -rf /opt/haxorport-server /etc/systemd/system/haxorport-server.service
sudo systemctl daemon-reload
```

## Troubleshooting

### Masalah: Service tidak dapat dijalankan

Periksa log service untuk informasi lebih lanjut:
```
sudo journalctl -u haxorport-server.service
```

Atau periksa log aplikasi:
```
sudo cat /opt/haxorport-server/logs/haxor-server.log
sudo cat /opt/haxorport-server/logs/haxor-server-error.log
```

### Masalah: Port yang digunakan sudah terpakai

Secara default, Haxorport Server menggunakan port 8081. Jika port ini sudah digunakan oleh aplikasi lain, Anda perlu mengubah konfigurasi.

1. Hentikan service:
   ```
   sudo systemctl stop haxorport-server.service
   ```

2. Edit file konfigurasi (jika ada):
   ```
   sudo nano /opt/haxorport-server/config.yaml
   ```

3. Jalankan kembali service:
   ```
   sudo systemctl start haxorport-server.service
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

## Informasi Tambahan

- Direktori instalasi: `/opt/haxorport-server`
- Binary aplikasi: `/opt/haxorport-server/bin/haxor-server`
- File service: `/etc/systemd/system/haxorport-server.service`
- Direktori log: `/opt/haxorport-server/logs`

## Lisensi

Haxorport dikembangkan oleh Alwan Putra Andriansyah dan terinspirasi oleh ngrok.
