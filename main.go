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
	"os"
	"text/template"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/boj/redistore.v1"
)

/* Application constants */
const ()

var (
	/* Extract env variables */
	staticLocation = os.Getenv("BS_STATIC_LOCATION")
	logFile        = os.Getenv("BS_LOGFILE")
	msgFile        = os.Getenv("BS_MSGFILE")
	serverPort     = os.Getenv("BS_SERVER_PORT")
	/* templating */
	tmpl = template.Must(template.ParseGlob(staticLocation + "/templ/*"))
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

		if r.Form["validEntry"][0] == "false" {
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

				_, _ = app.mfile.WriteString(fmt.Sprintf("     Name: %s\n", r.Form["contactName"][0]))
				_, _ = app.mfile.WriteString(fmt.Sprintf("  Company: %s\n", r.Form["contactCompany"][0]))
				_, _ = app.mfile.WriteString(fmt.Sprintf("    Email: %s\n", r.Form["contactEmail"][0]))
				_, _ = app.mfile.WriteString(fmt.Sprintf("    Phone: %s\n", r.Form["contactPhone"][0]))
				_, _ = app.mfile.WriteString(fmt.Sprintf("  Message:\n%s\n", r.Form["contactMessage"][0]))
				_, _ = app.mfile.WriteString("----------------------------------------------------------------------\n")
			}
		}

		msgStatus := MsgStatus{ValidToSend: validToRecord, Name: r.Form["contactName"][0]}
		tmpl.ExecuteTemplate(w, "contactussent.gotmpl.html", msgStatus)
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

func (app *App) registerUser(r *http.Request) error {
	app.user.Username = r.PostFormValue("username")
	app.user.Password = r.PostFormValue("password")
	app.user.Realname = "Yves Hoebeke"
	app.user.Title = "Owner"
	app.user.LoginTime = time.Now()
	app.log.Printf("Registering user %s as %s with username: %s and password: %s", app.user.Realname, app.user.Title, app.user.Username, app.user.Password)

	return nil
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
	r.HandleFunc("/expertise", app.expertise).Methods("GET")
	r.HandleFunc("/terms", app.terms).Methods("GET")
	r.HandleFunc("/privacy", app.privacy).Methods("GET")
	r.HandleFunc("/getlog", app.getlog).Methods("GET")
	r.HandleFunc("/getmsg", app.getmsg).Methods("GET")
	r.HandleFunc("/request", app.request).Methods("POST")

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
		https://bytesupply.com/request

		****************************************************
	*/
}
