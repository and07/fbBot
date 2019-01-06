package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var token = os.Getenv("TOKEN")

// const ...
const (
	EndPoint = "https://graph.facebook.com/v2.6/me/messages"
)

// ReceivedMessage ...
type ReceivedMessage struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

// Entry ...
type Entry struct {
	ID        string      `json:"id"`
	Time      int64       `json:"time"`
	Messaging []Messaging `json:"messaging"`
}

// Messaging ...
type Messaging struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   Message   `json:"message"`
}

// Sender ...
type Sender struct {
	ID string `json:"id"`
}

// Recipient ...
type Recipient struct {
	ID string `json:"id"`
}

// Message ...
type Message struct {
	MID  string `json:"mid"`
	Seq  int64  `json:"seq"`
	Text string `json:"text"`
}

// SendMessage ...
type SendMessage struct {
	Recipient Recipient `json:"recipient"`
	Message   struct {
		Text string `json:"text"`
	} `json:"message"`
}

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
		messagingEvents := receivedMessage.Entry[0].Messaging
		for _, event := range messagingEvents {
			senderID := event.Sender.ID
			if &event.Message != nil && event.Message.Text != "" {
				if event.Message.Text == "generic" {
					sendGenericMessage(senderID)
				} else {
					sentTextMessage(senderID, event.Message.Text)
				}

			}
		}
		fmt.Fprintf(w, "Success")
	}
}

func sentTextMessage(senderID string, text string) {
	recipient := new(Recipient)
	recipient.ID = senderID
	m := new(SendMessage)
	m.Recipient = *recipient
	m.Message.Text = text
	b, err := json.Marshal(m)
	if err != nil {
		log.Print(err)
		return
	}

	res, err := request(string(b))
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("res %#v", res)
}

func sendGenericMessage(senderID string) {
	messageData := fmt.Sprintf(`{
	  recipient: {
		id: "%s"
	  },
	  message: {
		attachment: {
		  type: "template",
		  payload: {
			template_type: "generic",
			elements: [{
			  title: "rift",
			  subtitle: "Next-generation virtual reality",
			  item_url: "https://www.oculus.com/en-us/rift/",
			  image_url: "http://messengerdemo.parseapp.com/img/rift.png",
			  buttons: [{
				type: "web_url",
				url: "https://www.oculus.com/en-us/rift/",
				title: "Open Web URL"
			  }, {
				type: "postback",
				title: "Call Postback",
				payload: "Payload for first bubble",
			  }],
			}, {
			  title: "touch",
			  subtitle: "Your Hands, Now in VR",
			  item_url: "https://www.oculus.com/en-us/touch/",
			  image_url: "http://messengerdemo.parseapp.com/img/touch.png",
			  buttons: [{
				type: "web_url",
				url: "https://www.oculus.com/en-us/touch/",
				title: "Open Web URL"
			  }, {
				type: "postback",
				title: "Call Postback",
				payload: "Payload for second bubble",
			  }]
			}]
		  }
		}
	  }
	}`, senderID)

	res, err := request(messageData)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("res %#v", res)

}

func request(messageData string) (map[string]interface{}, error) {
	data := []byte(messageData)

	req, err := http.NewRequest("POST", EndPoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Add("access_token", token)
	req.URL.RawQuery = values.Encode()
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{Timeout: time.Duration(30 * time.Second)}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var result map[string]interface{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
