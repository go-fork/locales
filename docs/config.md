# Locales Configuration

Configuration cho Locales Service Provider được thiết kế đơn giản và tập trung vào việc quản lý loader và truy xuất message. Các cấu hình liên quan đến HTTP request handling được tách ra cho middleware.

## Config Structure

```go
type Config struct {
    // DefaultLocale là ngôn ngữ mặc định để khởi tạo bundle.
    DefaultLocale string `yaml:"default_locale" json:"default_locale" mapstructure:"default_locale"`

    // LocalesDirectory là thư mục chứa các tệp ngôn ngữ.
    LocalesDirectory string `yaml:"locales_directory" json:"locales_directory" mapstructure:"locales_directory"`

    // SupportedLocales là danh sách các ngôn ngữ được hỗ trợ.
    SupportedLocales []string `yaml:"supported_locales" json:"supported_locales" mapstructure:"supported_locales"`

    // Bundle là cấu hình bundle cho các tệp ngôn ngữ.
    Bundle string `yaml:"bundle" json:"bundle" mapstructure:"bundle"`

    // EmbeddedFS cho phép tải bản dịch từ embed.FS thay vì thư mục.
    EmbeddedFS *embed.FS `yaml:"-" json:"-" mapstructure:"-"`

    // EmbeddedPath chỉ định đường dẫn trong EmbeddedFS nơi chứa file bản dịch.
    EmbeddedPath string `yaml:"embedded_path" json:"embedded_path" mapstructure:"embedded_path"`
}
```

## Configuration Fields

### `DefaultLocale`
- **Type**: `string`
- **Required**: Yes
- **Default**: `"en"`
- **Purpose**: Ngôn ngữ mặc định để khởi tạo i18n Bundle
- **Example**: `"en"`, `"vi"`, `"fr"`

### `LocalesDirectory`
- **Type**: `string`  
- **Required**: Yes (khi không sử dụng EmbeddedFS)
- **Default**: `"locales"`
- **Purpose**: Đường dẫn đến thư mục chứa file ngôn ngữ
- **Example**: `"locales"`, `"assets/translations"`, `"./i18n"`

### `SupportedLocales`
- **Type**: `[]string`
- **Required**: Yes
- **Default**: `["en"]`
- **Purpose**: Danh sách các ngôn ngữ được hỗ trợ bởi ứng dụng
- **Example**: `["en", "vi", "fr", "es"]`

### `Bundle`
- **Type**: `string`
- **Required**: Yes
- **Default**: `"message"`
- **Purpose**: Tên bundle cho i18n package
- **Example**: `"message"`, `"app"`, `"translations"`

### `EmbeddedFS`
- **Type**: `*embed.FS`
- **Required**: No
- **Default**: `nil`
- **Purpose**: Embedded filesystem để load translations từ binary
- **Usage**: Chỉ dùng khi muốn embed translation files vào binary

### `EmbeddedPath`
- **Type**: `string`
- **Required**: No (required khi sử dụng EmbeddedFS)
- **Default**: `""`
- **Purpose**: Đường dẫn trong EmbeddedFS chứa translation files
- **Example**: `"locales"`, `"assets/i18n"`

## YAML Configuration

### Basic Configuration
```yaml
locales:
  default_locale: "en"
  locales_directory: "locales"
  supported_locales:
    - "en"
    - "vi"
    - "fr"
    - "es"
  bundle: "message"
```

### With Embedded Path
```yaml
locales:
  default_locale: "en"
  locales_directory: "locales"
  supported_locales:
    - "en"
    - "vi"
  bundle: "message"
  embedded_path: "assets/locales"
```

### Production Configuration
```yaml
locales:
  default_locale: "en"
  locales_directory: "/app/locales"
  supported_locales:
    - "en"
    - "vi"
    - "fr"
    - "es"
    - "de"
    - "it"
    - "pt"
    - "zh"
    - "ja"
    - "ru"
    - "ko"
  bundle: "app_messages"
  embedded_path: "locales"
```

## Environment Variables

Configuration có thể được override bằng environment variables:

```bash
# Override default locale
LOCALES_DEFAULT_LOCALE=vi

# Override locales directory
LOCALES_LOCALES_DIRECTORY=/custom/path

# Override bundle name
LOCALES_BUNDLE=custom_bundle

# Override embedded path
LOCALES_EMBEDDED_PATH=custom/locales
```

## Default Configuration

```go
func DefaultConfig() *Config {
    return &Config{
        DefaultLocale:    "en",
        LocalesDirectory: "locales",
        Bundle:           "message",
        SupportedLocales: []string{"en"},
    }
}
```

## Configuration Examples

### 1. File System Loading
```go
config := &locales.Config{
    DefaultLocale:    "vi",
    LocalesDirectory: "./translations",
    SupportedLocales: []string{"en", "vi", "fr"},
    Bundle:           "app",
}

manager := locales.NewManager(config)
```

### 2. Embedded FS Loading
```go
//go:embed locales/*.json
var localesFS embed.FS

config := &locales.Config{
    DefaultLocale:    "en",
    SupportedLocales: []string{"en", "vi"},
    Bundle:           "message",
    EmbeddedFS:       &localesFS,
    EmbeddedPath:     "locales",
}

manager := locales.NewManager(config)
```

### 3. Hybrid Loading
```go
//go:embed locales/*.json
var localesFS embed.FS

config := &locales.Config{
    DefaultLocale:    "en",
    LocalesDirectory: "./custom_locales", // Fallback to file system
    SupportedLocales: []string{"en", "vi", "fr"},
    Bundle:           "message",
    EmbeddedFS:       &localesFS,         // Primary loading
    EmbeddedPath:     "locales",
}

manager := locales.NewManager(config)
```

## Loading Priority

Manager sẽ load translations theo thứ tự ưu tiên:

1. **EmbeddedFS** (nếu được cấu hình)
2. **LocalesDirectory** (nếu EmbeddedFS không có hoặc thất bại)
3. **Runtime Registration** (qua `RegisterTranslations`)

## Validation Rules

### Required Fields
- `DefaultLocale` phải được cung cấp
- `SupportedLocales` phải chứa ít nhất một ngôn ngữ
- `Bundle` phải được cung cấp
- Ít nhất một trong `LocalesDirectory` hoặc `EmbeddedFS` phải được cấu hình

### Language Code Format
- Sử dụng ISO 639-1 codes (2 chữ cái): `en`, `vi`, `fr`
- Hoặc BCP 47 tags: `en-US`, `zh-CN`, `pt-BR`

### Directory Structure
```
locales/
├── en.json
├── vi.json
├── fr.json
└── es.json
```

### File Format
Chỉ hỗ trợ JSON format:
```json
{
  "welcome": "Welcome",
  "hello_user": "Hello {{.name}}",
  "nested": {
    "message": "This is nested"
  }
}
```

## Configuration Best Practices

### 1. Use Environment-specific Configs
```yaml
# config/development.yaml
locales:
  default_locale: "en"
  locales_directory: "./locales"
  
# config/production.yaml  
locales:
  default_locale: "en"
  locales_directory: "/app/locales"
  embedded_path: "locales"
```

### 2. Validate Supported Locales
```go
func ValidateConfig(config *locales.Config) error {
    if len(config.SupportedLocales) == 0 {
        return errors.New("supported_locales cannot be empty")
    }
    
    found := false
    for _, lang := range config.SupportedLocales {
        if lang == config.DefaultLocale {
            found = true
            break
        }
    }
    
    if !found {
        return errors.New("default_locale must be in supported_locales")
    }
    
    return nil
}
```

### 3. Use Meaningful Bundle Names
```yaml
# Good ✅
locales:
  bundle: "app_messages"      # Clear purpose
  bundle: "user_interface"    # Specific scope
  bundle: "api_responses"     # Domain specific

# Avoid ❌  
locales:
  bundle: "messages"          # Too generic
  bundle: "data"              # Unclear purpose
```

### 4. Organize by Environment
```bash
configs/
├── app.yaml           # Base config
├── development.yaml   # Dev overrides
├── staging.yaml       # Staging overrides
└── production.yaml    # Prod overrides
```

## Migration from Old Config

### Removed Fields
Các field sau đã được loại bỏ khỏi locales config và chuyển sang middleware:

```yaml
# ❌ Removed from locales config
locales:
  accept_language: true      # -> http.middleware.i18n.accept_language
  query_param: "lang"        # -> http.middleware.i18n.query_param  
  cookie: "lang"             # -> http.middleware.i18n.cookie
  cookie_max_age: 2592000    # -> http.middleware.i18n.cookie_max_age
```

### Migration Example
```yaml
# Old config ❌
locales:
  default_locale: "en"
  accept_language: true
  query_param: "lang"
  cookie: "lang"

# New config ✅
locales:
  default_locale: "en"
  locales_directory: "locales"
  supported_locales: ["en", "vi"]
  
http:
  middleware:
    i18n:
      accept_language: true
      query_param: "lang"
      cookie: "lang"
```
