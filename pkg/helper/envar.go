package helper

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func GetEnvAsInt(name string, defaultValue int) int {
	valueStr := GetEnvString(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}

func GetEnvAsBool(name string, defaultValue bool) bool {
	valStr := GetEnvString(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultValue
}

func GetEnvAsSlice(name string, defaultValue []string, sep string) []string {
	valStr := GetEnvString(name, "")

	if valStr == "" {
		return defaultValue

	}

	val := strings.Split(valStr, sep)

	return val
}
