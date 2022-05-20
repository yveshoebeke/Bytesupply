package main

/*
	Bytesupply.com - Web Server Pages App
	=====================================

	Complete documentation and user guides are available here:
	https://https://github.com/yveshoebeke/bytesupply/blob/master/README.md

	@author	yves.hoebeke@accds.com - 1011001.1110110.1100101.1110011

	@version 1.0.0

	(c) 2020 - Bytesupply, LLC - All Rights Reserved.
*/

import (
	"net/http"
	"os"
	"time"

	"bytesupply.com/packages/app"
	"bytesupply.com/packages/router"
	"bytesupply.com/packages/utilities"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/* Middleware */
func malcolm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.AppStruct.Log.Printf("User: %s | URL: %s | Method: %s | IP: %s", app.AppStruct.User.Username, r.URL.Path, r.Method, utilities.GetIP(r))
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
 *	- Routings						*
 *	- Serve and Listen.				*
 ************************************

*/
func main() {
	app.AppStruct.Log.Println("Starting service.")

	/* Close DB and Logfile connections when done */
	defer app.AppStruct.DB.Close()
	defer app.LogF.Close()

	/* Routers definitions */
	r := mux.NewRouter()

	/* Middleware */
	r.Use(malcolm)

	/* Allow static content */
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(router.StaticLocation))))

	/* Handlers --> Router */
	r.HandleFunc("/", router.Homepage).Methods(http.MethodGet)
	r.HandleFunc("/home", router.Home).Methods(http.MethodGet)
	r.HandleFunc("/company", router.Company).Methods(http.MethodGet)

	r.HandleFunc("/staff", router.Staff).Methods(http.MethodGet)
	r.HandleFunc("/history", router.History).Methods(http.MethodGet)
	r.HandleFunc("/expertise", router.Expertise).Methods(http.MethodGet)
	r.HandleFunc("/terms", router.Terms).Methods(http.MethodGet)
	r.HandleFunc("/privacy", router.Privacy).Methods(http.MethodGet)

	r.HandleFunc("/search", router.Search).Methods(http.MethodGet, http.MethodPost)

	r.HandleFunc("/product/{item:[a-zA-Z]+}", router.Product).Methods(http.MethodGet)
	r.HandleFunc("/products", router.Products).Methods(http.MethodGet)

	r.HandleFunc("/contactus", router.Contactus).Methods(http.MethodGet, http.MethodPost)

	r.HandleFunc("/login", router.Login).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/logout", router.Logout).Methods(http.MethodGet, http.MethodPost)

	r.HandleFunc("/admin", router.Admin).Methods(http.MethodGet)
	r.HandleFunc("/profile", router.Profile).Methods(http.MethodGet)

	r.HandleFunc("/getmessages", router.Getmessages).Methods(http.MethodGet)
	r.HandleFunc("/changemessagestatus/{id:[0-9]+}/{status:[0-9]}/{referer:[a-z]+}", router.Changemessagestatus).Methods(http.MethodGet)

	r.HandleFunc("/user", router.User).Methods(http.MethodGet)
	r.HandleFunc("/getusers", router.Getusers).Methods(http.MethodGet)
	r.HandleFunc("/updateuser", router.Updateuser).Methods(http.MethodPost)

	/* Server setup and start */
	BytesupplyServer := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, r),
		Addr:         ":80",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	/*
	**************************************
	* Setup and initialization completed *
	*                                    *
	*         Launch the server!         *
	**************************************
	 */
	app.AppStruct.Log.Fatal(BytesupplyServer.ListenAndServe())
}
