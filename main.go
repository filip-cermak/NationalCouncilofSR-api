// Sample helloworld is an App Engine app.
package main

// [START import]
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"strconv"
	"github.com/gocolly/colly"

)

type voting_info struct{
	Names []string
	Parties []string
	Votes  []string
}

func scrape_meeting(meeting_id int)([]byte){
	// Instantiate default collector
	c := colly.NewCollector()

	var name_slc []string
	var party_slc []string
	var vote_slc []string

	// On every a element which has href attribute call callback
	party:= ""

    c.OnHTML("td", func(e *colly.HTMLElement) {
    	// name and vote description

    	if len(e.Text) > 0 && string(e.Text[0]) !=  "[" && string(e.Text[0]) !=  "\n" {
	    	party = e.Text
	    	//fmt.Printf("%q", party)
		}

		if len(e.Text) > 0 && string(e.Text[0]) ==  "[" && string(e.Text[0]) !=  "\n"{
			vote_type := e.Text [1]
			name := e.Text[4:]

			name_slc = append(name_slc, name)
			party_slc = append(party_slc, party)
			vote_slc = append(vote_slc, string(vote_type))
			}
		//generate the output string
		//fmt.Printf("%q", name_str)

         })

    s := strconv.Itoa(meeting_id)
     // Start scraping
    c.Visit("https://www.nrsr.sk/web/Default.aspx?sid=schodze/hlasovanie/hlasklub&ID=" + s)

    //output

	vote_obj:= &voting_info{ Names: name_slc, Parties: party_slc, Votes: vote_slc}
	b, err := json.Marshal(vote_obj)

	if err != nil {
		log.Fatal(err)
	}

	return b
}

// [END import]
// [START main_func]

func main() {
	http.HandleFunc("/", indexHandler)
	

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

// [END main_func]

// [START indexHandler]

// indexHandler responds to requests with our greeting.
func indexHandler_2(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}

func indexHandler(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprint(w, string(scrape_meeting(43535)))
}

// [END indexHandler]
// [END gae_go111_app]