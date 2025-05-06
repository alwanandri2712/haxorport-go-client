package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// Hubungkan ke server WebSocket
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:9090/ws", nil)
	if err != nil {
		log.Fatalf("Gagal terhubung ke server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Terhubung ke server WebSocket di ws://localhost:9090/ws")

	// Channel untuk menangani sinyal interupsi
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Channel untuk mengirim pesan
	done := make(chan struct{})

	// Goroutine untuk membaca pesan dari server
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Gagal membaca pesan: %v", err)
				return
			}
			fmt.Printf("Server: %s\n", message)
		}
	}()

	// Goroutine untuk mengirim ping secara berkala
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	// Reader untuk membaca input dari terminal
	reader := bufio.NewReader(os.Stdin)

	// Loop utama
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// Kirim ping
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("Gagal mengirim ping: %v", err)
				return
			}
		case <-interrupt:
			// Kirim pesan penutupan
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("Gagal mengirim pesan penutupan: %v", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		default:
			// Baca input dari terminal
			fmt.Print("Pesan: ")
			message, _ := reader.ReadString('\n')
			message = message[:len(message)-1] // Hapus newline

			if message == "exit" {
				return
			}

			// Kirim pesan ke server
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Printf("Gagal mengirim pesan: %v", err)
				return
			}
		}
	}
}
