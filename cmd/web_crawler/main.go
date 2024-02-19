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
    // TODO: Implement in args
    url := "https://www.google.com/search?sca_esv=d220b328ad44e505&sxsrf=ACQVn0-N_PSr68nLqelTse_qZtFFHdrn5w:1708180318217&q=random+image&tbm=isch&source=lnms&prmd=ivsnbmtz&sa=X&ved=2ahUKEwiTldHIy7KEAxXnCbkGHeOvAfAQ0pQJegQIFBAB&biw=960&bih=959&dpr=1"
    links, h1Text, err := Crawl(url)

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

    links, h1Text, err := extractElements(resp.Body)

    return links, h1Text, nil
}

func extractElements(body io.Reader, tags ...string) ([]string, string, error) {
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
