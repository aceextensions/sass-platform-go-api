package db

import (
	"fmt"
	"strings"
)

// FilterOptions defines a single filter criteria
type FilterOptions struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // eq, contains, gt, lt, gte, lte, in
	Value    interface{} `json:"value"`
}

// QueryOptions defines sorting, filtering and pagination parameters
type QueryOptions struct {
	Search       string          `json:"search"`
	SearchFields []string        `json:"searchFields"`
	Filters      []FilterOptions `json:"filters"`
	SortBy       string          `json:"sortBy"`    // camelCase from frontend
	SortOrder    string          `json:"sortOrder"` // asc, desc
	Page         int             `json:"page"`
	Limit        int             `json:"limit"`
}

// BuiltQuery contains the generated SQL fragments and arguments
type BuiltQuery struct {
	WhereClause string        // e.g. "WHERE name ILIKE $1 AND role = $2"
	Args        []interface{} // values for $1, $2, etc.
	OrderBy     string        // e.g. "ORDER BY created_at DESC"
	Limit       int
	Offset      int
}

// BuildQuery generates SQL fragments and arguments based on QueryOptions
func BuildQuery(options QueryOptions, startArgIndex int) BuiltQuery {
	var conditions []string
	args := make([]interface{}, 0)
	argIdx := startArgIndex

	// 1. Handle Global Search
	if options.Search != "" && len(options.SearchFields) > 0 {
		var searchConditions []string
		searchTerm := "%" + options.Search + "%"
		for _, field := range options.SearchFields {
			snakeField := camelToSnake(field)
			searchConditions = append(searchConditions, fmt.Sprintf("%s ILIKE $%d", snakeField, argIdx))
		}
		if len(searchConditions) > 0 {
			conditions = append(conditions, "("+strings.Join(searchConditions, " OR ")+")")
			args = append(args, searchTerm)
			argIdx++
		}
	}

	// 2. Handle Filters
	for _, filter := range options.Filters {
		snakeField := camelToSnake(filter.Field)
		var op string
		var val interface{} = filter.Value

		switch filter.Operator {
		case "eq":
			op = "="
		case "contains":
			op = "ILIKE"
			val = "%" + fmt.Sprint(filter.Value) + "%"
		case "gt":
			op = ">"
		case "lt":
			op = "<"
		case "gte":
			op = ">="
		case "lte":
			op = "<="
		default:
			continue
		}

		conditions = append(conditions, fmt.Sprintf("%s %s $%d", snakeField, op, argIdx))
		args = append(args, val)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 3. Handle Sorting
	orderBy := ""
	if options.SortBy != "" {
		snakeSort := camelToSnake(options.SortBy)
		order := "ASC"
		if strings.ToUpper(options.SortOrder) == "DESC" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("ORDER BY %s %s", snakeSort, order)
	} else {
		orderBy = "ORDER BY created_at DESC" // Default
	}

	// 4. Handle Pagination
	limit := options.Limit
	if limit <= 0 {
		limit = 25
	}
	page := options.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	return BuiltQuery{
		WhereClause: whereClause,
		Args:        args,
		OrderBy:     orderBy,
		Limit:       limit,
		Offset:      offset,
	}
}

// camelToSnake converts camelCase to snake_case
func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(valToLower(r))
	}
	return result.String()
}

func valToLower(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + ('a' - 'A')
	}
	return r
}
