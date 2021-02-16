package app

import (
	"database/sql"
	// "log"
	"os"

	log "github.com/sirupsen/logrus"
)

// App - application structure */
type App struct {
	Log   *log.Logger
	Lfile *os.File
	User  *User
	DB    *sql.DB
	// Messages Messages
}

// User - info */
type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Realname  string `json:"realname"`
	Title     string `json:"title"`
	LastLogin string `jason:"lastlogin"`
	LoginTime string `json:"logintime"`
}
