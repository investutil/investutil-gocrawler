package main

import (
    "fmt"
    "github.com/gocolly/colly/v2"
    "regexp"
    "strings"
    "time"
)

func main() {
    // Create a new collector with a timeout
    c := colly.NewCollector(
        colly.Async(true),
        colly.MaxDepth(1),
        colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
    )

    // Set a request timeout
    c.SetRequestTimeout(30 * time.Second)

    // Add a callback function to the collector
    c.OnHTML("div[slot='marketTimeNotice'] span", func(e *colly.HTMLElement) {
        // Extracted content
        marketTimeNotice := e.Text
        fmt.Println("Extracted content:", marketTimeNotice)

        // Check if market is open or closed
        if strings.Contains(marketTimeNotice, "Market Open") {
            fmt.Println("Market is open")
        } else if strings.Contains(marketTimeNotice, "At close") {
            re := regexp.MustCompile(`At close: (.*) at (.*) PM EDT`)
            match := re.FindStringSubmatch(marketTimeNotice)

            if len(match) > 2 {
                date := match[1]
                time := match[2]
                fmt.Println("Date:", date)
                fmt.Println("Time:", time)
            }
            fmt.Println("Market is closed")
        } else {
            fmt.Println("Market status not found")
        }
    })

    // Handle request errors
    c.OnError(func(r *colly.Response, err error) {
        fmt.Println("Request error:", err)
    })

    // Start visiting the page
    fmt.Println("Visiting page...")
    c.Visit("https://finance.yahoo.com/quote/%5EGSPC/")

    // Wait until all asynchronous jobs are finished
    c.Wait()

    // Prevent the program from exiting too early
    select {}
}

