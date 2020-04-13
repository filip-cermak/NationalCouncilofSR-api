package main

import (
	"strconv"
	"encoding/json"
	"github.com/gocolly/colly"
	"log"
	"strings"

)

type voting_info struct{
	Names []string
	Parties []string
	Votes  []string
}

type voting_session struct{
 	Session_ID string
 	Timestamp string
 	Text string
}

type sessions struct{
 	VotingSessions []voting_session
}

func scrape_meeting_id()([]byte){
	// Instantiate default collector
	c := colly.NewCollector()

	var allSessions []voting_session

	// On every a element which has href attribute call callback

    c.OnHTML("tr[class^='tab_zoznam']", func(e *colly.HTMLElement) {

		//timestamp
		timestamp := e.ChildText("td:nth-child(2)")
		timestamp = strings.Replace(timestamp, " ", "", -1)
		timestamp = strings.Replace(timestamp, "\n", " ", -1)
		//text
		text := e.ChildText("td:nth-child(5)")
		//sessionID
		s := strings.Split(e.ChildAttr("a[href]", "href"), "=")
		session_ID := s[len(s)-1]

		allSessions = append(allSessions, voting_session{Session_ID: session_ID, Timestamp: timestamp, Text: text})
		
		})

    c.Visit("https://www.nrsr.sk/web/default.aspx?SectionId=108")

    CurrentVotingSessions := &sessions{VotingSessions:allSessions}
	b, err := json.Marshal(CurrentVotingSessions)	

	if err != nil {
		log.Fatal(err)
	}

	return b
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
