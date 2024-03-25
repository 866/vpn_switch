package main

import (
	"log"
	"net/http"

	"github.com/866/vpn_switch/handlers"
)

// The server's port to be listening on
const Port = ":80"

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Serve favicon
	http.HandleFunc("/favicon.ico", faviconHandler)

	// Handle main entry
	http.HandleFunc("/", handlers.HandleRoot)

	// Handle VPN requests
	http.HandleFunc("/vpn", handlers.HandleVPN)

	// Define usb off route
	http.HandleFunc("/usboff", handlers.HandleUSB)

	// Define vpn configuration upload requests
	http.HandleFunc("/upload", handlers.UploadHandler)

	// Define reboot request
	http.HandleFunc("/reboot", handlers.HandleReboot)

	// Define shutdown request
	http.HandleFunc("/shutdown", handlers.HandleShutdown)

	// Login handlers
	http.HandleFunc("/login", handlers.Signin)
	http.HandleFunc("/signup", handlers.Signup)
	http.HandleFunc("/refresh", handlers.Refresh)
	http.HandleFunc("/logout", handlers.Logout)

	// Run the server
	log.Printf("Server listening on port %v...", Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}
