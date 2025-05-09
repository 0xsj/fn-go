// pkg/repository/options.go
package repository

// QueryOptions contains parameters for filtering, sorting, and pagination
type QueryOptions struct {
	Filters map[string]interface{}
	Sort    []SortOption
	Limit   int
	Offset  int
}

// SortOption represents a field to sort by and its direction
type SortOption struct {
	Field     string
	Direction SortDirection
}

// SortDirection indicates ascending or descending sort
type SortDirection int

const (
	SortAscending SortDirection = iota
	SortDescending
)

// QueryOption is a function that modifies QueryOptions
type QueryOption func(*QueryOptions)

// WithFilter adds a filter condition
func WithFilter(field string, value interface{}) QueryOption {
	return func(opts *QueryOptions) {
		if opts.Filters == nil {
			opts.Filters = make(map[string]interface{})
		}
		opts.Filters[field] = value
	}
}

// WithFilters adds multiple filter conditions
func WithFilters(filters map[string]interface{}) QueryOption {
	return func(opts *QueryOptions) {
		if opts.Filters == nil {
			opts.Filters = make(map[string]interface{})
		}
		for k, v := range filters {
			opts.Filters[k] = v
		}
	}
}

// WithSort adds a sort condition
func WithSort(field string, direction SortDirection) QueryOption {
	return func(opts *QueryOptions) {
		opts.Sort = append(opts.Sort, SortOption{Field: field, Direction: direction})
	}
}

// WithPagination adds pagination
func WithPagination(limit, offset int) QueryOption {
	return func(opts *QueryOptions) {
		opts.Limit = limit
		opts.Offset = offset
	}
}

// ApplyOptions applies the provided query options
func ApplyOptions(opts ...QueryOption) QueryOptions {
	options := QueryOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}