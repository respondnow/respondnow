package api

type Pagination struct {
	All        bool  `json:"all"`
	Index      int64 `json:"index"`
	Limit      int64 `json:"limit"`
	TotalPages int64 `json:"totalPages"`
	TotalItems int64 `json:"totalItems"`
}

func GetPagination(page, limit, total int64, all bool) Pagination {
	if all || limit == 0 {
		return Pagination{
			All: all,
		}
	}
	rem := total % limit
	var pageCount int64 = total / limit
	if rem != 0 {
		pageCount = pageCount + 1
	}
	return Pagination{
		All:        all,
		Index:      page,
		Limit:      limit,
		TotalPages: pageCount,
		TotalItems: total,
	}
}
