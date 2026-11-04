package main

import (
	"os"
	"testing"
)

// unsetSettingsEnv unsets the variable during the test and restores it afterwards.
func unsetSettingsEnv(t *testing.T) {
	t.Setenv(settingsEnvVar, "")
	if err := os.Unsetenv(settingsEnvVar); err != nil {
		t.Fatal(err)
	}
}

func TestReadServerSettings(t *testing.T) {
	t.Run("Reads settings from the environment and unsets the variable", func(t *testing.T) {
		t.Setenv(settingsEnvVar, `{"port":1}`)
		settings, err := readServerSettings([]string{"server"})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if settings != `{"port":1}` {
			t.Errorf("Expected settings from the environment, got %q", settings)
		}
		if _, ok := os.LookupEnv(settingsEnvVar); ok {
			t.Errorf("Expected %s to be unset", settingsEnvVar)
		}
	})

	t.Run("Prefers the environment over the argument", func(t *testing.T) {
		t.Setenv(settingsEnvVar, `{"port":1}`)
		settings, err := readServerSettings([]string{"server", `{"port":2}`})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if settings != `{"port":1}` {
			t.Errorf("Expected settings from the environment, got %q", settings)
		}
	})

	t.Run("Falls back to the first argument", func(t *testing.T) {
		unsetSettingsEnv(t)
		settings, err := readServerSettings([]string{"server", `{"port":2}`})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if settings != `{"port":2}` {
			t.Errorf("Expected settings from the argument, got %q", settings)
		}
	})

	t.Run("Returns an error without settings", func(t *testing.T) {
		unsetSettingsEnv(t)
		if _, err := readServerSettings([]string{"server"}); err == nil {
			t.Error("Expected an error")
		}
	})
}
