package config

import "testing"

func TestGetEnv(t *testing.T) {
	variable := "GO111MODULE"
	expected := "on"
	actual := GetEnv(variable)
	if actual != expected {
		t.Errorf("value %s of environment variable %s does not match the expected value %s", actual, variable, expected)
	}
}
