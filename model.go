package main

// ReceivedMessage ...
type ReceivedMessage struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

// Entry ...
type Entry struct {
	ID        string       `json:"id"`
	Time      int64        `json:"time"`
	Messaging []*Messaging `json:"messaging"`
}

// Messaging ...
type Messaging struct {
	Sender    *Sender    `json:"sender"`
	Recipient *Recipient `json:"recipient"`
	Timestamp int64      `json:"timestamp"`
	Message   *Message   `json:"message"`
	Delivery  *Delivery  `json:"delivery,omitempty"`
	Postback  *Postback  `json:"postback,omitempty"`
	Optin     *Optin     `json:"optin,empty"`
}

// Sender ...
type Sender struct {
	ID string `json:"id"`
}

// Recipient ...
type Recipient struct {
	ID string `json:"id"`
}

// AttachmentType ...
type AttachmentType string

const (
	AttachmentTypeTemplate AttachmentType = "template"
	AttachmentTypeImage    AttachmentType = "image"
	AttachmentTypeVideo    AttachmentType = "video"
	AttachmentTypeAudio    AttachmentType = "audio"
	AttachmentTypeLocation AttachmentType = "location"
)

// Attachment ...
type Attachment struct {
	Type    AttachmentType `json:"type"`
	Payload interface{}    `json:"payload,omitempty"`
}

// Element ...
type Element struct {
	Title    string    `json:"title"`
	Subtitle string    `json:"subtitle"`
	ItemURL  string    `json:"item_url"`
	ImageURL string    `json:"image_url"`
	Buttons  []*Button `json:"buttons"`
}

// Button ...
type Button struct {
	Type    string `json:"type"`
	URL     string `json:"url,omitempty"`
	Title   string `json:"title"`
	Payload string `json:"payload,omitempty"`
}

// Payload ...
type Payload struct {
	TemplateType string     `json:"template_type"`
	Elements     []*Element `json:"elements"`
}

// Message ...
type Message struct {
	ID          string        `json:"mid"`
	Text        string        `json:"text,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
	Seq         int           `json:"seq"`
}

// SendMessageText ...
type SendMessageText struct {
	Recipient Recipient `json:"recipient"`
	Message   struct {
		Text string `json:"text,omitempty"`
	} `json:"message"`
}

// SendMessageGeneric ...
type SendMessageGeneric struct {
	Recipient Recipient   `json:"recipient"`
	Message   MessageData `json:"message"`
}

// MessagesData ...
type MessagesData struct {
	Messages []MessageData `json:"messages"`
}

// MessageData ...
type MessageData struct {
	Attachment Attachment `json:"attachment,omitempty"`
}

// Delivery ...
type Delivery struct {
	MessageIDS []string `json:"mids"`
	Watermark  int64    `json:"watermark"`
	Seq        int      `json:"seq"`
}

// Postback ...
type Postback struct {
	Payload string `json:"payload"`
}

// Optin ...
type Optin struct {
	Ref string `json:"ref"`
}
