package repository

type QueryOptions struct {
	Filters  map[string]any
	Sort     []SortOption
	Limit    int
	Offset   int
}

type SortOption struct {
	Field     string
	Direction SortDirection
}

type SortDirection int

const (
	SortAscending SortDirection = iota
	SortDescending
)

type QueryOption func(*QueryOptions)

func WithFilter(field string, value any) QueryOption {
	return func(opts *QueryOptions) {
		if opts.Filters == nil {
			opts.Filters = make(map[string]any)
		}
		opts.Filters[field] = value
	}
}

// WithFilters adds multiple filter conditions
func WithFilters(filters map[string]any) QueryOption {
	return func(opts *QueryOptions) {
		if opts.Filters == nil {
			opts.Filters = make(map[string]any)
		}
		for k, v := range filters {
			opts.Filters[k] = v
		}
	}
}

// WithSort adds a sort condition
func WithSort(field string, direction SortDirection) QueryOption {
	return func(opts *QueryOptions) {
		opts.Sort = append(opts.Sort, SortOption{
			Field:     field,
			Direction: direction,
		})
	}
}

// WithPagination adds pagination
func WithPagination(limit, offset int) QueryOption {
	return func(opts *QueryOptions) {
		opts.Limit = limit
		opts.Offset = offset
	}
}

// ApplyOptions applies all query options
func ApplyOptions(opts ...QueryOption) QueryOptions {
	options := QueryOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}