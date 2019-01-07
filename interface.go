package main

// Rsser ...
type Rsser interface {
	GetRssData() PostPageData
	Name() string
}

// Post ...
type Post struct {
	Title       string   `protobuf:"bytes,1,opt,name=Title" json:"Title,omitempty"`
	Slug        string   `protobuf:"bytes,2,opt,name=Slug" json:"Slug,omitempty"`
	Description string   `protobuf:"bytes,3,opt,name=Description" json:"Description,omitempty"`
	Link        string   `protobuf:"bytes,4,opt,name=Link" json:"Link,omitempty"`
	Image       string   `protobuf:"bytes,5,opt,name=Image" json:"Image,omitempty"`
	SourceImage string   `protobuf:"bytes,6,opt,name=SourceImage" json:"SourceImage,omitempty"`
	Published   int64    `protobuf:"varint,7,opt,name=Published" json:"Published,omitempty"`
	Categories  []string `protobuf:"bytes,8,rep,name=Categories" json:"Categories,omitempty"`
	SourceTitle string   `protobuf:"bytes,9,opt,name=SourceTitle" json:"SourceTitle,omitempty"`
}

// PostPageData ...
type PostPageData struct {
	PageTitle string           `protobuf:"bytes,1,opt,name=PageTitle" json:"PageTitle,omitempty"`
	PageImage string           `protobuf:"bytes,2,opt,name=PageImage" json:"PageImage,omitempty"`
	Pages     map[string]*Post `protobuf:"bytes,3,rep,name=Pages" json:"Pages,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}
