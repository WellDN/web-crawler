package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/welldn/web-crawler/cmd/crawler"
)

func main() {
    args := os.Args[1:]
    if len(args) != 1 {
        fmt.Println("Give me one url to crawl")
    }

    url := args[0]

    crawler.HandleDownloaded(url)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        crawler.HandleRoot(w, r, url)
    })

    fmt.Println("Server started on: \thttp://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
