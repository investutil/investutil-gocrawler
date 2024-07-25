# web-scraping
The service of Investutil needs data. Data from web scraping of publicly accessible sites is a very important source, besides the data gathered from free or paid APIs.

## Architecture

The web scraping framwork, we will use  Colly (Golang)

Why not rust:
Go are easier to code, and web scraping requires more agility.

Reference project layout: https://github.com/golang-standards/project-layout

```
.
├── cmd
│   └── main
│       └── main.go       // Entry point of the application
├── internal
│   ├── collector         // Crawler logic
│   │   └── collector.go  // Implementation of the crawler
│   ├── database          // Database operations
│   │   └── mongodb.go    // MongoDB related code
│   ├── queue             // Queue operations
│   │   └── rabbitmq.go   // RabbitMQ related code
│   ├── config            // Configuration files
│   │   └── config.go     // Configuration loading and parsing
│   ├── models            // Data models
│   │   └── models.go     // Definitions of data structures
│   └── utils             // Utility functions
│       └── utils.go      // Common helper functions
└── go.mod                // Go module file
```
