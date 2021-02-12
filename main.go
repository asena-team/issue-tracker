package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/gorilla/mux"
)

var jiraClient *jira.Client

var FileServerPaths = []string{"css", "js", "images"}

func main() {
	jiraClient = JiraClient()

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static/"))
	r.Handle("/", fs)
	for _, path := range FileServerPaths {
		r.PathPrefix("/" + path).Handler(fs)
	}

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	})

	api := r.PathPrefix("/api/v" + strconv.Itoa(APIVersion)).Subrouter()
	api.Use(APIRateLimiter)
	api.HandleFunc("/issue", IssueHandler).Methods("POST")
	i := 0
	api.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		i++
		_, _ = writer.Write([]byte(strconv.Itoa(i)))
	}).Methods("GET")

	srv := &http.Server{
		Addr:         "0.0.0.0:" + Port(),
		Handler:      &Server{r},
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 30,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	go CleanupVisitors()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_ = srv.Shutdown(ctx)

	log.Println("APIServer shut down.")
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
	if origin := r.Header.Get("Origin"); origin != "" {
		r.Header.Set("Access-Control-Allow-Origin", origin)
		r.Header.Set("Access-Control-Allow-Methods", "GET, POST")
	}

	s.r.ServeHTTP(w, r)
}

func APIRateLimiter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal APIServer Error", http.StatusInternalServerError)
			return
		}

		limiter := GetVisitor(ip)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func IssueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")
	var issue Issue
	err := json.NewDecoder(r.Body).Decode(&issue)
	if err != nil {
		// TODO::Handle Error
	}

	// Duplicated shit codes
	if issue.Title == "" || !Compare(len(issue.Title), 150, 1) {
		// TODO::Handle Error
	}

	if issue.Description == "" || !Compare(len(issue.Description), 2000, 30) {
		// TODO::Handle Error
	}

	if issue.Reporter == "" || !Compare(len(issue.Reporter), 50, 1) {
		// TODO::Handle Error
	}

	if !Contains(issue.Type, IssueTypes) {
		// TODO::Handle Error
	}

	if !Contains(issue.Priority, Priorities) {
		// TODO::Handle Error
	}

	i := NewIssue(issue)
	_, _, err = jiraClient.Issue.Create(&i)
	if err != nil {
		// TODO::Handle Error
	}
}
