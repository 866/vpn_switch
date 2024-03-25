package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// The path where to store the uploaded file
const vpnConfPath = "/etc/wireguard/wg0.conf"

// Display the named template
func display(w http.ResponseWriter, data interface{}) {
	UploadTemplate.Execute(w, data)
}

func createNewVPN(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, _, err := r.FormFile("myFile")
	if err != nil {
		log.Println("Error Retrieving the File: ", err)
		return
	}

	defer file.Close()

	// Create file
	dst, err := os.Create(vpnConfPath)
	if err != nil {
		log.Printf("Error when uploading file to %v. Error: %v", vpnConfPath, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("The file has been written to ", vpnConfPath)

	// Run the command to create the connection
	output, err := setNewVPN()
	if err != nil {
		fmt.Fprintf(w, "The connection can not be created. Command output: %v, Error: %v\n", output, err)
		return
	}
	fmt.Fprintf(w, "The new VPN connection has been created successfully\n")
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	if !CheckLoginAndRedirect(w, r) {
		return
	}
	// Handles /upload route
	switch r.Method {
	case http.MethodPost:
		// Uploads a file
		log.Println("Handling the file upload and creating a new connection.")
		createNewVPN(w, r)
	default:
		// Displays a page for uploading
		log.Println("Upload page is accessed.")
		display(w, nil)		
	}
}
