/*
	Bytesupply.com - Web Server Pages App
	=====================================

	Complete documentation and user guides are available here:
	https://https://github.com/yveshoebeke/bytesupply/blob/master/README.md

	@author	yves.hoebeke@accds.com - 1011001.1110110.1100101.1110011

	@version 1.0.0

	(c) 2020 - Bytesupply, LLC - All Rights Reserved.
*/

package main

/* System libraries */
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	googleapi "bytesupply.com/googleapi"
	utilities "bytesupply.com/utilities"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	/* Extract env variables */
	getlog         = utilities.Getlog
	staticLocation = os.Getenv("BS_STATIC_LOCATION")
	logFile        = os.Getenv("BS_LOGFILE")
	msgFile        = os.Getenv("BS_MSGFILE")
	serverPort     = os.Getenv("BS_SERVER_PORT")
	dbHost         = os.Getenv("BS_MYSQL_HOST")
	dbPort         = os.Getenv("BS_MYSQL_PORT")
	dbUser         = os.Getenv("BS_MYSQL_USERNAME")
	dbPassword     = os.Getenv("BS_MYSQL_PASSWORD")
	dbDatabase     = os.Getenv("BS_MYSQL_DB")
	/* sql statements */
	sqlAddMessage      = `INSERT INTO messages (user,name,company,email,phone,url,message) VALUES (?, ?, ?, ?, ?, ?, ?)`
	sqlAddUser         = `INSERT INTO users (name,password,company,email,phone,url) VALUES (?, ?, ?, ?, ?, ?)`
	sqlUserLogin       = `SELECT name, password, title, lastlogin FROM users WHERE email = ?`
	sqlUpdateLastlogin = `UPDATE users SET lastlogin=NOW() WHERE email=?`
	/* templating */
	tmpl    = template.Must(template.New("").Funcs(funcMap).ParseGlob(staticLocation + "/templates/*"))
	funcMap = template.FuncMap{
		"hasHTTP": func(myUrl string) string {
			if strings.Contains(myUrl, "://") {
				return myUrl
			}

			return "https://" + myUrl
		},
	}
)

// type getlog = utilities.Getlog

// User - info */
type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Realname  string `json:"realname"`
	Title     string `json:"title"`
	LastLogin string `jason:"lastlogin"`
	LoginTime string `json:"logintime"`
}

// App - application structure */
type App struct {
	log   *log.Logger
	lfile *os.File
	User  User
	db    *sql.DB
}

// Data - database structure */
type Data struct {
	ReqType   string    `json:"reqtype"`
	ReqCmd    string    `json:"reqcmd"`
	Timestamp time.Time `json:"Timestamp"`
}

/* Routers */
func (app *App) homepage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.go.html", app)
}

func (app *App) home(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.go.html", app)
}

func (app *App) company(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticLocation+"/html/company.html")
}

func (app *App) staff(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticLocation+"/html/staff.html")
}

func (app *App) history(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticLocation+"/html/history.html")
}

func (app *App) admin(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "admin.go.html", app)
}

func (app *App) profile(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "profile.go.html", app)
}

func (app *App) logout(w http.ResponseWriter, r *http.Request) {
	// Set app user to default values
	app.User.Username = "WWW"
	app.User.Password = "*"
	app.User.Realname = "Visitor"
	app.User.Title = "visitor"
	app.User.LastLogin = time.Now().Format(time.RFC3339)
	app.User.LoginTime = time.Now().Format(time.RFC3339)

	tmpl.ExecuteTemplate(w, "home.go.html", app)
}

func (app *App) login(w http.ResponseWriter, r *http.Request) {
	type Login struct {
		SigninErrors   []string
		RegisterErrors []string
	}
	if r.Method == http.MethodGet {
		// Get - present form(s)
		var login Login
		// d := Data{Nothing: "Nothing"}
		tmpl.ExecuteTemplate(w, "login.go.html", login)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		var login Login
		var user User
		t := time.Now().Format(time.RFC3339)

		if r.FormValue("submitLoginRegister") == "Login" {
			if !utilities.IsEmailAddress(r.FormValue("loginName"), true) {
				login.RegisterErrors = append(login.RegisterErrors, fmt.Sprintf("Login must be email."))
				tmpl.ExecuteTemplate(w, "login.go.html", login)
				return
			}

			err := app.db.QueryRow(sqlUserLogin, r.FormValue("loginName")).Scan(&user.Realname, &user.Password, &user.Title, &user.LastLogin)
			if err != nil {
				app.log.Println("User login query failed:", err.Error()) // proper error handling instead of panic in your app
				login.SigninErrors = append(login.SigninErrors, fmt.Sprintf("'%s' is not registered.", r.FormValue("loginName")))
				tmpl.ExecuteTemplate(w, "login.go.html", login)
				return
			}
			// Check password hashes
			pwdMatch := utilities.ComparePasswords(user.Password, []byte(r.FormValue("loginPassword")))

			// If matched update last login time and update app user data
			if pwdMatch {
				_, err := app.db.Exec(`UPDATE users SET lastlogin=NOW() WHERE email=?`, r.FormValue("loginName"))
				if err != nil {
					app.log.Println("Login lastlogin update sql err:", err.Error())
					login.SigninErrors = append(login.SigninErrors, fmt.Sprintf("Report SQL error: %s", err.Error()))
					tmpl.ExecuteTemplate(w, "login.go.html", login)
				}
				// Register user into App
				app.User.Username = r.FormValue("loginName")
				app.User.Password = user.Password
				app.User.Realname = user.Realname
				app.User.Title = user.Title
				app.User.LastLogin = user.LastLogin
				app.User.LoginTime = t
				app.log.Printf("User %s logged in", r.FormValue("loginName"))
				http.Redirect(w, r, "/home", http.StatusSeeOther)
			} else {
				app.log.Printf("Login for %s with %s failed to match.", r.FormValue("loginName"), r.FormValue("loginPassword"))
				login.SigninErrors = append(login.SigninErrors, fmt.Sprintf("Wrong Rmail or Password."))
				tmpl.ExecuteTemplate(w, "login.go.html", login)
			}
		} else if r.FormValue("submitLoginRegister") == "Register" {
			// Hash and Verify password
			pwdGiven := utilities.HashAndSalt([]byte(r.FormValue("registerPassword")))
			pwdMatch := utilities.ComparePasswords(pwdGiven, []byte(r.FormValue("registerVerifyPassword")))

			// run through validation/vaccination filter function to be added to the utilities
			if !utilities.IsEmailAddress(r.FormValue("registerEmail"), true) {
				login.RegisterErrors = append(login.RegisterErrors, fmt.Sprintf("Invalid Email address."))
			}
			if !utilities.IsAlphaNumeric(r.FormValue("registerName"), true) {
				login.RegisterErrors = append(login.RegisterErrors, fmt.Sprintf("Invalid Name entry."))
			}
			if !utilities.IsAlphaNumeric(r.FormValue("registerCompany"), false) {
				login.RegisterErrors = append(login.RegisterErrors, fmt.Sprintf("Invalid Company entry."))
			}
			if !utilities.IsPhoneNumber(r.FormValue("registerPhone"), false) {
				login.RegisterErrors = append(login.RegisterErrors, fmt.Sprintf("Invalid Phone number."))
			}
			if !utilities.IsURLAddress(r.FormValue("registerURL"), false) {
				login.RegisterErrors = append(login.RegisterErrors, fmt.Sprintf("Invalid URL address."))
			}
			if !pwdMatch {
				login.RegisterErrors = append(login.RegisterErrors, fmt.Sprintf("Verify Passwords failed."))
			}

			if len(login.RegisterErrors) > 0 {
				tmpl.ExecuteTemplate(w, "login.go.html", login)
			} else {
				_, err := app.db.Exec(sqlAddUser, r.FormValue("registerName"), pwdGiven, r.FormValue("registerCompany"), r.FormValue("registerEmail"), r.FormValue("registerPhone"), r.FormValue("registerURL"))
				if err != nil {
					app.log.Println("Register INSERT sql err:", err.Error())
					http.Redirect(w, r, "/home", http.StatusExpectationFailed)
				}

				app.log.Printf("User %s registered", r.FormValue("registerName"))
				app.User.Username = r.FormValue("registerEmail")
				app.User.Password = pwdGiven
				app.User.Realname = r.FormValue("registerName")
				app.User.Title = "user"
				app.User.LastLogin = t
				app.User.LoginTime = t

				http.Redirect(w, r, "/home", http.StatusSeeOther)
			}
		} else {
			app.log.Println("Wrong login/register switch value")
			http.Redirect(w, r, "/home", http.StatusBadRequest)
		}
	}
}

func (app *App) contactus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, staticLocation+"/html/contactus.html")
	} else if r.Method == http.MethodPost {
		// process contact us info
		type MsgStatus struct {
			ValidToSend bool   `json:"validtosend"`
			Name        string `json:"name"`
		}

		r.ParseForm()

		var validToRecord = false

		if r.FormValue("validEntry") == "false" {
			validToRecord = false
		} else {
			validToRecord = true
			// Validate (name, email and message are mandatory)
			for varName, varValue := range r.Form {
				switch varName {
				case "contactName":
				case "contactEmail":
				case "contactMessage":
					if len(varValue[0]) == 0 {
						validToRecord = false
					}
					break
				default:
					break
				}
			}

			if validToRecord {
				_, err := app.db.Exec(sqlAddMessage, app.User.Username, r.FormValue("contactName"), r.FormValue("contactCompany"), r.FormValue("contactEmail"), r.FormValue("contactPhone"), r.FormValue("contactURL"), r.FormValue("contactMessage"))
				if err != nil {
					app.log.Println("ContactUs INSERT sql err:", err.Error())
				}
			}
		}

		msgStatus := MsgStatus{ValidToSend: validToRecord, Name: r.FormValue("contactName")}
		tmpl.ExecuteTemplate(w, "contactussent.go.html", msgStatus)
	}
}

func (app *App) search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	searchKey := url.QueryEscape(r.FormValue("searchKey"))

	if len(searchKey) != 0 {
		searchResults, err := googleapi.GetSearchResults(searchKey)
		if err != nil {
			app.log.Println("Google API Err:", err)
		} else {
			tmpl.ExecuteTemplate(w, "search.go.html", searchResults)
		}
	} else {
		http.Redirect(w, r, r.FormValue("referer"), http.StatusSeeOther)
	}
}

func (app *App) expertise(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticLocation+"/html/expertise.html")
}

func (app *App) terms(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticLocation+"/html/terms.html")
}

func (app *App) privacy(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticLocation+"/html/privacy.html")
}

func (app *App) products(w http.ResponseWriter, r *http.Request) {
	type Item struct {
		ItemToShow string `json:"itemtoshow"`
	}
	item := Item{ItemToShow: "all"}
	tmpl.ExecuteTemplate(w, "product.go.html", item)
}

func (app *App) product(w http.ResponseWriter, r *http.Request) {
	type Item struct {
		ItemToShow string `json:"itemtoshow"`
	}
	vars := mux.Vars(r)
	itemtoshow := vars["item"]
	item := Item{ItemToShow: itemtoshow}
	app.log.Println("Item:", vars["item"])
	tmpl.ExecuteTemplate(w, "product.go.html", item)
}

func (app *App) test(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app.log.Println("Object:", vars["object"])
	http.ServeFile(w, r, staticLocation+"/html/"+vars["object"]+".html")
}

func (app *App) api(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	version := vars["version"]
	request := vars["request"]

	app.log.Println("@api with version:", version, "and request:", request)

	switch version {
	default:
	case "v1":
		switch request {
		case "qTurHm":
			app.qTurHm(w, r)
		case "request":
			app.request(w, r)
		}
	}
}

func (app *App) qTurHm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, staticLocation+"/html/contactus.html")
	} else if r.Method == http.MethodPost {
		type Target struct {
			Top    int `json:"top"`
			Left   int `json:"left"`
			Width  int `json:"width"`
			Height int `json:"height"`
		}

		type Move struct {
			T int `json:"t"`
			X int `json:"x"`
			Y int `json:"y"`
		}

		type QTurHm struct {
			Key           string `json:"userkey"`
			TimeCreated   int    `json:"timestamp"`
			ResultContent string `json:"resultcontent"`
			URL           string `json:"origURL"`
			Mobile        bool   `json:"mobile"`
			Target        Target `json:"target"`
			Receiver      string `json:"receiver"`
			SampleCount   int    `json:"samples"`
			Moves         []Move `json:"moves"`
		}

		var q QTurHm

		// Try to decode the request body into the struct.
		err := json.NewDecoder(r.Body).Decode(&q)
		if err != nil {
			app.log.Println("API error (qTurHm):", err.Error())
			return
		}

		app.log.Printf("%v", q)
		app.log.Printf("Key: %s Time: %d", q.Key, q.TimeCreated)
		rfn := q.Key + "_" + strconv.Itoa(q.TimeCreated)
		app.log.Printf("Result File Name: %s should be: %s", rfn, q.ResultContent)

		res := []byte("8")
		werr := ioutil.WriteFile("/go/bin/data/qTurHm/"+rfn, res, 0644)
		if werr != nil {
			app.log.Printf("Error writing result file /go/bin/data/qTurHm/%s: %v", rfn, werr)
		}
	}
}

func (app *App) request(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.log.Println("Error parsing Body:", err)
	}
	var data Data
	json.Unmarshal(reqBody, &data)
	data.Timestamp = time.Now()

	json.NewEncoder(w).Encode(data)

	app.log.Printf("Request command received: %s", data.ReqType)
}

/* Middleware */
func (app *App) inMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.log.Printf("User: %s | URL: %s | Method: %s | IP: %s", app.User.Username, r.URL.Path, r.Method, utilities.GetIP(r))
		next.ServeHTTP(w, r)
	})
}

/*
       ^ ^
      (o O)
 ___oOO(.)OOo___
 _______________

 ************************************
 *	Execution start point!!!!!!!!!	*
 *	Structure and Methods 			*
 *	- Setup and start of app.		*
 *	- Serve and Listen.				*
 ************************************

*/
func init() {
}

func main() {
	// Logging
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:     false,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	logger.SetLevel(log.InfoLevel)

	// log file set up
	lf, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening logfile: %s -> %v", logFile, err)
	}
	defer lf.Close()

	mw := io.MultiWriter(os.Stdout, lf)
	logger.SetOutput(mw)

	// mysql connectivity
	dbConnectData := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbDatabase)
	db, err := sql.Open("mysql", dbConnectData)
	if err != nil {
		fmt.Println("db connect issue:", err.Error())
	}
	defer db.Close()

	// Initial user data (before actual login)
	user := User{
		Username:  "WWW",
		Password:  "*",
		Realname:  "Visitor",
		Title:     "visitor",
		LastLogin: time.Now().Format(time.RFC3339),
		LoginTime: time.Now().Format(time.RFC3339),
	}

	// Set app values
	app := &App{
		log:   logger,
		lfile: lf,
		User:  user,
		db:    db,
	}

	app.log.Println("Starting service.")

	/* Routers definitions */
	r := mux.NewRouter()

	/* Middleware */
	r.Use(app.inMiddleWare)

	/* Allow static content */
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticLocation))))

	/* Handlers */
	r.HandleFunc("/", app.homepage).Methods(http.MethodGet)
	r.HandleFunc("/home", app.home).Methods(http.MethodGet)
	r.HandleFunc("/company", app.company).Methods(http.MethodGet)
	r.HandleFunc("/login", app.login).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/logout", app.logout).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/staff", app.staff).Methods(http.MethodGet)
	r.HandleFunc("/history", app.history).Methods(http.MethodGet)
	r.HandleFunc("/contactus", app.contactus).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/search", app.search).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/expertise", app.expertise).Methods(http.MethodGet)
	r.HandleFunc("/terms", app.terms).Methods(http.MethodGet)
	r.HandleFunc("/admin", app.admin).Methods(http.MethodGet)
	r.HandleFunc("/profile", app.profile).Methods(http.MethodGet)
	r.HandleFunc("/privacy", app.privacy).Methods(http.MethodGet)
	r.HandleFunc("/product/{item:[a-zA-Z]+}", app.product).Methods(http.MethodGet)
	r.HandleFunc("/products", app.products).Methods(http.MethodGet)
	r.HandleFunc("/getlog", getlog).Methods(http.MethodGet)
	r.HandleFunc("/request", app.request).Methods("POST")
	r.HandleFunc("/test/{object:[a-z]+}", app.test).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/api/{version:[a-z0-9]+}/{request:[a-zA-Z]+}", app.api).Methods(http.MethodGet, http.MethodPost)

	/* Server setup and start */
	BytesupplyServer := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		Addr:         serverPort,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	/*
	**************************************
	* Setup and initialization completed *
	*                                    *
	*         Launch the server!         *
	**************************************
	 */
	app.log.Fatal(BytesupplyServer.ListenAndServe())

	/*
		****************************************************
		POST request test:

		curl --header "Content-Type: application/json" \
		--request POST \
		--data '{"reqtype":"test", "reqcmd":"Here is some requested data"}' \
		https://bytesupply.com/api/v1/request

		****************************************************
	*/
}
