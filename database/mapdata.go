package database

var (
	UserPortfolios = make(map[int64]*Portfolio)
	UserStates     = make(map[int64]string)
)
