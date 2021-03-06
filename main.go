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
	EndPoint            = "https://graph.facebook.com/v2.11/me/"
	FOXNEWS             = "FOXNEWS"
	OZILOO              = "OZILOO"
	GET_STARTED_PAYLOAD = "GET_STARTED_PAYLOAD"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	foxnews := NewFoxnews()
	defer foxnews.Closed()

	oziloo := NewOziloo()
	defer oziloo.Closed()

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/webhook", webhookHandler(foxnews, oziloo))
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)

	m := msg{}
	q := m.SendMSG(oziloo)
	defer func() {
		close(q)
	}()
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Facebook Bot")
}

func webhookHandler(rssFN, rssOZ Rsser) func(w http.ResponseWriter, r *http.Request) {

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
					senderID := message.Sender.ID
					if message.Message != nil {
						messageText := message.Message.Text
						if messageText != "" {

							switch messageText {
							case "generic":
								sendGenericMessage(senderID)
							case foxnews:
								sendGenericRssMessage(senderID, rssFN)
							case oziloo:
								sendGenericRssMessage(senderID, rssOZ)
							default:
								sentTextMessage(senderID, message.Message.Text)
							}
						}
					} else if message.Postback != nil {
						switch message.Postback.Payload {
						case FOXNEWS:
							sendGenericRssMessage(senderID, rssFN)
						case OZILOO:
							sendGenericRssMessage(senderID, rssOZ)
						case GET_STARTED_PAYLOAD:
							sentTextMessage(senderID, "Get started")
						default:
							sentTextMessage(senderID, fmt.Sprintf("Postback called with payload: %s", message.Postback.Payload))
						}

					}
				}
			}
			fmt.Fprintf(w, "Success")
		}
	}
}
