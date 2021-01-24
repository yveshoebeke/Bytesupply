package utilities

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	logFile = os.Getenv("BS_LOGFILE")
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
