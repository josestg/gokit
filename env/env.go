package env

import (
	"os"
	"strconv"
	"time"
)

// String returns the env string value if the key exists.
// Otherwise, returns fallback value.
func String(key, fallback string) string {
	v, exists := getEnv(key)
	if !exists {
		return fallback
	}

	return v
}

// Int returns the env integer value if the key exists.
// Otherwise, returns fallback value.
func Int(key string, fallback int) int {
	return Parse(key, fallback, strconv.Atoi)
}

// Int64 returns the env int64 value if the key exists.
// Otherwise, returns fallback value.
func Int64(key string, fallback int64) int64 {
	return Parse(key, fallback, func(v string) (int64, error) {
		return strconv.ParseInt(v, 10, 64)
	})
}

// Float64 returns the env float64 value if the key exists.
// Otherwise, returns fallback value.
func Float64(key string, fallback float64) float64 {
	return Parse(key, fallback, func(v string) (float64, error) {
		return strconv.ParseFloat(v, 64)
	})
}

// Duration returns the env duration value if the key exists.
// Otherwise, returns fallback value.
func Duration(key string, fallback time.Duration) time.Duration {
	return Parse(key, fallback, time.ParseDuration)
}

// Bool returns the env boolean value if the key exists.
// Otherwise, returns fallback value.
func Bool(key string, fallback bool) bool {
	return Parse(key, fallback, strconv.ParseBool)
}

// Parse the env value to the given type.
func Parse[T any](key string, fallback T, parser func(v string) (T, error)) T {
	v, exists := getEnv(key)
	if !exists {
		return fallback
	}

	return must(parser(v))
}

// getEnv returns the env value and exists flag.
func getEnv(key string) (string, bool) {
	v, exists := os.LookupEnv(key)
	return v, exists
}

// must panics if the error is not nil.
func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}
