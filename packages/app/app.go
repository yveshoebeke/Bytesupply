package app

import (
	sql "database/sql"
	"fmt"
	"io"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	log "github.com/sirupsen/logrus"
)

const (
	StaticLocation = "./static/"
	TemplatePath   = "./static/templates/*.go.html"
)

var (
	LogF                                           *os.File
	AppStruct                                      *App
	LogFile, ServerPort                            string
	DbHost, DbPort, DbUser, DbPassword, DbDatabase string
)

// App - application structure */
type App struct {
	Log   *log.Logger
	Lfile *os.File
	User  *User
	DB    *sql.DB
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

func init() {
	ServerPort = os.Getenv("BS_SERVER_PORT")
	LogFile = "/go/bin/log/bytesupply.log" //os.Getenv("BS_LOGFILE")
	DbHost = os.Getenv("BS_MYSQL_HOST")
	DbPort = os.Getenv("BS_MYSQL_PORT")
	DbUser = os.Getenv("BS_MYSQL_USERNAME")
	DbPassword = os.Getenv("BS_MYSQL_PASSWORD")
	DbDatabase = os.Getenv("BS_MYSQL_DB")

	InitApp()
}

func InitApp() {
	// Logging
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:     false,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	logger.SetLevel(log.InfoLevel)

	// log file set up
	LogF, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening logfile: %s -> %v", LogFile, err)
	}
	// Note: LogF.Close() --> in: main.go

	mw := io.MultiWriter(os.Stdout, LogF)
	logger.SetOutput(mw)

	// mysql connectivity
	dbConnectData := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DbUser, DbPassword, DbHost, DbPort, DbDatabase)
	db, err := sql.Open("mysql", dbConnectData)
	if err != nil {
		fmt.Println("db connect issue:", err.Error())
	}
	// Note: db.Close() --> in: main.go

	// Initial user data (before actual login)
	user := &User{
		Username:  "WWW",
		Password:  "*",
		Realname:  "Visitor",
		Title:     "visitor",
		LastLogin: time.Now().Format(time.RFC3339),
		LoginTime: time.Now().Format(time.RFC3339),
	}

	// Set app values
	AppStruct = &App{
		Log:   logger,
		User:  user,
		Lfile: LogF,
		DB:    db,
	}
}
