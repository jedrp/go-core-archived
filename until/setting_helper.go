package until

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Setting struct {
	EnvVar string
	items  map[string]string
}

func (s *Setting) IsConfigured() bool {
	if s == nil || len(s.items) > 0 {
		return true
	}
	return false
}

func (s *Setting) GetBoolValue(key string, defaultValue bool) bool {
	v := s.getValue(key)
	if v == "" {
		return defaultValue
	}
	defaultValue, err := strconv.ParseBool(v)
	if err != nil {
		panic(err)
	}
	return defaultValue
}

func (s *Setting) GetIntValue(key string, defaultValue int) int {
	v := s.getValue(key)
	if v == "" {
		return defaultValue
	}
	defaultValue, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return defaultValue
}

func (s *Setting) GetStringValue(key string, defaultValue string) string {
	v := s.getValue(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func (s *Setting) getValue(key string) string {
	if !s.IsConfigured() {
		panic(s.EnvVar + " Setting was not configured")
	}
	if v, ok := s.items[key]; ok {
		return v
	}
	return ""
}

func GetSettings(envName string) (*Setting, error) {
	str := os.Getenv(envName)
	if str == "" {
		return nil, nil
	}
	setting := &Setting{
		EnvVar: envName,
	}
	values := strings.Split(str, ";")
	for _, v := range values {
		v = strings.Trim(v, " ")
		s := strings.Split(v, "=")
		if len(s) != 2 {
			return nil, fmt.Errorf("Setting is not correct - %s", v)
		}
		setting.items[s[0]] = s[1]
	}
	return setting, nil
}
