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

// // Message -
// type Message struct {
// 	ID      int    //INT AUTO_INCREMENT PRIMARY KEY,
// 	User    string //VARCHAR(20) NOT NULL DEFAULT 'Unknown',
// 	Name    string //VARCHAR(100) NOT NULL,
// 	Company string //VARCHAR(100) DEFAULT '',
// 	Email   string //VARCHAR(100) NOT NULL,
// 	Phone   string //VARCHAR(20) DEFAULT '',
// 	URL     string //VARCHAR(200) DEFAULT '',
// 	Message string //TEXT NOT NULL,
// 	Status  int    //INT DEFAULT 0,
// 	Qturhm  int    //INT DEFAULT -1,
// 	Created string //TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// }

// // Messages -
// type Messages struct {
// 	Messages []Message
// }

// MessageData -
// type MessageData struct {
// 	App      *App
// 	Messages Messages
// }
