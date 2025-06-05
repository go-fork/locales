# Locales Service Provider

**Locales Service Provider** cung cấp giải pháp quản lý đa ngôn ngữ (i18n) chuyên nghiệp cho ứng dụng Go, tập trung vào **loader và truy xuất message** với hiệu suất cao và kiến trúc modular.

## 🎯 Đặc điểm chính

- ✅ **Single Responsibility**: Chỉ quản lý loading và translation, tách biệt HTTP logic
- ✅ **Multiple Sources**: Hỗ trợ file system, embedded FS, và runtime registration
- ✅ **Thread Safe**: Concurrent operations với RWMutex optimization  
- ✅ **Module Based**: Plugin-friendly architecture cho microservices
- ✅ **High Performance**: Localizer caching và lazy loading
- ✅ **Graceful Fallback**: Robust error handling và fallback strategies

## 🚀 Quick Start

### Installation

```bash
go get github.com/Fork/providers/locales
```

### Basic Usage

```go
package main

import (
    "github.com/Fork/providers/locales"
)

func main() {
    // Setup
    config := locales.DefaultConfig()
    config.SupportedLocales = []string{"en", "vi", "fr"}
    
    manager := locales.NewManager(config)
    manager.LoadTranslations("./locales", nil)
    
    // Direct translation
    msg := manager.Translate("vi", "welcome", map[string]interface{}{
        "name": "John",
    })
    
    // Function for specific language
    viTranslate := manager.TranslateFunc("vi")
    greeting := viTranslate("hello", nil)
    farewell := viTranslate("goodbye", nil)
}
```

## ⚙️ Configuration

```yaml
# config/app.yaml
locales:
  default_locale: "en"           # Language for bundle initialization
  locales_directory: "locales"   # Directory containing translation files
  supported_locales:             # List of supported languages
    - "en"
    - "vi"
    - "fr"
  bundle: "message"              # Bundle name
  embedded_path: "locales"       # Path in embedded FS (optional)
```

## 📁 Translation Files

```json
// locales/vi.json
{
  "welcome": "Chào mừng {{.name}}",
  "auth": {
    "login": "Đăng nhập",
    "logout": "Đăng xuất"
  },
  "validation": {
    "required": "Trường này là bắt buộc",
    "email": "Email không hợp lệ"
  }
}
```

## 🔧 API Reference

### Core Methods

```go
// Direct translation
msg := manager.Translate("vi", "welcome", map[string]interface{}{
    "name": "John",
})

// Language-specific function
viTranslate := manager.TranslateFunc("vi")
msg1 := viTranslate("auth.login", nil)
msg2 := viTranslate("validation.required", nil)
```

### Loading Sources

```go
// 1. File System
err := manager.LoadTranslations("./locales", nil)

// 2. Embedded FS
//go:embed locales/*.json
var localesFS embed.FS
err := manager.LoadEmbeddedTranslations(localesFS, "locales")

// 3. Runtime Registration (for modules/plugins)
manager.RegisterTranslations("auth", map[string]map[string]string{
    "en": {"login": "Login"},
    "vi": {"login": "Đăng nhập"},
})
```

## 🌐 Web Integration

### With Gin Framework

```go
// Middleware
func I18nMiddleware(manager locales.Manager) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Detect language from Accept-Language, query param, cookie
        lang := detectLanguage(c)
        
        // Inject translate function
        translateFunc := manager.TranslateFunc(lang)
        c.Set("translate", translateFunc)
        
        c.Next()
    }
}

// Handler
func GetUser(c *gin.Context) {
    translate := c.MustGet("translate").(func(string, map[string]interface{}) string)
    
    if user == nil {
        message := translate("user.not_found", map[string]interface{}{
            "id": userID,
        })
        c.JSON(404, gin.H{"error": message})
        return
    }
    
    c.JSON(200, user)
}
```

## 🏗️ Architecture

### Separation of Concerns

- **Locales Service Provider**: Manages loading and message retrieval
- **HTTP Middleware**: Handles language detection from requests
- **Clean API**: Explicit language parameter, no hidden dependencies

### Loading Priority

1. **Runtime Registration** (module-specific translations)
2. **Bundle Translations** (from files or embedded FS)  
3. **Fallback** (returns key if not found)

## 📦 Multiple Loading Strategies

### 1. File System Loading
```
locales/
├── en.json
├── vi.json
├── fr.json
└── modules/
    ├── auth/
    └── admin/
```

### 2. Embedded FS Loading
```go
//go:embed locales/*.json
var localesFS embed.FS

config := &locales.Config{
    EmbeddedFS:   &localesFS,
    EmbeddedPath: "locales",
}
```

### 3. Module Registration
```go
// For plugin systems
authPlugin.RegisterTranslations(manager)
adminPlugin.RegisterTranslations(manager)
```

## 🚀 Performance Features

- **Localizer Caching**: One cached localizer per language
- **Lazy Loading**: Localizers created only when needed
- **Thread Safe**: Concurrent access with RWMutex
- **Memory Efficient**: Only cache what's actually used
- **Fallback Strategy**: Graceful degradation on missing translations

## 📚 Documentation

- **[Overview](docs/overview.md)** - Architecture and design principles
- **[Manager](docs/manager.md)** - API reference and examples
- **[Config](docs/config.md)** - Configuration guide and best practices  
- **[Loader](docs/loader.md)** - Loading strategies and file formats

## 🔄 Migration from v1

### API Changes
```go
// Old ❌
manager.Get(key, args, locales...)
manager.MustGet(key, args, locales...)

// New ✅
manager.Translate(lang, key, args)
manager.TranslateFunc(lang)(key, args)
```

### Config Changes
```yaml
# Old ❌ - Mixed HTTP concerns
locales:
  accept_language: true
  query_param: "lang"
  cookie: "lang"

# New ✅ - Separated concerns
locales:
  default_locale: "en"
  locales_directory: "locales"
  
http:
  middleware:
    i18n:
      accept_language: true
      query_param: "lang" 
      cookie: "lang"
```

## 🧪 Testing

```go
func TestTranslation(t *testing.T) {
    config := locales.DefaultConfig()
    manager := locales.NewManager(config)
    
    // Register test translations
    manager.RegisterTranslations("test", map[string]map[string]string{
        "en": {"hello": "Hello {{.name}}"},
        "vi": {"hello": "Xin chào {{.name}}"},
    })
    
    result := manager.Translate("vi", "test.hello", map[string]interface{}{
        "name": "John",
    })
    
    assert.Equal(t, "Xin chào John", result)
}
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Related Projects

- **[go-i18n](https://github.com/nicksnyder/go-i18n)** - Underlying i18n library
- **[Fork/providers](github.com/go-fork/providers)** - Collection of service providers

---

*Made with ❤️ by the Fork team*
