package locales

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/locales/*
var testFS embed.FS

func TestDefaultConfig(t *testing.T) {
	t.Run("should return default config with expected values", func(t *testing.T) {
		config := DefaultConfig()

		assert.NotNil(t, config)
		assert.Equal(t, "en", config.DefaultLocale)
		assert.Equal(t, "locales", config.LocalesDirectory)
		assert.Equal(t, "message", config.Bundle)
		assert.Equal(t, []string{"en"}, config.SupportedLocales)
		assert.Nil(t, config.EmbeddedFS)
		assert.Empty(t, config.EmbeddedPath)
	})
}

func TestConfig_Struct(t *testing.T) {
	t.Run("should create config with custom values", func(t *testing.T) {
		config := &Config{
			DefaultLocale:    "vi",
			LocalesDirectory: "translations",
			SupportedLocales: []string{"en", "vi", "fr"},
			Bundle:           "messages",
			EmbeddedFS:       &testFS,
			EmbeddedPath:     "testdata/locales",
		}

		assert.Equal(t, "vi", config.DefaultLocale)
		assert.Equal(t, "translations", config.LocalesDirectory)
		assert.Equal(t, []string{"en", "vi", "fr"}, config.SupportedLocales)
		assert.Equal(t, "messages", config.Bundle)
		assert.NotNil(t, config.EmbeddedFS)
		assert.Equal(t, "testdata/locales", config.EmbeddedPath)
	})

	t.Run("should support empty supported locales", func(t *testing.T) {
		config := &Config{
			DefaultLocale:    "en",
			LocalesDirectory: "locales",
			SupportedLocales: []string{},
			Bundle:           "message",
		}

		assert.Equal(t, "en", config.DefaultLocale)
		assert.Empty(t, config.SupportedLocales)
	})

	t.Run("should support multiple supported locales", func(t *testing.T) {
		supportedLangs := []string{"en", "vi", "fr", "de", "ja", "ko"}
		config := &Config{
			DefaultLocale:    "en",
			LocalesDirectory: "locales",
			SupportedLocales: supportedLangs,
			Bundle:           "message",
		}

		assert.Equal(t, supportedLangs, config.SupportedLocales)
		assert.Len(t, config.SupportedLocales, 6)
	})
}

func TestConfig_EmbeddedFS(t *testing.T) {
	t.Run("should handle embedded filesystem", func(t *testing.T) {
		config := &Config{
			DefaultLocale:    "en",
			LocalesDirectory: "locales",
			SupportedLocales: []string{"en"},
			Bundle:           "message",
			EmbeddedFS:       &testFS,
			EmbeddedPath:     "testdata/locales",
		}

		assert.NotNil(t, config.EmbeddedFS)
		assert.Equal(t, "testdata/locales", config.EmbeddedPath)
	})

	t.Run("should handle nil embedded filesystem", func(t *testing.T) {
		config := &Config{
			DefaultLocale:    "en",
			LocalesDirectory: "locales",
			SupportedLocales: []string{"en"},
			Bundle:           "message",
			EmbeddedFS:       nil,
			EmbeddedPath:     "",
		}

		assert.Nil(t, config.EmbeddedFS)
		assert.Empty(t, config.EmbeddedPath)
	})
}
