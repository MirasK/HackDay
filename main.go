package main

import (
	"encoding/json"
	"fmt"
	"hackday/api"
	"hackday/app"
	"net/http"
	"os"
)

func main() {
	defaultPort := "8080"
	port := os.Getenv("PORT")
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

	mux := http.NewServeMux()
	// static files define
	static := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", static))

	// api on graphql
	mux.HandleFunc("/auth", api.CreateTokenEndpoint)
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		result := api.ExecuteQuery(r.URL.Query().Get("query"), api.APISchema, r.URL.Query().Get("token"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// app hanlders
	mux.HandleFunc("/", app.Hsign)
	mux.HandleFunc("/forgot", app.Hcontact)
	mux.HandleFunc("/logout", app.Hlogout)
	mux.HandleFunc("/profile", app.Hprofile)
	mux.HandleFunc("/profile/settings", app.Hsettings)
	mux.HandleFunc("/verification", app.Hverification)
	mux.HandleFunc("/contact", app.Hcontact)

	fmt.Println("listening on: " + host + ":" + port)
	app.WriteLog("listening on: " + host + ":" + port)
	e = http.ListenAndServe(":"+port, mux)
	if e != nil {
		app.WriteLog(e.Error())
		return
	}
}
