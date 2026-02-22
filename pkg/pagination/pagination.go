package pagination

import (
	"honda-leasing-api/pkg/response"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

// Normalize ensures page and limit values are within valid ranges.
func Normalize(page, limit int) (int, int) {
	if page <= 0 {
		page = DefaultPage
	}
	if limit <= 0 || limit > MaxLimit {
		limit = DefaultLimit
	}
	return page, limit
}

// GetOffset calculates the database offset from page and limit.
func GetOffset(page, limit int) int {
	if page <= 0 {
		page = DefaultPage
	}
	if limit <= 0 {
		limit = DefaultLimit
	}
	return (page - 1) * limit
}

// BuildMeta constructs PaginationMeta from the given parameters.
func BuildMeta(page, limit int, total int64) response.PaginationMeta {
	hasMore := (int64(page) * int64(limit)) < total
	return response.PaginationMeta{
		Page:    page,
		Limit:   limit,
		Total:   total,
		HasMore: hasMore,
	}
}
