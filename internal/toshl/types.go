package toshl

// Currency represents a currency with its code and exchange rate.
type Currency struct {
	Code  string  `json:"code"`
	Rate  float64 `json:"rate"`
	Fixed bool    `json:"fixed"`
}

// Entry represents a Toshl financial entry (expense or income).
type Entry struct {
	ID       string   `json:"id"`
	Amount   float64  `json:"amount"`
	Currency Currency `json:"currency"`
	Date     string   `json:"date"`
	Desc     string   `json:"desc,omitempty"`
	Account  string   `json:"account"`
	Category string   `json:"category"`
	Tags     []string `json:"tags,omitempty"`
	Modified string   `json:"modified"`
	Deleted  bool     `json:"deleted,omitempty"`
	Readonly bool     `json:"readonly,omitempty"`
}

// Account represents a Toshl financial account.
type Account struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Balance        float64  `json:"balance"`
	InitialBalance float64  `json:"initial_balance"`
	Currency       Currency `json:"currency"`
	Status         string   `json:"status"`
	Order          int      `json:"order"`
	Modified       string   `json:"modified"`
	Deleted        bool     `json:"deleted,omitempty"`
}

// Category represents a Toshl category.
type Category struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Modified string `json:"modified"`
	Deleted  bool   `json:"deleted,omitempty"`
}

// Tag represents a Toshl tag.
type Tag struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Modified string `json:"modified"`
	Deleted  bool   `json:"deleted,omitempty"`
}

// Budget represents a Toshl budget.
type Budget struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Amount   float64  `json:"amount"`
	Currency Currency `json:"currency"`
	From     string   `json:"from"`
	To       string   `json:"to"`
	Planned  float64  `json:"planned,omitempty"`
	Spending float64  `json:"spending,omitempty"`
	Modified string   `json:"modified"`
	Deleted  bool     `json:"deleted,omitempty"`
}

// ListEntriesParams holds parameters for listing entries.
type ListEntriesParams struct {
	From     string
	To       string
	Account  string
	Category string
	PerPage  int
	Page     int
}

// ListParams holds common pagination parameters.
type ListParams struct {
	PerPage int
	Page    int
}
