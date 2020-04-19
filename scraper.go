package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type votingInfo struct {
	Names   []string
	Parties []string
	Votes   []string
}

type votingSession struct {
	SessionID string
	Timestamp string
	Text      string
}

type sessions struct {
	VotingSessions []votingSession
}

func scrapeMeetingID() []byte {
	// Instantiate default collector
	c := colly.NewCollector()

	var allSessions []votingSession

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
		sessionID := s[len(s)-1]

		allSessions = append(allSessions, votingSession{SessionID: sessionID, Timestamp: timestamp, Text: text})

	})

	c.Visit("https://www.nrsr.sk/web/default.aspx?SectionId=108")

	CurrentVotingSessions := &sessions{VotingSessions: allSessions}
	b, err := json.Marshal(CurrentVotingSessions)

	if err != nil {
		log.Fatal(err)
	}

	return b
}

func scrapeMeeting(meetingID int) []byte {
	// Instantiate default collector
	c := colly.NewCollector()

	var nameSlc []string
	var partySlc []string
	var voteSlc []string

	// On every a element which has href attribute call callback
	party := ""

	c.OnHTML("td", func(e *colly.HTMLElement) {
		// name and vote description

		if len(e.Text) > 0 && string(e.Text[0]) != "[" && string(e.Text[0]) != "\n" {
			party = e.Text
			//fmt.Printf("%q", party)
		}

		if len(e.Text) > 0 && string(e.Text[0]) == "[" && string(e.Text[0]) != "\n" {
			voteType := e.Text[1]
			name := e.Text[4:]

			nameSlc = append(nameSlc, name)
			partySlc = append(partySlc, party)
			voteSlc = append(voteSlc, string(voteType))
		}

	})

	s := strconv.Itoa(meetingID)

	// Start scraping
	c.Visit("https://www.nrsr.sk/web/Default.aspx?sid=schodze/hlasovanie/hlasklub&ID=" + s)

	//output
	voteObj := &votingInfo{Names: nameSlc, Parties: partySlc, Votes: voteSlc}
	b, err := json.Marshal(voteObj)

	if err != nil {
		log.Fatal(err)
	}

	return b
}
