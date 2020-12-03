/*
Bytesupply.com - Web Server Pages App
	=====================================

	Complete documentation and user guides are available here:
	https://github.com/AccuityDeliverySystems/ACCDS-2.0/blob/master/README.md

	@author	yves.hoebeke@accds.com - 1011001.1110110.1100101.1110011

	@version 1.0.0

	(c) 2020 - Bytesupply, LLC - All Rights Reserved.
*/

package main

/* System libraries */
import (
	"bufio"
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

	"bytesupply.com/googleapi"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/boj/redistore.v1"
)

var (
	/* Extract env variables */
	staticLocation = os.Getenv("BS_STATIC_LOCATION")
	logFile        = os.Getenv("BS_LOGFILE")
	msgFile        = os.Getenv("BS_MSGFILE")
	serverPort     = os.Getenv("BS_SERVER_PORT")

	/* templating */
	tmpl    = template.Must(template.New("").Funcs(funcMap).ParseGlob(staticLocation + "/templ/*"))
	funcMap = template.FuncMap{
		"hasHTTP": func(myUrl string) string {
			if strings.Contains(myUrl, "://") {
				return myUrl
			}

			return "https://" + myUrl
		},
	}
)

// User - info */
type User struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Realname  string    `json:"realname"`
	Title     string    `json:"title"`
	LoginTime time.Time `json:"logintime"`
}

// Data - database structure */
type Data struct {
	ReqType   string    `json:"reqtype"`
	ReqCmd    string    `json:"reqcmd"`
	Timestamp time.Time `json:"Timestamp"`
}

// AppDatabase - application db */
type AppDatabase struct {
	name string
	//	db   *gorm.DB
}

// App - application structure */
type App struct {
	databases []AppDatabase
	log       *log.Logger
	mfile     *os.File
	lfile     *os.File
	store     *redistore.RediStore
	user      User
}

// GetIP - IP address retriever */
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARD-FOR")
	if forwarded != "" {
		return forwarded
	}

	return r.RemoteAddr
}

/* Routers */
func (app *App) homepage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticLocation+"/html/index.html")
}

func (app *App) bytesupply(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticLocation+"/html/bytesupply.html")
}

func (app *App) staff(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticLocation+"/html/staff.html")
}

func (app *App) history(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, staticLocation+"/html/history.html")
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
				// record data -> db table or -> txt file ... here ---> revisit.
				_, err := app.mfile.WriteString(fmt.Sprintf("Timestamp: %s\n", time.Now().Format(time.RFC3339)))
				if err != nil {
					app.log.Printf("Error writing %v: %v", msgFile, err)
				}

				_, _ = app.mfile.WriteString(fmt.Sprintf("     Name: %s\n", r.FormValue("contactName")))
				_, _ = app.mfile.WriteString(fmt.Sprintf("  Company: %s\n", r.FormValue("contactCompany")))
				_, _ = app.mfile.WriteString(fmt.Sprintf("    Email: %s\n", r.FormValue("contactEmail")))
				_, _ = app.mfile.WriteString(fmt.Sprintf("    Phone: %s\n", r.FormValue("contactPhone")))
				_, _ = app.mfile.WriteString(fmt.Sprintf("  Message:\n%s\n", r.FormValue("contactMessage")))
				_, _ = app.mfile.WriteString("----------------------------------------------------------------------\n")
				_, _ = app.mfile.WriteString(fmt.Sprintf(" Response:\n%s\n", r.FormValue("g-recaptcha-response")))
				_, _ = app.mfile.WriteString("======================================================================\n")
			}
		}

		msgStatus := MsgStatus{ValidToSend: validToRecord, Name: r.FormValue("contactName")}
		tmpl.ExecuteTemplate(w, "contactussent.gotmpl.html", msgStatus)
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
			tmpl.ExecuteTemplate(w, "search.gotmpl.html", searchResults)
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

func (app *App) test(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app.log.Println("Object:", vars["object"])
	http.ServeFile(w, r, staticLocation+"/html/"+vars["object"]+".html")
}

func (app *App) getlog(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<p style=\"color:blue;\"><a href=\"/\">Bytesupply</a></p><p>Access log</p>")

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

func (app *App) getmsg(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<p style=\"color:blue;\"><a href=\"/\">Bytesupply</a></p><p>Access messages</p>")

	msgfile, err := os.Open(msgFile)
	if err != nil {
		fmt.Fprintf(w, "<p style=\"color:blue;\">%s failed to open: %s</p>", msgFile, err)
	} else {
		scanner := bufio.NewScanner(msgfile)
		scanner.Split(bufio.ScanLines)

		fmt.Fprintf(w, "<ul>")
		for scanner.Scan() {
			fmt.Fprintf(w, "<li>%s</li>", scanner.Text())
		}
		fmt.Fprintf(w, "</ul>")
		msgfile.Close()
	}
}

func (app *App) registerUser(r *http.Request) error {
	app.user.Username = r.PostFormValue("username")
	app.user.Password = r.PostFormValue("password")
	app.user.Realname = "Yves Hoebeke"
	app.user.Title = "Owner"
	app.user.LoginTime = time.Now()
	app.log.Printf("Registering user %s as %s with username: %s and password: %s", app.user.Realname, app.user.Title, app.user.Username, app.user.Password)

	return nil
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
		app.log.Printf("User: %s | URL: %s | Method: %s | IP: %s", app.user.Username, r.URL.Path, r.Method, GetIP(r))
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

	// message file set up
	mf, err := os.OpenFile(msgFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening msgFile:", err)
	}
	defer mf.Close()

	lf, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error opening logfile: %s -> %v", logFile, err)
	}
	defer lf.Close()

	mw := io.MultiWriter(os.Stdout, lf)
	logger.SetOutput(mw)

	// Set app values
	app := &App{
		log:   logger,
		lfile: lf,
		mfile: mf,
	}

	app.log.Println("Starting service.")

	/* Routers definitions */
	r := mux.NewRouter()

	/* Middleware */
	r.Use(app.inMiddleWare)

	/* Allow static content */
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticLocation))))

	/* Handlers */
	r.HandleFunc("/", app.homepage).Methods("GET")
	r.HandleFunc("/bytesupply", app.bytesupply).Methods("GET")
	r.HandleFunc("/staff", app.staff).Methods("GET")
	r.HandleFunc("/history", app.history).Methods("GET")
	r.HandleFunc("/contactus", app.contactus).Methods("GET", "POST")
	r.HandleFunc("/search", app.search).Methods("GET", "POST")
	r.HandleFunc("/expertise", app.expertise).Methods("GET")
	r.HandleFunc("/terms", app.terms).Methods("GET")
	r.HandleFunc("/privacy", app.privacy).Methods("GET")
	r.HandleFunc("/getlog", app.getlog).Methods("GET")
	r.HandleFunc("/getmsg", app.getmsg).Methods("GET")
	r.HandleFunc("/request", app.request).Methods("POST")
	r.HandleFunc("/test/{object:[a-z]+}", app.test).Methods("GET", "POST")
	r.HandleFunc("/api/{version:[a-z0-9]+}/{request:[a-zA-Z]+}", app.api).Methods("GET", "POST")

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
	fmt.Println("Starting server")
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
