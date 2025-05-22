// pkg/common/config/loader.go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// EnvProvider loads configuration from environment variables
type EnvProvider struct {
	prefix       string
	requiredVars []string
	missingVars  []string
}

// NewEnvProvider creates a new environment variable provider
func NewEnvProvider(prefix string) *EnvProvider {
	return &EnvProvider{
		prefix:       prefix,
		requiredVars: []string{},
		missingVars:  []string{},
	}
}

// WithRequiredVars adds required variables to be validated
func (p *EnvProvider) WithRequiredVars(vars ...string) Provider {
	p.requiredVars = append(p.requiredVars, vars...)
	return p
}

// Validate checks that all required variables are present
func (p *EnvProvider) Validate() error {
	p.missingVars = []string{}
	
	for _, v := range p.requiredVars {
		key := p.formatKey(v)
		if val := os.Getenv(key); val == "" {
			p.missingVars = append(p.missingVars, key)
		}
	}
	
	if len(p.missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(p.missingVars, ", "))
	}
	
	return nil
}

// MissingVars returns the list of missing required variables
func (p *EnvProvider) MissingVars() []string {
	return p.missingVars
}

// Get retrieves a string value
func (p *EnvProvider) Get(key string) string {
	return os.Getenv(p.formatKey(key))
}

// GetDefault retrieves a string value with a default
func (p *EnvProvider) GetDefault(key string, defaultValue string) string {
	val := p.Get(key)
	if val == "" {
		return defaultValue
	}
	return val
}

// GetInt retrieves an integer value
func (p *EnvProvider) GetInt(key string) int {
	str := p.Get(key)
	if str == "" {
		return 0
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}

// GetIntDefault retrieves an integer value with a default
func (p *EnvProvider) GetIntDefault(key string, defaultValue int) int {
	str := p.Get(key)
	if str == "" {
		return defaultValue
	}
	
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	
	return val
}

// GetBool retrieves a boolean value
func (p *EnvProvider) GetBool(key string) bool {
	str := strings.ToLower(p.Get(key))
	return str == "true" || str == "1" || str == "yes"
}

// GetBoolDefault retrieves a boolean value with a default
func (p *EnvProvider) GetBoolDefault(key string, defaultValue bool) bool {
	str := p.Get(key)
	if str == "" {
		return defaultValue
	}
	str = strings.ToLower(str)
	return str == "true" || str == "1" || str == "yes"
}

// GetDuration retrieves a duration value
func (p *EnvProvider) GetDuration(key string) (time.Duration, error) {
	str := p.Get(key)
	if str == "" {
		return 0, fmt.Errorf("environment variable %s not set", p.formatKey(key))
	}
	
	return time.ParseDuration(str)
}

// GetDurationDefault retrieves a duration value with a default
func (p *EnvProvider) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	str := p.Get(key)
	if str == "" {
		return defaultValue
	}
	
	duration, err := time.ParseDuration(str)
	if err != nil {
		return defaultValue
	}
	
	return duration
}

// GetSlice retrieves a slice of strings
func (p *EnvProvider) GetSlice(key string, separator string) []string {
	str := p.Get(key)
	if str == "" {
		return []string{}
	}
	
	items := strings.Split(str, separator)
	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}
	
	return items
}

// GetMap retrieves a map from a string
func (p *EnvProvider) GetMap(key string, itemSeparator, keyValueSeparator string) map[string]string {
	str := p.Get(key)
	if str == "" {
		return map[string]string{}
	}
	
	result := make(map[string]string)
	items := strings.Split(str, itemSeparator)
	
	for _, item := range items {
		parts := strings.SplitN(item, keyValueSeparator, 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	
	return result
}

// formatKey formats a key with the prefix
func (p *EnvProvider) formatKey(key string) string {
	if p.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s_%s", p.prefix, key)
}

// FileProvider loads configuration from a file
type FileProvider struct {
	data         map[string]any
	requiredVars []string
	missingVars  []string
}

// NewJSONFileProvider creates a new JSON file provider
func NewJSONFileProvider(filepath string) (*FileProvider, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()
	
	var data map[string]any
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	
	return &FileProvider{
		data:         data,
		requiredVars: []string{},
		missingVars:  []string{},
	}, nil
}

// WithRequiredVars adds required variables to be validated
func (p *FileProvider) WithRequiredVars(vars ...string) Provider {
	p.requiredVars = append(p.requiredVars, vars...)
	return p
}

// Validate checks that all required variables are present
func (p *FileProvider) Validate() error {
	p.missingVars = []string{}
	
	for _, path := range p.requiredVars {
		if !p.hasPath(path) {
			p.missingVars = append(p.missingVars, path)
		}
	}
	
	if len(p.missingVars) > 0 {
		return fmt.Errorf("missing required configuration properties: %s", strings.Join(p.missingVars, ", "))
	}
	
	return nil
}

// MissingVars returns the list of missing required variables
func (p *FileProvider) MissingVars() []string {
	return p.missingVars
}

// Get retrieves a string value
func (p *FileProvider) Get(key string) string {
	value, _ := p.getPath(key)
	if value == nil {
		return ""
	}
	
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetDefault retrieves a string value with a default
func (p *FileProvider) GetDefault(key string, defaultValue string) string {
	value := p.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetInt retrieves an integer value
func (p *FileProvider) GetInt(key string) int {
	value, _ := p.getPath(key)
	if value == nil {
		return 0
	}
	
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return i
	default:
		return 0
	}
}

// GetIntDefault retrieves an integer value with a default
func (p *FileProvider) GetIntDefault(key string, defaultValue int) int {
	value, found := p.getPath(key)
	if !found {
		return defaultValue
	}
	
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return defaultValue
		}
		return i
	default:
		return defaultValue
	}
}

// GetBool retrieves a boolean value
func (p *FileProvider) GetBool(key string) bool {
	value, _ := p.getPath(key)
	if value == nil {
		return false
	}
	
	switch v := value.(type) {
	case bool:
		return v
	case string:
		str := strings.ToLower(v)
		return str == "true" || str == "1" || str == "yes"
	case int:
		return v != 0
	case float64:
		return v != 0
	default:
		return false
	}
}

// GetBoolDefault retrieves a boolean value with a default
func (p *FileProvider) GetBoolDefault(key string, defaultValue bool) bool {
	value, found := p.getPath(key)
	if !found {
		return defaultValue
	}
	
	switch v := value.(type) {
	case bool:
		return v
	case string:
		str := strings.ToLower(v)
		return str == "true" || str == "1" || str == "yes"
	case int:
		return v != 0
	case float64:
		return v != 0
	default:
		return defaultValue
	}
}

// GetDuration retrieves a duration value
func (p *FileProvider) GetDuration(key string) (time.Duration, error) {
	value := p.Get(key)
	if value == "" {
		return 0, fmt.Errorf("configuration property %s not found", key)
	}
	
	return time.ParseDuration(value)
}

// GetDurationDefault retrieves a duration value with a default
func (p *FileProvider) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	value := p.Get(key)
	if value == "" {
		return defaultValue
	}
	
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	
	return duration
}

// GetSlice retrieves a slice of strings
func (p *FileProvider) GetSlice(key string, separator string) []string {
	// First try to get as an array
	value, found := p.getPath(key)
	if !found {
		return []string{}
	}
	
	// If it's already an array in the JSON, use that
	if arr, ok := value.([]any); ok {
		result := make([]string, len(arr))
		for i, v := range arr {
			result[i] = fmt.Sprintf("%v", v)
		}
		return result
	}
	
	// Otherwise, treat as a string and split
	str := p.Get(key)
	if str == "" {
		return []string{}
	}
	
	items := strings.Split(str, separator)
	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}
	
	return items
}

// getPath retrieves a value from a nested path (e.g., "database.connection.host")
func (p *FileProvider) getPath(path string) (any, bool) {
	parts := strings.Split(path, ".")
	var current any = p.data
	
	for _, part := range parts {
		// Try to convert to map
		currentMap, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		
		// Try to get the next level
		current, ok = currentMap[part]
		if !ok {
			return nil, false
		}
	}
	
	return current, true
}

// hasPath checks if a path exists in the configuration
func (p *FileProvider) hasPath(path string) bool {
	_, found := p.getPath(path)
	return found
}

// LoadJSONFile loads a configuration from a JSON file
func LoadJSONFile(filepath string, cfg any) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	
	return nil
}

// CompositeProvider combines multiple providers
type CompositeProvider struct {
	providers    []Provider
	requiredVars []string
	missingVars  []string
}

// NewCompositeProvider creates a new composite provider
func NewCompositeProvider(providers ...Provider) *CompositeProvider {
	return &CompositeProvider{
		providers:    providers,
		requiredVars: []string{},
		missingVars:  []string{},
	}
}

// WithRequiredVars adds required variables to be validated
func (p *CompositeProvider) WithRequiredVars(vars ...string) Provider {
	p.requiredVars = append(p.requiredVars, vars...)
	return p
}

// Validate checks that all required variables are present
func (p *CompositeProvider) Validate() error {
	p.missingVars = []string{}
	
	for _, v := range p.requiredVars {
		found := false
		for _, provider := range p.providers {
			if value := provider.Get(v); value != "" {
				found = true
				break
			}
		}
		
		if !found {
			p.missingVars = append(p.missingVars, v)
		}
	}
	
	if len(p.missingVars) > 0 {
		return fmt.Errorf("missing required configuration properties: %s", strings.Join(p.missingVars, ", "))
	}
	
	return nil
}

// MissingVars returns the list of missing required variables
func (p *CompositeProvider) MissingVars() []string {
	return p.missingVars
}

// Get retrieves a string value from the first provider that has it
func (p *CompositeProvider) Get(key string) string {
	for _, provider := range p.providers {
		if value := provider.Get(key); value != "" {
			return value
		}
	}
	return ""
}

// GetDefault retrieves a string value with a default
func (p *CompositeProvider) GetDefault(key string, defaultValue string) string {
	value := p.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetInt retrieves an integer value
func (p *CompositeProvider) GetInt(key string) int {
	for _, provider := range p.providers {
		if value := provider.Get(key); value != "" {
			i, err := strconv.Atoi(value)
			if err == nil {
				return i
			}
		}
	}
	return 0
}

// GetIntDefault retrieves an integer value with a default
func (p *CompositeProvider) GetIntDefault(key string, defaultValue int) int {
	for _, provider := range p.providers {
		if value := provider.Get(key); value != "" {
			i, err := strconv.Atoi(value)
			if err == nil {
				return i
			}
		}
	}
	return defaultValue
}

// GetBool retrieves a boolean value
func (p *CompositeProvider) GetBool(key string) bool {
	for _, provider := range p.providers {
		if value := provider.Get(key); value != "" {
			value = strings.ToLower(value)
			if value == "true" || value == "1" || value == "yes" {
				return true
			}
			return false
		}
	}
	return false
}

// GetBoolDefault retrieves a boolean value with a default
func (p *CompositeProvider) GetBoolDefault(key string, defaultValue bool) bool {
	for _, provider := range p.providers {
		if value := provider.Get(key); value != "" {
			value = strings.ToLower(value)
			if value == "true" || value == "1" || value == "yes" {
				return true
			}
			return false
		}
	}
	return defaultValue
}

// GetDuration retrieves a duration value
func (p *CompositeProvider) GetDuration(key string) (time.Duration, error) {
	for _, provider := range p.providers {
		if value := provider.Get(key); value != "" {
			return time.ParseDuration(value)
		}
	}
	return 0, fmt.Errorf("configuration property %s not found", key)
}

// GetDurationDefault retrieves a duration value with a default
func (p *CompositeProvider) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	for _, provider := range p.providers {
		if value := provider.Get(key); value != "" {
			duration, err := time.ParseDuration(value)
			if err == nil {
				return duration
			}
		}
	}
	return defaultValue
}

// GetSlice retrieves a slice of strings
func (p *CompositeProvider) GetSlice(key string, separator string) []string {
	for _, provider := range p.providers {
		if value := provider.Get(key); value != "" {
			items := strings.Split(value, separator)
			for i, item := range items {
				items[i] = strings.TrimSpace(item)
			}
			return items
		}
	}
	return []string{}
}