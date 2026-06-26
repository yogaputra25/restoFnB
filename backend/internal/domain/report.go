package domain

type SalesReportItem struct {
	MenuItemID string  `json:"menu_item_id"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Revenue    float64 `json:"revenue"`
}

type SalesReportCategory struct {
	CategoryID   string            `json:"category_id"`
	CategoryName string            `json:"category_name"`
	Items        []SalesReportItem `json:"items"`
}

type SalesReportDaily struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Orders  int     `json:"orders"`
}

type SalesReportSummary struct {
	TotalRevenue float64 `json:"total_revenue"`
	TotalOrders  int     `json:"total_orders"`
	TotalItems   int     `json:"total_items"`
}

type SalesReportResponse struct {
	Period         string               `json:"period"`
	From           string               `json:"from"`
	To             string               `json:"to"`
	Summary        SalesReportSummary   `json:"summary"`
	Categories     []SalesReportCategory `json:"categories"`
	DailyBreakdown []SalesReportDaily   `json:"daily_breakdown"`
}
