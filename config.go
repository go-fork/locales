package locales

import (
	"embed"
)

// Config định nghĩa cấu hình cho service locales.
//
// Config chứa các thông tin cần thiết để khởi tạo và cấu hình
// loader và truy xuất message cho hệ thống đa ngôn ngữ (i18n).
// Việc xác định ngôn ngữ từ request thuộc về middleware.
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

// DefaultConfig trả về cấu hình mặc định cho service locales.
//
// Hàm này tạo và trả về một đối tượng Config với các giá trị mặc định
// cho việc load và quản lý message. Việc xác định ngôn ngữ từ request
// thuộc về middleware.
//
// Returns:
//   - *Config: Con trỏ tới đối tượng Config với các giá trị mặc định
func DefaultConfig() *Config {
	return &Config{
		DefaultLocale:    "en",
		LocalesDirectory: "locales",
		Bundle:           "message",
		SupportedLocales: []string{"en"},
	}
}
