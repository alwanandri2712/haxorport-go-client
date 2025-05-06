package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade koneksi HTTP ke WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Gagal upgrade koneksi: %v", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Klien terhubung: %s\n", conn.RemoteAddr())

	for {
		// Baca pesan
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Gagal membaca pesan: %v", err)
			break
		}

		// Tampilkan pesan
		fmt.Printf("Pesan dari %s: %s\n", conn.RemoteAddr(), message)

		// Kirim balasan
		response := fmt.Sprintf("Server menerima: %s", message)
		if err := conn.WriteMessage(messageType, []byte(response)); err != nil {
			log.Printf("Gagal mengirim pesan: %v", err)
			break
		}
	}

	fmt.Printf("Klien terputus: %s\n", conn.RemoteAddr())
}

func main() {
	// Daftarkan handler WebSocket
	http.HandleFunc("/ws", handleWebSocket)

	// Mulai server HTTP
	port := 9090
	fmt.Printf("Server WebSocket berjalan di :%d/ws\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Gagal memulai server: %v", err)
	}
}
