package locales

import (
	"embed"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/locales/*
var testEmbedFS embed.FS

func TestLoaderFunc(t *testing.T) {
	t.Run("should implement Loader interface", func(t *testing.T) {
		loaderFunc := LoaderFunc(func(path string) ([]byte, error) {
			return []byte("test content"), nil
		})

		// Test that LoaderFunc implements Loader interface
		var _ Loader = loaderFunc

		content, err := loaderFunc.LoadMessage("test.txt")
		assert.NoError(t, err)
		assert.Equal(t, []byte("test content"), content)
	})

	t.Run("should return error when function returns error", func(t *testing.T) {
		expectedErr := errors.New("load error")
		loaderFunc := LoaderFunc(func(path string) ([]byte, error) {
			return nil, expectedErr
		})

		content, err := loaderFunc.LoadMessage("test.txt")
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, content)
	})

	t.Run("should pass filepath to underlying function", func(t *testing.T) {
		var receivedPath string
		loaderFunc := LoaderFunc(func(path string) ([]byte, error) {
			receivedPath = path
			return []byte("content"), nil
		})

		_, err := loaderFunc.LoadMessage("/path/to/file.txt")
		assert.NoError(t, err)
		assert.Equal(t, "/path/to/file.txt", receivedPath)
	})
}

func TestEmbedLoader(t *testing.T) {
	t.Run("should implement Loader interface", func(t *testing.T) {
		loader := &EmbedLoader{FS: testEmbedFS}

		// Test that EmbedLoader implements Loader interface
		var _ Loader = loader
	})

	t.Run("should load existing file from embedded FS", func(t *testing.T) {
		loader := &EmbedLoader{FS: testEmbedFS}

		content, err := loader.LoadMessage("testdata/locales/message.en.yaml")
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
		assert.Contains(t, string(content), "welcome:")
		assert.Contains(t, string(content), "Hello {{.Name}}")
	})

	t.Run("should return error for non-existing file", func(t *testing.T) {
		loader := &EmbedLoader{FS: testEmbedFS}

		content, err := loader.LoadMessage("non-existing-file.yaml")
		assert.Error(t, err)
		assert.Nil(t, content)
	})

	t.Run("should load Vietnamese translation file", func(t *testing.T) {
		loader := &EmbedLoader{FS: testEmbedFS}

		content, err := loader.LoadMessage("testdata/locales/message.vi.yaml")
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
		assert.Contains(t, string(content), "Chào mừng")
		assert.Contains(t, string(content), "Xin chào {{.Name}}")
	})
}

func TestFileSystemLoader(t *testing.T) {
	// Create temporary directory and files for testing
	tempDir := t.TempDir()

	testFile := filepath.Join(tempDir, "test.yaml")
	testContent := "test: content"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	t.Run("should load file from filesystem using os.ReadFile", func(t *testing.T) {
		content, err := os.ReadFile(testFile)
		assert.NoError(t, err)
		assert.Equal(t, testContent, string(content))
	})

	t.Run("should return error for non-existing file", func(t *testing.T) {
		nonExistingFile := filepath.Join(tempDir, "non-existing.yaml")
		content, err := os.ReadFile(nonExistingFile)
		assert.Error(t, err)
		assert.Empty(t, content)
	})
}

func TestLoader_Integration(t *testing.T) {
	t.Run("should work with different loader implementations", func(t *testing.T) {
		// Test with EmbedLoader
		embedLoader := &EmbedLoader{FS: testEmbedFS}
		content1, err1 := embedLoader.LoadMessage("testdata/locales/message.en.yaml")
		assert.NoError(t, err1)
		assert.NotEmpty(t, content1)

		// Test with LoaderFunc
		funcLoader := LoaderFunc(func(path string) ([]byte, error) {
			return testEmbedFS.ReadFile(path)
		})
		content2, err2 := funcLoader.LoadMessage("testdata/locales/message.en.yaml")
		assert.NoError(t, err2)
		assert.Equal(t, content1, content2)
	})

	t.Run("should handle different file formats", func(t *testing.T) {
		// Create a temporary JSON file
		tempDir := t.TempDir()
		jsonFile := filepath.Join(tempDir, "test.json")
		jsonContent := `{"test": "content", "number": 123}`
		err := os.WriteFile(jsonFile, []byte(jsonContent), 0644)
		require.NoError(t, err)

		// Use LoaderFunc with os.ReadFile
		loader := LoaderFunc(os.ReadFile)
		content, err := loader.LoadMessage(jsonFile)
		assert.NoError(t, err)
		assert.Equal(t, jsonContent, string(content))
	})
}

func TestFSLoader(t *testing.T) {
	t.Run("should implement Loader interface", func(t *testing.T) {
		fsLoader := FSLoader{FS: testEmbedFS}

		// Test that FSLoader implements Loader interface
		var _ Loader = fsLoader

		content, err := fsLoader.LoadMessage("testdata/locales/message.en.yaml")
		assert.NoError(t, err)
		assert.Contains(t, string(content), "welcome")
	})

	t.Run("should return error for non-existing file", func(t *testing.T) {
		fsLoader := FSLoader{FS: testEmbedFS}

		content, err := fsLoader.LoadMessage("testdata/locales/nonexistent.yaml")
		assert.Error(t, err)
		assert.Nil(t, content)
	})

	t.Run("should load Vietnamese translation file", func(t *testing.T) {
		fsLoader := FSLoader{FS: testEmbedFS}

		content, err := fsLoader.LoadMessage("testdata/locales/message.vi.yaml")
		assert.NoError(t, err)
		assert.Contains(t, string(content), "Chào mừng")
	})
}

func TestNewModuleLoader(t *testing.T) {
	t.Run("should create new module loader", func(t *testing.T) {
		manager := NewManager(DefaultConfig())

		loader := NewModuleLoader("test-module", manager)

		assert.NotNil(t, loader)
		assert.Equal(t, "test-module", loader.moduleName)
		assert.Equal(t, manager, loader.manager)
	})

	t.Run("should create loader with different module names", func(t *testing.T) {
		manager := NewManager(DefaultConfig())

		loader1 := NewModuleLoader("auth", manager)
		loader2 := NewModuleLoader("user", manager)

		assert.Equal(t, "auth", loader1.moduleName)
		assert.Equal(t, "user", loader2.moduleName)
		assert.Equal(t, manager, loader1.manager)
		assert.Equal(t, manager, loader2.manager)
	})
}

func TestModuleLoader_RegisterMessage(t *testing.T) {
	t.Run("should register single message", func(t *testing.T) {
		manager := NewManager(DefaultConfig())
		loader := NewModuleLoader("test", manager)

		// This should not panic
		loader.RegisterMessage("en", "hello", "Hello World")

		// Verify the message was registered by trying to translate it
		result := manager.Translate("en", "test.hello", nil)
		assert.Equal(t, "Hello World", result)
	})

	t.Run("should register messages for different languages", func(t *testing.T) {
		manager := NewManager(DefaultConfig())
		loader := NewModuleLoader("test", manager)

		loader.RegisterMessage("en", "greeting", "Hello")
		loader.RegisterMessage("vi", "greeting", "Xin chào")

		// Verify both languages
		resultEn := manager.Translate("en", "test.greeting", nil)
		assert.Equal(t, "Hello", resultEn)

		resultVi := manager.Translate("vi", "test.greeting", nil)
		assert.Equal(t, "Xin chào", resultVi)
	})

	t.Run("should handle empty message", func(t *testing.T) {
		manager := NewManager(DefaultConfig())
		loader := NewModuleLoader("test", manager)

		// This should not panic
		loader.RegisterMessage("en", "empty", "")

		result := manager.Translate("en", "test.empty", nil)
		assert.Equal(t, "", result)
	})
}

func TestModuleLoader_RegisterMessages(t *testing.T) {
	t.Run("should register multiple messages", func(t *testing.T) {
		manager := NewManager(DefaultConfig())
		loader := NewModuleLoader("test", manager)

		messages := map[string]string{
			"hello":   "Hello World",
			"goodbye": "Goodbye",
			"thanks":  "Thank you",
		}

		loader.RegisterMessages("en", messages)

		// Verify all messages were registered
		assert.Equal(t, "Hello World", manager.Translate("en", "test.hello", nil))
		assert.Equal(t, "Goodbye", manager.Translate("en", "test.goodbye", nil))
		assert.Equal(t, "Thank you", manager.Translate("en", "test.thanks", nil))
	})

	t.Run("should handle empty messages map", func(t *testing.T) {
		manager := NewManager(DefaultConfig())
		loader := NewModuleLoader("test", manager)

		messages := map[string]string{}

		// This should not panic
		loader.RegisterMessages("en", messages)
	})

	t.Run("should handle nil messages map", func(t *testing.T) {
		manager := NewManager(DefaultConfig())
		loader := NewModuleLoader("test", manager)

		// This should not panic
		loader.RegisterMessages("en", nil)
	})

	t.Run("should register for multiple languages", func(t *testing.T) {
		manager := NewManager(DefaultConfig())
		loader := NewModuleLoader("test", manager)

		enMessages := map[string]string{
			"welcome": "Welcome",
			"error":   "Error occurred",
		}

		viMessages := map[string]string{
			"welcome": "Chào mừng",
			"error":   "Đã xảy ra lỗi",
		}

		loader.RegisterMessages("en", enMessages)
		loader.RegisterMessages("vi", viMessages)

		// Verify both languages
		assert.Equal(t, "Welcome", manager.Translate("en", "test.welcome", nil))
		assert.Equal(t, "Chào mừng", manager.Translate("vi", "test.welcome", nil))
		assert.Equal(t, "Error occurred", manager.Translate("en", "test.error", nil))
		assert.Equal(t, "Đã xảy ra lỗi", manager.Translate("vi", "test.error", nil))
	})
}

// Test for FSLoader LoadMessage method (line 81)
func TestFSLoader_LoadMessage(t *testing.T) {
	// Create test filesystem
	testFS := fstest.MapFS{
		"translations/en.json": &fstest.MapFile{
			Data: []byte(`{"hello": "Hello", "world": "World"}`),
		},
		"translations/vi.json": &fstest.MapFile{
			Data: []byte(`{"hello": "Xin chào", "world": "Thế giới"}`),
		},
	}

	loader := FSLoader{FS: testFS}

	t.Run("should load existing file", func(t *testing.T) {
		data, err := loader.LoadMessage("translations/en.json")
		require.NoError(t, err)
		assert.JSONEq(t, `{"hello": "Hello", "world": "World"}`, string(data))
	})

	t.Run("should return error for non-existent file", func(t *testing.T) {
		_, err := loader.LoadMessage("non-existent.json")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file does not exist")
	})
}

// Test for ModuleLoader LoadJSON method
func TestModuleLoader_LoadJSON(t *testing.T) {
	// Create temporary files for testing
	tempDir := t.TempDir()

	// Create test JSON files
	enFile := filepath.Join(tempDir, "en.json")
	viFile := filepath.Join(tempDir, "vi.json")
	invalidFile := filepath.Join(tempDir, "invalid.json")

	err := os.WriteFile(enFile, []byte(`{"hello": "Hello", "goodbye": "Goodbye"}`), 0644)
	require.NoError(t, err)

	err = os.WriteFile(viFile, []byte(`{"hello": "Xin chào", "goodbye": "Tạm biệt"}`), 0644)
	require.NoError(t, err)

	err = os.WriteFile(invalidFile, []byte(`{"invalid": json}`), 0644)
	require.NoError(t, err)

	manager := NewManager(nil)
	loader := NewModuleLoader("testmodule", manager)

	t.Run("should load valid JSON file", func(t *testing.T) {
		err := loader.LoadJSON(enFile)
		require.NoError(t, err)

		// Verify translation was registered
		translation := manager.Translate("en", "testmodule.hello", nil)
		assert.Equal(t, "Hello", translation)
	})

	t.Run("should load Vietnamese JSON file", func(t *testing.T) {
		err := loader.LoadJSON(viFile)
		require.NoError(t, err)

		// Verify translation was registered
		translation := manager.Translate("vi", "testmodule.hello", nil)
		assert.Equal(t, "Xin chào", translation)
	})

	t.Run("should return error for invalid JSON", func(t *testing.T) {
		err := loader.LoadJSON(invalidFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse translation file")
	})

	t.Run("should return error for non-existent file", func(t *testing.T) {
		err := loader.LoadJSON("/non/existent/file.json")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read translation file")
	})
}

// Test for ModuleLoader LoadJSONDirectory method
func TestModuleLoader_LoadJSONDirectory(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	localesDir := filepath.Join(tempDir, "locales")
	err := os.MkdirAll(localesDir, 0755)
	require.NoError(t, err)

	// Create test files
	enFile := filepath.Join(localesDir, "en.json")
	viFile := filepath.Join(localesDir, "vi.json")
	txtFile := filepath.Join(localesDir, "readme.txt")
	invalidFile := filepath.Join(localesDir, "invalid.json")

	err = os.WriteFile(enFile, []byte(`{"welcome": "Welcome", "error": "Error"}`), 0644)
	require.NoError(t, err)

	err = os.WriteFile(viFile, []byte(`{"welcome": "Chào mừng", "error": "Lỗi"}`), 0644)
	require.NoError(t, err)

	err = os.WriteFile(txtFile, []byte("This is not JSON"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(invalidFile, []byte(`{"broken": json}`), 0644)
	require.NoError(t, err)

	manager := NewManager(nil)
	loader := NewModuleLoader("testmodule", manager)

	t.Run("should load all valid JSON files from directory", func(t *testing.T) {
		// Capture stdout to check warnings
		err := loader.LoadJSONDirectory(localesDir)
		require.NoError(t, err) // Verify English translations were loaded
		translation := manager.Translate("en", "testmodule.welcome", nil)
		assert.Equal(t, "Welcome", translation)

		// Verify Vietnamese translations were loaded
		translation = manager.Translate("vi", "testmodule.welcome", nil)
		assert.Equal(t, "Chào mừng", translation)

		// Verify non-JSON files were ignored (no error)
		// Invalid JSON files should show warning but not fail the process
	})

	t.Run("should return error for non-existent directory", func(t *testing.T) {
		err := loader.LoadJSONDirectory("/non/existent/directory")
		assert.Error(t, err)
	})
}

// Test for ModuleLoader LoadFromFS method
func TestModuleLoader_LoadFromFS(t *testing.T) {
	manager := NewManager(nil)
	loader := NewModuleLoader("testmodule", manager)

	t.Run("should load all valid JSON files from filesystem", func(t *testing.T) {
		// Create test filesystem with only valid files
		testFS := fstest.MapFS{
			"locales/en.json": &fstest.MapFile{
				Data: []byte(`{"home": "Home", "about": "About"}`),
			},
			"locales/vi.json": &fstest.MapFile{
				Data: []byte(`{"home": "Trang chủ", "about": "Giới thiệu"}`),
			},
			"locales/config.txt": &fstest.MapFile{
				Data: []byte("Not a JSON file"),
			},
		}

		err := loader.LoadFromFS(testFS, "locales")
		require.NoError(t, err)

		// Verify translations were loaded
		translation := manager.Translate("en", "testmodule.home", nil)
		assert.Equal(t, "Home", translation)

		translation = manager.Translate("vi", "testmodule.home", nil)
		assert.Equal(t, "Trang chủ", translation)
	})

	t.Run("should return error for invalid JSON in filesystem", func(t *testing.T) {
		invalidFS := fstest.MapFS{
			"locales/invalid.json": &fstest.MapFile{
				Data: []byte(`{"broken": json}`),
			},
		}

		err := loader.LoadFromFS(invalidFS, "locales")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse translation file")
	})

	t.Run("should return error for non-existent directory", func(t *testing.T) {
		testFS := fstest.MapFS{
			"locales/en.json": &fstest.MapFile{
				Data: []byte(`{"test": "Test"}`),
			},
		}

		err := loader.LoadFromFS(testFS, "non-existent")
		assert.Error(t, err)
	})
}

// Test for RegisterModule function
func TestRegisterModule(t *testing.T) {
	// Create temporary directory with translation files
	tempDir := t.TempDir()
	localesDir := filepath.Join(tempDir, "locales")
	err := os.MkdirAll(localesDir, 0755)
	require.NoError(t, err)

	// Create test files
	enFile := filepath.Join(localesDir, "en.json")
	viFile := filepath.Join(localesDir, "vi.json")

	err = os.WriteFile(enFile, []byte(`{"signup": "Sign Up", "login": "Login"}`), 0644)
	require.NoError(t, err)

	err = os.WriteFile(viFile, []byte(`{"signup": "Đăng ký", "login": "Đăng nhập"}`), 0644)
	require.NoError(t, err)

	manager := NewManager(nil)

	t.Run("should register module translations successfully", func(t *testing.T) {
		err := RegisterModule(manager, "auth", localesDir)
		require.NoError(t, err) // Verify translations were registered
		translation := manager.Translate("en", "auth.signup", nil)
		assert.Equal(t, "Sign Up", translation)

		translation = manager.Translate("vi", "auth.login", nil)
		assert.Equal(t, "Đăng nhập", translation)
	})

	t.Run("should return error for non-existent directory", func(t *testing.T) {
		err := RegisterModule(manager, "nonexistent", "/non/existent/directory")
		assert.Error(t, err)
	})
}

// Test for RegisterModuleFromFS function
func TestRegisterModuleFromFS(t *testing.T) {
	// Create test filesystem
	testFS := fstest.MapFS{
		"locales/en.json": &fstest.MapFile{
			Data: []byte(`{"dashboard": "Dashboard", "settings": "Settings"}`),
		},
		"locales/vi.json": &fstest.MapFile{
			Data: []byte(`{"dashboard": "Bảng điều khiển", "settings": "Cài đặt"}`),
		},
	}

	manager := NewManager(nil)

	t.Run("should register module translations from filesystem", func(t *testing.T) {
		err := RegisterModuleFromFS(manager, "admin", testFS, "locales")
		require.NoError(t, err) // Verify translations were registered
		translation := manager.Translate("en", "admin.dashboard", nil)
		assert.Equal(t, "Dashboard", translation)

		translation = manager.Translate("vi", "admin.settings", nil)
		assert.Equal(t, "Cài đặt", translation)
	})

	t.Run("should return error for non-existent directory in filesystem", func(t *testing.T) {
		err := RegisterModuleFromFS(manager, "nonexistent", testFS, "non-existent")
		assert.Error(t, err)
	})

	t.Run("should handle invalid JSON in filesystem", func(t *testing.T) {
		invalidFS := fstest.MapFS{
			"locales/invalid.json": &fstest.MapFile{
				Data: []byte(`{"broken": json}`),
			},
		}

		err := RegisterModuleFromFS(manager, "invalid", invalidFS, "locales")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse translation file")
	})
}
