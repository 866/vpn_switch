package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func vpnCommandStatus() bool {
	// Check the VPN status
	// Run the shell command
	cmd := exec.Command("bash", "-c", "nmcli con show --active | grep ua-vpn")
	output, err := cmd.Output()
	// Convert output to string
	response := string(output)
	if err != nil && err.Error() != "exit status 1" {
		log.Println("Error when checking VPN connection status:", err)
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

func HandleVPN(w http.ResponseWriter, r *http.Request) {
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
