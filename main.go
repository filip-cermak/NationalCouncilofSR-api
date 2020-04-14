package main

// [START import]
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gocloud.dev/server"
)

var storage *Storage

func init() {
	storage = NewStorage()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandlerWebsite)
	r.HandleFunc("/meetings/", cached("24h", indexHandlerMeetings))
	r.HandleFunc("/voting/{id:[0-9]+}", cached("24h", indexHandlerVoting))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := server.New(r, nil)
	log.Printf("Listening on port %s", port)
	log.Fatal(srv.ListenAndServe(fmt.Sprintf(":%s", port)))

}

func indexHandlerWebsite(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(content))

}

func indexHandlerMeetings(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	fmt.Fprint(w, string(scrape_meeting_id()))
}

func indexHandlerVoting(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)
	varID := vars["id"]
	//fmt.Fprint(w, varId)
	i, err := strconv.Atoi(varID)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, string(scrape_meeting(i)))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
