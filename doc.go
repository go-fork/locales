// Package locales cung cấp service provider và các công cụ để hỗ trợ đa ngôn ngữ (i18n) trong ứng dụng Go.
//
// Package này được thiết kế để tích hợp với framework go.fork.vn/di thông qua service provider pattern,
// nhưng cũng có thể được sử dụng độc lập. Nó phụ thuộc vào thư viện go-i18n/v2 để quản lý bản dịch
// và cung cấp các API đơn giản để làm việc với bản dịch đa ngôn ngữ.
//
// # Tính năng
//
// Package locales cung cấp các tính năng sau:
//   - Tải bản dịch từ file JSON trong thư mục hoặc từ embed.FS
//   - Đăng ký và sử dụng bản dịch theo module
//   - Hỗ trợ thay thế biến trong chuỗi bản dịch
//   - Tự động phát hiện ngôn ngữ từ header Accept-Language
//   - Lưu trữ lựa chọn ngôn ngữ của người dùng trong cookie
//   - API để thêm, tải và truy xuất bản dịch
//
// # Sử dụng với Service Provider
//
// Để sử dụng locales với framework go.fork.vn/core, bạn cần đăng ký ServiceProvider vào container:
//
//	app := core.NewApplication()
//	app.Register(config.NewServiceProvider())
//	app.Register(locales.NewServiceProvider())
//	app.Boot()
//
// # Sử dụng trực tiếp
//
// Bạn cũng có thể sử dụng Manager trực tiếp mà không cần ServiceProvider:
//
//	config := locales.DefaultConfig()
//	config.DefaultLocale = "vi"
//	config.LocalesDirectory = "locales"
//
//	manager := locales.NewManager(config)
//	message := manager.Translate("fr", "welcome", map[string]interface{}{"name": "John"})
//
// # Đăng ký bản dịch cho module
//
// Để đăng ký bản dịch cho một module cụ thể:
//
//	// Đăng ký từ thư mục
//	locales.RegisterModule(manager, "mymodule", "path/to/mymodule/locales")
//
//	// Hoặc đăng ký từ embed.FS
//	//go:embed locales/*.json
//	var localesFS embed.FS
//	locales.RegisterModuleFromFS(manager, "mymodule", localesFS, "locales")
package locales
