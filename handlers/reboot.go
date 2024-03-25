package handlers

import (
	"log"
	"net/http"
	"encoding/json"
	"os/exec"
)

// Handles /reboot route
func HandleReboot(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	if !CheckLoginAndRedirect(w, r) {
		return
	}
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

// Handles /shutdown route
func HandleShutdown(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	if !CheckLoginAndRedirect(w, r) {
		return
	}
	// Handles a USB request
	log.Println("Send command \"Shutdown -P 0\"")
	err := cmdShutdown()
	// Create a response
	response := make(map[string]string)
	if err == nil {
		response["message"] = "Shutdown has been executed successfully."
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
func cmdShutdown() (err error) {
	// Runs usb networking off command in shell
	cmd := exec.Command("shutdown", "-P", "0")
	err = cmd.Run()
	if err != nil {
		log.Println("The error occured while running 'reboot", err)
	}
	return
}