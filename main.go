package main

import (
	"encoding/json"
	"fmt"
	"hackday/api"
	"hackday/app"
	"net/http"
	"os"
)

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	// api on graphql
	mux.HandleFunc("/auth", api.CreateTokenEndpoint)
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		result := api.ExecuteQuery(r.URL.Query().Get("query"), api.APISchema, r.URL.Query().Get("token"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// app hanlders
	mux.HandleFunc("/", app.Hsign)                        // 100%
	mux.HandleFunc("/s/", app.HsaveUser)                  // 100%
	mux.HandleFunc("/forgot", app.Hforgot)                // 100%
	mux.HandleFunc("/verification", app.Hverification)    // 100%
	mux.HandleFunc("/restore", app.Hrestore)              // 100%
	mux.HandleFunc("/logout", app.Hlogout)                // 100%
	mux.HandleFunc("/profile", app.Hprofile)              // 100%
	mux.HandleFunc("/profile/settings/", app.Hsettings)   // 100%
	mux.HandleFunc("/profile/ch/", app.HphotoAndSocials)  // 0%
	mux.HandleFunc("/profile/my-vacantions", app.Hworks)  // 100%
	mux.HandleFunc("/profile/my-vacantions/", app.Hwork)  // 100%
	mux.HandleFunc("/profile/my-subscription", app.Hsubs) // 100%
	mux.HandleFunc("/user/", app.Hprofile)                // 100%
	mux.HandleFunc("/user/settings/", app.Hsettings)      // 100%
	mux.HandleFunc("/contact", app.Hcontact)              // 100%
	mux.HandleFunc("/create-work", app.HworkCreate)       // 100%
	mux.HandleFunc("/filter", app.Hfilter)                // 100%
	mux.HandleFunc("/works", app.Hworks)                  // 100%
	mux.HandleFunc("/works/", app.Hwork)                  // 100%

	// static files define
	static := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", static))

	return mux
}

func main() {
	defaultPort := "8080"
	port := os.Getenv("PORT")
	// host := "https://hackday2020.herokuapp.com"
	host := "http://localhost"
	if host == "http://localhost" {
		port = defaultPort
	}

	e := app.InitProg()
	if e != nil {
		app.WriteLog(e.Error())
		return
	}
	// check sessions expire per minute
	go app.CheckPerMin()

	mux := routes()

	fmt.Println("listening on: " + host + ":" + port)
	app.WriteLog("listening on: " + host + ":" + port)
	e = http.ListenAndServe(":"+port, mux)
	if e != nil {
		app.WriteLog(e.Error())
		return
	}
}
