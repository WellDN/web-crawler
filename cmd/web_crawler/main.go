package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"golang.org/x/net/html"
)

func main() {
	// URL to start crawling
	url := "https://upload.wikimedia.org/wikipedia/commons/b/b6/Image_created_with_a_mobile_phone.png"

	links, h1Text, errImg, err := crawl(url)
	if err != nil {
		log.Fatal(err)
	}

    errImg = downloadImage(url)
    if errImg != nil {
        log.Fatal(errImg)
    }
    fmt.Printf("Image downloaded successfully.")

	fmt.Println("Links:")
	for _, link := range links {
		fmt.Println(link)
	}

	fmt.Println("\n<h1> tag Text:")
	fmt.Println(h1Text)
}

func crawl(url string) ([]string, string, error, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", nil, err
	}
	defer resp.Body.Close()

	links, h1Text, errImg := extractElements(resp.Body, "a", "h1", "img")

	return links, h1Text, errImg, nil
}

func extractElements(body io.Reader, tags ...string) ([]string, string, error) {
	tokenizer := html.NewTokenizer(body)
	links := []string{}
	var h1Text string
	var errImg error 

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			// End of the document
			return links, h1Text, errImg 
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
							for _, attr := range token.Attr {
								if attr.Key == "src" {
									imageUrl := strings.TrimSpace(attr.Val)
									errImg = downloadImage(imageUrl)
									if errImg != nil {
                                        log.Printf("Error downloading image from %s: %v\n", imageUrl, errImg)
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
        log.Println(err)
    }
    defer resp.Body.Close()


    downloadsDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }
    downloadsDir = filepath.Join(downloadsDir, "Downloads")
    // If folder /Downloads doesn't exist, create new folder
    if _, err := os.Stat(downloadsDir); os.IsNotExist(err) {
        err := os.Mkdir(downloadsDir, 0750)
        if err != nil {
            log.Println(err)
        }
    }

    filename := filepath.Base(url)
    filepath := filepath.Join(downloadsDir, filename)

    file, err := os.Create(filepath)

    if err != nil {
        log.Println(err)
    }
    defer file.Close()

    _, err = io.Copy(file, resp.Body)

    if err != nil {
        log.Println(err)
    }

    fmt.Printf("Image downloaded to the %s path.", filepath)

    return nil
}
