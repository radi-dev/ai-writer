package main

import (
	// "context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/radi-dev/ai-writer/services/conversation"
)

var fileServ = http.FileServer(http.Dir("templates/static/"))

func main() {
	// ctx := context.Background()

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServ)).Methods("GET")
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/assistant", assistant).Methods("GET")
	r.HandleFunc("/analytics", analytics).Methods("GET")
	r.HandleFunc("/media", media).Methods("GET")
	r.HandleFunc("/scheduler", scheduler).Methods("GET")
	r.HandleFunc("/settings", settings).Methods("GET")
	r.HandleFunc("/team", team).Methods("GET")
	r.HandleFunc("/templates", templates).Methods("GET")

	r.HandleFunc("/generate", generateResponse).Methods("POST")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", "3000"), r); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home page")
	t, err := template.ParseFiles("templates/index.html", "templates/pages/content_generator.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func assistant(w http.ResponseWriter, r *http.Request) {
	fmt.Println("templates")
	t, err := template.ParseFiles("templates/index.html", "templates/pages/ai_assistant.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func analytics(w http.ResponseWriter, r *http.Request) {
	fmt.Println("analytics")
	t, err := template.ParseFiles("templates/index.html", "templates/pages/analytics.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func media(w http.ResponseWriter, r *http.Request) {
	fmt.Println("media")
	t, err := template.ParseFiles("templates/index.html", "templates/pages/media_library.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func scheduler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("scheduler")
	t, err := template.ParseFiles("templates/index.html", "templates/pages/scheduler.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func settings(w http.ResponseWriter, r *http.Request) {
	fmt.Println("settings")
	t, err := template.ParseFiles("templates/index.html", "templates/pages/settings.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func team(w http.ResponseWriter, r *http.Request) {
	fmt.Println("team")
	t, err := template.ParseFiles("templates/index.html", "templates/pages/team.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func templates(w http.ResponseWriter, r *http.Request) {
	fmt.Println("templates")
	t, err := template.ParseFiles("templates/index.html", "templates/pages/templates.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func generateResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("generateResponsePost")
	ctx := r.Context()
	fmt.Println("\n\nformVal", r.FormValue("length-slider"))
	topic := r.FormValue("topic-input")
	lengthStr := r.FormValue("length-slider")
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		http.Error(w, "Invalid length value", http.StatusBadRequest)
		return
	}

	response := conversation.WriteLinkedInArticle(ctx, topic, length)
	t, err := template.ParseFiles("templates/forms/content_generator_form.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	type responseData struct {
		Response string
	}
	data := responseData{
		Response: response,
	}

	t.Execute(w, data)
}
