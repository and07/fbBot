package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var token = os.Getenv("TOKEN")
var rss = os.Getenv("RSS")

// const ...
const (
	EndPoint = "https://graph.facebook.com/v2.6/me/messages"
)

func main() {

	foxnews := NewFoxnews()
	defer foxnews.Closed()

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/webhook", webhookHandler(foxnews))
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Facebook Bot")
}

func webhookHandler(rss Rsser) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if r.URL.Query().Get("hub.verify_token") == token {
				fmt.Fprintf(w, r.URL.Query().Get("hub.challenge"))
			} else {
				fmt.Fprintf(w, "Error, wrong validation token")
			}
		}
		if r.Method == "POST" {
			var receivedMessage ReceivedMessage
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Print(err)
			}
			if err = json.Unmarshal(b, &receivedMessage); err != nil {
				log.Print(err)
			}

			for _, entry := range receivedMessage.Entry {
				for _, message := range entry.Messaging {
					if message.Message != nil {
						messageText := message.Message.Text
						if messageText != "" {
							senderID := message.Sender.ID
							switch messageText {
							case "generic":
								sendGenericMessage(senderID)
							case "foxnews":
								sendGenericRssMessage(senderID, rss)
							default:
								sentTextMessage(senderID, message.Message.Text)
							}
						}
					}
				}
			}
			fmt.Fprintf(w, "Success")
		}
	}
}
