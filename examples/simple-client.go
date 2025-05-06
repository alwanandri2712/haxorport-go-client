package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Hubungkan ke server
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		fmt.Printf("Gagal terhubung ke server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Terhubung ke server localhost:9090")

	// Buat reader untuk membaca input dari terminal
	reader := bufio.NewReader(os.Stdin)

	// Buat reader untuk membaca respons dari server
	serverReader := bufio.NewReader(conn)

	// Goroutine untuk membaca respons dari server
	go func() {
		for {
			response, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("Koneksi terputus")
				os.Exit(1)
			}
			fmt.Print("Server: ", response)
		}
	}()

	// Loop utama untuk mengirim pesan
	for {
		fmt.Print("Pesan: ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "exit" {
			break
		}

		// Kirim pesan ke server
		fmt.Fprintf(conn, "%s\n", message)
	}
}
