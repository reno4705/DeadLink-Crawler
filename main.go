package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

var wg sync.WaitGroup
var mu sync.Mutex

func checkLink(link string) {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Printf("Error: %s -> %v\n", link, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {

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

func parseLinks(body io.Reader) ([]string, error) {
	var links []string
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}
	var visit func(*html.Node) 
	visit = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := strings.TrimSpace(attr.Val)
					if href != "" && !strings.HasPrefix(href, "#") && !strings.HasPrefix(href, "mailto:") { // skip the anchor tag which does not have redirect link
						links = append(links, href)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
	}
	visit(doc)
	return links, nil
}

func classifyLinks(baseURL string, rawLinks []string) ([]string, []string) {
	var internal, external []string 
	seen := make(map[string]bool)   

	base, err := url.Parse(baseURL) 
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

func crawl(link string, visited map[string]bool) {
	defer wg.Done()

	mu.Lock() 
	if visited[link] {
		mu.Unlock()
		return
	}

	visited[link] = true
	mu.Unlock()

	checkLink(link) 

	links, err := extractLinks(link)
	if err != nil {
		fmt.Printf("Failed to extract links from %s: %v\n", link, err)
		return
	}

	internal, external := classifyLinks(link, links)

	for _, ex := range external {
		checkLink(ex) 
	}

	for _, in := range internal {
		
		wg.Add(1)            
		go crawl(in, visited) 
	}
}

func main() {
	startURL := os.Args[1] // test link = https://scrape-me.dreamsofcode.io
	visited := make(map[string]bool)

	wg.Add(1)
	go crawl(startURL, visited)

	wg.Wait() 
}
