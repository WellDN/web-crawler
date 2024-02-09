package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	// URL to start crawling
	url := "https://upload.wikimedia.org/wikipedia/commons/b/b6/Image_created_with_a_mobile_phone.png"

	links, h1Text, err := crawl(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Links:")
	for _, link := range links {
		fmt.Println(link)
	}

	fmt.Println("\n<h1> tag Text:")
	fmt.Println(h1Text)
}

func crawl(url string) ([]string, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	links, h1Text := extractElements(resp.Body, "a", "h1", "img")

	return links, h1Text, nil
}

func extractElements(body io.Reader, tags ...string) ([]string, string) {
	tokenizer := html.NewTokenizer(body)
	links := []string{}
	var h1Text string

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			// End of the document
			return links, h1Text
		case html.StartTagToken, html.EndTagToken:
			token := tokenizer.Token()
			for _, tag := range tags {
				if token.Data == tag {
					// Grab the link inside of <a href="example.com">
					// TODO: make it so you grab the whole page inside instead of just the link.
					if token.Data == "a" {
						for _, attr := range token.Attr {
							if attr.Key == "href" {
								link := strings.TrimSpace(attr.Val)
								links = append(links, link)
							}
						}
						// Grab the text inside of <h1> tag
					} else if tag == "h1" {
						if tag == "h1" && tokenType == html.StartTagToken {
							tokenType = tokenizer.Next()
							h1Text = strings.TrimSpace(tokenizer.Token().Data)
						}
						// Download image
					} else if tag == "img" {
						if tokenType == html.StartTagToken {
							for _, attr := range token.Attr { // This is downloading every img, gonna make proper later but just to see if it downloads correctly
								if attr.Key == "src" {
									imageUrl := strings.TrimSpace(attr.Val)
									err := downloadImage(imageUrl)
									if err != nil {
                                        log.Printf("Error downloading image from %s: %v\n", imageUrl, err)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func downloadImage(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

    // This is just printing which one was suppose to be downloaded.
    // TODO: download the image instead of just printing which one is going to be downloaded
	log.Printf("Image downloaded from %s\n", url)

	return nil
}
