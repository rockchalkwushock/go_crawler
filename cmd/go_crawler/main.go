package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Crawler struct{}

func NewCrawler() *Crawler {
	return &Crawler{}
}

func (c *Crawler) Fetch(url string) (string, error) {
	resp, err := http.Get(url)

	fmt.Println("#############################")
	fmt.Printf("Response: %v\n", resp)
	fmt.Println("#############################")

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Crawler) ParseLinks(htmlContent string, keyword string) ([]string, error) {
	// Implement link parsing
	doc, err := html.Parse(strings.NewReader(htmlContent))

	if err != nil {
		return nil, err
	}

	var links []string
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if keyword != "" {
						if strings.Contains(n.FirstChild.Data, keyword) {
							links = append(links, a.Val)
							break
						}
					} else {
						links = append(links, a.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links, nil
}

func main() {
	var urls []string
	keywordPtr := flag.String("keyword", "", "Keyword to search for")
	flag.Parse()

	urls = flag.Args()
	keyword := *keywordPtr
	results := make(chan []string)

	fmt.Println("#############################")
	fmt.Printf("Querying keyword (%v) & URL(s) (%v)\n", keyword, urls)
	fmt.Println("#############################")

	for _, url := range urls {
		go func(url string) {
			crawler := NewCrawler()
			content, err := crawler.Fetch(url)

			if err != nil {
				fmt.Printf("Error fetching URL %s: %v\n", url, err)
				results <- []string{}
				return
			}

			links, err := crawler.ParseLinks(content, keyword)

			if err != nil {
				fmt.Printf("Error parsing links from URL %s: %v\n", url, err)
				results <- []string{}
				return
			}

			results <- links
		}(url)
	}

	for range urls {
		links := <-results
		fmt.Println("Found links:", links)
	}
}
