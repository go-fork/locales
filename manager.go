package locales

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// Manager định nghĩa interface cho việc quản lý bản dịch i18n.
//
// Manager cung cấp các phương thức để làm việc với bản dịch,
// bao gồm việc cấu hình, tải bản dịch từ nhiều nguồn, và phản hồi
// các yêu cầu bản dịch từ ứng dụng.
type Manager interface {
	// Translate trả về bản dịch của khóa message trong ngôn ngữ đã chỉ định.
	//
	// Tham số:
	//   - lang: Mã ngôn ngữ.
	//   - key: Khóa của message cần dịch.
	//   - args: Các tham số thay thế trong message.
	//
	// Trả về:
	//   - string: Bản dịch của message.
	Translate(lang, key string, args map[string]interface{}) string

	// TranslateFunc trả về hàm dịch cho ngôn ngữ cụ thể.
	//
	// Tham số:
	//   - lang: Mã ngôn ngữ.
	//
	// Trả về:
	//   - func(string, map[string]interface{}) string: Hàm dịch cho ngôn ngữ đã chỉ định.
	TranslateFunc(lang string) func(key string, args map[string]interface{}) string

	// LoadTranslations tải bản dịch từ thư mục chỉ định.
	//
	// Tham số:
	//   - directory: Đường dẫn đến thư mục chứa file bản dịch.
	//   - loaderFunc: Hàm tùy chỉnh để tải file bản dịch. Nếu nil, sử dụng os.ReadFile.
	//
	// Trả về:
	//   - error: Lỗi nếu có khi tải bản dịch.
	LoadTranslations(directory string, loaderFunc func(string) ([]byte, error)) error

	// LoadEmbeddedTranslations tải bản dịch từ filesystem nhúng.
	//
	// Tham số:
	//   - embeddedFS: Filesystem nhúng chứa file bản dịch.
	//   - directory: Đường dẫn trong filesystem nhúng.
	//
	// Trả về:
	//   - error: Lỗi nếu có khi tải bản dịch.
	LoadEmbeddedTranslations(embeddedFS embed.FS, directory string) error

	// RegisterTranslations đăng ký bản dịch cho một module cụ thể.
	//
	// Tham số:
	//   - module: Tên của module.
	//   - translations: Map các bản dịch theo ngôn ngữ.
	RegisterTranslations(module string, translations map[string]map[string]string)

	// SetLocale thiết lập ngôn ngữ mặc định.
	//
	// Tham số:
	//   - locale: Mã ngôn ngữ cần thiết lập.
	SetLocale(locale string)

	// Bundle trả về đối tượng i18n.Bundle đang sử dụng.
	//
	// Trả về:
	//   - *i18n.Bundle: Đối tượng bundle i18n.
	Bundle() *i18n.Bundle
}

// manager là implementation của Manager.
//
// manager quản lý việc tải và truy xuất bản dịch cho ứng dụng.
type manager struct {
	// config chứa cấu hình cho manager.
	config *Config

	// bundle là đối tượng i18n.Bundle chính để quản lý bản dịch.
	bundle *i18n.Bundle

	// localizer là map lưu trữ các đối tượng i18n.Localizer cho mỗi ngôn ngữ.
	localizer map[string]*i18n.Localizer

	// mutex bảo vệ truy cập đồng thời vào các biến chia sẻ.
	mutex sync.RWMutex

	// messageMap lưu trữ bản dịch theo module và ngôn ngữ.
	messageMap map[string]map[string]string // For storing module-specific messages
}

// NewManager tạo một đối tượng manager mới với cấu hình cho trước.
//
// Tham số:
//   - config: Cấu hình cho manager.
//
// Trả về:
//   - Manager: Đối tượng Manager đã được khởi tạo.
func NewManager(config *Config) Manager {
	// Handle nil config
	if config == nil {
		config = DefaultConfig()
	}

	m := &manager{
		config:     config,
		localizer:  make(map[string]*i18n.Localizer),
		messageMap: make(map[string]map[string]string),
	}

	// Khởi tạo bundle với ngôn ngữ mặc định
	m.bundle = i18n.NewBundle(language.Make(config.DefaultLocale))
	m.bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	m.bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	m.bundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)

	// Tải bản dịch từ cấu hình
	if config.EmbeddedFS != nil && config.EmbeddedPath != "" {
		_ = m.LoadEmbeddedTranslations(*config.EmbeddedFS, config.EmbeddedPath)
	} else if config.LocalesDirectory != "" {
		_ = m.LoadTranslations(config.LocalesDirectory, nil)
	}

	return m
}

// Translate trả về bản dịch của khóa message trong ngôn ngữ đã chỉ định.
//
// Tham số:
//   - lang: Mã ngôn ngữ.
//   - key: Khóa của message cần dịch.
//   - args: Các tham số thay thế trong message.
//
// Trả về:
//   - string: Bản dịch của message.
func (m *manager) Translate(lang, key string, args map[string]interface{}) string {
	localizeConfig := &i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: args,
	}

	// Kiểm tra trong module-specific messages trước
	m.mutex.RLock()
	if moduleMsgs, ok := m.messageMap[lang]; ok {
		if msg, ok := moduleMsgs[key]; ok {
			// Thay thế tham số đơn giản
			result := msg
			for k, val := range args {
				placeholder := "{{." + k + "}}"
				replacement := fmt.Sprintf("%v", val)
				result = strings.Replace(result, placeholder, replacement, -1)
			}
			m.mutex.RUnlock()
			return result
		}
	}
	m.mutex.RUnlock()

	// Sử dụng bundle làm fallback
	loc := m.getLocalizer(lang)
	result, err := loc.Localize(localizeConfig)
	if err != nil {
		// Nếu không tìm thấy, trả về key làm fallback
		return key
	}

	return result
}

// TranslateFunc trả về hàm dịch cho ngôn ngữ cụ thể.
//
// Tham số:
//   - lang: Mã ngôn ngữ.
//
// Trả về:
//   - func(string, map[string]interface{}) string: Hàm dịch cho ngôn ngữ đã chỉ định.
func (m *manager) TranslateFunc(lang string) func(key string, args map[string]interface{}) string {
	return func(key string, args map[string]interface{}) string {
		return m.Translate(lang, key, args)
	}
}

// LoadTranslations tải bản dịch từ thư mục chỉ định.
//
// Tham số:
//   - directory: Đường dẫn đến thư mục chứa file bản dịch.
//   - loaderFunc: Hàm tùy chỉnh để tải file bản dịch. Nếu nil, sử dụng os.ReadFile.
//
// Trả về:
//   - error: Lỗi nếu có khi tải bản dịch.
func (m *manager) LoadTranslations(directory string, loaderFunc func(string) ([]byte, error)) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".json" || filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			var data []byte
			if loaderFunc != nil {
				data, err = loaderFunc(path)
			} else {
				data, err = os.ReadFile(path)
			}

			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", path, err)
				return nil // Tiếp tục với file khác
			}

			_, err = m.bundle.ParseMessageFileBytes(data, path)
			if err != nil {
				fmt.Printf("Error parsing message file %s: %v\n", path, err)
				return nil // Tiếp tục với file khác
			}
		}
		return nil
	})
}

// LoadEmbeddedTranslations tải bản dịch từ filesystem nhúng.
//
// Tham số:
//   - embeddedFS: Filesystem nhúng chứa file bản dịch.
//   - directory: Đường dẫn trong filesystem nhúng.
//
// Trả về:
//   - error: Lỗi nếu có khi tải bản dịch.
func (m *manager) LoadEmbeddedTranslations(embeddedFS embed.FS, directory string) error {
	return fs.WalkDir(embeddedFS, directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".json" || filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			data, err := embeddedFS.ReadFile(path)
			if err != nil {
				fmt.Printf("Error reading embedded file %s: %v\n", path, err)
				return nil // Tiếp tục với file khác
			}

			_, err = m.bundle.ParseMessageFileBytes(data, path)
			if err != nil {
				fmt.Printf("Error parsing embedded message file %s: %v\n", path, err)
				return nil // Tiếp tục với file khác
			}
		}
		return nil
	})
}

// RegisterTranslations đăng ký bản dịch cho một module cụ thể.
//
// Tham số:
//   - module: Tên của module.
//   - translations: Map các bản dịch theo ngôn ngữ.
func (m *manager) RegisterTranslations(module string, translations map[string]map[string]string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for lang, messages := range translations {
		if _, ok := m.messageMap[lang]; !ok {
			m.messageMap[lang] = make(map[string]string)
		}

		for key, message := range messages {
			// Prefix với module name để tránh conflict
			moduleKey := fmt.Sprintf("%s.%s", module, key)
			m.messageMap[lang][moduleKey] = message
		}
	}
}

// SetLocale thiết lập ngôn ngữ mặc định.
//
// Tham số:
//   - locale: Mã ngôn ngữ cần thiết lập.
func (m *manager) SetLocale(locale string) {
	m.config.DefaultLocale = locale
}

// Bundle trả về đối tượng i18n.Bundle đang sử dụng.
//
// Trả về:
//   - *i18n.Bundle: Đối tượng bundle i18n.
func (m *manager) Bundle() *i18n.Bundle {
	return m.bundle
}

// getLocalizer trả về đối tượng localizer cho ngôn ngữ cụ thể.
// Nếu localizer cho ngôn ngữ được yêu cầu chưa tồn tại, một localizer mới
// sẽ được tạo và lưu vào cache cho các lần gọi sau.
//
// Tham số:
//   - lang: Mã ngôn ngữ cần lấy localizer.
//
// Trả về:
//   - *i18n.Localizer: Localizer cho ngôn ngữ đã chỉ định.
func (m *manager) getLocalizer(lang string) *i18n.Localizer {
	m.mutex.RLock()
	loc, ok := m.localizer[lang]
	m.mutex.RUnlock()

	if !ok {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		// Kiểm tra lại trong trường hợp đã được tạo trong lúc chờ write lock
		loc, ok = m.localizer[lang]
		if !ok {
			loc = i18n.NewLocalizer(m.bundle, lang)
			m.localizer[lang] = loc
		}
	}

	return loc
}
