package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// Template stucture
type mainPageInfo struct {
	Checkbox string
}

// Define the templates which are inititialized
// when the package is loaded
var (
	HomeTemplate   *template.Template
	UploadTemplate *template.Template
)

func init() {
	// Load the templates
	HomeTemplate = template.Must(template.ParseFiles("./templates/index.html"))
	UploadTemplate = template.Must(template.ParseFiles("./templates/upload.html"))
}

// Checks if the user is logged in and redirects to the login page if necessary
func CheckLoginAndRedirect(w http.ResponseWriter, r *http.Request) bool {
	log.Println("Check if the user is logged in.")
	loggedin, err := UserIsLoggedIn(r)
	// Check for the unexpected error
	if err != http.ErrNoCookie {
		log.Println("Unexpected error occurred while loggin in: ", err)
		return false
	}
	// Check if user is not logged int
	if !loggedin {
		log.Println("The user is not logged in. Redirecting")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}
	log.Println("The user is logged in.")
	return true
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		log.Printf("Wrong address is accessed: %v", r.URL.Path)
		return
	}
	// Check if the user is logged in
	if !CheckLoginAndRedirect(w, r) {
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
