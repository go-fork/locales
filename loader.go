package locales

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Loader là interface cho việc tải file bản dịch.
// Interface này định nghĩa phương thức để tải nội dung của một file bản dịch
// từ bất kỳ nguồn nào (filesystem, embed.FS, v.v.).
type Loader interface {
	// LoadMessage tải nội dung của một file bản dịch.
	//
	// Tham số:
	//   - filepath (string): Đường dẫn đến file bản dịch.
	//
	// Trả về:
	//   - []byte: Nội dung của file bản dịch.
	//   - error: Lỗi nếu có khi tải file.
	LoadMessage(filepath string) ([]byte, error)
}

// LoaderFunc là kiểu adapter cho phép sử dụng
// các hàm thông thường như message loaders.
// Kiểu này triển khai interface Loader.
type LoaderFunc func(string) ([]byte, error)

// LoadMessage gọi f(filepath).
// Phương thức này cho phép LoaderFunc được sử dụng như một Loader.
//
// Tham số:
//   - filepath (string): Đường dẫn đến file bản dịch.
//
// Trả về:
//   - []byte: Nội dung của file bản dịch.
//   - error: Lỗi nếu có khi tải file.
func (f LoaderFunc) LoadMessage(filepath string) ([]byte, error) {
	return f(filepath)
}

// EmbedLoader cung cấp hỗ trợ tải bản dịch từ filesystem nhúng.
// Cấu trúc này triển khai interface Loader và cho phép tải
// file bản dịch từ embed.FS.
type EmbedLoader struct {
	// FS là filesystem nhúng chứa các file bản dịch.
	FS embed.FS
}

// LoadMessage tải một file bản dịch từ filesystem nhúng.
//
// Tham số:
//   - filepath (string): Đường dẫn đến file bản dịch trong filesystem nhúng.
//
// Trả về:
//   - []byte: Nội dung của file bản dịch.
//   - error: Lỗi nếu có khi tải file.
func (l EmbedLoader) LoadMessage(filepath string) ([]byte, error) {
	return l.FS.ReadFile(filepath)
}

// FSLoader cung cấp hỗ trợ tải bản dịch từ fs.FS.
// Cấu trúc này triển khai interface Loader và cho phép tải
// file bản dịch từ bất kỳ triển khai fs.FS nào.
type FSLoader struct {
	// FS là filesystem chứa các file bản dịch.
	FS fs.FS
}

// LoadMessage tải một file bản dịch từ fs.FS.
//
// Tham số:
//   - filepath (string): Đường dẫn đến file bản dịch trong filesystem.
//
// Trả về:
//   - []byte: Nội dung của file bản dịch.
//   - error: Lỗi nếu có khi tải file.
func (l FSLoader) LoadMessage(filepath string) ([]byte, error) {
	return fs.ReadFile(l.FS, filepath)
}

// ModuleTranslations lưu trữ bản dịch cho một module cụ thể.
// Cấu trúc này chứa tên module và các bản dịch cho module đó
// theo từng ngôn ngữ.
type ModuleTranslations struct {
	// Module là tên của module.
	Module string

	// Translations là map từ mã ngôn ngữ đến map bản dịch thông báo.
	Translations map[string]map[string]string
}

// ModuleLoader cung cấp chức năng đăng ký bản dịch riêng cho module.
// Cấu trúc này cho phép các module đăng ký bản dịch riêng của mình với manager.
type ModuleLoader struct {
	// moduleName là tên của module.
	moduleName string

	// manager là đối tượng Manager để đăng ký bản dịch.
	manager Manager
}

// NewModuleLoader tạo một module loader mới.
// Hàm này tạo một đối tượng ModuleLoader để quản lý bản dịch
// cho một module cụ thể.
//
// Tham số:
//   - moduleName: Tên của module.
//   - manager: Manager interface để đăng ký bản dịch.
//
// Trả về:
//   - *ModuleLoader: Đối tượng ModuleLoader mới.
func NewModuleLoader(moduleName string, manager Manager) *ModuleLoader {
	return &ModuleLoader{
		moduleName: moduleName,
		manager:    manager,
	}
}

// RegisterMessage đăng ký một thông báo bản dịch cho một module.
//
// Tham số:
//   - lang: Mã ngôn ngữ của thông báo.
//   - messageID: ID của thông báo.
//   - translation: Bản dịch của thông báo.
func (ml *ModuleLoader) RegisterMessage(lang string, messageID string, translation string) {
	translations := make(map[string]map[string]string)
	if _, ok := translations[lang]; !ok {
		translations[lang] = make(map[string]string)
	}
	translations[lang][messageID] = translation
	ml.manager.RegisterTranslations(ml.moduleName, translations)
}

// RegisterMessages đăng ký nhiều thông báo bản dịch cho một module.
//
// Tham số:
//   - lang: Mã ngôn ngữ của các thông báo.
//   - messages: Map chứa các cặp ID và bản dịch của thông báo.
func (ml *ModuleLoader) RegisterMessages(lang string, messages map[string]string) {
	translations := make(map[string]map[string]string)
	translations[lang] = messages
	ml.manager.RegisterTranslations(ml.moduleName, translations)
}

// LoadJSON tải bản dịch từ một file JSON cho một module.
// Phương thức này đọc file JSON, xác định ngôn ngữ từ tên file,
// phân tích nội dung thành map bản dịch và đăng ký chúng cho module.
//
// Tham số:
//   - path (string): Đường dẫn đến file JSON bản dịch.
//
// Trả về:
//   - error: Lỗi nếu có khi tải hoặc phân tích file.
func (ml *ModuleLoader) LoadJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read translation file %s: %w", path, err)
	}

	// Try to determine language from filename
	lang := filepath.Base(path)
	lang = lang[:len(lang)-len(filepath.Ext(lang))]

	// Parse JSON into messages map
	var messages map[string]string
	if err := json.Unmarshal(data, &messages); err != nil {
		return fmt.Errorf("failed to parse translation file %s: %w", path, err)
	}

	// Register messages
	ml.RegisterMessages(lang, messages)
	return nil
}

// LoadJSONDirectory tải tất cả file JSON bản dịch từ một thư mục.
// Phương thức này quét tất cả file JSON trong thư mục và tải chúng
// sử dụng LoadJSON. Các lỗi khi tải file sẽ được ghi log nhưng
// không dừng quá trình tải.
//
// Tham số:
//   - dirPath (string): Đường dẫn đến thư mục chứa các file JSON bản dịch.
//
// Trả về:
//   - error: Lỗi nếu có khi duyệt thư mục.
func (ml *ModuleLoader) LoadJSONDirectory(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".json" {
			if err := ml.LoadJSON(path); err != nil {
				fmt.Printf("Warning: failed to load translation file %s: %v\n", path, err)
			}
		}
		return nil
	})
}

// LoadFromFS tải bản dịch từ một fs.FS (hữu ích cho embed.FS) cho một module.
// Phương thức này quét tất cả file JSON trong filesystem và tải chúng
// vào module. Các lỗi khi tải file sẽ được trả về.
//
// Tham số:
//   - fsys (fs.FS): Filesystem chứa các file bản dịch.
//   - dirPath (string): Đường dẫn trong filesystem đến thư mục chứa các file bản dịch.
//
// Trả về:
//   - error: Lỗi nếu có khi duyệt filesystem.
func (ml *ModuleLoader) LoadFromFS(fsys fs.FS, dirPath string) error {
	return fs.WalkDir(fsys, dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".json" {
			data, err := fs.ReadFile(fsys, path)
			if err != nil {
				return fmt.Errorf("failed to read translation file %s: %w", path, err)
			}

			// Try to determine language from filename
			lang := filepath.Base(path)
			lang = lang[:len(lang)-len(filepath.Ext(lang))]

			// Parse JSON into messages map
			var messages map[string]string
			if err := json.Unmarshal(data, &messages); err != nil {
				return fmt.Errorf("failed to parse translation file %s: %w", path, err)
			}

			// Register messages
			ml.RegisterMessages(lang, messages)
		}
		return nil
	})
}

// RegisterModule đăng ký bản dịch i18n cho module từ thư mục locales.
// Hàm này tạo một ModuleLoader mới với tên module được cung cấp và tải
// tất cả các file JSON bản dịch từ thư mục đã chỉ định.
//
// Tham số:
//   - manager (Manager): Manager interface để đăng ký bản dịch.
//   - name (string): Tên của module (được sử dụng làm namespace cho các khóa dịch).
//   - localeDir (string): Đường dẫn đến thư mục chứa các file JSON bản dịch.
//
// Trả về:
//   - error: Lỗi nếu có khi tải bản dịch.
//
// Ví dụ:
//
//	err := locales.RegisterModule(manager, "totp", "internal/modules/totp/locales")
func RegisterModule(manager Manager, name, localeDir string) error {
	loader := NewModuleLoader(name, manager)
	return loader.LoadJSONDirectory(localeDir)
}

// RegisterModuleFromFS đăng ký bản dịch i18n cho module từ filesystem nhúng.
// Hàm này tạo một ModuleLoader mới với tên module được cung cấp và tải
// tất cả các file JSON bản dịch từ filesystem nhúng.
//
// Tham số:
//   - manager (Manager): Manager interface để đăng ký bản dịch.
//   - name (string): Tên của module (được sử dụng làm namespace cho các khóa dịch).
//   - fsys (fs.FS): Filesystem chứa các file bản dịch.
//   - dirPath (string): Đường dẫn trong filesystem đến thư mục chứa các file bản dịch.
//
// Trả về:
//   - error: Lỗi nếu có khi tải bản dịch.
//
// Ví dụ:
//
//	//go:embed locales/*.json
//	var localesFS embed.FS
//	err := locales.RegisterModuleFromFS(manager, "totp", localesFS, "locales")
func RegisterModuleFromFS(manager Manager, name string, fsys fs.FS, dirPath string) error {
	loader := NewModuleLoader(name, manager)
	return loader.LoadFromFS(fsys, dirPath)
}
