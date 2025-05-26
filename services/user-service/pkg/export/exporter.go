// services/user-service/pkg/export/exporter.go
package export

import (
	"encoding/csv"
	"encoding/json"
	"io"

	"github.com/0xsj/fn-go/pkg/models"
)

// Format represents supported export formats
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Exporter exports user data in different formats
type Exporter struct {
	format Format
}

// NewExporter creates a new data exporter
func NewExporter(format Format) *Exporter {
	return &Exporter{
		format: format,
	}
}

// ExportUsers exports users to the specified writer
func (e *Exporter) ExportUsers(users []*models.User, w io.Writer) error {
	switch e.format {
	case FormatJSON:
		return e.exportJSON(users, w)
	case FormatCSV:
		return e.exportCSV(users, w)
	default:
		return e.exportJSON(users, w)
	}
}

// exportJSON exports users in JSON format
func (e *Exporter) exportJSON(users []*models.User, w io.Writer) error {
	// Create a sanitized version of users without sensitive data
	exportUsers := make([]map[string]any, len(users))
	
	for i, user := range users {
		exportUsers[i] = map[string]any{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"phone":     user.Phone,
			"role":      user.Role,
			"status":    user.Status,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		}
	}
	
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(exportUsers)
}

// exportCSV exports users in CSV format
func (e *Exporter) exportCSV(users []*models.User, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()
	
	// Write header
	header := []string{"ID", "Username", "Email", "First Name", "Last Name", "Phone", "Role", "Status", "Created At", "Updated At"}
	if err := writer.Write(header); err != nil {
		return err
	}
	
	// Write user data
	for _, user := range users {
		record := []string{
			user.ID,
			user.Username,
			user.Email,
			user.FirstName,
			user.LastName,
			user.Phone,
			string(user.Role),
			string(user.Status),
			user.CreatedAt.Format("2006-01-02 15:04:05"),
			user.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	
	return nil
}