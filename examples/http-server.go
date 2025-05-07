package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Menerima permintaan dari %s: %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
	fmt.Fprintf(w, "Hello from Haxorport HTTP Server!")
}

func main() {
	// Daftarkan handler HTTP
	http.HandleFunc("/", handleRoot)

	// Mulai server HTTP
	port := 9090
	fmt.Printf("Server HTTP berjalan di 0.0.0.0:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil); err != nil {
		log.Fatalf("Gagal memulai server: %v", err)
	}
}
