package main

import (
	"log"
	"sync"
	"time"

	"github.com/gosimple/slug"
	"github.com/mmcdole/gofeed"
)

type Foxnews struct {
	ticker *time.Ticker
	mx     sync.RWMutex
	exit   chan struct{}
	data   map[string]*Post
	urlRss string
	title  string
	image  string
}

const urlRss = "http://feeds.foxnews.com/foxnews/latest"

func NewFoxnews() *Foxnews {

	f := Foxnews{}
	f.urlRss = urlRss
	f.data = make(map[string]*Post)
	f.exit = make(chan struct{})
	f.data = f.getData()
	f.Start()
	return &f
}

func (f *Foxnews) Start() {

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

func (f *Foxnews) Closed() {
	f.exit <- struct{}{}
	close(f.exit)
}

func (f *Foxnews) Name() string {
	return "foxnews"
}

func (f *Foxnews) GetRssData() PostPageData {
	f.mx.RLock()
	defer f.mx.RUnlock()

	return PostPageData{
		PageTitle: f.title,
		PageImage: f.image,
		Pages:     f.data,
	}
}

func (f *Foxnews) getData() map[string]*Post {
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
	log.Printf("%#v \n", feed.Items[0].Categories)
	for _, v := range feed.Items {

		if v.Extensions["media"]["group"] != nil {
			log.Printf("%#v \n", v.Extensions["media"]["group"][0].Children["content"][0].Attrs["url"])

			if _, ok := data[slug.Make(v.Title)]; !ok {

				t1, _ := time.Parse(time.RFC1123, v.Published)

				data[slug.Make(v.Title)] = &Post{
					Published:   t1.Unix(),
					Categories:  v.Categories,
					Title:       v.Title,
					Slug:        slug.Make(v.Title),
					Link:        v.Link,
					Description: v.Description,
					Image:       v.Extensions["media"]["group"][0].Children["content"][0].Attrs["url"],
					SourceImage: "//global.fncstatic.com/static/orion/styles/img/fox-news/favicons/apple-touch-icon-60x60.png",
					SourceTitle: feed.Title,
				}

			}

		}

	}
	f.title = feed.Title
	f.image = feed.Image.URL
	return data
}
