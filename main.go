package main

// [START import]
import (
	"fmt"
	"log"
	"os"

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
	r.HandleFunc("/meetings/", cached("24h", indexHandlerMeetingsEN))
	r.HandleFunc("/meetings/sl", cached("24h", indexHandlerMeetingsSL))
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
