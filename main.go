// Sample helloworld is an App Engine app.
package main

// [START import]
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"gocloud.dev/server"
	)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", HandlerWebsite)
	r.HandleFunc("/meetings", indexHandlerMeetings)
	r.HandleFunc("/voting/{id:[0-9]+}", indexHandlerVoting)

    port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := server.New(r, nil)
	log.Printf("Listening on port %s", port)
	log.Fatal(srv.ListenAndServe(fmt.Sprintf(":%s", port)))

}

func HandlerWebsite(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, "An unofficial API for the National Council of the Slovak Republic voting statistics, work in progress")

}

func indexHandlerMeetings(w http.ResponseWriter, r *http.Request){

	fmt.Fprint(w, string(scrape_meeting_id()))
}

func indexHandlerVoting(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	varId := vars["id"]
	fmt.Fprint(w, varId)
}

// [END indexHandler]
// [END gae_go111_app]