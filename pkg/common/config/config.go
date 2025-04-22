package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Provider interface {
	Get(key string) string 
	GetInt(key string) int 
	GetBool(key string) bool
}

type EnvProvider struct {
	prefix string
}

func NewEnvProvider(prefix string) *EnvProvider {
	return &EnvProvider{
		prefix: prefix,
	}
}

func (p *EnvProvider) Get(key string) string {
	return os.Getenv(p.formatKey(key))
}

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


func (p *EnvProvider) GetBool(key string) bool {
	str := strings.ToLower(p.Get(key))
	return str == "true" || str == "1" || str == "yes"
}


func (p *EnvProvider) formatKey(key string) string {
	if p.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s_%s", p.prefix, key)
}

func LoadJsonFile(filepath string, cfg interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(cfg)
}