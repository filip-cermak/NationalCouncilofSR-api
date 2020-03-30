// Sample helloworld is an App Engine app.
package main

// [START import]
import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", indexHandlerWebsite)
	http.HandleFunc("/meetings", indexHandlerMeetings)
	http.HandleFunc("/voting", indexHandlerVoting)
	

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	// [END setting_port]
}

// indexHandler responds to requests with our greeting.

func indexHandlerWebsite(w http.ResponseWriter, r *http.Request){
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

	fmt.Fprint(w, "test2")
}
// [END indexHandler]
// [END gae_go111_app]