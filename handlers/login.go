package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Set the token expire value
const TokenExpireMin = 60 * 24 * 7

// Login token name to set the coookie
const LoginToken = "login_token"

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string
	Username string
}

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var sessions = map[string]session{}

// each session contains the username of the user and the time at which it expires
type session struct {
	username string
	expiry   time.Time
}

// The "db" package level variable will hold the reference to our database instance
var db *sql.DB

// Initialization function
func init() {
	// Initialize the database
	initDB()
}

// Open the database and setup the database pointer
func initDB() {
	var err error
	// Connect to the postgres db
	//you might have to change the connection string to add your database credentials
	db, err = sql.Open("sqlite3", "users.sqlite")
	if err != nil {
		panic(err)
	}
}

// we'll use this method later to determine if the session has expired
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

func Signup(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	loggedin, _ := UserIsLoggedIn(r)
	if !loggedin {
		return
	}
	// Parse and decode the request body into a new `Credentials` instance
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Error occurred when decoding request body: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Println("Error occurred when hashing the password: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Next, insert the username, along with the hashed password into the database
	if _, err = db.Exec("insert into users values ($1, $2)", creds.Username, string(hashedPassword)); err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		log.Println("Error occurred when inserting rows into db: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
}

// Returns true if the user is logged in, else otherwise
// Can return http.ErrNoCookie if there is no cookie
func UserIsLoggedIn(r *http.Request) (authorized bool, err error) {
	// Explicitely defin authorized to false
	authorized = false
	// Look for the login cookie
	c, err := r.Cookie(LoginToken)
	if err != nil {
		if err == http.ErrNoCookie {
			return
		}
		// For any other type of error, return a bad request status
		return
	}
	sessionToken := c.Value
	// We then get the session from our session map
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		return
	}
	// If the session is present, but has expired, we can delete the session, and return
	// an unauthorized status
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		return
	}
	// The checks are passed so the user is authorized
	authorized = true
	return
}

// Refresh user login cookie and extend its expiration period
func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code from this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie(LoginToken)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	userSession, exists := sessions[sessionToken]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// (END) The code until this point is the same as the first part of the `Welcome` route

	// If the previous session is valid, create a new session token for the current user
	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(TokenExpireMin * time.Minute)

	// Set the token in the session map, along with the user whom it represents
	sessions[newSessionToken] = session{
		username: userSession.username,
		expiry:   expiresAt,
	}

	// Delete the older session token
	delete(sessions, sessionToken)

	// Set the new token as the users LoginToken cookie
	http.SetCookie(w, &http.Cookie{
		Name:    LoginToken,
		Value:   newSessionToken,
		Expires: expiresAt,
	})
}

// Handles logout procedure
func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(LoginToken)
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// remove the users session from the session map
	delete(sessions, sessionToken)

	// We need to let the client know that the cookie is expired
	// In the response, we set the session token to an empty
	// value and set its expiry as the current time
	http.SetCookie(w, &http.Cookie{
		Name:    LoginToken,
		Value:   "",
		Expires: time.Now(),
	})

	log.Println("The user is logged out. The cookie is destroyed.")
	r.Method = http.MethodGet
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Handles login post
func loginPOST(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse the form data: ", err)
		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
		return
	}
	// Parse and decode the request body into a new `Credentials` instance
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	creds := &Credentials{username, password}
	// Get the existing entry present in the database for the given username
	log.Println(creds)
	result := db.QueryRow("select password from users where username=$1", creds.Username)
	// We create another instance of `Credentials` to store the credentials we get from the database
	storedCreds := &Credentials{}
	// Store the obtained password in `storedCreds`
	err = result.Scan(&storedCreds.Password)
	if err != nil {
		log.Println("Error while requesting the sql table data: ", err)
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Wrong login name and password.")
			return
		}
		// If the error is of any other type, send a 500 status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		log.Println("Wrong login name and password.")
		w.WriteHeader(http.StatusUnauthorized)
	}

	// If we reach this point, that means the users password was correct, and that they are authorized
	// The default 200 status is sent

	// Create a new random session token
	// we use the "github.com/google/uuid" library to generate UUIDs
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(TokenExpireMin * time.Minute)

	// Set the token in the session map, along with the session information
	sessions[sessionToken] = session{
		username: creds.Username,
		expiry:   expiresAt,
	}

	// Finally, we set the client cookie for LoginToken as the session token we just generated
	// we also set an expiry time of TokenExpireMin minutes
	http.SetCookie(w, &http.Cookie{
		Name:    LoginToken,
		Value:   sessionToken,
		Expires: expiresAt,
	})
	log.Println("The user is logged in. The cookie is set.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handles signing in procedure
func Signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Uploads a file
		log.Println("Handling the /login POST.")
		loginPOST(w, r)
	default:
		log.Println("Handling the /login GET.")
		LoginTemplate.Execute(w, nil)
	}

}
