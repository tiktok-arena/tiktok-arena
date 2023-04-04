package models

type PaginationQueries struct {
	Page       int    `query:"page"`
	Count      int    `query:"count"`
	SearchText string `query:"search"`
}

func ValidatePaginationQueries(queries *PaginationQueries) {
	if queries.Page <= 0 {
		queries.Page = 1
	}

	switch {
	case queries.Count > 50:
		queries.Count = 50
	case queries.Count <= 0:
		queries.Count = 20
	}
}
