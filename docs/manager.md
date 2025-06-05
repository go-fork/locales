# Locales Manager

Manager là thành phần chính của Locales Service Provider, chịu trách nhiệm quản lý và truy xuất các bản dịch đa ngôn ngữ.

## Interface Definition

```go
type Manager interface {
    // Translation Methods
    Translate(lang, key string, args map[string]interface{}) string
    TranslateFunc(lang string) func(string, map[string]interface{}) string
    
    // Loading Methods
    LoadTranslations(directory string, loaderFunc func(string) ([]byte, error)) error
    LoadEmbeddedTranslations(embeddedFS embed.FS, directory string) error
    RegisterTranslations(module string, translations map[string]map[string]string)
    
    // Utility Methods
    SetLocale(locale string)
    Bundle() *i18n.Bundle
}
```

## Core Methods

### Translation Methods

#### `Translate(lang, key string, args map[string]interface{}) string`

Dịch một message key thành ngôn ngữ cụ thể.

**Tham số:**
- `lang`: Mã ngôn ngữ (ví dụ: "vi", "en", "fr")
- `key`: Khóa của message cần dịch
- `args`: Map các tham số để thay thế trong template

**Trả về:** String đã được dịch, hoặc key gốc nếu không tìm thấy bản dịch.

**Ví dụ:**
```go
// Basic usage
msg := manager.Translate("vi", "welcome", nil)
// Kết quả: "Chào mừng"

// With parameters
msg := manager.Translate("vi", "hello_user", map[string]interface{}{
    "name": "John",
    "role": "admin",
})
// Template: "Xin chào {{.name}}, bạn là {{.role}}"
// Kết quả: "Xin chào John, bạn là admin"
```

#### `TranslateFunc(lang string) func(string, map[string]interface{}) string`

Tạo một function closure để dịch cho ngôn ngữ cụ thể. Hữu ích khi cần dịch nhiều message cho cùng một ngôn ngữ.

**Tham số:**
- `lang`: Mã ngôn ngữ

**Trả về:** Function có signature `func(string, map[string]interface{}) string`

**Ví dụ:**
```go
// Tạo function cho tiếng Việt
viTranslate := manager.TranslateFunc("vi")

// Sử dụng multiple times
title := viTranslate("page.title", nil)
desc := viTranslate("page.description", map[string]interface{}{
    "count": 10,
})
error := viTranslate("validation.required", map[string]interface{}{
    "field": "email",
})
```

### Loading Methods

#### `LoadTranslations(directory string, loaderFunc func(string) ([]byte, error)) error`

Load bản dịch từ thư mục chứa các file JSON.

**Tham số:**
- `directory`: Đường dẫn đến thư mục chứa file ngôn ngữ
- `loaderFunc`: Function tùy chỉnh để đọc file. Nếu `nil`, sử dụng `os.ReadFile`

**Cấu trúc thư mục:**
```
locales/
├── en.json
├── vi.json
├── fr.json
└── es.json
```

**Format file JSON:**
```json
{
  "welcome": "Chào mừng",
  "hello_user": "Xin chào {{.name}}",
  "validation": {
    "required": "Trường {{.field}} là bắt buộc",
    "email": "Email không hợp lệ"
  }
}
```

**Ví dụ:**
```go
// Load từ thư mục
err := manager.LoadTranslations("./locales", nil)
if err != nil {
    log.Fatal(err)
}

// Load với custom loader
err := manager.LoadTranslations("./locales", func(path string) ([]byte, error) {
    // Custom logic to read file
    return customReadFile(path)
})
```

#### `LoadEmbeddedTranslations(embeddedFS embed.FS, directory string) error`

Load bản dịch từ embedded filesystem. Hữu ích cho việc đóng gói translation files vào binary.

**Tham số:**
- `embeddedFS`: Embedded filesystem chứa file bản dịch
- `directory`: Đường dẫn trong embedded FS

**Ví dụ:**
```go
//go:embed locales/*.json
var localesFS embed.FS

func init() {
    manager := locales.NewManager(config)
    err := manager.LoadEmbeddedTranslations(localesFS, "locales")
    if err != nil {
        log.Fatal(err)
    }
}
```

#### `RegisterTranslations(module string, translations map[string]map[string]string)`

Đăng ký bản dịch cho module cụ thể tại runtime. Hữu ích cho plugin systems hoặc dynamic modules.

**Tham số:**
- `module`: Tên module (sẽ được prefix vào key)
- `translations`: Map translation theo ngôn ngữ

**Ví dụ:**
```go
// Đăng ký translation cho auth module
manager.RegisterTranslations("auth", map[string]map[string]string{
    "en": {
        "login_required": "Login required",
        "invalid_credentials": "Invalid username or password",
        "token_expired": "Your session has expired",
    },
    "vi": {
        "login_required": "Yêu cầu đăng nhập",
        "invalid_credentials": "Tên đăng nhập hoặc mật khẩu không hợp lệ", 
        "token_expired": "Phiên làm việc đã hết hạn",
    },
})

// Sử dụng với prefix
msg := manager.Translate("vi", "auth.login_required", nil)
// Kết quả: "Yêu cầu đăng nhập"
```

### Utility Methods

#### `SetLocale(locale string)`

Thiết lập ngôn ngữ mặc định cho bundle. Thay đổi `DefaultLocale` trong config.

**Ví dụ:**
```go
manager.SetLocale("vi")
```

#### `Bundle() *i18n.Bundle`

Trả về đối tượng `i18n.Bundle` bên trong để truy cập advanced features.

**Ví dụ:**
```go
bundle := manager.Bundle()
// Sử dụng bundle trực tiếp cho advanced operations
```

## Loading Strategy

Manager sử dụng chiến lược loading theo thứ tự ưu tiên:

1. **Module-specific translations** (được đăng ký qua `RegisterTranslations`)
2. **Bundle translations** (được load từ files hoặc embedded FS)
3. **Fallback** (trả về key gốc)

## Caching & Performance

- **Localizer Caching**: Mỗi ngôn ngữ có một `i18n.Localizer` được cache
- **Lazy Loading**: Localizer chỉ được tạo khi cần thiết
- **Thread Safe**: Sử dụng `sync.RWMutex` cho concurrent access
- **Memory Efficient**: Chỉ cache những localizer thực sự được sử dụng

## Error Handling

Manager sử dụng graceful error handling:

- Nếu không tìm thấy bản dịch, trả về key gốc
- Load errors được log nhưng không gây crash
- File parsing errors được skip và tiếp tục với file khác

## Best Practices

### 1. Sử dụng TranslateFunc cho cùng ngôn ngữ
```go
// ✅ Tốt - tạo function một lần
translate := manager.TranslateFunc(userLang)
title := translate("page.title", nil)
desc := translate("page.desc", nil)

// ❌ Không tối ưu - tạo lại localizer mỗi lần
title := manager.Translate(userLang, "page.title", nil)  
desc := manager.Translate(userLang, "page.desc", nil)
```

### 2. Organize keys theo hierarchy
```json
{
  "auth": {
    "login": "Đăng nhập",
    "logout": "Đăng xuất",
    "register": "Đăng ký"
  },
  "validation": {
    "required": "Trường này là bắt buộc",
    "email": "Email không hợp lệ"  
  }
}
```

### 3. Sử dụng meaningful fallbacks
```go
// Template có placeholder rõ ràng
msg := manager.Translate(lang, "error.not_found", map[string]interface{}{
    "resource": "User",
    "id": userID,
})
// "User with ID {{.id}} not found"
```

### 4. Module registration pattern
```go
// Trong package auth
func RegisterTranslations(manager locales.Manager) {
    manager.RegisterTranslations("auth", authTranslations)
}

// Trong main application
authModule.RegisterTranslations(localesManager)
```
