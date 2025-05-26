# Crypto Data Collector

A Go-based cryptocurrency data collection system that fetches data from various sources and stores it in MongoDB.

## Features

- Collects cryptocurrency price data from CoinGecko API
- Uses RabbitMQ for reliable data processing
- Stores data in MongoDB
- Supports automatic error retry
- Separates data collection and processing

## Prerequisites

Before running the application, make sure you have the following services installed and running:

1. MongoDB (version 4.0 or later)
2. RabbitMQ (version 3.8 or later)

## Configuration

The application uses a YAML configuration file (`config.yaml`) with the following structure:

```yaml
database:
  mongodb:
    uri: "mongodb://localhost:27017"
    database: "investutil"

queue:
  rabbitmq:
    uri: "amqp://guest:guest@localhost:5672/"
    exchange: "crypto"
    exchange_type: "direct"
    queue: "crypto_data"
    routing_key: "crypto.prices"

collector:
  schedule: "0 0 * * *"  # Run at midnight every day
```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/investutil-gocrawler.git
cd investutil-gocrawler
```

2. Install dependencies:
```bash
go mod download
```

## Running the Application

The application can run in two modes:

### 1. Collector Mode

This mode fetches data from external sources and sends it to RabbitMQ:

```bash
go run cmd/main/main.go -mode collect
```

### 2. Processor Mode

This mode processes data from RabbitMQ and stores it in MongoDB:

```bash
go run cmd/main/main.go -mode process
```

### Additional Options

- `-config`: Specify a custom config file path (default: "config.yaml")
```bash
go run cmd/main/main.go -mode collect -config /path/to/config.yaml
```

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

## Error Handling

The system includes automatic retry mechanisms:
- Collection errors: The collector will log errors and exit
- Processing errors: Failed messages will be retried up to 3 times with increasing delays (1s, 2s, 3s)
- After 3 failed attempts, messages will be requeued

## Monitoring

To check if the system is running properly:

1. RabbitMQ Management Console:
   - Open `http://localhost:15672` (default credentials: guest/guest)
   - Check queue status under "Queues" tab

2. MongoDB:
   - Use MongoDB Compass or CLI to check stored data:
   ```bash
   mongosh
   use investutil
   db.bitcoin_prices.find().sort({timestamp: -1}).limit(1)
   ```

## Troubleshooting

Common issues and solutions:

1. Connection errors:
   - Ensure MongoDB is running: `systemctl status mongodb`
   - Ensure RabbitMQ is running: `systemctl status rabbitmq-server`

2. Data not being processed:
   - Check if collector is running in collect mode
   - Check if processor is running in process mode
   - Check RabbitMQ queue status

## License

[MIT License](LICENSE)
