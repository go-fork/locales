# Changelog

Tất cả các thay đổi đáng chú ý đối với dự án này sẽ được ghi lại trong file này.

Định dạng này dựa trên [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
và dự án này tuân theo [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-01-09

### Added

- Triển khai Manager interface với các phương thức dịch cốt lõi:
  - `Translate()` cho việc dịch trực tiếp với ngôn ngữ và key cụ thể
  - `TranslateFunc()` tạo function dịch cho ngôn ngữ cụ thể
  - `SetLocale()` và `Bundle()` cho quản lý cấu hình
- Hệ thống loading đa nguồn với 3 chiến lược chính:
  - File System loading từ thư mục chứa file JSON/YAML
  - Embedded FS loading cho việc đóng gói translation vào binary
  - Runtime Registration cho plugin systems và dynamic modules
- ModuleLoader interface hỗ trợ:
  - `RegisterMessage()` đăng ký một bản dịch cho module
  - `RegisterMessages()` đăng ký nhiều bản dịch cho module
  - `LoadJSON()` và `LoadJSONDirectory()` load từ file JSON
  - `LoadYAML()` và `LoadYAMLDirectory()` load từ file YAML
- Config struct cho quản lý cấu hình locales:
  - DefaultLocale, SupportedLocales configuration
  - LocalesDirectory và Bundle settings
  - EmbeddedFS và EmbeddedPath cho embedded loading
- Service Provider tích hợp với go.fork.vn/di framework:
  - Register() và Boot() lifecycle methods
  - Dependency injection support cho applications
- Các utility functions cho module registration:
  - `RegisterModule()` đăng ký module từ thư mục
  - `RegisterModuleFromFS()` đăng ký module từ embedded filesystem
  - `NewModuleLoader()` tạo loader cho module cụ thể
- Thread-safe concurrent operations với RWMutex optimization
- Localizer caching và lazy loading cho performance tối ưu
- Template parameter substitution với Go template syntax
- Graceful error handling và fallback strategies khi thiếu translation
- Hỗ trợ multiple file formats: JSON và YAML
- Comprehensive test suite với 96.2% code coverage
- Extensive documentation trong thư mục docs/
- Mock implementations cho testing
- Sample configurations và build scripts
- Release templates và migration guides

[Unreleased]: https://github.com/go-fork/locales/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/go-fork/locales/releases/tag/v0.1.0