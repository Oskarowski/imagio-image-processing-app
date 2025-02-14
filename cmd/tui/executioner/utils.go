package executioner

import (
	"fmt"
	"strconv"
)

func parseIntArg(args map[string]string, key string) (int, error) {
	value, exists := args[key]
	if !exists {
		return 0, fmt.Errorf("missing required argument: %s", key)
	}
	return strconv.Atoi(value)
}

func parseBoolArg(args map[string]string, key string) (bool, error) {
	value, exists := args[key]
	if !exists {
		return false, fmt.Errorf("missing required argument: %s", key)
	}
	return strconv.ParseBool(value)
}

func parseFloatArg(args map[string]string, key string) (float64, error) {
	value, exists := args[key]
	if !exists {
		return 0, fmt.Errorf("missing required argument: %s", key)
	}
	return strconv.ParseFloat(value, 64)
}
