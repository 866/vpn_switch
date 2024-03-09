package handlers

import (
	"net/http"
	"html/template"
	"log"
)

// Template stucture
type mainPageInfo struct {
	Checkbox string
}

// Define the templates which are inititialized
// when the package is loaded
var (
    HomeTemplate *template.Template
    UploadTemplate *template.Template
)

func init() {
	// Load the templates
	HomeTemplate = template.Must(template.ParseFiles("./templates/index.html"))
	UploadTemplate = template.Must(template.ParseFiles("./templates/upload.html"))
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		log.Printf("Wrong address is accessed: %v", r.URL.Path)
		return
	}
	
	log.Println("Main page is entered.")
	data := mainPageInfo{""}
	// Check status to change the slider position
	if vpnCommandStatus() {
		data.Checkbox = "checked"
	}
	HomeTemplate.Execute(w, data)
}