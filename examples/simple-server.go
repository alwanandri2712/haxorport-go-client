package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	// Mulai server TCP
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Gagal memulai server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server berjalan di :9090")

	for {
		// Terima koneksi
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Gagal menerima koneksi: %v", err)
			continue
		}

		// Tangani koneksi dalam goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Klien terhubung: %s\n", conn.RemoteAddr())

	// Buat reader
	reader := bufio.NewReader(conn)

	for {
		// Baca pesan
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Klien terputus: %s\n", conn.RemoteAddr())
			break
		}

		// Tampilkan pesan
		message = strings.TrimSpace(message)
		fmt.Printf("Pesan dari %s: %s\n", conn.RemoteAddr(), message)

		// Kirim balasan
		response := fmt.Sprintf("Server menerima: %s\n", message)
		conn.Write([]byte(response))
	}
}
