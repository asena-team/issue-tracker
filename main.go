package main

import (
	"context"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/gorilla/mux"
)

var (
	jiraClient *jira.Client
	views      *template.Template
)

var FileServerPaths = []string{"css", "js", "images"}

func main() {
	jiraClient = JiraClient()

	var err error
	views, err = template.ParseGlob(filepath.Join("views", "*.html"))
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	for _, path := range FileServerPaths {
		r.PathPrefix("/static/" + path).Handler(fs)
	}

	r.HandleFunc("/", Handler).Methods("POST", "GET")
	r.Use(APIRateLimiter)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	})

	srv := &http.Server{
		Addr:         net.JoinHostPort("", Port()),
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
	Mail        string `json:"mail"`
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
		if r.Method != "GET" {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Println(err)
				executeResult(&w, false, "IP adresiniz işlenemdi. Lütfen daha sonra tekrar deneyin.")

				return
			}

			limiter := GetVisitor(ip)
			if !limiter.Allow() {
				executeResult(&w, false, "Hey! Lütfen biraz yavaşla.\nBu kadar kısa sürede, bu kadar fazla istek gönderemezsin.\nEğer bunun bir hata olduğunu düşünüyorsan lütfen bizimle iletişime geç.")

				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if err := views.ExecuteTemplate(w, "index", map[string]interface{}{}); err != nil {
			log.Printf("failed to render template: %v", err)
		}

		return
	}

	if err := r.ParseForm(); err != nil {
		executeResult(&w, false, "Bilinmeyen bir hata oluştu. Lütfen daha sonra tekrar deneyin.")

		return
	}

	issue, ok := ValidateForm(&r.Form)
	if !ok {
		executeResult(&w, false, "Göndermeye çalıştığınız bilgiler istenen şartlara uymuyor.")

		return
	}

	i := NewIssue(issue)
	res, _, err := jiraClient.Issue.Create(i)
	if err != nil {
		log.Printf("issue creation error: %s", err)
		executeResult(&w, false, "Bilinmeyen bir hata oluştu. Lütfen daha sonra tekrar deneyin.")

		return
	}

	executeResult(&w, true, "Gönderiniz tarafımıza başarıyla ulaştı.\nGeri bildiriminiz için teşekkür ederiz.\nIssue Key: "+res.Key)
}

func executeResult(w *http.ResponseWriter, result bool, message string) {
	arr := map[string]interface{}{
		"result":   result,
		"messages": strings.Split(message, "\n"),
	}

	if err := views.ExecuteTemplate(*w, "index", arr); err != nil {
		log.Printf("failed to render template: %v", err)
	}
}
