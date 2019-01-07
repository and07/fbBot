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

// const ...
const (
	EndPoint = "https://graph.facebook.com/v2.6/me/messages"
)

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/webhook", webhookHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Facebook Bot")
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
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
					if message.Message.Text != "" {
						senderID := message.Sender.ID
						if message.Message.Text == "generic" {
							sendGenericMessage(senderID)
						} else {
							sentTextMessage(senderID, message.Message.Text)
						}
					}
				}
			}
		}
		fmt.Fprintf(w, "Success")
	}
}
