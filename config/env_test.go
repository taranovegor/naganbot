package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	key := "NAGANBOT_TEST_ENV"
	expected := "test_value"
	os.Setenv(key, expected)
	defer os.Unsetenv(key)

	actual := GetEnv(key)
	if actual != expected {
		t.Errorf("value %s of environment variable %s does not match the expected value %s", actual, key, expected)
	}
}
