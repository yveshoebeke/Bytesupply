package messages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yveshoebeke/Bytesupply/packages/app"
	"github.com/yveshoebeke/Bytesupply/packages/dbsql"
	"github.com/yveshoebeke/Bytesupply/packages/utilities"
)

type Contact struct {
	Errors []string
}

type MsgStatus struct {
	ValidToSend bool   `json:"validtosend"`
	Name        string `json:"name"`
}

func GetContactUsMessagesProcess(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	var contact Contact
	r.ParseForm()

	if !utilities.IsAlphaNumeric(r.FormValue("contactName"), true) || len(r.FormValue("contactName")) < 3 {
		contact.Errors = append(contact.Errors, "Invalid (Mandatory) Name entry.")
	}
	if !utilities.IsAlphaNumeric(r.FormValue("contactCompany"), false) {
		contact.Errors = append(contact.Errors, "Invalid Company entry.")
	}
	if !utilities.IsPhoneNumber(r.FormValue("contactPhone"), false) {
		contact.Errors = append(contact.Errors, "Invalid Phone number.")
	}
	if !utilities.IsEmailAddress(r.FormValue("contactEmail"), true) {
		contact.Errors = append(contact.Errors, "Invalid (Mandatory) Email address.")
	}
	if !utilities.IsURLAddress(r.FormValue("contactURL"), false) {
		contact.Errors = append(contact.Errors, "Invalid URL address.")
	}
	if !utilities.IsAlphaNumeric(r.FormValue("contactMessage"), true) || len(r.FormValue("contactMessage")) < 3 {
		contact.Errors = append(contact.Errors, "Invalid (Mandatory) Message entry.")
	}

	if len(contact.Errors) > 0 {
		tmpl.ExecuteTemplate(w, "contactus.go.html", contact)
	} else {
		_, err := app.AppStruct.DB.Exec(dbsql.AddMessage, app.AppStruct.User.Username, r.FormValue("contactName"), r.FormValue("contactCompany"), r.FormValue("contactEmail"), r.FormValue("contactPhone"), r.FormValue("contactURL"), r.FormValue("contactMessage"))
		if err != nil {
			fmt.Println("ContactUs INSERT sql err:", err.Error())
			app.AppStruct.Log.Println("ContactUs INSERT sql err:", err.Error())
		}

		msgStatus := MsgStatus{ValidToSend: true, Name: r.FormValue("contactName")}
		tmpl.ExecuteTemplate(w, "contactussent.go.html", msgStatus)
	}
}

type MessageData struct {
	App      *app.App
	Messages utilities.Messages
}

func GetMessagesProcess() (MessageData, error) {
	var mm utilities.Messages
	var m utilities.Message

	messages, err := app.AppStruct.DB.Query(dbsql.GetAllMessagesByStatus, "%")
	if err != nil {
		return MessageData{}, err
	}
	// defer messages.Close()

	for messages.Next() {
		err := messages.Scan(&m.ID, &m.User, &m.Name, &m.Company, &m.Email, &m.Phone, &m.URL, &m.Message, &m.Status, &m.Qturhm, &m.Created)
		if err != nil {
			return MessageData{}, err
		}

		mm.Messages = append(mm.Messages, m)
	}

	data := MessageData{
		App:      app.AppStruct,
		Messages: mm,
	}

	return data, nil
}

func ChangeMessageStatusProcess(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	status := vars["status"]
	referer := vars["referer"]

	_, err := app.AppStruct.DB.Exec(dbsql.UpdateMessageStatus, status, id)
	if err != nil {
		return "", err
	}

	return referer, nil
}
