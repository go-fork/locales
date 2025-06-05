# Migration Guide - v0.1.0

## Overview
Đây là hướng dẫn setup và sử dụng Go Locales Service Provider v0.1.0 - phiên bản đầu tiên của thư viện.

## Prerequisites
- Go 1.23.9 hoặc cao hơn
- Hiểu biết cơ bản về Go modules và i18n concepts

## Installation

### Step 1: Install Package
```bash
go get go.fork.vn/locales
```

### Step 2: Import in Your Project
```go
import (
    "go.fork.vn/locales"
)
```

## Quick Setup Guide

### Basic Usage
```go
package main

import (
    "go.fork.vn/locales"
)

func main() {
    // Tạo configuration
    config := locales.DefaultConfig()
    config.SupportedLocales = []string{"en", "vi", "fr"}
    config.DefaultLocale = "en"
    
    // Khởi tạo manager
    manager := locales.NewManager(config)
    
    // Load translations từ thư mục
    err := manager.LoadTranslations("./locales", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Sử dụng translation
    message := manager.Translate("vi", "welcome", map[string]interface{}{
        "name": "John",
    })
    fmt.Println(message)
}
```

### Directory Structure
Tạo thư mục locales với cấu trúc như sau:
```
locales/
├── en.json
├── vi.json
└── fr.json
```

**Example en.json:**
```json
{
  "welcome": "Welcome {{.name}}",
  "auth": {
    "login": "Login",
    "logout": "Logout"
  }
}
```

**Example vi.json:**
```json
{
  "welcome": "Chào mừng {{.name}}",
  "auth": {
    "login": "Đăng nhập", 
    "logout": "Đăng xuất"
  }
}
```

## Configuration Options

### Basic Config
```go
type Config struct {
    DefaultLocale     string   // Default language (e.g., "en")
    SupportedLocales  []string // Supported languages
    LocalesDirectory  string   // Directory containing translation files
    Bundle           string   // Bundle name for i18n
    EmbeddedFS       *embed.FS // For embedded translations
    EmbeddedPath     string   // Path in embedded FS
}
```

### With Service Provider (go.fork.vn/di)
```go
package main

import (
    "go.fork.vn/core"
    "go.fork.vn/config"
    "go.fork.vn/locales"
)

func main() {
    app := core.NewApplication()
    
    // Register service providers
    app.Register(config.NewServiceProvider())
    app.Register(locales.NewServiceProvider())
    
    // Boot application
    app.Boot()
    
    // Get locales manager from container
    manager := app.Make("locales.manager").(locales.Manager)
}
```

## Loading Strategies

### 1. File System Loading
```go
manager := locales.NewManager(config)
err := manager.LoadTranslations("./locales", nil)
```

### 2. Embedded FS Loading
```go
//go:embed locales/*.json
var localesFS embed.FS

config := &locales.Config{
    EmbeddedFS:   &localesFS,
    EmbeddedPath: "locales",
}
manager := locales.NewManager(config)
err := manager.LoadEmbeddedTranslations(localesFS, "locales")
```

### 3. Runtime Registration (For Plugins)
```go
manager.RegisterTranslations("auth", map[string]map[string]string{
    "en": {
        "login": "Login",
        "logout": "Logout",
    },
    "vi": {
        "login": "Đăng nhập",
        "logout": "Đăng xuất", 
    },
})
```

## Module System

### Basic Module Registration
```go
// Register module từ directory
err := locales.RegisterModule(manager, "auth", "./modules/auth/locales")

// Register module từ embedded FS
//go:embed modules/auth/locales/*.json
var authLocalesFS embed.FS
err := locales.RegisterModuleFromFS(manager, "auth", authLocalesFS, "modules/auth/locales")
```

### Using ModuleLoader
```go
loader := locales.NewModuleLoader("auth", manager)

// Load single file
err := loader.LoadJSON("./auth/en.json")

// Load entire directory
err := loader.LoadJSONDirectory("./auth/locales")

// Register individual messages
loader.RegisterMessage("en", "welcome", "Welcome to Auth Module")
```

## Best Practices

### 1. Key Organization
```json
{
  "module": {
    "feature": {
      "action": "Translation text"
    }
  }
}
```

### 2. Parameter Usage
```json
{
  "user_message": "Hello {{.name}}, you have {{.count}} notifications"
}
```

```go
message := manager.Translate("en", "user_message", map[string]interface{}{
    "name":  "John",
    "count": 5,
})
```

### 3. Error Handling
```go
message := manager.Translate("vi", "some.key", nil)
if message == "some.key" {
    // Translation not found, handle fallback
    message = manager.Translate("en", "some.key", nil)
}
```

## Common Patterns

### Language-specific Functions
```go
viTranslate := manager.TranslateFunc("vi")
msg1 := viTranslate("welcome", nil)
msg2 := viTranslate("goodbye", nil)
```

### Module-based Organization
```go
// Auth module
manager.RegisterTranslations("auth", authTranslations)

// Use with module prefix
loginMsg := manager.Translate("vi", "auth.login", nil)
```

## Troubleshooting

### Translation Not Found
- Kiểm tra file path và format
- Verify key exists trong translation file
- Check language code chính xác
- Ensure file được load thành công

### Performance Issues
- Use `TranslateFunc()` cho repeated translations với cùng language
- Leverage localizer caching
- Consider pre-warming localizers cho common languages

### Module Conflicts
- Use descriptive module names
- Organize keys hierarchically
- Avoid key conflicts giữa modules
    Field1 string
    Field2 int64 // Changed from int
    Field3 bool  // New field
}
```

### Configuration Changes
If you're using configuration files:

```yaml
# Old configuration format
old_setting: value
deprecated_option: true

# New configuration format
new_setting: value
# deprecated_option removed
new_option: false
```

## Step-by-Step Migration

### Step 1: Update Dependencies
```bash
go get go.fork.vn/scheduler@v0.1.0
go mod tidy
```

### Step 2: Update Import Statements
```go
// If import paths changed
import (
    "go.fork.vn/scheduler" // Updated import
)
```

### Step 3: Update Code
Replace deprecated function calls:

```go
// Before
result := scheduler.OldFunction(param)

// After
result := scheduler.NewFunction(param, defaultValue)
```

### Step 4: Update Configuration
Update your configuration files according to the new schema.

### Step 5: Run Tests
```bash
go test ./...
```

## Common Issues and Solutions

### Issue 1: Function Not Found
**Problem**: `undefined: scheduler.OldFunction`  
**Solution**: Replace with `scheduler.NewFunction`

### Issue 2: Type Mismatch
**Problem**: `cannot use int as int64`  
**Solution**: Cast the value or update variable type

## Getting Help
- Check the [documentation](https://pkg.go.dev/go.fork.vn/scheduler@v0.1.0)
- Search [existing issues](https://github.com/go-fork/scheduler/issues)
- Create a [new issue](https://github.com/go-fork/scheduler/issues/new) if needed

## Rollback Instructions
If you need to rollback:

```bash
go get go.fork.vn/scheduler@previous-version
go mod tidy
```

Replace `previous-version` with your previous version tag.

---
**Need Help?** Feel free to open an issue or discussion on GitHub.
