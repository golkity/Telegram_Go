package database

type Stock struct {
	Name    string
	Price   float64
	Percent float64
}

type Portfolio struct {
	Stocks map[string]Stock
	Total  float64
}

type UserState struct {
	State        string
	StockName    string
	StockPrice   float64
	StockPercent float64
}
