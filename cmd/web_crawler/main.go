package main

import (
    "net/http"
    "fmt"
    "bufio"
)

// Thats gonna be used to just a example our crawler its gonna take the html/metadata, download and display somewhere else
func main() {
    resp, err := http.Get("https://gobyexample.com")
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

    // Print the HTTP response status
    fmt.Println("Response status", resp.Status)

    // Print the first 5 lines of the response body
    scanner := bufio.NewScanner(resp.Body)
    for i := 0; scanner.Scan() && i < 5; i++ {
        fmt.Println(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }
}
