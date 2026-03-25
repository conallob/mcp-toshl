package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/conallob/mcp-toshl/internal/toshl"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var (
	// Build-time variables set by GoReleaser.
	version = "dev"
	commit  = "none"
	date    = "unknown"

	versionFlag = flag.Bool("version", false, "print version and exit")
	tokenFlag   = flag.String("token", "", "Toshl personal access token (or set TOSHL_TOKEN env var)")
)

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Printf("mcp-toshl %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", date)
		os.Exit(0)
	}

	token := *tokenFlag
	if token == "" {
		token = os.Getenv("TOSHL_TOKEN")
	}
	if token == "" {
		fmt.Fprintln(os.Stderr, "Error: Toshl API token required. Use -token flag or set TOSHL_TOKEN environment variable.")
		os.Exit(1)
	}

	log.SetPrefix("[mcp-toshl] ")
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(os.Stderr)

	client := toshl.NewClient(token)

	s := server.NewMCPServer("mcp-toshl", version,
		server.WithToolCapabilities(false),
	)

	registerTools(s, client)

	log.Printf("Starting mcp-toshl %s", version)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func registerTools(s *server.MCPServer, client *toshl.Client) {
	// list_entries
	s.AddTool(
		mcp.NewTool("list_entries",
			mcp.WithDescription("List Toshl financial entries (expenses and incomes). Negative amounts are expenses, positive are incomes."),
			mcp.WithString("from", mcp.Description("Start date in YYYY-MM-DD format (e.g. 2024-01-01)")),
			mcp.WithString("to", mcp.Description("End date in YYYY-MM-DD format (e.g. 2024-12-31)")),
			mcp.WithString("account", mcp.Description("Filter by account ID")),
			mcp.WithString("category", mcp.Description("Filter by category ID")),
			mcp.WithNumber("per_page", mcp.Description("Number of entries per page (default 200)")),
			mcp.WithNumber("page", mcp.Description("Page number for pagination (default 0)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			params := toshl.ListEntriesParams{
				From:     req.GetString("from", ""),
				To:       req.GetString("to", ""),
				Account:  req.GetString("account", ""),
				Category: req.GetString("category", ""),
				PerPage:  req.GetInt("per_page", 0),
				Page:     req.GetInt("page", 0),
			}
			entries, err := client.ListEntries(params)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(entries)
		},
	)

	// get_entry
	s.AddTool(
		mcp.NewTool("get_entry",
			mcp.WithDescription("Get a single Toshl entry by ID."),
			mcp.WithString("id", mcp.Required(), mcp.Description("The entry ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := req.GetString("id", "")
			entry, err := client.GetEntry(id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(entry)
		},
	)

	// list_accounts
	s.AddTool(
		mcp.NewTool("list_accounts",
			mcp.WithDescription("List all Toshl financial accounts with their balances and currencies."),
			mcp.WithNumber("per_page", mcp.Description("Number of accounts per page (default 200)")),
			mcp.WithNumber("page", mcp.Description("Page number for pagination (default 0)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			accounts, err := client.ListAccounts(toshl.ListParams{
				PerPage: req.GetInt("per_page", 0),
				Page:    req.GetInt("page", 0),
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(accounts)
		},
	)

	// get_account
	s.AddTool(
		mcp.NewTool("get_account",
			mcp.WithDescription("Get a single Toshl account by ID."),
			mcp.WithString("id", mcp.Required(), mcp.Description("The account ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id := req.GetString("id", "")
			account, err := client.GetAccount(id)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(account)
		},
	)

	// list_categories
	s.AddTool(
		mcp.NewTool("list_categories",
			mcp.WithDescription("List all Toshl categories used to classify entries."),
			mcp.WithNumber("per_page", mcp.Description("Number of categories per page (default 200)")),
			mcp.WithNumber("page", mcp.Description("Page number for pagination (default 0)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			categories, err := client.ListCategories(toshl.ListParams{
				PerPage: req.GetInt("per_page", 0),
				Page:    req.GetInt("page", 0),
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(categories)
		},
	)

	// list_tags
	s.AddTool(
		mcp.NewTool("list_tags",
			mcp.WithDescription("List all Toshl tags used to label entries."),
			mcp.WithNumber("per_page", mcp.Description("Number of tags per page (default 200)")),
			mcp.WithNumber("page", mcp.Description("Page number for pagination (default 0)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			tags, err := client.ListTags(toshl.ListParams{
				PerPage: req.GetInt("per_page", 0),
				Page:    req.GetInt("page", 0),
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(tags)
		},
	)

	// list_budgets
	s.AddTool(
		mcp.NewTool("list_budgets",
			mcp.WithDescription("List all active Toshl budgets with their spending totals."),
			mcp.WithNumber("per_page", mcp.Description("Number of budgets per page (default 200)")),
			mcp.WithNumber("page", mcp.Description("Page number for pagination (default 0)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			budgets, err := client.ListBudgets(toshl.ListParams{
				PerPage: req.GetInt("per_page", 0),
				Page:    req.GetInt("page", 0),
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(budgets)
		},
	)
}

func jsonResult(v interface{}) (*mcp.CallToolResult, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}
