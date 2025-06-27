# WPP Analytics CLI

A command-line tool for fetching analytics and template analytics data from Facebook Graph API.

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

The CLI supports two modes: **analytics** (default) and **template analytics**.

### Basic Analytics (Default Mode)

```bash
./wppanalyticscli -wbaid=<WBA_ID> -start=<ISO_8601_DATE> -end=<ISO_8601_DATE> [-granularity=<GRANULARITY>]
```

### Template Analytics

```bash
./wppanalyticscli -mode=template -wbaid=<WBA_ID> -start=<ISO_8601_DATE> -end=<ISO_8601_DATE> -templates=<TEMPLATE_IDS> -metrics=<METRIC_TYPES> [-granularity=daily]
```

### Parameters

#### Common Parameters
- `-wbaid`: WhatsApp Business Account ID (required)
- `-start`: Start date in ISO-8601 format (required)
- `-end`: End date in ISO-8601 format (required)
- `-timezone`: Timezone for date display (optional, default: America/Sao_Paulo)
- `-mode`: Mode selection (optional, default: analytics)
  - Valid values: `analytics`, `template`

#### Analytics Mode Parameters
- `-granularity`: Data granularity (optional, default: DAY)
  - Valid values: `HALF_HOUR`, `DAY`, `MONTH`

#### Template Analytics Parameters
- `-templates`: Comma-separated template IDs (required for template mode)
- `-metrics`: Comma-separated metric types (required for template mode)
  - Valid values: `cost`, `clicked`, `delivered`, `read`, `sent`
- `-granularity`: Data granularity (default: daily)
  - Valid values: `daily`

### Examples

#### Basic Analytics

```bash
export FB_ACCESS_TOKEN="your_access_token_here"

# Daily analytics (default)
./wppanalyticscli -wbaid=932157148829117 -start=2025-06-20 -end=2025-06-24

# Monthly granularity
./wppanalyticscli -wbaid=932157148829117 -start=2025-01-01 -end=2025-06-30 -granularity=MONTH

# Half-hour granularity
./wppanalyticscli -wbaid=932157148829117 -start=2025-06-24T00:00:00Z -end=2025-06-24T23:59:59Z -granularity=HALF_HOUR
```

#### Template Analytics

```bash
# Template analytics with all metrics
./wppanalyticscli -mode=template -wbaid=932157148829117 -start=2025-06-20 -end=2025-06-24 -templates=1026573095658757 -metrics=cost,clicked,delivered,read,sent

# Multiple templates
./wppanalyticscli -mode=template -wbaid=932157148829117 -start=2025-06-20 -end=2025-06-24 -templates=1026573095658757,1234567890123456 -metrics=delivered,read,clicked

# Specific metrics only
./wppanalyticscli -mode=template -wbaid=932157148829117 -start=2025-06-20 -end=2025-06-24 -templates=1026573095658757 -metrics=cost,clicked
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

