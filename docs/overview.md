# Locales Service Provider - Overview

Locales Service Provider cung cấp giải pháp quản lý đa ngôn ngữ (i18n) cho ứng dụng Go. Service này tập trung vào việc **load và truy xuất message** từ các nguồn khác nhau, không bao gồm logic xác định ngôn ngữ từ HTTP request.

## Kiến trúc và Nguyên tắc

### Tách biệt trách nhiệm

Service Provider này tuân theo nguyên tắc **Single Responsibility Principle** với các trách nhiệm rõ ràng:

- ✅ **Locales Service Provider**: Quản lý loader và truy xuất message
- ✅ **HTTP Middleware**: Xác định ngôn ngữ từ request (Accept-Language, query param, cookie)

### API Interface

```go
type Manager interface {
    // Dịch message với ngôn ngữ cụ thể
    Translate(lang, key string, args map[string]interface{}) string
    
    // Tạo function dịch cho ngôn ngữ cụ thể
    TranslateFunc(lang string) func(string, map[string]interface{}) string
    
    // Load từ thư mục
    LoadTranslations(directory string, loaderFunc func(string) ([]byte, error)) error
    
    // Load từ embedded filesystem
    LoadEmbeddedTranslations(embeddedFS embed.FS, directory string) error
    
    // Đăng ký translation cho module
    RegisterTranslations(module string, translations map[string]map[string]string)
    
    // Utility methods
    SetLocale(locale string)
    Bundle() *i18n.Bundle
}
```

## Tính năng chính

### 1. Multiple Loading Sources
- **File System**: Load từ thư mục chứa file JSON
- **Embedded FS**: Load từ `embed.FS` cho binary distribution  
- **Runtime Registration**: Đăng ký translation tại runtime cho module

### 2. Flexible Translation API
```go
// Dịch trực tiếp
msg := manager.Translate("vi", "welcome", map[string]interface{}{
    "name": "John",
})

// Tạo function cho ngôn ngữ cụ thể  
viTranslate := manager.TranslateFunc("vi")
msg1 := viTranslate("hello", nil)
msg2 := viTranslate("goodbye", map[string]interface{}{"name": "John"})
```

### 3. Module-specific Translations
```go
manager.RegisterTranslations("auth", map[string]map[string]string{
    "en": {
        "login_required": "Login required",
        "invalid_credentials": "Invalid credentials",
    },
    "vi": {
        "login_required": "Yêu cầu đăng nhập", 
        "invalid_credentials": "Thông tin đăng nhập không hợp lệ",
    },
})

// Sử dụng: manager.Translate("vi", "auth.login_required", nil)
```

## Cấu hình

```yaml
locales:
  default_locale: "en"           # Ngôn ngữ khởi tạo bundle
  locales_directory: "locales"   # Thư mục chứa file ngôn ngữ
  supported_locales:             # Danh sách ngôn ngữ hỗ trợ
    - "en"
    - "vi"
  bundle: "message"              # Tên bundle
  embedded_path: "locales"       # Đường dẫn trong embedded FS
```

## Tích hợp với Middleware

Service này hoạt động độc lập và được tích hợp với HTTP middleware để xác định ngôn ngữ:

```go
// Trong middleware
func (m *I18nMiddleware) Handle(c *gin.Context) {
    lang := m.detectLanguage(c) // Từ Accept-Language, query, cookie
    
    // Inject translate function vào context
    translateFunc := localesManager.TranslateFunc(lang)
    c.Set("translate", translateFunc)
    
    c.Next()
}

// Trong handler
func (h *Handler) GetUser(c *gin.Context) {
    translate := c.MustGet("translate").(func(string, map[string]interface{}) string)
    
    message := translate("user.not_found", map[string]interface{}{
        "id": userID,
    })
    
    c.JSON(404, gin.H{"error": message})
}
```

## Performance Features

- **Localizer Caching**: Cache `i18n.Localizer` instances per language
- **Concurrent Safe**: Thread-safe operations với `sync.RWMutex`
- **Fallback Strategy**: Trả về key nếu không tìm thấy translation
- **Lazy Loading**: Chỉ tạo localizer khi cần thiết

## Use Cases

1. **Web Applications**: Kết hợp với HTTP middleware cho multi-language websites
2. **APIs**: Trả về error message theo ngôn ngữ client
3. **CLI Applications**: Support multiple languages cho command line tools
4. **Microservices**: Shared translation service across services
5. **Module Systems**: Mỗi module có thể đăng ký translation riêng

## Migration từ phiên bản cũ

Phiên bản mới loại bỏ các method liên quan đến HTTP request:

```go
// Cũ ❌
manager.Get(key, args, locales...)
manager.MustGet(key, args, locales...)

// Mới ✅  
manager.Translate(lang, key, args)
manager.TranslateFunc(lang)(key, args)
```

Các cấu hình middleware được tách ra khỏi service config:

```yaml
# Cũ ❌ - Trong locales config
locales:
  accept_language: true
  query_param: "lang"
  cookie: "lang"

# Mới ✅ - Trong middleware config  
http:
  middleware:
    i18n:
      accept_language: true
      query_param: "lang"
      cookie: "lang"
```
