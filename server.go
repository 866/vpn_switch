package main

import (
	"log"
	"net/http"

	"github.com/866/vpn_switch/handlers"
)

// The server's port to be listening on
const Port = ":80"


func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Handle main entry
	http.HandleFunc("/", handlers.HandleRoot)

	// Handle VPN requests
	http.HandleFunc("/vpn", handlers.HandleVPN)

	// Define usb off route
	http.HandleFunc("/usboff", handlers.HandleUSB)

	// Define vpn configuration upload requests
	http.HandleFunc("/upload", handlers.UploadHandler)

	// Run the server
	log.Printf("Server listening on port %v...", Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}
