package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
)

func HandleUSB(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	if !CheckLoginAndRedirect(w, r) {
		return
	}
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
