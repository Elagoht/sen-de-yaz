package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Development env var is preferred over .env
var envFiles = []string{".env.development", ".env"}

// Loads environment variables from .env files
func LoadEnvFile(dir string) (string, error) {
	for _, name := range envFiles {
		content, err := os.ReadFile(filepath.Join(dir, name))
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", err
		}
		if err := applyEnv(string(content)); err != nil {
			return "", fmt.Errorf("%s: %w", name, err)
		}
		return name, nil
	}
	return "", nil
}

// Returns env as string and panics if not found
func EnvString(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	panic(fmt.Sprintf("Missing environment variable: %s", key))
}

// Returns env as integer and panics if not found
func EnvInt(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		panic(fmt.Sprintf("Missing environment variable: %s", key))
	}
	return value
}

// Parses environment variables
func applyEnv(content string) error {
	for number, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		key, value, found := strings.Cut(line, "=")
		if !found {
			return fmt.Errorf("line %d: want KEY=value", number+1)
		}
		key = strings.TrimSpace(key)
		if key == "" || strings.ContainsAny(key, " \t") {
			return fmt.Errorf("line %d: %q is not a variable name", number+1, key)
		}
		value = strings.TrimSpace(value)
		if value != "" && (value[0] == '"' || value[0] == '\'') {
			quote := value[0]
			end := strings.IndexByte(value[1:], quote)
			if end < 0 {
				return fmt.Errorf("line %d: unterminated quote", number+1)
			}
			rest := strings.TrimSpace(value[end+2:])
			if rest != "" && !strings.HasPrefix(rest, "#") {
				return fmt.Errorf("line %d: text after the closing quote", number+1)
			}
			value = value[1 : end+1]
		}
		if _, ok := os.LookupEnv(key); !ok {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("line %d: set %s: %w", number+1, key, err)
			}
		}
	}
	return nil
}
