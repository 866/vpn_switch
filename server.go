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
const Port = ":80"

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

func cmdVPN(command string) (err error) {
	// Runs vpn command in shell
	cmd := exec.Command("nmcli", "connection", command, "ua-vpn")
	err = cmd.Run()
	if err != nil {
		log.Printf("VPN %v NMCLI error output: %v", command, err)
	}
	return
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
		var vpnErr error
		if vpn_status, ok := requestData["vpn"]; !ok {
			http.Error(w, "Invalid JSON. 'vpn' field missing.", http.StatusBadRequest)
			return
		} else if vpn_status {
			vpnErr = cmdVPN("up")
		} else {
			vpnErr = cmdVPN("down")
		}

		// Respond with a message
		response := make(map[string]string)
		if vpnErr == nil {
			response["message"] = "Received VPN set command successfully."
		} else {
			response["message"] = "Error: " + vpnErr.Error()
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Errors with sending the response to /vpn. Error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("The response has been sent successfully: %v", response)
	default:
		// For unsupported methods, return Method Not Allowed status
		log.Printf("The HTTP method %v for /vpn is not allowed", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleUSB(w http.ResponseWriter, r *http.Request) {
	// Handles a USB request
	log.Println("Send command \"USB Off\"")
	err := cmdUSBOff()
	// Create a response
	response := make(map[string]string)
	if err == nil {
		response["message"] = "USB down has been executed successfully."
	} else {
		response["message"] = "Error: " + err.Error()
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Errors with sending the response%v. Error: %v", response, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("The response has been sent successfully: %v", response)
}

func cmdUSBOff() (err error) {
	// Runs usb networking off command in shell
	cmd := exec.Command("ifconfig", "usb0", "down")
	err = cmd.Run()
	if err != nil {
		log.Println("The error occured while running 'ifconfig usb0 down'", err)
	}
	return
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

	// Define usb off route
	http.HandleFunc("/usboff", handleUSB)

	// Run the server
	log.Printf("Server listening on port %v...", Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}
