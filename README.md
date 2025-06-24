# WPP Analytics CLI

A command-line tool for fetching analytics data from Facebook Graph API.

## Prerequisites

- Go 1.19 or higher
- Facebook Graph API access token

## Installation

### Build from source

```bash
make build
```

This will create the `wppanalyticscli` binary in the current directory.

### Clean build artifacts

```bash
make clean
```

## Configuration

Set your Facebook Graph API access token as an environment variable:

```bash
export FB_ACCESS_TOKEN="your_access_token_here"
```

## Usage

```bash
./wppanalyticscli -wbaid=<WBA_ID> -start=<ISO_8601_DATE> -end=<ISO_8601_DATE> [-granularity=<GRANULARITY>]
```

### Parameters

- `-wbaid`: WhatsApp Business Account ID (required)
- `-start`: Start date in ISO-8601 format (required)
- `-end`: End date in ISO-8601 format (required)  
- `-granularity`: Data granularity (optional, default: DAY)
  - Valid values: `HALF_HOUR`, `DAY`, `MONTH`

### Examples

#### Basic usage with daily granularity
```bash
export FB_ACCESS_TOKEN="EAAOTgr7rkAoBO2E8lZAPCXM7CtNDBseZCODnxVmif1HzfZAIdTx1BH06KBhOuoNR8ZCxxfRZCXYv30QiOG96qA6bZBNJqZBIkYLCg5m1tK1J50sC4dLECnXD5dEKdPbJOanZAJvo5SF6i1ljuHOkS3cNIWbH2BUscSZARMhu2Of43pGv8UIvqq5n26nMQK3ZAI4GlzFUWREYQZCjUG7mlJ1UGEbMgnNkAvq8z4ZAfY7ZAfVYLXwJhgAZDZD"

./wppanalyticscli -wbaid=932157148829117 -start=2025-06-20T00:00:00Z -end=2025-06-24T00:00:00Z
```

#### With monthly granularity
```bash
./wppanalyticscli -wbaid=932157148829117 -start=2025-01-01T00:00:00Z -end=2025-06-30T00:00:00Z -granularity=MONTH
```

#### With half-hour granularity
```bash
./wppanalyticscli -wbaid=932157148829117 -start=2025-06-24T00:00:00Z -end=2025-06-24T23:59:59Z -granularity=HALF_HOUR
```

## Date Format

All dates must be in ISO-8601 format with timezone information:

- `2025-06-24T00:00:00Z` (UTC)
- `2025-06-24T10:30:00-05:00` (with timezone offset)
- `2025-06-24T15:45:30+02:00` (with timezone offset)

## Output

The tool outputs the API response as pretty-printed JSON to stdout. Error messages are sent to stderr.

## Development

### Running tests
```bash
make test
```

### Building for different platforms
```bash
make build-linux
make build-windows
make build-darwin
```

### Release
The project is configured with GoReleaser for automated releases. Push a git tag to trigger a release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Error Handling

The tool will exit with status code 1 and display an error message if:

- Required parameters are missing
- FB_ACCESS_TOKEN environment variable is not set
- Date parsing fails
- Invalid granularity is specified
- API request fails