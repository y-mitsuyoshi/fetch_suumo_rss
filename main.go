package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// RSS structure to unmarshal the XML data
type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}

func main() {
	// URL of the RSS feed
	url := "https://suumo.jp/jj/bukken/ichiran/JJ012FC001/?ar=030&bs=020&ekInputCd=17640&ekInputNk=-1&ekInputNm=%E6%B8%8B%E8%B0%B7&ekInputTj=45&et=20&kb=1&km=0&kr=A&kt=6500&md=5&pc=100&pj=2&po=1&ta=13&rssFlg=1"

	// Get the RSS feed
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching the RSS feed:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the RSS feed:", err)
		return
	}

	// Unmarshal the XML data
	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		fmt.Println("Error unmarshalling the XML data:", err)
		return
	}

	// Regular expressions to extract information
	rePrice := regexp.MustCompile(`(\d+(?:,\d+)*万円)`)
	reSize := regexp.MustCompile(`(\d+\.\d+m&sup2;|\d+坪)`)
	reLayout := regexp.MustCompile(`(\d+LDK(?:\+\S*)?)`)
	reLocation := regexp.MustCompile(`東京都\S+`)

	// Process each item
	for _, item := range rss.Channel.Items {
		fmt.Println("リンク:", item.Link)

		// Extract information using regular expressions
		price := rePrice.FindAllString(item.Description, -1)
		size := reSize.FindAllString(item.Description, -1)
		layout := reLayout.FindAllString(item.Description, -1)
		location := reLocation.FindAllString(item.Description, -1)

		fmt.Println("価格:", strings.Join(price, ", "))
		fmt.Println("サイズ:", strings.Join(size, ", "))
		fmt.Println("間取り:", strings.Join(layout, ", "))
		fmt.Println("場所:", strings.Join(location, ", "))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
