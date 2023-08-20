package dtos

type PaginationQueries struct {
	Page       int    `query:"page" json:"page"`
	Count      int    `query:"count" json:"count"`
	SearchText string `query:"search" json:"searchText"`
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
