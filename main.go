package main

import (
	"encoding/json"
	"fmt"
	"hackday/api"
	"hackday/app"
	"hackday/db"
	"net/http"
	"os"
)

var defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	host := "http://localhost"
	if host == "http://localhost" {
		port = defaultPort
	}
	client, e := db.Conn()
	if e != nil {
		panic(e)
	}

	// static files define
	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", static))

	// api on graphql
	api.SetClient(client)
	http.HandleFunc("/auth", api.CreateTokenEndpoint)
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		result := api.ExecuteQuery(r.URL.Query().Get("query"), api.APISchema, r.URL.Query().Get("token"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// app on graphql
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := api.ExecuteQuery(r.URL.Query().Get("query"), api.AppSchema, r.URL.Query().Get("token"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// app hanlders
	http.HandleFunc("/", app.Hsign)
	http.HandleFunc("/profile", app.Hprofile)
	http.HandleFunc("/verification", app.Hverification)
	http.HandleFunc("/contact", app.Hcontact)

	// check sessions expire per minute
	// go forum.CheckPerMin()

	fmt.Println("listening on: " + host + ":" + port)
	e = http.ListenAndServe(":"+port, nil)
	if e != nil {
		panic(e)
	}
}
