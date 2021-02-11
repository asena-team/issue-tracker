package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/gorilla/mux"
)

var jiraClient *jira.Client

func main() {
	jiraClient = JiraClient()

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static/"))
	r.Handle("/", &Server{r}).Handler(fs)
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	})

	r.HandleFunc("/issue", IssueHandler).Methods("POST")

	srv := &http.Server{
		Addr:         "0.0.0.0:" + Port(),
		Handler:      r,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 30,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_ = srv.Shutdown(ctx)

	log.Println("Server shut down.")
	os.Exit(0)
}

type Server struct {
	r *mux.Router
}

type Issue struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Priority    string `json:"priority"`
	Reporter    string `json:"reporter"`
}

func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "80"
	}

	return port
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("origin"); origin != "" {
		r.Header.Set("Access-Control-Allow-Origin", origin)
		r.Header.Set("Access-Control-Allow-Methods", "GET, POST")
	}

	s.r.ServeHTTP(w, r)
}

func IssueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")
	var issue Issue
	_ = json.NewDecoder(r.Body).Decode(&issue)

	i := NewIssue(issue)
	_, _, _ = jiraClient.Issue.Create(&i) // TODO::Handle Errors
}
