# Locales Loader

Loader component chịu trách nhiệm tải bản dịch từ các nguồn khác nhau và parsing chúng thành format mà i18n bundle có thể sử dụng.

## Loading Strategies

Locales Manager hỗ trợ ba chiến lược loading chính:

1. **File System Loading** - Từ thư mục trên disk
2. **Embedded FS Loading** - Từ embedded filesystem trong binary  
3. **Runtime Registration** - Đăng ký translations tại runtime

## File System Loading

### Cấu trúc thư mục

```
locales/
├── en.json          # English translations
├── vi.json          # Vietnamese translations  
├── fr.json          # French translations
├── es.json          # Spanish translations
├── de.json          # German translations
└── subdirectory/    # Nested directories supported
    ├── en.json
    └── vi.json
```

### API Usage

```go
// Basic loading
err := manager.LoadTranslations("./locales", nil)
if err != nil {
    log.Fatal("Failed to load translations:", err)
}

// Custom loader function
customLoader := func(filePath string) ([]byte, error) {
    // Custom logic: decrypt, validate, transform, etc.
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, err
    }
    
    // Example: decrypt file
    decryptedData, err := decrypt(data)
    if err != nil {
        return nil, err
    }
    
    return decryptedData, nil
}

err := manager.LoadTranslations("./encrypted_locales", customLoader)
```

### File Formats

Hiện tại chỉ hỗ trợ **JSON format**:

```json
{
  "welcome": "Welcome to our application",
  "hello_user": "Hello {{.name}}, welcome back!",
  "auth": {
    "login": "Login",
    "logout": "Logout",
    "register": "Register"
  },
  "validation": {
    "required": "This field is required",
    "email": "Please enter a valid email address",
    "min_length": "Must be at least {{.min}} characters long"
  },
  "errors": {
    "not_found": "{{.resource}} with ID {{.id}} not found",
    "unauthorized": "You are not authorized to perform this action",
    "server_error": "An unexpected error occurred. Please try again."
  }
}
```

### Error Handling

File System Loading sử dụng graceful error handling:

```go
// Errors are logged but don't stop the loading process
func (m *manager) LoadTranslations(directory string, loaderFunc func(string) ([]byte, error)) error {
    return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err // Directory access error - stop
        }
        
        if info.IsDir() {
            return nil // Skip directories
        }

        if filepath.Ext(path) == ".json" {
            // Load and parse file
            data, err := loadFile(path, loaderFunc)
            if err != nil {
                fmt.Printf("Error reading file %s: %v\n", path, err)
                return nil // Log error but continue with other files
            }

            _, err = m.bundle.ParseMessageFileBytes(data, path)
            if err != nil {
                fmt.Printf("Error parsing message file %s: %v\n", path, err)
                return nil // Log error but continue with other files
            }
        }
        return nil
    })
}
```

## Embedded FS Loading

### Setup

```go
//go:embed locales/*.json
var localesFS embed.FS

func main() {
    config := &locales.Config{
        DefaultLocale:    "en",
        SupportedLocales: []string{"en", "vi", "fr"},
        Bundle:           "message",
        EmbeddedFS:       &localesFS,
        EmbeddedPath:     "locales",
    }
    
    manager := locales.NewManager(config)
    // Embedded files are automatically loaded in NewManager
}
```

### Manual Loading

```go
//go:embed assets/translations/*.json
var translationsFS embed.FS

err := manager.LoadEmbeddedTranslations(translationsFS, "assets/translations")
if err != nil {
    log.Fatal("Failed to load embedded translations:", err)
}
```

### Directory Structure trong Embedded FS

```
// Trong source code
assets/
└── translations/
    ├── en.json
    ├── vi.json
    ├── fr.json
    └── modules/
        ├── auth/
        │   ├── en.json
        │   └── vi.json
        └── admin/
            ├── en.json
            └── vi.json
```

### Binary Distribution Benefits

- **Single Binary**: Tất cả translations được đóng gói trong binary
- **No External Dependencies**: Không cần external files hoặc directories
- **Performance**: Faster loading, no disk I/O sau khi binary load
- **Security**: Translations không thể bị modify externally

## Runtime Registration

### Basic Registration

```go
// Đăng ký translations cho auth module
manager.RegisterTranslations("auth", map[string]map[string]string{
    "en": {
        "login_required":      "Authentication required",
        "invalid_credentials": "Invalid username or password",
        "token_expired":       "Your session has expired",
        "access_denied":       "Access denied",
    },
    "vi": {
        "login_required":      "Yêu cầu xác thực",
        "invalid_credentials": "Tên đăng nhập hoặc mật khẩu không đúng",
        "token_expired":       "Phiên làm việc đã hết hạn",
        "access_denied":       "Truy cập bị từ chối",
    },
    "fr": {
        "login_required":      "Authentification requise",
        "invalid_credentials": "Nom d'utilisateur ou mot de passe invalide",
        "token_expired":       "Votre session a expiré",
        "access_denied":       "Accès refusé",
    },
})

// Sử dụng với module prefix
msg := manager.Translate("vi", "auth.login_required", nil)
// Kết quả: "Yêu cầu xác thực"
```

### Plugin System Integration

```go
// Plugin interface
type Plugin interface {
    Name() string
    RegisterTranslations(manager locales.Manager)
}

// Auth plugin implementation
type AuthPlugin struct{}

func (p *AuthPlugin) Name() string {
    return "auth"
}

func (p *AuthPlugin) RegisterTranslations(manager locales.Manager) {
    manager.RegisterTranslations("auth", map[string]map[string]string{
        "en": authEnTranslations,
        "vi": authViTranslations,
        "fr": authFrTranslations,
    })
}

// Plugin loader
func LoadPlugins(manager locales.Manager, plugins []Plugin) {
    for _, plugin := range plugins {
        plugin.RegisterTranslations(manager)
        log.Printf("Loaded translations for plugin: %s", plugin.Name())
    }
}
```

### Dynamic Module Registration

```go
// API endpoint để đăng ký translations
func RegisterModuleTranslations(c *gin.Context) {
    var request struct {
        Module       string                            `json:"module"`
        Translations map[string]map[string]string     `json:"translations"`
    }
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Validate request
    if request.Module == "" {
        c.JSON(400, gin.H{"error": "module name is required"})
        return
    }
    
    // Register translations
    localesManager.RegisterTranslations(request.Module, request.Translations)
    
    c.JSON(200, gin.H{
        "message": fmt.Sprintf("Translations registered for module: %s", request.Module),
    })
}
```

## Loading Priority & Fallback

Manager sử dụng loading strategy theo thứ tự ưu tiên:

```go
func (m *manager) Translate(lang, key string, args map[string]interface{}) string {
    // 1. Kiểm tra trong Runtime Registration trước (module-specific)
    m.mutex.RLock()
    if moduleMsgs, ok := m.messageMap[lang]; ok {
        if msg, ok := moduleMsgs[key]; ok {
            // Found in runtime registration
            result := processTemplate(msg, args)
            m.mutex.RUnlock()
            return result
        }
    }
    m.mutex.RUnlock()

    // 2. Fallback đến Bundle (từ file system hoặc embedded FS)
    loc := m.getLocalizer(lang)
    result, err := loc.Localize(&i18n.LocalizeConfig{
        MessageID:    key,
        TemplateData: args,
    })
    
    if err != nil {
        // 3. Ultimate fallback - trả về key
        return key
    }

    return result
}
```

## Advanced Loading Scenarios

### 1. Multi-source Loading

```go
func setupMultiSourceLoading(manager locales.Manager) error {
    // 1. Load base translations từ embedded FS
    //go:embed base_locales/*.json
    var baseFS embed.FS
    
    err := manager.LoadEmbeddedTranslations(baseFS, "base_locales")
    if err != nil {
        return fmt.Errorf("failed to load base translations: %w", err)
    }
    
    // 2. Load custom translations từ file system (override base)
    err = manager.LoadTranslations("./custom_locales", nil)
    if err != nil {
        log.Printf("Custom translations not found, using base only: %v", err)
    }
    
    // 3. Load plugin translations
    loadPluginTranslations(manager)
    
    return nil
}
```

### 2. Hot Reload cho Development

```go
func watchTranslationFiles(manager locales.Manager, directory string) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    err = watcher.Add(directory)
    if err != nil {
        log.Fatal(err)
    }

    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            
            if event.Op&fsnotify.Write == fsnotify.Write {
                log.Printf("Translation file modified: %s", event.Name)
                
                // Reload translations
                err := manager.LoadTranslations(directory, nil)
                if err != nil {
                    log.Printf("Failed to reload translations: %v", err)
                }
            }
            
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            log.Printf("Watcher error: %v", err)
        }
    }
}
```

### 3. Lazy Loading với Cache

```go
type LazyLoader struct {
    manager   locales.Manager
    directory string
    loaded    map[string]bool
    mutex     sync.Mutex
}

func (l *LazyLoader) LoadLanguage(lang string) error {
    l.mutex.Lock()
    defer l.mutex.Unlock()
    
    if l.loaded[lang] {
        return nil // Already loaded
    }
    
    filePath := filepath.Join(l.directory, lang+".json")
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return fmt.Errorf("translation file not found: %s", filePath)
    }
    
    // Load specific language file
    data, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }
    
    _, err = l.manager.Bundle().ParseMessageFileBytes(data, filePath)
    if err != nil {
        return err
    }
    
    l.loaded[lang] = true
    return nil
}
```

## Best Practices

### 1. File Organization

```
locales/
├── common/           # Shared translations
│   ├── en.json
│   ├── vi.json
│   └── fr.json
├── modules/          # Module-specific
│   ├── auth/
│   ├── admin/
│   └── user/
└── errors/           # Error messages
    ├── en.json
    ├── vi.json
    └── fr.json
```

### 2. Key Naming Conventions

```json
{
  "common": {
    "buttons": {
      "save": "Save",
      "cancel": "Cancel",
      "delete": "Delete"
    },
    "messages": {
      "success": "Operation completed successfully",
      "loading": "Loading..."
    }
  },
  "auth": {
    "forms": {
      "login": {
        "title": "Login",
        "email_placeholder": "Enter your email",
        "password_placeholder": "Enter your password"
      }
    },
    "errors": {
      "invalid_credentials": "Invalid email or password"
    }
  }
}
```

### 3. Template Parameters

```json
{
  "user": {
    "profile": {
      "welcome": "Welcome back, {{.name}}!",
      "last_login": "Last login: {{.date}}",
      "messages_count": "You have {{.count}} new {{.count | plural \"message\" \"messages\"}}"
    }
  }
}
```

### 4. Error Handling Strategy

```go
// Graceful degradation
func safeTranslate(manager locales.Manager, lang, key string, args map[string]interface{}) string {
    result := manager.Translate(lang, key, args)
    
    // If translation failed, try fallback language
    if result == key && lang != "en" {
        result = manager.Translate("en", key, args)
    }
    
    // If still failed, return user-friendly message
    if result == key {
        return fmt.Sprintf("Translation missing: %s", key)
    }
    
    return result
}
```

### 5. Performance Optimization

```go
// Pre-warm localizers cho common languages
func preWarmLocalizers(manager locales.Manager, languages []string) {
    for _, lang := range languages {
        // Trigger localizer creation
        _ = manager.TranslateFunc(lang)
    }
}

// Use trong application startup
func init() {
    manager := locales.NewManager(config)
    preWarmLocalizers(manager, []string{"en", "vi", "fr", "es"})
}
```
