package main

import (
	"encoding/json"
	"fmt"
	"os"
	"parse-azh-colly/lib/cleaner"

	"github.com/gocolly/colly"
)

type NewsItem struct {
	Title    string `json:"title"`
	Announce string `json:"announce"`
}

func main() {
	scrapeURL := "https://azh.kz/ru/news/in-atyrau"

	c := colly.NewCollector(colly.AllowedDomains("https://azh.kz", "azh.kz"))

	var newsItems []NewsItem

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting %s\n", r.URL)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error while scraping %s: %s\n", r.Request.URL, err.Error())
	})

	c.OnHTML("h3.news-list__title, span.news-list__announce", func(h *colly.HTMLElement) {
		switch {
		case h.Attr("class") == "news-list__title":
			newsItems = append(newsItems, NewsItem{Title: cleaner.CleanResult(h.Text)})
		case h.Attr("class") == "news-list__announce":
			if len(newsItems) > 0 {
				newsItems[len(newsItems)-1].Announce = cleaner.CleanResult(h.Text)
			}

		}
	})

	c.OnScraped(func(r *colly.Response) {
		jsonData, err := json.MarshalIndent(newsItems, "", "    ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}

		file, err := os.Create("news.json")
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		_, err = file.Write(jsonData)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		fmt.Println("Scraped data written to news.json")
	})

	c.Visit(scrapeURL)
}
