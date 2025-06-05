# Release Notes - v0.1.0

## Overview
Phiên bản đầu tiên của Go Locales Service Provider - một thư viện quản lý đa ngôn ngữ (i18n) chuyên nghiệp cho ứng dụng Go với kiến trúc modular và hiệu suất cao.

## What's New
### 🚀 Features
- **Manager Interface**: Core translation management với `Translate()`, `TranslateFunc()`, `SetLocale()`, và `Bundle()`
- **Multiple Loading Strategies**: Hỗ trợ File System, Embedded FS, và Runtime Registration
- **ModuleLoader System**: Plugin-friendly architecture với `RegisterMessage()`, `RegisterMessages()`, `LoadJSON()`, `LoadYAML()`
- **Service Provider Integration**: Tích hợp với go.fork.vn/di framework thông qua Register() và Boot() lifecycle
- **Thread-Safe Operations**: Concurrent access với RWMutex optimization
- **Template Parameter Substitution**: Hỗ trợ Go template syntax với map[string]interface{} parameters
- **Multi-Format Support**: JSON và YAML file formats
- **Graceful Fallback**: Robust error handling và fallback strategies

### 🔧 Improvements
- **Localizer Caching**: Mỗi ngôn ngữ có cached localizer cho performance tối ưu
- **Lazy Loading**: Localizer chỉ được tạo khi cần thiết
- **Memory Efficient**: Chỉ cache những gì thực sự được sử dụng
- **Module-based Architecture**: Plugin system support cho microservices

### 📚 Documentation
- Complete documentation suite trong docs/ directory
- API reference cho Manager interface
- Configuration guide với examples
- Loader strategies documentation
- Best practices và usage patterns
- Sample configurations và templates

## Breaking Changes
### ⚠️ Important Notes
Đây là phiên bản đầu tiên (v0.1.0), không có breaking changes từ phiên bản trước.

## Migration Guide
See [MIGRATION.md](./MIGRATION.md) for detailed setup instructions.

## Dependencies
### Added
- github.com/nicksnyder/go-i18n/v2 v2.6.0: Core i18n functionality
- github.com/stretchr/testify v1.10.0: Testing framework
- go.fork.vn/config v0.1.3: Configuration management
- go.fork.vn/di v0.1.3: Dependency injection framework
- golang.org/x/text v0.25.0: Text processing utilities
- gopkg.in/yaml.v3 v3.0.1: YAML parsing support

## Performance
- **Thread-Safe**: RWMutex optimization cho concurrent operations
- **Caching Strategy**: Localizer caching giảm overhead khi translate
- **Memory Optimization**: Lazy loading và efficient memory usage
- **Loading Priority**: Runtime Registration > Bundle Translations > Fallback

## Testing
- **Comprehensive Test Suite**: 96.2% code coverage
- **Mock Implementations**: Complete mocks cho testing
- **Error Case Coverage**: Extensive error scenario testing
- **Concurrent Testing**: Thread safety validation
- **Integration Tests**: End-to-end testing scenarios

## Contributors
Thanks to all contributors who made this release possible:
- @zinzinday
- @nghiant0921

## Download
- Source code: [go.fork.vn/scheduler@v0.1.0]
- Documentation: [pkg.go.dev/go.fork.vn/scheduler@v0.1.0]

---
Release Date: 2025-06-05
