// pkg/common/db/util.go
package db

import (
	"database/sql"
	"fmt"
	"strings"
)

// BuildWhereClause builds a WHERE clause from filters
func BuildWhereClause(filters map[string]any) (string, []any) {
	if len(filters) == 0 {
		return "", nil
	}
	
	var clauses []string
	var args []any
	
	for field, value := range filters {
		// Handle nil values (IS NULL)
		if value == nil {
			clauses = append(clauses, fmt.Sprintf("%s IS NULL", field))
			continue
		}
		
		// Handle slices (IN)
		if values, ok := value.([]any); ok && len(values) > 0 {
			placeholders := make([]string, len(values))
			for i := range placeholders {
				placeholders[i] = "?"
			}
			clauses = append(clauses, fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ", ")))
			args = append(args, values...)
			continue
		}
		
		// Handle regular equality
		clauses = append(clauses, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}
	
	return "WHERE " + strings.Join(clauses, " AND "), args
}

// BuildOrderByClause builds an ORDER BY clause
func BuildOrderByClause(sorts []SortOption) string {
	if len(sorts) == 0 {
		return ""
	}
	
	var clauses []string
	
	for _, sort := range sorts {
		direction := "ASC"
		if sort.Direction == SortDescending {
			direction = "DESC"
		}
		clauses = append(clauses, fmt.Sprintf("%s %s", sort.Field, direction))
	}
	
	return "ORDER BY " + strings.Join(clauses, ", ")
}

// BuildLimitOffsetClause builds a LIMIT/OFFSET clause
func BuildLimitOffsetClause(limit, offset int) string {
	if limit <= 0 {
		return ""
	}
	
	if offset > 0 {
		return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
	}
	
	return fmt.Sprintf("LIMIT %d", limit)
}

// BuildQueryClauses builds complete query clauses
func BuildQueryClauses(filters map[string]any, sorts []SortOption, limit, offset int) (string, []any) {
	var clauses []string
	var args []any
	
	// WHERE clause
	whereClause, whereArgs := BuildWhereClause(filters)
	if whereClause != "" {
		clauses = append(clauses, whereClause)
		args = append(args, whereArgs...)
	}
	
	// ORDER BY clause
	orderByClause := BuildOrderByClause(sorts)
	if orderByClause != "" {
		clauses = append(clauses, orderByClause)
	}
	
	// LIMIT/OFFSET clause
	limitOffsetClause := BuildLimitOffsetClause(limit, offset)
	if limitOffsetClause != "" {
		clauses = append(clauses, limitOffsetClause)
	}
	
	return strings.Join(clauses, " "), args
}

// ScanRowsToMap scans SQL rows into maps
func ScanRowsToMap(rows *sql.Rows) ([]map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	
	var results []map[string]any
	
	for rows.Next() {
		// Create a slice of any to hold the values
		values := make([]any, len(columns))
		// Create a slice of pointers to those values
		scanArgs := make([]any, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		
		// Scan the result into the composite values
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		
		// Create the map and store values
		result := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			
			// Convert sql.RawBytes to string
			if b, ok := val.([]byte); ok {
				result[col] = string(b)
			} else {
				result[col] = val
			}
		}
		
		results = append(results, result)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return results, nil
}

// SortOption defines a sort option
type SortOption struct {
	Field     string
	Direction SortDirection
}

// SortDirection indicates sort direction
type SortDirection int

const (
	SortAscending SortDirection = iota
	SortDescending
)