package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"golang.org/x/net/html"
)

var wg sync.WaitGroup
var mu sync.Mutex

// to check if a link is active or dead by verifying its status code
func checkLink(link string) {
	resp, err := http.Get(link)
	if (err != nil) {
		fmt.Printf("Error: %s -> %v\n", link, err)
		return 
	}
	defer resp.Body.Close()

	if (resp.StatusCode >= 400) {
		fmt.Printf("DEAD LINK: %s -> %d\n", link, resp.StatusCode)
	} else {
		fmt.Printf("Ok LINK: %s -> %d\n", link, resp.StatusCode)
	}
}


func extractLinks(link string) ([]string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseLinks(resp.Body)
}

// extract all links that are avaliable in the html page
func parseLinks(body io.Reader) ([]string, error) {
	var links []string
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}
	var visit func(*html.Node)  // fun to iterate through the html code to identify the anchor tag
	visit = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := strings.TrimSpace(attr.Val)
					if (href != "" && !strings.HasPrefix(href, "#") && !strings.HasPrefix(href, "mailto:")) { // skip the anchor tag which does not have redorect link
						links = append(links,href)
					}
				}
			}
		}
		// check all children tags of each tag
		for c := n.FirstChild; c!= nil; c=c.NextSibling {
			visit(c)
		}
	}
	visit(doc)
	return links, nil
}

func classifyLinks(baseURL string, rawLinks []string) ([]string, []string) {
	var internal, external []string  // to store url of same domain and diff domain
	seen := make(map[string]bool) // keep track of visited urls

	base, err := url.Parse(baseURL) // parses a raw url into a [URL] structure or Object
	if err != nil {
		return nil, nil
	}

	for _, href := range rawLinks {
		link, err := url.Parse(href)
		if err != nil {
			continue
		}

		resolved := base.ResolveReference(link)

		fullURL := resolved.String()
		if seen[fullURL] {
			continue
		}
		seen[fullURL] = true

		if resolved.Host == base.Host {
			internal = append(internal, fullURL)
		} else {
			external = append(external, fullURL)
		}
	}

	return internal, external
}

// recursive web crawler
func crawl(link string, visited map[string]bool) {
	defer wg.Done() // defer - executed when the function is completed
	// once the function is finished, the crawl is marked as completed

	mu.Lock() // mutex lock
	if visited[link] {
		mu.Unlock()
		return
	}

	visited[link] = true
	mu.Unlock()

	checkLink(link) // to check wheather the link is dead or not

	links, err := extractLinks(link)
	if err != nil {
		fmt.Printf("Failed to extract links from %s: %v\n", link, err)
		return
	}

	internal, external := classifyLinks(link, links)

	for _, ex := range external {
		checkLink(ex) // just check, not crawling for external links
	}

	for _, in := range internal {
		// for each internal link
		wg.Add(1) // increase wait group count
		go crawl(in, visited) // recursively crawl each link
	}
}

func main() {
	startURL := "https://scrape-me.dreamsofcode.io" // input url
	visited := make(map[string]bool)

	wg.Add(1)
	go crawl(startURL, visited)

	wg.Wait() // blocks main fun from exiting until all goroutines have called wg.Done()
}