package utilities

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	logFile           = os.Getenv("BS_LOGFILE")
	emailRegex        = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	phoneRegex        = regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	alphaNumericRegex = regexp.MustCompile("^[a-zA-Z0-9_,.!? ]*$")
)

// User - info */
type User struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Realname  string    `json:"realname"`
	Title     string    `json:"title"`
	LoginTime time.Time `json:"logintime"`
}

// App - application structure */
type App struct {
	log   *log.Logger
	lfile *os.File
	user  User
	db    *sql.DB
}

// Getlog - func (app *App) Getlog(w http.ResponseWriter, r *http.Request) {
func Getlog(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<p style=\"color:blue;\"><a href=\"/home\">Bytesupply</a></p><p>Access log</p>")

	logfile, err := os.Open(logFile)
	if err != nil {
		fmt.Fprintf(w, "<p style=\"color:blue;\">%s failed to open: %s</p>", logFile, err)
	} else {
		scanner := bufio.NewScanner(logfile)
		scanner.Split(bufio.ScanLines)

		fmt.Fprintf(w, "<ul>")
		for scanner.Scan() {
			fmt.Fprintf(w, "<li>%s</li>", scanner.Text())
		}
		fmt.Fprintf(w, "</ul>")
		logfile.Close()
	}
}

// GetIP - IP address retriever */
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARD-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

//HashAndSalt - Encrypt password.
func HashAndSalt(pwd []byte) string {
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// ComparePasswords - Just what it says - true -> ok else false
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// IsEmailAddress - check if valid email address
func IsEmailAddress(emailAddress string, mandatory bool) bool {
	if !mandatory && len(emailAddress) == 0 {
		return true
	}
	if len(emailAddress) < 3 && len(emailAddress) > 254 {
		return false
	}
	return emailRegex.MatchString(emailAddress)
}

// IsPhoneNumber - check if valid phone number
// Allows for:
// "1(234)5678901x1234"
// "(+351) 282 43 50 50"
// "90191919908"
// "555-8909"
// "001 6867684"
// "001 6867684x1"
// "1 (234) 567-8901"
// "1-234-567-8901 ext1234"
func IsPhoneNumber(phoneNumber string, mandatory bool) bool {
	if !mandatory && len(phoneNumber) == 0 {
		return true
	}
	return phoneRegex.MatchString(phoneNumber)
}

// IsURLAddress - check if valid URL
func IsURLAddress(urlAddress string, mandatory bool) bool {
	if !mandatory && len(urlAddress) == 0 {
		return true
	}
	_, err := url.ParseRequestURI(urlAddress)
	if err != nil {
		return false
	}

	u, err := url.Parse(urlAddress)
	// if u.Schemw = "" -> prepend "http://" ? -> revisit later.
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// IsAlphaNumeric - Only alpha numeric characters allowed
func IsAlphaNumeric(stringValue string, mandatory bool) bool {
	if !mandatory && len(stringValue) == 0 {
		return true
	}
	return alphaNumericRegex.MatchString(stringValue)
}
