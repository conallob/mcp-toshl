package toshl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURL        = "https://api.toshl.com"
	defaultPerPage = 200
)

// Client is a Toshl API client.
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates a new Toshl API client using a personal access token.
// The token is sent as the username in HTTP Basic Auth with an empty password,
// as per the Toshl API documentation.
func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) get(path string, params url.Values) ([]byte, error) {
	reqURL := baseURL + path
	if len(params) > 0 {
		reqURL += "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(c.token, "")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// ListEntries returns a list of entries filtered by the given parameters.
func (c *Client) ListEntries(p ListEntriesParams) ([]Entry, error) {
	params := url.Values{}
	if p.From != "" {
		params.Set("from", p.From)
	}
	if p.To != "" {
		params.Set("to", p.To)
	}
	if p.Account != "" {
		params.Set("account", p.Account)
	}
	if p.Category != "" {
		params.Set("category", p.Category)
	}
	perPage := defaultPerPage
	if p.PerPage > 0 {
		perPage = p.PerPage
	}
	params.Set("per_page", strconv.Itoa(perPage))
	if p.Page > 0 {
		params.Set("page", strconv.Itoa(p.Page))
	}

	body, err := c.get("/entries", params)
	if err != nil {
		return nil, err
	}

	var entries []Entry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("parsing entries: %w", err)
	}
	return entries, nil
}

// GetEntry returns a single entry by ID.
func (c *Client) GetEntry(id string) (*Entry, error) {
	body, err := c.get("/entries/"+id, nil)
	if err != nil {
		return nil, err
	}

	var entry Entry
	if err := json.Unmarshal(body, &entry); err != nil {
		return nil, fmt.Errorf("parsing entry: %w", err)
	}
	return &entry, nil
}

// ListAccounts returns all accounts.
func (c *Client) ListAccounts(p ListParams) ([]Account, error) {
	params := url.Values{}
	perPage := defaultPerPage
	if p.PerPage > 0 {
		perPage = p.PerPage
	}
	params.Set("per_page", strconv.Itoa(perPage))
	if p.Page > 0 {
		params.Set("page", strconv.Itoa(p.Page))
	}

	body, err := c.get("/accounts", params)
	if err != nil {
		return nil, err
	}

	var accounts []Account
	if err := json.Unmarshal(body, &accounts); err != nil {
		return nil, fmt.Errorf("parsing accounts: %w", err)
	}
	return accounts, nil
}

// GetAccount returns a single account by ID.
func (c *Client) GetAccount(id string) (*Account, error) {
	body, err := c.get("/accounts/"+id, nil)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := json.Unmarshal(body, &account); err != nil {
		return nil, fmt.Errorf("parsing account: %w", err)
	}
	return &account, nil
}

// ListCategories returns all categories.
func (c *Client) ListCategories(p ListParams) ([]Category, error) {
	params := url.Values{}
	perPage := defaultPerPage
	if p.PerPage > 0 {
		perPage = p.PerPage
	}
	params.Set("per_page", strconv.Itoa(perPage))
	if p.Page > 0 {
		params.Set("page", strconv.Itoa(p.Page))
	}

	body, err := c.get("/categories", params)
	if err != nil {
		return nil, err
	}

	var categories []Category
	if err := json.Unmarshal(body, &categories); err != nil {
		return nil, fmt.Errorf("parsing categories: %w", err)
	}
	return categories, nil
}

// ListTags returns all tags.
func (c *Client) ListTags(p ListParams) ([]Tag, error) {
	params := url.Values{}
	perPage := defaultPerPage
	if p.PerPage > 0 {
		perPage = p.PerPage
	}
	params.Set("per_page", strconv.Itoa(perPage))
	if p.Page > 0 {
		params.Set("page", strconv.Itoa(p.Page))
	}

	body, err := c.get("/tags", params)
	if err != nil {
		return nil, err
	}

	var tags []Tag
	if err := json.Unmarshal(body, &tags); err != nil {
		return nil, fmt.Errorf("parsing tags: %w", err)
	}
	return tags, nil
}

// ListBudgets returns all budgets.
func (c *Client) ListBudgets(p ListParams) ([]Budget, error) {
	params := url.Values{}
	perPage := defaultPerPage
	if p.PerPage > 0 {
		perPage = p.PerPage
	}
	params.Set("per_page", strconv.Itoa(perPage))
	if p.Page > 0 {
		params.Set("page", strconv.Itoa(p.Page))
	}

	body, err := c.get("/budgets", params)
	if err != nil {
		return nil, err
	}

	var budgets []Budget
	if err := json.Unmarshal(body, &budgets); err != nil {
		return nil, fmt.Errorf("parsing budgets: %w", err)
	}
	return budgets, nil
}
