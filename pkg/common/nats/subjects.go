// pkg/common/nats/subjects.go
package nats

import (
	"strings"
)

// Service names
const (
	ServiceUser         = "user"
	ServiceAuth         = "auth"
	ServiceNotification = "notification"
	ServiceIncident     = "incident"
	ServiceLocation     = "location"
	ServiceTest         = "test" // Added for testing
)

// Operation types
const (
	OpCreate = "create"
	OpUpdate = "update"
	OpDelete = "delete"
	OpGet    = "get"
	OpList   = "list"
	OpEvent  = "event"
)

// Test subjects for simple ping/pong and fizz/buzz testing
const (
	SubjectPing = "ping"
	SubjectPong = "pong"
	SubjectFizz = "fizz"
	SubjectBuzz = "buzz"
)

// BuildSubject builds a subject string from components
// Format: service.operation.entity.id
func BuildSubject(service, operation, entity string, id ...string) string {
	parts := []string{service, operation, entity}
	
	// Add ID if provided
	if len(id) > 0 && id[0] != "" {
		parts = append(parts, id[0])
	}
	
	return strings.Join(parts, ".")
}

// TestSubjects provides subject helpers for testing
type TestSubjects struct{}

// Ping returns the ping subject
func (s TestSubjects) Ping() string {
	return SubjectPing
}

// Pong returns the pong subject
func (s TestSubjects) Pong() string {
	return SubjectPong
}

// Fizz returns the fizz subject
func (s TestSubjects) Fizz() string {
	return SubjectFizz
}

// Buzz returns the buzz subject
func (s TestSubjects) Buzz() string {
	return SubjectBuzz
}

// UserSubjects provides subject helpers for user service
type UserSubjects struct{}

// Create returns subject for creating users
func (s UserSubjects) Create() string {
	return BuildSubject(ServiceUser, OpCreate, "user")
}

// Get returns subject for getting a user
func (s UserSubjects) Get(id string) string {
	return BuildSubject(ServiceUser, OpGet, "user", id)
}

// List returns subject for listing users
func (s UserSubjects) List() string {
	return BuildSubject(ServiceUser, OpList, "user")
}

// Update returns subject for updating a user
func (s UserSubjects) Update(id string) string {
	return BuildSubject(ServiceUser, OpUpdate, "user", id)
}

// Delete returns subject for deleting a user
func (s UserSubjects) Delete(id string) string {
	return BuildSubject(ServiceUser, OpDelete, "user", id)
}

// Event returns subject for user events
func (s UserSubjects) Event(event string) string {
	return BuildSubject(ServiceUser, OpEvent, event)
}

// AuthSubjects provides subject helpers for auth service
type AuthSubjects struct{}

// Login returns subject for login
func (s AuthSubjects) Login() string {
	return BuildSubject(ServiceAuth, "command", "login")
}

// Validate returns subject for token validation
func (s AuthSubjects) Validate() string {
	return BuildSubject(ServiceAuth, "command", "validate")
}

// NotificationSubjects provides subject helpers for notification service
type NotificationSubjects struct{}

// Send returns subject for sending notifications
func (s NotificationSubjects) Send() string {
	return BuildSubject(ServiceNotification, "command", "send")
}

// Subjects provides access to all subject helpers
type Subjects struct{}

// Test returns test subjects for ping/pong testing
func (s Subjects) Test() TestSubjects {
	return TestSubjects{}
}

// User returns user service subjects
func (s Subjects) User() UserSubjects {
	return UserSubjects{}
}

// Auth returns auth service subjects
func (s Subjects) Auth() AuthSubjects {
	return AuthSubjects{}
}

// Notification returns notification service subjects
func (s Subjects) Notification() NotificationSubjects {
	return NotificationSubjects{}
}

// GetSubjects returns a subjects helper
func GetSubjects() Subjects {
	return Subjects{}
}