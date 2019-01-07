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

	res, err := request(string(b))
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("res %#v", res)
}

func sendGenericRssMessage(senderID string, rss Rsser) {
	recipient := new(Recipient)
	recipient.ID = senderID
	m := new(SendMessageGeneric)
	m.Recipient = *recipient

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
	m.Message.Attachment = attachment
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
