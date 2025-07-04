# Dead Link Checker

A fast, reliable command-line tool written in Go for detecting broken links across websites. The tool crawls websites recursively, following internal links to specified depths, and identifies HTTP errors, timeouts, and unreachable resources.

## Features

- **Recursive Website Crawling**: Automatically discovers and follows internal links to map entire website structures
- **Configurable Depth Control**: Set maximum crawl depth to control scope and execution time
- **Comprehensive Link Detection**: Extracts and validates all types of links including relative paths, absolute URLs, and parent directory references
- **Dead Link Identification**: Detects HTTP errors (404, 500, etc.), timeouts, and unreachable resources
- **Professional CLI Interface**: Built with Cobra CLI for intuitive command-line usage
- **Fast and Concurrent**: Efficient HTTP client with proper error handling and timeout management
- **Duplicate Prevention**: Smart tracking to avoid crawling the same URLs multiple times

## Installation

### Prerequisites

- Go 1.19 or higher
- Internet connection for testing external links

### Build from Source

```bash
# Clone the repository
git clone https://github.com/your-username/dead-link-checker.git
cd dead-link-checker

# Build the executable
go build -o dead-link-checker .

# (Optional) Install globally
go install .
```

### Download Binary

Download the latest release from the [releases page](https://github.com/your-username/dead-link-checker/releases).

## Usage

### Basic Usage

```bash
# Check a website with default depth (2 levels)
./dead-link-checker check https://example.com

# Check with custom depth
./dead-link-checker check https://example.com -d 5

# Check with maximum depth of 1 (homepage only)
./dead-link-checker check https://example.com --depth 1
```

### Command Options

```bash
# Get general help
./dead-link-checker --help

# Get help for the check command
./dead-link-checker check --help
```

### Available Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--depth` | `-d` | `2` | Maximum crawl depth (0 = homepage only) |
| `--help` | `-h` | - | Show help information |

## How It Works

### 1. **Website Crawling**
The tool starts from the provided URL and recursively discovers internal links:
- Parses HTML content to extract all `<a href="">` links
- Resolves relative URLs (`/about`, `../contact`) to absolute URLs
- Identifies internal vs external links based on domain matching
- Respects the specified depth limit to prevent infinite crawling

### 2. **Link Classification**
- **Internal Links**: Same domain as the starting URL (followed recursively)
- **External Links**: Different domains (checked but not crawled)
- **Relative Paths**: Converted to absolute URLs using the base page URL
- **Fragment Links**: Links with anchors (`#section`) are preserved

### 3. **Dead Link Detection**
Each discovered link is validated using HTTP HEAD requests:
- **HTTP 2xx**: Link is alive and accessible
- **HTTP 4xx/5xx**: Link is broken (404 Not Found, 500 Server Error, etc.)
- **Timeout/Network Error**: Link is unreachable
- **Invalid URL**: Malformed URLs are reported as broken

## Example Output

```bash
$ ./dead-link-checker check https://example.com -d 2

Checking https://example.com
Collecting dead URLs:
https://example.com/broken-page
https://example.com/missing-image.jpg
https://external-site.com/dead-link
check called
```

## Architecture

The project follows clean architecture principles with clear separation of concerns:

```
├── cmd/           # CLI interface (Cobra commands)
├── internal/      # Core business logic
│   ├── crawler.go    # Website crawling and link discovery
│   ├── parser.go     # HTML parsing and link extraction
│   ├── checker.go    # Dead link detection and validation
│   └── scraper.go    # HTTP client and content fetching
└── main.go        # Application entry point
```

### Key Components

- **Crawler**: Recursive website traversal with depth control and duplicate prevention
- **Parser**: HTML parsing using Go's `golang.org/x/net/html` package
- **Checker**: HTTP validation using Go's standard `net/http` client
- **CLI**: Professional command-line interface built with Cobra

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test package
go test ./internal -v

# Run benchmarks
go test -bench=. ./internal
```

### Test Coverage

The project includes comprehensive test suites:
- **Unit Tests**: All core functions have dedicated test cases
- **Integration Tests**: End-to-end workflow validation
- **Edge Case Testing**: Invalid URLs, network errors, malformed HTML
- **Benchmark Tests**: Performance measurement for critical functions

### Code Quality

- **Go Modules**: Dependency management with `go.mod`
- **Error Handling**: Comprehensive error handling throughout
- **Documentation**: All exported functions include documentation
- **Type Safety**: Strong typing with minimal use of `interface{}`

## Dependencies

- **[Cobra](https://github.com/spf13/cobra)**: Modern CLI framework for Go
- **[golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html)**: HTML parsing library

## Performance Considerations

- **Memory Efficient**: Uses maps for tracking visited URLs to prevent memory leaks
- **Network Optimized**: HTTP HEAD requests minimize bandwidth usage
- **Timeout Handling**: Configurable timeouts prevent hanging requests
- **Duplicate Prevention**: Smart caching prevents redundant crawling

## Limitations

- **JavaScript-Rendered Content**: Only parses static HTML, does not execute JavaScript
- **Authentication**: Does not handle authenticated pages or login-protected content
- **Rate Limiting**: No built-in rate limiting (relies on Go's HTTP client defaults)
- **Robots.txt**: Does not respect robots.txt restrictions

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Go](https://golang.org/)
- CLI powered by [Cobra](https://github.com/spf13/cobra)
- HTML parsing by [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html)
