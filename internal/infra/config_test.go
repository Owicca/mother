package infra

import (
	"testing"
	"reflect"
)

func TestNewConfig(t *testing.T) {
	t.Run("return a Config", func(t *testing.T) {
		want := "test"
		cfg := NewConfig(want)
		if reflect.TypeOf(Config{}) != reflect.TypeOf(cfg) {
			t.Errorf("NewConfig(%q) is not of type Config", want)
		}
	})
	t.Run("path is correct", func(t *testing.T) {
		want := "test"
		cfg := NewConfig(want)
		if cfg.CfgPath != want {
			t.Errorf("NewConfig(%q) = %q", want, cfg.CfgPath)
		}
	})
}

func TestString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := Config{
		}

		want := `""`
		got := cfg.String()
		if got != want {
			t.Errorf("Config.String() = %q, but wanted an empty string", got)
		}
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("file does not exist", func(t *testing.T) {
		path := "null"

		_, err := LoadConfig(path)
		if err == nil {
			t.Errorf("LoadConfig(%q) did not err on invalid path", path)
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		path := "./testdata/malformed_config.json"

		_, err := LoadConfig(path)
		if err == nil {
			t.Errorf("LoadConfig(%q) did not err on malformed json", path)
		}
	})
}