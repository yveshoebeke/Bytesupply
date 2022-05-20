package users

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/yveshoebeke/Bytesupply/packages/app"
	"github.com/yveshoebeke/Bytesupply/packages/dbsql"
	"github.com/yveshoebeke/Bytesupply/packages/utilities"
)

type Login struct {
	SigninErrors   []string
	RegisterErrors []string
}

func ProcessLogin(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	// r.ParseForm()
	r.ParseMultipartForm(10 << 20)
	var login Login
	var user app.User
	t := time.Now().Format(time.RFC3339)

	if r.FormValue("submitLoginRegister") == "Login" {
		if !utilities.IsEmailAddress(r.FormValue("loginName"), true) {
			login.RegisterErrors = append(login.RegisterErrors, "Login must be email.")
			tmpl.ExecuteTemplate(w, "login.go.html", login)
			return
		}

		err := app.AppStruct.DB.QueryRow(dbsql.UserLogin, r.FormValue("loginName")).Scan(&user.Realname, &user.Password, &user.Title, &user.LastLogin)
		if err != nil {
			app.AppStruct.Log.Println("User login query failed:", err.Error()) // proper error handling instead of panic in your app
			login.SigninErrors = append(login.SigninErrors, fmt.Sprintf("'%s' is not registered.", r.FormValue("loginName")))
			tmpl.ExecuteTemplate(w, "login.go.html", login)
			return
		}
		// Check password hashes
		pwdMatch := utilities.ComparePasswords(user.Password, []byte(r.FormValue("loginPassword")))

		// If matched update last login time and update app user data
		if pwdMatch {
			_, err := app.AppStruct.DB.Exec(`UPDATE users SET lastlogin=NOW() WHERE email=?`, r.FormValue("loginName"))
			if err != nil {
				app.AppStruct.Log.Println("Login lastlogin update sql err:", err.Error())
				login.SigninErrors = append(login.SigninErrors, fmt.Sprintf("Report SQL error: %s", err.Error()))
				tmpl.ExecuteTemplate(w, "login.go.html", login)
			}
			// Register user into App
			u := utilities.User{}
			u.Username = r.FormValue("loginName")
			u.Password = user.Password
			u.Realname = user.Realname
			u.Title = user.Title
			u.LastLogin = user.LastLogin
			u.LoginTime = t

			app.AppStruct.User.Username = r.FormValue("loginName")
			app.AppStruct.User.Password = user.Password
			app.AppStruct.User.Realname = user.Realname
			app.AppStruct.User.Title = user.Title
			app.AppStruct.User.LastLogin = user.LastLogin
			app.AppStruct.User.LoginTime = t

			app.AppStruct.Log.Printf("User %s logged in", r.FormValue("loginName"))
			fmt.Println(u)

			tmpl.ExecuteTemplate(w, "welcome.go.html", u)
		} else {
			app.AppStruct.Log.Printf("Login for %s with %s failed to match.", r.FormValue("loginName"), r.FormValue("loginPassword"))
			login.SigninErrors = append(login.SigninErrors, "Wrong Email or Password.")

			tmpl.ExecuteTemplate(w, "login.go.html", login)
		}
	} else if r.FormValue("submitLoginRegister") == "Register" {
		// Hash and Verify password
		pwdGiven := utilities.HashAndSalt([]byte(r.FormValue("registerPassword")))
		pwdMatch := utilities.ComparePasswords(pwdGiven, []byte(r.FormValue("registerVerifyPassword")))

		// run through validation/vaccination filter function to be added to the utilities
		if !utilities.IsEmailAddress(r.FormValue("registerEmail"), true) {
			login.RegisterErrors = append(login.RegisterErrors, "Invalid Email address.")
		}
		if !utilities.IsAlphaNumeric(r.FormValue("registerName"), true) {
			login.RegisterErrors = append(login.RegisterErrors, "Invalid Name entry.")
		}
		if !utilities.IsAlphaNumeric(r.FormValue("registerCompany"), false) {
			login.RegisterErrors = append(login.RegisterErrors, "Invalid Company entry.")
		}
		if !utilities.IsPhoneNumber(r.FormValue("registerPhone"), false) {
			login.RegisterErrors = append(login.RegisterErrors, "Invalid Phone number.")
		}
		if !utilities.IsURLAddress(r.FormValue("registerURL"), false) {
			login.RegisterErrors = append(login.RegisterErrors, "Invalid URL address.")
		}
		if !pwdMatch {
			login.RegisterErrors = append(login.RegisterErrors, "Verify Passwords failed.")
		}
		uploadFilename, uploadError := utilities.UploadProfilePicture(r)
		if uploadError != nil {
			login.RegisterErrors = append(login.RegisterErrors, fmt.Sprintf("Picture upload failed: %s.", uploadError))
		}

		if len(login.RegisterErrors) > 0 {
			tmpl.ExecuteTemplate(w, "login.go.html", login)
		} else {
			_, err := app.AppStruct.DB.Exec(dbsql.AddUser, r.FormValue("registerName"), pwdGiven, r.FormValue("registerCompany"), r.FormValue("registerEmail"), r.FormValue("registerPhone"), r.FormValue("registerURL"), uploadFilename)
			if err != nil {
				app.AppStruct.Log.Println("Register INSERT sql err:", err.Error())
				http.Redirect(w, r, "/home", http.StatusExpectationFailed)
			}

			app.AppStruct.Log.Printf("User %s registered", r.FormValue("registerName"))
			app.AppStruct.User.Username = r.FormValue("registerEmail")
			app.AppStruct.User.Password = pwdGiven
			app.AppStruct.User.Realname = r.FormValue("registerName")
			app.AppStruct.User.Title = "user"
			app.AppStruct.User.LastLogin = t
			app.AppStruct.User.LoginTime = t

			http.Redirect(w, r, "/home", http.StatusSeeOther)
		}
	} else {
		app.AppStruct.Log.Println("Wrong login/register switch value")
		http.Redirect(w, r, "/home", http.StatusBadRequest)
	}

}

func ProcessLogout() {
	// Set app user to default values
	app.AppStruct.User.Username = "WWW"
	app.AppStruct.User.Password = "*"
	app.AppStruct.User.Realname = "Visitor"
	app.AppStruct.User.Title = "visitor"
	app.AppStruct.User.LastLogin = time.Now().Format(time.RFC3339)
	app.AppStruct.User.LoginTime = time.Now().Format(time.RFC3339)
}

type MessageData struct {
	App             *app.App
	TotalUserCount  int
	ActiveUserCount int
	MessageCount    int
}

func AdminData() (*MessageData, error) {
	data := MessageData{
		App:             app.AppStruct,
		TotalUserCount:  0,
		ActiveUserCount: 0,
		MessageCount:    0,
	}

	messagecounterr := app.AppStruct.DB.QueryRow(dbsql.CountUnreadMessages).Scan(&data.MessageCount)
	if messagecounterr != nil {
		return nil, messagecounterr
	}

	totalusercounterr := app.AppStruct.DB.QueryRow(dbsql.CountUsersByStatus, "%").Scan(&data.TotalUserCount)
	if totalusercounterr != nil {
		return nil, totalusercounterr
	}

	activeusercounterr := app.AppStruct.DB.QueryRow(dbsql.CountUsersByStatus, "1").Scan(&data.ActiveUserCount)
	if activeusercounterr != nil {
		return nil, activeusercounterr
	}

	return &data, nil
}

type UserData struct {
	App   *app.App
	Users utilities.Users
}

func GetUsersProcess() (UserData, error) {
	var uu utilities.Users
	var u utilities.UserRecord

	users, err := app.AppStruct.DB.Query(dbsql.GetAllUsersByStatus, "%")
	if err != nil {
		return UserData{}, err
	}
	defer users.Close()

	for users.Next() {
		err := users.Scan(&u.Name, &u.Title, &u.Password, &u.Company, &u.Email, &u.Phone, &u.URL, &u.Comment, &u.Picture, &u.Lastlogin, &u.Status, &u.Qturhm, &u.Created)
		if err != nil {
			return UserData{}, err
		}
		uu.Users = append(uu.Users, u)
	}

	data := UserData{
		App:   app.AppStruct,
		Users: uu,
	}

	return data, nil
}

func UpdateUsersProcess(r *http.Request) (string, string, error) {
	var sqlQuery string
	var val interface{}
	var ok bool

	r.ParseForm()
	email := r.FormValue("email")
	field := r.FormValue("field")
	value := r.FormValue("value")
	referer := r.FormValue("referer")

	switch field {
	case "status":
		val, ok = utilities.AllowedUserStatus[value]
		fmt.Printf("Status value -> %v %T\n", value, value)
		if !ok {
			app.AppStruct.Log.Printf("Wrong Status Value %s given for %s", value, email)
		}
		sqlQuery = fmt.Sprintf(dbsql.UpdateUser, field)
	case "title":
		val, ok = utilities.AllowedUserTitles[value]
		fmt.Printf("Title value -> %v %T\n", value, value)
		if !ok {
			app.AppStruct.Log.Printf("Wrong Title Value %s given for %s", value, email)
		}
		sqlQuery = fmt.Sprintf(dbsql.UpdateUser, field)
	case "comment":
		val = value
		sqlQuery = fmt.Sprintf(dbsql.UpdateUser, field)
	default:
		app.AppStruct.Log.Printf("Error changing Field: %s with Value: %s for User: %s", field, value, email)
	}

	fmt.Printf("--> Status value -> %v %T\n", val, val)

	_, err := app.AppStruct.DB.Exec(sqlQuery, val, email)
	if err != nil {
		return "", email, err
	}

	return referer, "", nil
}
