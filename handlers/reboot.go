package handlers

import (
	"log"
	"net/http"
	"encoding/json"
	"os/exec"
)

// Handles /reboot route
func HandleReboot(w http.ResponseWriter, r *http.Request) {
	// Handles a USB request
	log.Println("Send command \"Reboot\"")
	err := cmdReboot()
	// Create a response
	response := make(map[string]string)
	if err == nil {
		response["message"] = "Reboot has been executed successfully."
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

// Sends the reboot command
func cmdReboot() (err error) {
	// Runs usb networking off command in shell
	cmd := exec.Command("reboot")
	err = cmd.Run()
	if err != nil {
		log.Println("The error occured while running 'reboot", err)
	}
	return
}