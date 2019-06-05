package main

import (
	"log"
	"sync"
	"time"

	"github.com/gosimple/slug"
	"github.com/mmcdole/gofeed"
)

// Oziloo ...
type Oziloo struct {
	ticker *time.Ticker
	mx     sync.RWMutex
	exit   chan struct{}
	data   map[string]*Post
	urlRss string
	title  string
	image  string
}

const (
	urlOzilooRss = "https://oziloo.com/es/rss"
	oziloo       = "oziloo"
)

// NewOziloo ...
func NewOziloo() *Oziloo {

	f := Oziloo{}
	f.urlRss = urlOzilooRss
	f.data = make(map[string]*Post)
	f.exit = make(chan struct{})
	f.data = f.getData()
	f.Start()
	return &f
}

// Start update cache
func (f *Oziloo) Start() {

	f.ticker = time.NewTicker(30 * time.Minute)
	go func() {

		for {
			select {
			case <-f.ticker.C:
				go func() {
					f.mx.Lock()
					defer f.mx.Unlock()
					defer log.Println("Tick")

					f.data = f.getData()

					//log.Printf("coin pair = %#v\n", marketName)
					//log.Printf("book = %#v\n", book)
					//log.Printf("history = %#v\n", history)
				}()
			case <-f.exit:
				log.Println("Exit")
				f.ticker.Stop()
				return
			}
		}
	}()
}

// Closed update cache
func (f *Oziloo) Closed() {
	f.exit <- struct{}{}
	close(f.exit)
}

// Name ...
func (f *Oziloo) Name() string {
	return oziloo
}

// GetRssData ...
func (f *Oziloo) GetRssData() PostPageData {
	f.mx.RLock()
	defer f.mx.RUnlock()

	return PostPageData{
		PageTitle: f.title,
		PageImage: f.image,
		Pages:     f.data,
	}
}

func (f *Oziloo) getData() map[string]*Post {
	fp := gofeed.NewParser()

	data := make(map[string]*Post)

	feed, errFeed := fp.ParseURL(f.urlRss)
	if errFeed != nil {
		log.Println(errFeed)
		return data
	}
	log.Println(feed.Title)
	log.Println(feed.Image.URL)

	log.Printf("%#v \n", feed.Items[0].Published)

	for _, v := range feed.Items {

		if len(v.Enclosures) > 0 {
			if _, ok := data[slug.Make(v.Title)]; !ok {

				t1, _ := time.Parse(time.RFC1123, v.Published)

				data[slug.Make(v.Title)] = &Post{
					Published:   t1.Unix(),
					Categories:  []string{"test"},
					Title:       v.Title,
					Slug:        slug.Make(v.Title),
					Link:        v.Link,
					Description: v.Description,
					Image:       v.Enclosures[0].URL,
					SourceImage: feed.Image.URL,
					SourceTitle: feed.Title,
				}
				log.Printf("%#v", data[slug.Make(v.Title)])
			}
		}

	}
	f.title = feed.Title
	f.image = feed.Image.URL
	return data
}
