// services/user-service/pkg/search/indexer.go
package search

import (
	"strings"
	"sync"

	"github.com/0xsj/fn-go/pkg/models"
)

// Indexer provides simple in-memory search capabilities for users
type Indexer struct {
	userIndex map[string][]string // Maps search terms to user IDs
	mutex     sync.RWMutex
}

// NewIndexer creates a new search indexer
func NewIndexer() *Indexer {
	return &Indexer{
		userIndex: make(map[string][]string),
	}
}

// IndexUser indexes a user for searching
func (i *Indexer) IndexUser(user *models.User) {
	if user == nil {
		return
	}

	i.mutex.Lock()
	defer i.mutex.Unlock()

	// Extract searchable terms
	terms := i.extractTerms(user)

	// Add user ID to each term's entry
	for _, term := range terms {
		// Skip very short terms
		if len(term) < 2 {
			continue
		}

		// Convert to lowercase
		term = strings.ToLower(term)

		// Check if the term already exists
		userIDs, exists := i.userIndex[term]
		if !exists {
			i.userIndex[term] = []string{user.ID}
			continue
		}

		// Check if user is already indexed for this term
		found := false
		for _, id := range userIDs {
			if id == user.ID {
				found = true
				break
			}
		}

		if !found {
			i.userIndex[term] = append(userIDs, user.ID)
		}
	}
}

// SearchUsers searches for users matching the query
func (i *Indexer) SearchUsers(query string, limit int) []string {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	// Process query
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return []string{}
	}

	// Split query into terms
	queryTerms := strings.Fields(query)

	// Find matching user IDs for each term
	matchesPerTerm := make([]map[string]bool, 0, len(queryTerms))

	for _, term := range queryTerms {
		termMatches := make(map[string]bool)

		// Look for exact matches
		if userIDs, exists := i.userIndex[term]; exists {
			for _, id := range userIDs {
				termMatches[id] = true
			}
		}

		// Look for prefix matches
		for indexedTerm, userIDs := range i.userIndex {
			if strings.HasPrefix(indexedTerm, term) {
				for _, id := range userIDs {
					termMatches[id] = true
				}
			}
		}

		if len(termMatches) > 0 {
			matchesPerTerm = append(matchesPerTerm, termMatches)
		}
	}

	// Find users that match all terms
	if len(matchesPerTerm) == 0 {
		return []string{}
	}

	// Start with all matches from the first term
	finalMatches := matchesPerTerm[0]

	// Intersect with matches from other terms
	for i := 1; i < len(matchesPerTerm); i++ {
		termMatches := matchesPerTerm[i]
		for id := range finalMatches {
			if !termMatches[id] {
				delete(finalMatches, id)
			}
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(finalMatches))
	for id := range finalMatches {
		result = append(result, id)
		if len(result) >= limit && limit > 0 {
			break
		}
	}

	return result
}

// RemoveUser removes a user from the index
func (i *Indexer) RemoveUser(userID string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	for term, userIDs := range i.userIndex {
		// Create a new slice without the user ID
		newIDs := make([]string, 0, len(userIDs))
		for _, id := range userIDs {
			if id != userID {
				newIDs = append(newIDs, id)
			}
		}

		if len(newIDs) == 0 {
			// If no users left for this term, remove it
			delete(i.userIndex, term)
		} else {
			i.userIndex[term] = newIDs
		}
	}
}

// extractTerms extracts searchable terms from a user
func (i *Indexer) extractTerms(user *models.User) []string {
	terms := []string{
		user.ID,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.FirstName + " " + user.LastName,
	}
	
	// Add phone if not empty
	if user.Phone != "" {
		terms = append(terms, user.Phone)
	}

	return terms
}