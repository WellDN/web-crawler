package main

import (
    "fmt"
    "log"
    "net/http"
    "golang.org/x/net/html"
    "strings"
    "io"
)

func main() {
    // URL to start crawling
    url := "https://example.com" 

    // Call the crawl function
    links, err := crawl(url)
    if err != nil {
        log.Fatal(err)
    }

    // Print the titles of the crawled pages
    for _, link := range links {
        fmt.Println(link)
    }
}

func crawl(url string) ([]string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    links := extractLinks(resp.Body)

    return links, nil
}

func extractLinks(body io.Reader) []string {
    tokenizer := html.NewTokenizer(body)
    links := []string{}

    for {
        tokenType := tokenizer.Next()
        switch tokenType {
        case html.ErrorToken:
            // End of the document
            return links
        case html.StartTagToken, html.EndTagToken:
            token := tokenizer.Token()
            if token.Data == "a" {
                for _, attr := range token.Attr {
                    if attr.Key == "href" {
                        link := strings.TrimSpace(attr.Val)
                        links = append(links, link)
                    }
                }
            }
        }
    }
}

