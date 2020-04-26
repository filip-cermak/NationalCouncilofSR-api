package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

func indexHandlerMeetingsSL(w http.ResponseWriter, r *http.Request) error {
	enableCors(&w)
	s, err := scrapeMeetingID("sl")
	fmt.Fprint(w, string(s))

	if err != nil {
		return err
	}

	return nil
}

func indexHandlerMeetingsEN(w http.ResponseWriter, r *http.Request) error {
	enableCors(&w)
	s, err := scrapeMeetingID("en")
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
