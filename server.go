package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

// The server's port to be listening on
const Port = ":8080"

type MainPageInfo struct {
	Checkbox string
}

func vpnCommandStatus() bool {
	// Check the VPN status
	// Run the shell command
	cmd := exec.Command("bash", "-c", "nmcli con show --active | grep ua-vpn")
	output, err := cmd.Output()
	// Convert output to string
	response := string(output)
	if err != nil && err.Error() != "exit status 1" {
		fmt.Println("Error:", err)
		return false
	}
	// Check if there is an active vpn connection in the list
	return !(strings.TrimSpace(response) == "")
}

func cmdVPN(command string) {
	// Runs vpn command in shell
	cmd := exec.Command("nmcli", "connection", command, "ua-vpn")
	err := cmd.Run()
	if err != nil {
		log.Printf("VPN %v NMCLI error output: %v", command, err)
		return
	}
}

func handleVPN(w http.ResponseWriter, r *http.Request) {
	// Check for other addresses
	if r.URL.Path != "/vpn" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		log.Printf("Wrong address is accessed: %v", r.URL.Path)
		return
	}

	// We expect JSON communication
	w.Header().Set("Content-Type", "application/json")
	// Handle GET and POST differently
	switch r.Method {

	// GET request
	case http.MethodGet:
		// Create a map representing the JSON response
		jsonResponse := map[string]bool{"vpn": true}

		// Encode the JSON response and write it to the response writer
		err := json.NewEncoder(w).Encode(jsonResponse)
		if err != nil {
			// If there's an error encoding JSON, return an internal server error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("GET from /vpn response: %v", jsonResponse)

	// POST request
	case http.MethodPost:

		// For POST requests, decode the JSON body
		var requestData map[string]bool
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Log the payload
		log.Printf("POST to /vpn with payload: %v", requestData)

		// Validate the received JSON
		if _, ok := requestData["vpn"]; !ok {
			http.Error(w, "Invalid JSON. 'vpn' field missing.", http.StatusBadRequest)
			return
		}

		// Respond with success message
		successResponse := map[string]string{"message": "Received VPN set command successfully"}
		if err := json.NewEncoder(w).Encode(successResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		// For unsupported methods, return Method Not Allowed status
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func cmdUSBOff() {
	// Runs usb networking off command in shell
	cmd := exec.Command("ifconfig", "usb0", "down")
	err := cmd.Run()
	if err != nil {
		log.Println("The error occured while running 'ifconfig usb0 down'", err)
		return
	}
}

func main() {

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Handle main entry
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Main page is entered.")
		data := MainPageInfo{""}
		// Check status to change the slider position
		if vpnCommandStatus() {
			data.Checkbox = "checked"
		}
		tmpl.Execute(w, data)
	})

	// Handle VPN requests
	http.HandleFunc("/vpn", handleVPN)

	// Define VPN manipulation routes
	http.HandleFunc("/vpnoff", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Send command \"VPN Off\"")
		fmt.Fprintf(w, "Sending a command to up a VPN...")
		cmdVPN("down")
	})

	http.HandleFunc("/vpnon", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Send command \"VPN On\"")
		fmt.Fprintf(w, "Sending a command to up a VPN...")
		cmdVPN("up")
	})

	// Define usb off route
	http.HandleFunc("/usboff", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Send command \"USB Off\"")
		fmt.Fprintf(w, "Sending a command to down a USB tethering...")
		cmdUSBOff()
	})

	// Run the server
	log.Printf("Server listening on port %v...", Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}
