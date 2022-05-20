package router

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"bytesupply.com/packages/app"
	"bytesupply.com/packages/googleapi"
	"bytesupply.com/packages/messages"
	"bytesupply.com/packages/users"
	"bytesupply.com/packages/utilities"
	"github.com/gorilla/mux"
)

const (
	StaticLocation = "./static/"
	TemplatePath   = "./static/templates/*.go.html"
)

var (
	/* templating */
	tmpl    = template.Must(template.New("main").Funcs(funcMap).ParseGlob(TemplatePath))
	funcMap = template.FuncMap{
		"hasHTTP": func(myUrl string) string {
			if strings.Contains(myUrl, "://") {
				return myUrl
			}

			return "https://" + myUrl
		},
		"userStatus": func(myStatus int) string {
			return utilities.AllowedUserStatusByInt[myStatus]
		},
	}
)

// Index page.
func Homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(app.AppStruct)
	tmpl.ExecuteTemplate(w, "index.go.html", app.AppStruct)
}

// Home page.
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println(app.AppStruct)
	tmpl.ExecuteTemplate(w, "home.go.html", app.AppStruct)
}

//------------------------------------------------------------
// About -> Company
func Company(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, StaticLocation+"html/company.html")
}

// About -> Staff
func Staff(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, StaticLocation+"html/staff.html")
}

// About -> History
func History(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, StaticLocation+"html/history.html")
}

//------------------------------------------------------------
// Expertise synopsis page.
func Expertise(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, StaticLocation+"html/expertise.html")
}

//------------------------------------------------------------
// Policies -> Terms and conditions.
func Terms(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, StaticLocation+"html/terms.html")
}

// Policies -> Privacy page.
func Privacy(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, StaticLocation+"html/privacy.html")
}

//------------------------------------------------------------
// Show Product accordions.
func Products(w http.ResponseWriter, r *http.Request) {
	type Item struct {
		ItemToShow string `json:"itemtoshow"`
	}
	item := Item{ItemToShow: "all"}
	tmpl.ExecuteTemplate(w, "product.go.html", item)
}

// Same as above but open specific product accordion.
func Product(w http.ResponseWriter, r *http.Request) {
	type Item struct {
		ItemToShow string `json:"itemtoshow"`
	}
	vars := mux.Vars(r)
	itemtoshow := vars["item"]
	item := Item{ItemToShow: itemtoshow}
	app.AppStruct.Log.Println("Item:", vars["item"])
	tmpl.ExecuteTemplate(w, "product.go.html", item)
}

//------------------------------------------------------------
// Contact us messaging.
func Contactus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var contact messages.Contact
		tmpl.ExecuteTemplate(w, "contactus.go.html", contact)
	} else if r.Method == http.MethodPost {
		// process contact us info
		messages.GetContactUsMessagesProcess(w, r, tmpl)
	}
}

// Read Contact US Messages
func Getmessages(w http.ResponseWriter, r *http.Request) {
	if app.AppStruct.User.Title != "admin" {
		http.Redirect(w, r, "/home", http.StatusForbidden)
		return
	}
	// MessageData -
	data, err := messages.GetMessagesProcess()
	if err != nil {
		app.AppStruct.Log.Println("Message retrieval query failed:", err.Error())
	}

	tmpl.ExecuteTemplate(w, "showMessages.go.html", data)
}

// Change message status.
func Changemessagestatus(w http.ResponseWriter, r *http.Request) {
	if app.AppStruct.User.Title != "admin" {
		http.Redirect(w, r, "/home", http.StatusForbidden)
		return
	}
	referer, err := messages.ChangeMessageStatusProcess(r)
	if err != nil {
		app.AppStruct.Log.Println("Message status update failed:", err.Error())
		return
	}

	http.Redirect(w, r, "/"+referer, http.StatusSeeOther)
}

//------------------------------------------------------------
// Logout.
func Logout(w http.ResponseWriter, r *http.Request) {
	users.ProcessLogout()
	tmpl.ExecuteTemplate(w, "home.go.html", app.AppStruct)
}

// Login
func Login(w http.ResponseWriter, r *http.Request) {
	type Login struct {
		SigninErrors   []string
		RegisterErrors []string
	}
	if r.Method == http.MethodGet {
		var login Login
		tmpl.ExecuteTemplate(w, "login.go.html", login)
	} else if r.Method == http.MethodPost {
		users.ProcessLogin(w, r, tmpl)
	}
}

//------------------------------------------------------------
// Admin functions.
func Admin(w http.ResponseWriter, r *http.Request) {
	if app.AppStruct.User.Title == "admin" {
		data, err := users.AdminData()
		if err != nil {
			app.AppStruct.Log.Println(err)
		}
		tmpl.ExecuteTemplate(w, "admin.go.html", data)
	} else {
		http.Redirect(w, r, "/home", http.StatusForbidden)
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "profile.go.html", app.AppStruct)
}

//------------------------------------------------------------
// Users
func User(w http.ResponseWriter, r *http.Request) {
	type MessageData struct {
		App          *app.App
		UserCount    int
		MessageCount int
	}

	data := MessageData{
		App:          app.AppStruct,
		UserCount:    0,
		MessageCount: 0,
	}

	tmpl.ExecuteTemplate(w, "user.go.html", data)
}

// Getusers.
func Getusers(w http.ResponseWriter, r *http.Request) {
	if app.AppStruct.User.Title != "admin" {
		http.Redirect(w, r, "/home", http.StatusForbidden)
		return
	}

	data, err := users.GetUsersProcess()
	if err != nil {
		app.AppStruct.Log.Println(err)
	}

	tmpl.ExecuteTemplate(w, "showUsers.go.html", data)
}

// Update user.
func Updateuser(w http.ResponseWriter, r *http.Request) {
	if app.AppStruct.User.Title != "admin" || r.Method != http.MethodPost {
		http.Redirect(w, r, "/home", http.StatusForbidden)
		return
	}

	referer, email, err := users.UpdateUsersProcess(r)
	if err != nil {
		app.AppStruct.Log.Printf("User update for User %s failed: %v", email, err.Error())
	}

	http.Redirect(w, r, "/"+referer, http.StatusSeeOther)
}

//------------------------------------------------------------<outlier>
// Google search
func Search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	searchKey := url.QueryEscape(r.FormValue("searchKey"))

	if len(searchKey) != 0 {
		searchResults, err := googleapi.GetSearchResults(searchKey)
		if err != nil {
			// app.Log.Println("Google API Err:", err)
			fmt.Println("Google API Err:", err)
		} else {
			tmpl.ExecuteTemplate(w, "search.go.html", searchResults)
		}
	} else {
		http.Redirect(w, r, r.FormValue("referer"), http.StatusSeeOther)
	}
}
