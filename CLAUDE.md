# mcp-toshl

## Overview

`mcp-toshl` is a Model Context Protocol (MCP) server written in Go that provides read-only access to [Toshl Finance](https://toshl.com) financial transaction data via the [Toshl API](https://developer.toshl.com/docs/).

## Architecture

```
mcp-toshl/
├── cmd/mcp-toshl/main.go       # Entry point: flags, server init, tool registration
├── internal/toshl/
│   ├── client.go               # HTTP client wrapping the Toshl REST API
│   └── types.go                # Go structs for Toshl API resources
├── .github/workflows/          # CI/CD: build, release, Claude integration
├── .goreleaser.yml             # Multi-platform release + Homebrew tap publishing
└── Makefile                    # build, test, fmt, vet, lint targets
```

## Essential Setup

1. Obtain a [Toshl personal access token](https://toshl.com/account/apps/) from your account settings.
2. Set the token via environment variable or flag:
   ```sh
   export TOSHL_TOKEN=your_token_here
   mcp-toshl
   # or
   mcp-toshl -token your_token_here
   ```

## Key Components

### MCP Server (`cmd/mcp-toshl/main.go`)
- Uses [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go) for MCP protocol implementation
- Communicates over stdio (JSON-RPC 2.0) — stdout is reserved for protocol messages, use stderr for logs
- Build-time version injection via `ldflags` (`main.version`, `main.commit`, `main.date`)

### Toshl Client (`internal/toshl/client.go`)
- Authentication: HTTP Basic Auth with the personal access token as username, empty password
- Base URL: `https://api.toshl.com`
- Default page size: 200 items
- Pagination via `per_page` and `page` query parameters

### Available MCP Tools

| Tool | Description |
|------|-------------|
| `list_entries` | List entries with optional date range, account, and category filters |
| `get_entry` | Get a single entry by ID |
| `list_accounts` | List all financial accounts with balances |
| `get_account` | Get a single account by ID |
| `list_categories` | List all categories |
| `list_tags` | List all tags |
| `list_budgets` | List all active budgets |

## Development

```sh
make build    # compile binary to bin/mcp-toshl
make test     # run tests
make vet      # run go vet
make fmt      # run go fmt
make lint     # fmt + vet
make deps     # go mod download + tidy
```

## Release

Releases are automated via GoReleaser triggered by pushing a `v*.*.*` tag. The release workflow:
1. Builds binaries for Linux, macOS, and FreeBSD (amd64 + arm64)
2. Creates a GitHub Release with archives and checksums
3. Publishes a Homebrew formula to [conallob/homebrew-tap](https://github.com/conallob/homebrew-tap)

Requires `HOMEBREW_TAP_GITHUB_TOKEN` secret in the repository settings.

## Style Conventions

- Keep tools read-only — this server does not create or modify Toshl data
- Return errors via `mcp.NewToolResultError()`, not as Go errors
- Serialize all results as indented JSON via `json.MarshalIndent`
- Log to stderr only; stdout is the MCP JSON-RPC channel
