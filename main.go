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
	r.HandleFunc("/deletecache", indexHandlerDeleteCache)

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

func indexHandlerMeetings(w http.ResponseWriter, r *http.Request) error {
	enableCors(&w)
	s, err := scrapeMeetingID()
	fmt.Fprint(w, string(s))

	if err != nil {
		return err
	}

	return nil
}

func indexHandlerVoting(w http.ResponseWriter, r *http.Request) error {
	enableCors(&w)
	vars := mux.Vars(r)
	varID := vars["id"]
	//fmt.Fprint(w, varId)
	i, err := strconv.Atoi(varID)

	if err != nil {
		http.NotFound(w, r)
		log.Fatal(err)
	}

	s, err := scrapeMeeting(i)

	fmt.Fprint(w, string(s))

	if err != nil {
		return err
	}

	return nil
}

func indexHandlerDeleteCache(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	storage = NewStorage()
	fmt.Print("Cache Deleted!\n")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
