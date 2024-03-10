package crawler 

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

// Crawl a web content from a website to another
func HandleRoot(w http.ResponseWriter, r *http.Request, url string) {
    // URL to start crawling
    links, h1Text, err := Crawl(url)

    // This is passing the crawled images to another website
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    htmlContent := "<html><head><title>Crawled Content</title></head><body><h1>Links:</h1><ul>"

    for _, link := range links {
        htmlContent += "<li>" + link + "</li>"
    }
    htmlContent += "</ul><h1><h1>Tag Text:</h1><p>" + h1Text + "</p></body></html>"

    w.Header().Set("Content-Type", "text/html")
    io.WriteString(w, htmlContent)
}

func HandleDownloaded(url string) {
    links, h1Text, err := Crawl(url)

    // This download the crawled images to the ~/Download directory
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

func Crawl(url string) ([]string, string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, "", err
    }
    defer resp.Body.Close()

    links, h1Text, err := ExtractElements(resp.Body)

    return links, h1Text, nil
}

// Read token data and returns accordingly with the html tag
func ExtractElements(body io.Reader) ([]string, string, error) {
    tokenizer := html.NewTokenizer(body)
    links := []string{}
    var h1Text string

    for {
        tokenType := tokenizer.Next()
        switch tokenType {
        case html.ErrorToken:
            // End of the document
            return links, h1Text, nil 
        case html.StartTagToken, html.SelfClosingTagToken:
            token := tokenizer.Token()
            // Grab the link inside of <a href="example.com">
            if token.Data == "a" {
                if tokenType != html.SelfClosingTagToken {
                    for _, attr := range token.Attr {
                        if attr.Key == "href" {
                            link := strings.TrimSpace(attr.Val)
                            links = append(links, link)
                        }
                    }
                }
                // Grab the text inside of <h1> tag
            } else if token.Data == "h1" && tokenType == html.StartTagToken {
                tokenType = tokenizer.Next()
                h1Text = strings.TrimSpace(tokenizer.Token().Data)
                // Download image
            } else if  token.Data == "img" {
                for _, attr := range token.Attr {
                    if attr.Key == "src" {
                        imageURL := strings.TrimSpace(attr.Val)
                        err := Download(imageURL)
                        if err != nil {
                            log.Printf("Error downloading image from %s: %v\n", imageURL, err)
                        }
                    }
                }
            }
        }
    }
}

func Download(url string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
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
        return err
    }
    defer file.Close()

    _, err = io.Copy(file, resp.Body)
    if err != nil {
        return err
    }

    fmt.Printf("Image downloaded to the %s path.", filepath)

    return nil
}
