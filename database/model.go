package database

type Stock struct {
	Name    string
	Price   float64
	Percent float64
}

type Portfolio struct {
	Stock map[string]Stock
	Total float64
}
