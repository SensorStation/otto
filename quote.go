package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"nhooyr.io/websocket"
)

var (
	quote          Quote
	quoteCallbacks map[*websocket.Conn]func(q Quote)
)

type Quote struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Action string `json:"action"`
}

func init() {
	quoteCallbacks = make(map[*websocket.Conn]func(q Quote))
	quote = Quote{
		Author: "Rodney King",
		Title:  "",
		Text:   "Why can't we all just get along?",
		Action: "setQuote",
	}
}

func (q Quote) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(quote)

	case "POST", "PUT":
		text, ok := r.URL.Query()["text"]
		if !ok || len(text[0]) < 1 {
			log.Println("Url Param 'text' is missing...")
			fmt.Fprint(w, "URL param 'text' is missing")
			return
		}
		author, ok := r.URL.Query()["author"]
		if !ok || len(author[0]) < 1 {
			log.Println("URL Param 'author' is missing")
		}
		title, ok := r.URL.Query()["title"]
		if !ok || len(title[0]) < 1 {
			log.Println("Url Param 'title' is missing")
		}

		switch len(author) {
		case 0:
			// nothing

		case 1:
			quote.Author = author[0]

		default:
			quote.Author = strings.Join(author, ",")
		}

		quote.Title = ""
		if len(title) > 0 {
			quote.Title = title[0]
		}
		quote.Text = ""
		if len(text) > 0 {
			quote.Text = text[0]
		}

		log.Println("Sending the quote to Websocket")
		for _, cb := range quoteCallbacks {
			cb(quote)
		}
	}
}
