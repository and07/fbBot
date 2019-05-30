package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type msg struct{}

func (m *msg) SendMSG(rss Rsser) chan struct{} {
	q := make(chan struct{})
	ticker := time.NewTicker(1 * time.Hour)
	log.Println("SendMSG Started at", time.Now())

	go func() {

		for {
			select {
			case <-ticker.C:
				//Call the periodic function here.
				fmt.Println("SendMSG tick")
				msg := getMessageRss(rss)
				var msgs MessagesData
				msgs.Messages = append(msgs.Messages, msg)
				b, err := json.Marshal(msgs)
				if err != nil {
					log.Println(err)
					return
				}

				res, errRequestMessageСreatives := requestMessageСreatives(string(b))
				if errRequestMessageСreatives != nil {
					log.Println(errRequestMessageСreatives)
					return
				}

				messageCreativeID := res["message_creative_id"].(string)
				data := fmt.Sprintf(`{    
					"message_creative_id": %s,
					"notification_type": "SILENT_PUSH",
					"messaging_type": "MESSAGE_TAG",
					"tag": "NON_PROMOTIONAL_SUBSCRIPTION"
				}`, messageCreativeID)
				res, errRequestBroadcastMessages := requestBroadcastMessages(data)
				if errRequestBroadcastMessages != nil {
					log.Println(errRequestBroadcastMessages)
					return
				}

				broadcastID := res["broadcast_id"].(string)
				log.Println(" ", data, " - ", broadcastID)
			case <-q:
				log.Println("SendMSG Stop", time.Now())
				ticker.Stop()
				return
			}
		}
	}()
	return q
}

func sentTextMessage(senderID string, text string) {
	recipient := new(Recipient)
	recipient.ID = senderID
	m := new(SendMessageText)
	m.Recipient = *recipient
	m.Message.Text = text
	b, err := json.Marshal(m)
	if err != nil {
		log.Print(err)
		return
	}

	res, err := requestMessages(string(b))
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("res %#v", res)
}

func getMessageRss(rss Rsser) MessageData {
	var m MessageData

	rssData := rss.GetRssData()
	var elements []*Element

	for _, v := range rssData.Pages {
		el := &Element{
			Title:    v.Title,
			Subtitle: v.Description,
			ItemURL:  v.Link,
			ImageURL: v.Image,
		}

		button := &Button{
			Type:  "web_url",
			URL:   v.Link,
			Title: "Open Web URL",
		}
		el.Buttons = append(el.Buttons, button)
		elements = append(elements, el)
	}

	var attachment Attachment
	attachment.Type = AttachmentTypeTemplate
	attachment.Payload = &Payload{
		TemplateType: "generic",
		Elements:     elements,
	}
	m.Attachment = attachment

	return m
}

func sendGenericRssMessage(senderID string, rss Rsser) {
	recipient := new(Recipient)
	recipient.ID = senderID
	m := new(SendMessageGeneric)
	m.Recipient = *recipient
	m.Message = getMessageRss(rss)
	b, err := json.Marshal(m)
	if err != nil {
		log.Print(err)
		return
	}

	res, err := requestMessages(string(b))
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
			  image_url: "https://multimedia.bbycastatic.ca/multimedia/products/500x500/112/11264/11264431.jpg",
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
			  image_url: "https://cf1.s3.souqcdn.com/item/2016/12/14/12/01/89/49/item_XL_12018949_18012245.jpg",
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

	res, err := requestMessages(messageData)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("res %#v", res)

}

func request(urlPath, messageData string) (map[string]interface{}, error) {
	data := []byte(messageData)

	req, err := http.NewRequest("POST", EndPoint+urlPath, bytes.NewBuffer(data))
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

func requestBroadcastMessages(messageData string) (map[string]interface{}, error) {
	return request("broadcast_messages", messageData)
}

func requestMessageСreatives(messageData string) (map[string]interface{}, error) {
	return request("message_creatives", messageData)
}

func requestMessages(messageData string) (map[string]interface{}, error) {
	return request("messages", messageData)
}
