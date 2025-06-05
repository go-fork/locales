package locales

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/locales/* testdata/error_cases/*
var testManagerFS embed.FS

func TestNewManager(t *testing.T) {
	t.Run("should create manager with default config", func(t *testing.T) {
		config := DefaultConfig()
		manager := NewManager(config)

		assert.NotNil(t, manager)
		assert.NotNil(t, manager.Bundle())
	})

	t.Run("should create manager with custom config", func(t *testing.T) {
		config := &Config{
			DefaultLocale:    "vi",
			LocalesDirectory: "translations",
			SupportedLocales: []string{"en", "vi"},
			Bundle:           "messages",
		}
		manager := NewManager(config)

		assert.NotNil(t, manager)
		assert.NotNil(t, manager.Bundle())
	})

	t.Run("should handle nil config gracefully", func(t *testing.T) {
		// This should not panic, manager should handle nil config
		manager := NewManager(nil)
		assert.NotNil(t, manager)
	})
}

func TestManager_SetLocale(t *testing.T) {
	manager := NewManager(DefaultConfig())

	t.Run("should set locale successfully", func(t *testing.T) {
		// This should not panic
		manager.SetLocale("vi")
		manager.SetLocale("en")
	})

	t.Run("should handle invalid locale", func(t *testing.T) {
		// This should not panic even with invalid locale
		manager.SetLocale("invalid-locale")
	})
}

func TestManager_LoadTranslations(t *testing.T) {
	manager := NewManager(DefaultConfig())

	t.Run("should load translations from directory", func(t *testing.T) {
		// Create temporary directory with test files
		tempDir := t.TempDir()

		// Create a test translation file
		enFile := filepath.Join(tempDir, "message.en.yaml")
		enContent := `welcome:
  one: "Welcome"
  other: "Welcome"
hello_name:
  one: "Hello {{.Name}}"
  other: "Hello {{.Name}}"`
		err := os.WriteFile(enFile, []byte(enContent), 0644)
		require.NoError(t, err)

		// Load translations
		err = manager.LoadTranslations(tempDir, nil)
		assert.NoError(t, err)

		// Test translation
		result := manager.Translate("en", "welcome", nil)
		assert.Equal(t, "Welcome", result)
	})

	t.Run("should load translations with custom loader function", func(t *testing.T) {
		loaderFunc := func(path string) ([]byte, error) {
			return testManagerFS.ReadFile(path)
		}

		err := manager.LoadTranslations("testdata/locales", loaderFunc)
		assert.NoError(t, err)

		// Test translation
		result := manager.Translate("en", "welcome", nil)
		assert.Equal(t, "Welcome", result)
	})

	t.Run("should handle non-existing directory", func(t *testing.T) {
		err := manager.LoadTranslations("/non/existing/directory", nil)
		assert.Error(t, err)
	})

	t.Run("should handle loader function error", func(t *testing.T) {
		loaderFunc := func(path string) ([]byte, error) {
			return nil, errors.New("loader error")
		}

		err := manager.LoadTranslations("any-path", loaderFunc)
		assert.Error(t, err)
	})
}

func TestManager_LoadEmbeddedTranslations(t *testing.T) {
	manager := NewManager(DefaultConfig())

	t.Run("should load embedded translations successfully", func(t *testing.T) {
		err := manager.LoadEmbeddedTranslations(testManagerFS, "testdata/locales")
		assert.NoError(t, err)

		// Test English translation
		result := manager.Translate("en", "welcome", nil)
		assert.Equal(t, "Welcome", result)

		// Test Vietnamese translation
		result = manager.Translate("vi", "welcome", nil)
		assert.Equal(t, "Chào mừng", result)
	})

	t.Run("should handle non-existing embedded directory", func(t *testing.T) {
		err := manager.LoadEmbeddedTranslations(testManagerFS, "non/existing/path")
		assert.Error(t, err)
	})

	t.Run("should handle parse errors gracefully", func(t *testing.T) {
		// Load from error_cases directory which contains malformed files
		// The function should continue loading other files even if some fail
		err := manager.LoadEmbeddedTranslations(testManagerFS, "testdata/error_cases")
		assert.NoError(t, err) // Should not fail completely due to error handling
	})

	t.Run("should skip non-translation files", func(t *testing.T) {
		// The error_cases directory contains a readme.txt file which should be skipped
		err := manager.LoadEmbeddedTranslations(testManagerFS, "testdata/error_cases")
		assert.NoError(t, err)

		// Should not have loaded any valid translations from error_cases
		// since they're all malformed or non-translation files
	})

	t.Run("should handle empty files gracefully", func(t *testing.T) {
		// Create a new manager for this test
		testManager := NewManager(DefaultConfig())

		// The error_cases directory contains empty.json which should be handled gracefully
		err := testManager.LoadEmbeddedTranslations(testManagerFS, "testdata/error_cases")
		assert.NoError(t, err) // Should continue with other files
	})

	t.Run("should handle multiple file extensions", func(t *testing.T) {
		// Load normal translations which include both .yaml files
		testManager := NewManager(DefaultConfig())
		err := testManager.LoadEmbeddedTranslations(testManagerFS, "testdata/locales")
		assert.NoError(t, err)

		// Verify both English and Vietnamese translations loaded
		assert.Equal(t, "Welcome", testManager.Translate("en", "welcome", nil))
		assert.Equal(t, "Chào mừng", testManager.Translate("vi", "welcome", nil))
	})

	t.Run("should handle directory traversal errors", func(t *testing.T) {
		// Test with a path that exists but might have permission issues
		// or other traversal problems - the function should handle fs.WalkDir errors
		manager := NewManager(DefaultConfig())

		// Try to load from a path that doesn't exist in the embedded FS
		err := manager.LoadEmbeddedTranslations(testManagerFS, "testdata/nonexistent")
		assert.Error(t, err) // fs.WalkDir should return an error for non-existent path
	})
}

func TestManager_Translate(t *testing.T) {
	manager := NewManager(DefaultConfig())
	err := manager.LoadEmbeddedTranslations(testManagerFS, "testdata/locales")
	require.NoError(t, err)

	t.Run("should translate simple message", func(t *testing.T) {
		result := manager.Translate("en", "welcome", nil)
		assert.Equal(t, "Welcome", result)

		result = manager.Translate("vi", "welcome", nil)
		assert.Equal(t, "Chào mừng", result)
	})

	t.Run("should translate message with parameters", func(t *testing.T) {
		args := map[string]interface{}{
			"Name": "John",
		}

		result := manager.Translate("en", "hello_name", args)
		assert.Equal(t, "Hello John", result)

		result = manager.Translate("vi", "hello_name", args)
		assert.Equal(t, "Xin chào John", result)
	})

	t.Run("should handle missing translation key", func(t *testing.T) {
		result := manager.Translate("en", "non_existing_key", nil)
		// Should return the key itself when translation not found
		assert.Equal(t, "non_existing_key", result)
	})

	t.Run("should handle unsupported language", func(t *testing.T) {
		result := manager.Translate("unsupported", "welcome", nil)
		// Should fallback or return the key
		assert.NotEmpty(t, result)
	})

	t.Run("should handle nil arguments", func(t *testing.T) {
		result := manager.Translate("en", "welcome", nil)
		assert.Equal(t, "Welcome", result)
	})
}

func TestManager_TranslateFunc(t *testing.T) {
	manager := NewManager(DefaultConfig())
	err := manager.LoadEmbeddedTranslations(testManagerFS, "testdata/locales")
	require.NoError(t, err)

	t.Run("should return translation function for English", func(t *testing.T) {
		translateFunc := manager.TranslateFunc("en")
		assert.NotNil(t, translateFunc)

		result := translateFunc("welcome", nil)
		assert.Equal(t, "Welcome", result)

		args := map[string]interface{}{"Name": "Alice"}
		result = translateFunc("hello_name", args)
		assert.Equal(t, "Hello Alice", result)
	})

	t.Run("should return translation function for Vietnamese", func(t *testing.T) {
		translateFunc := manager.TranslateFunc("vi")
		assert.NotNil(t, translateFunc)

		result := translateFunc("welcome", nil)
		assert.Equal(t, "Chào mừng", result)

		args := map[string]interface{}{"Name": "Alice"}
		result = translateFunc("hello_name", args)
		assert.Equal(t, "Xin chào Alice", result)
	})

	t.Run("should handle unsupported language", func(t *testing.T) {
		translateFunc := manager.TranslateFunc("unsupported")
		assert.NotNil(t, translateFunc)

		result := translateFunc("welcome", nil)
		assert.NotEmpty(t, result)
	})
}

func TestManager_RegisterTranslations(t *testing.T) {
	manager := NewManager(DefaultConfig())

	t.Run("should register translations for module", func(t *testing.T) {
		translations := map[string]map[string]string{
			"en": {
				"module.welcome": "Module Welcome",
				"module.goodbye": "Module Goodbye",
			},
			"vi": {
				"module.welcome": "Chào mừng Module",
				"module.goodbye": "Tạm biệt Module",
			},
		}

		// This should not panic
		manager.RegisterTranslations("test-module", translations)
	})

	t.Run("should handle empty translations", func(t *testing.T) {
		translations := map[string]map[string]string{}

		// This should not panic
		manager.RegisterTranslations("empty-module", translations)
	})

	t.Run("should handle nil translations", func(t *testing.T) {
		// This should not panic
		manager.RegisterTranslations("nil-module", nil)
	})
}

func TestManager_Bundle(t *testing.T) {
	manager := NewManager(DefaultConfig())

	t.Run("should return non-nil bundle", func(t *testing.T) {
		bundle := manager.Bundle()
		assert.NotNil(t, bundle)
	})

	t.Run("should return same bundle instance", func(t *testing.T) {
		bundle1 := manager.Bundle()
		bundle2 := manager.Bundle()
		assert.Same(t, bundle1, bundle2)
	})
}

func TestManager_ConcurrentAccess(t *testing.T) {
	manager := NewManager(DefaultConfig())
	err := manager.LoadEmbeddedTranslations(testManagerFS, "testdata/locales")
	require.NoError(t, err)

	t.Run("should handle concurrent translation requests", func(t *testing.T) {
		done := make(chan bool)

		// Start multiple goroutines making translation requests
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()

				for j := 0; j < 100; j++ {
					result := manager.Translate("en", "welcome", nil)
					assert.Equal(t, "Welcome", result)

					result = manager.Translate("vi", "welcome", nil)
					assert.Equal(t, "Chào mừng", result)
				}
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("should handle concurrent locale setting", func(t *testing.T) {
		done := make(chan bool)

		// Start multiple goroutines setting locale
		for i := 0; i < 5; i++ {
			go func(index int) {
				defer func() { done <- true }()

				for j := 0; j < 10; j++ {
					if index%2 == 0 {
						manager.SetLocale("en")
					} else {
						manager.SetLocale("vi")
					}
				}
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 5; i++ {
			<-done
		}
	})
}

// Additional rigorous test cases for edge scenarios and error handling

func TestManager_ErrorHandling(t *testing.T) {
	t.Run("should handle invalid YAML files gracefully", func(t *testing.T) {
		manager := NewManager(DefaultConfig())
		tempDir := t.TempDir()

		// Create invalid YAML file
		invalidFile := filepath.Join(tempDir, "invalid.yaml")
		invalidContent := `welcome:
  one: "Welcome"
  invalid yaml content without proper structure
    - broken
  other: "Welcome"`
		err := os.WriteFile(invalidFile, []byte(invalidContent), 0644)
		require.NoError(t, err)

		// Should handle error gracefully and continue
		err = manager.LoadTranslations(tempDir, nil)
		assert.NoError(t, err) // Should not fail, just skip invalid files
	})

	t.Run("should handle empty translation files", func(t *testing.T) {
		manager := NewManager(DefaultConfig())
		tempDir := t.TempDir()

		// Create empty YAML file
		emptyFile := filepath.Join(tempDir, "empty.yaml")
		err := os.WriteFile(emptyFile, []byte(""), 0644)
		require.NoError(t, err)

		err = manager.LoadTranslations(tempDir, nil)
		assert.NoError(t, err) // Should handle empty files gracefully
	})

	t.Run("should handle permission denied errors", func(t *testing.T) {
		manager := NewManager(DefaultConfig())

		// Try to access a directory that doesn't exist or has no permission
		err := manager.LoadTranslations("/root/secret", nil)
		assert.Error(t, err) // Should return error for inaccessible directories
	})
}

func TestManager_ComplexTranslations(t *testing.T) {
	manager := NewManager(DefaultConfig())
	err := manager.LoadEmbeddedTranslations(testManagerFS, "testdata/locales")
	require.NoError(t, err)

	t.Run("should handle multiple parameter substitutions", func(t *testing.T) {
		// Register a complex translation for testing
		translations := map[string]map[string]string{
			"en": {
				"complex": "Hello {{.Name}}, you have {{.Count}} messages from {{.Sender}}",
			},
			"vi": {
				"complex": "Xin chào {{.Name}}, bạn có {{.Count}} tin nhắn từ {{.Sender}}",
			},
		}
		manager.RegisterTranslations("test", translations)

		args := map[string]interface{}{
			"Name":   "John",
			"Count":  5,
			"Sender": "Alice",
		}

		result := manager.Translate("en", "test.complex", args)
		assert.Equal(t, "Hello John, you have 5 messages from Alice", result)

		result = manager.Translate("vi", "test.complex", args)
		assert.Equal(t, "Xin chào John, bạn có 5 tin nhắn từ Alice", result)
	})

	t.Run("should handle missing parameters gracefully", func(t *testing.T) {
		translations := map[string]map[string]string{
			"en": {
				"missing_params": "Hello {{.Name}}, you have {{.Count}} items",
			},
		}
		manager.RegisterTranslations("test", translations)

		// Only provide Name, missing Count
		args := map[string]interface{}{
			"Name": "John",
		}

		result := manager.Translate("en", "test.missing_params", args)
		// Should still work, leaving unreplaced placeholders
		assert.Contains(t, result, "Hello John")
		assert.Contains(t, result, "{{.Count}}")
	})

	t.Run("should handle nested module translations", func(t *testing.T) {
		translations := map[string]map[string]string{
			"en": {
				"auth.login":    "Login",
				"auth.logout":   "Logout",
				"user.profile":  "User Profile",
				"user.settings": "User Settings",
			},
		}
		manager.RegisterTranslations("admin", translations)

		result := manager.Translate("en", "admin.auth.login", nil)
		assert.Equal(t, "Login", result)

		result = manager.Translate("en", "admin.user.profile", nil)
		assert.Equal(t, "User Profile", result)
	})
}

func TestManager_LanguageFallback(t *testing.T) {
	config := &Config{
		DefaultLocale:    "en",
		LocalesDirectory: "locales",
		SupportedLocales: []string{"en", "vi", "fr"},
		Bundle:           "message",
	}
	manager := NewManager(config)
	err := manager.LoadEmbeddedTranslations(testManagerFS, "testdata/locales")
	require.NoError(t, err)

	t.Run("should fallback to default locale for unsupported language", func(t *testing.T) {
		// Request translation in unsupported language
		result := manager.Translate("fr", "welcome", nil)
		// Should fallback to English translation since we don't have fr translations
		assert.Equal(t, "Welcome", result)
	})

	t.Run("should handle partial language codes", func(t *testing.T) {
		// Test with language-region codes
		result := manager.Translate("en-US", "welcome", nil)
		assert.NotEmpty(t, result)

		result = manager.Translate("vi-VN", "welcome", nil)
		assert.NotEmpty(t, result)
	})
}

func TestManager_StateManagement(t *testing.T) {
	t.Run("should maintain state after locale changes", func(t *testing.T) {
		manager := NewManager(DefaultConfig())
		err := manager.LoadEmbeddedTranslations(testManagerFS, "testdata/locales")
		require.NoError(t, err)

		// Change locale multiple times
		manager.SetLocale("vi")
		manager.SetLocale("en")
		manager.SetLocale("fr")
		manager.SetLocale("en")

		// Should still work correctly
		result := manager.Translate("en", "welcome", nil)
		assert.Equal(t, "Welcome", result)
	})

	t.Run("should handle registration after translations loaded", func(t *testing.T) {
		manager := NewManager(DefaultConfig())
		err := manager.LoadEmbeddedTranslations(testManagerFS, "testdata/locales")
		require.NoError(t, err)

		// Register additional translations after loading
		newTranslations := map[string]map[string]string{
			"en": {"new_key": "New Message"},
			"vi": {"new_key": "Tin nhắn mới"},
		}
		manager.RegisterTranslations("runtime", newTranslations)

		result := manager.Translate("en", "runtime.new_key", nil)
		assert.Equal(t, "New Message", result)

		result = manager.Translate("vi", "runtime.new_key", nil)
		assert.Equal(t, "Tin nhắn mới", result)
	})
}

func TestManager_Performance(t *testing.T) {
	t.Run("should handle large number of translations efficiently", func(t *testing.T) {
		manager := NewManager(DefaultConfig())

		// Register a large number of translations
		translations := make(map[string]map[string]string)
		translations["en"] = make(map[string]string)
		translations["vi"] = make(map[string]string)

		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("key_%d", i)
			translations["en"][key] = fmt.Sprintf("English message %d", i)
			translations["vi"][key] = fmt.Sprintf("Tin nhắn tiếng Việt %d", i)
		}

		manager.RegisterTranslations("perf", translations)

		// Test that lookups are still fast
		start := time.Now()
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("perf.key_%d", i%1000)
			result := manager.Translate("en", key, nil)
			assert.NotEmpty(t, result)
		}
		duration := time.Since(start)

		// Should complete within reasonable time (adjust as needed)
		assert.Less(t, duration, time.Millisecond*100)
	})
}
