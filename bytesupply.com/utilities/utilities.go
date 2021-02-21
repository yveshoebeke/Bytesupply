package utilities

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	app "bytesupply.com/app"

	"golang.org/x/crypto/bcrypt"
)

var (
	logFile           = os.Getenv("BS_LOGFILE")
	emailRegex        = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	phoneRegex        = regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	alphaNumericRegex = regexp.MustCompile("^[a-zA-Z0-9_,.!? ]*$")
	// AllowedUserTitles - .
	AllowedUserTitles = map[string]string{
		"User":   "user",
		"Expert": "expert",
		"Admin":  "admin",
	}
	// AllowedUserStatus - .
	AllowedUserStatus = map[string]int{
		"Active":     1,
		"Deactivate": 2,
		"Onhold":     8,
		"Suspended":  9,
	}
	// AllowedUserStatusInt - .
	AllowedUserStatusByInt = []string{"err0", "Active", "Deactivate", "err3", "err4", "err5", "err6", "err7", "Onhold", "Suspended"}
	// AllowedImageFormats - .
	AllowedImageFormats = map[string]string{
		"image/png":  "png",
		"image/jpeg": "jpg",
		"image/gif":  "gif",
	}
	// ProfilePictureLocation - .
	ProfilePictureLocation = "/static/img/profile"
	// DefaultProfilePicture - .
	DefaultProfilePicture = ProfilePictureLocation + "/" + "defaultuser.png"
)

// App - .
type App app.App

// User - .
type User app.User

// Message -
type Message struct {
	ID      int    //INT AUTO_INCREMENT PRIMARY KEY,
	User    string //VARCHAR(20) NOT NULL DEFAULT 'Unknown',
	Name    string //VARCHAR(100) NOT NULL,
	Company string //VARCHAR(100) DEFAULT '',
	Email   string //VARCHAR(100) NOT NULL,
	Phone   string //VARCHAR(20) DEFAULT '',
	URL     string //VARCHAR(200) DEFAULT '',
	Message string //TEXT NOT NULL,
	Status  int    //INT DEFAULT 0,
	Qturhm  int    //INT DEFAULT -1,
	Created string //TIMESTAMP DEFAULT CURRENT_TIMESTAMP
}

// Messages -
type Messages struct {
	Messages []Message
}

// UserRecord -
type UserRecord struct {
	Name      string //VARCHAR(100) NOT NULL,
	Title     string //VARCHAR(100) NOT NULL DEFAULT 'user',
	Password  string //VARCHAR(100) NOT NULL,
	Company   string //VARCHAR(100) DEFAULT '',
	Email     string //VARCHAR(100) NOT NULL,
	Phone     string //VARCHAR(20) DEFAULT '',
	URL       string //VARCHAR(200) DEFAULT '',
	Comment   string
	Picture   string
	Lastlogin string //TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	Status    int    //INT DEFAULT 1,
	Qturhm    int    //INT DEFAULT -1,
	Created   string //TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
}

// Users -
type Users struct {
	Users []UserRecord
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

// GetMessages -
func (app *App) GetMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%T %v", app, app)
	fmt.Fprintf(w, "From utilities.GetMessages -> User: %v DB: %v\r", nil, nil)
}

// GetIP - IP address retriever
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

// UploadProfilePicture - .
func UploadProfilePicture(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("myFile")
	if errors.As(err, &http.ErrMissingFile) {
		fmt.Println("upload - file not found (empty entry?)")
		return DefaultProfilePicture, nil
	}
	if err != nil {
		return DefaultProfilePicture, err
	}
	defer file.Close()

	extension, allowed := AllowedImageFormats[handler.Header["Content-Type"][0]]
	if !allowed {
		return DefaultProfilePicture, fmt.Errorf("%s is not an allowed format", handler.Header["Content-Type"][0])
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(ProfilePictureLocation, fmt.Sprintf("upload-*.%s", extension))
	if err != nil {
		return DefaultProfilePicture, fmt.Errorf("Tempfile i/o -> %v", err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return DefaultProfilePicture, fmt.Errorf("ReadAll i/o -> %v", err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	// return name for db record and that we have successfully uploaded our file!
	return tempFile.Name(), nil

}
