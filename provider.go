package locales

import (
	"fmt"

	"go.fork.vn/config"
	"go.fork.vn/di"
)

// ServiceProvider định nghĩa interface cho Locales service provider.
//
// ServiceProvider kế thừa từ di.ServiceProvider và định nghĩa
// các phương thức cần thiết cho một Locales service provider.
type ServiceProvider interface {
	di.ServiceProvider
}

// serviceProvider là implementation của ServiceProvider.
//
// serviceProvider chịu trách nhiệm đăng ký các dịch vụ Locales vào DI container
// và cung cấp i18n manager cho các module khác trong ứng dụng.
type serviceProvider struct {
	providers []string
}

// NewServiceProvider tạo một Locales service provider mới.
//
// Hàm này khởi tạo và trả về một đối tượng ServiceProvider để sử dụng với DI container.
// ServiceProvider cho phép tự động đăng ký và cấu hình locales manager cho ứng dụng.
//
// Returns:
//   - ServiceProvider: Interface ServiceProvider đã được implement bởi serviceProvider
//
// Example:
//
//	app.Register(locales.NewServiceProvider())
func NewServiceProvider() ServiceProvider {
	return &serviceProvider{}
}

// Register đăng ký các dịch vụ Locales với DI container.
//
// Phương thức này đăng ký Locales manager vào container DI của ứng dụng.
// Nó khởi tạo một Locales manager mới và đăng ký nó dưới key "locales".
//
// Params:
//   - app: di.Application interface của ứng dụng, cung cấp phương thức Container() để lấy container DI
func (p *serviceProvider) Register(app di.Application) {
	// Lấy container từ app
	container := app.Container()
	if container == nil {
		return // Không làm gì khi không có container
	}

	// Kiểm tra và lấy config manager từ container
	configManager, ok := app.MustMake("config").(config.Manager)
	if !ok {
		panic("Locales provider requires config service to be registered")
	}

	// Tạo cấu hình mặc định và đọc cấu hình từ config
	localesConfig := DefaultConfig()
	err := configManager.UnmarshalKey("locales", localesConfig)
	if err != nil {
		fmt.Printf("Warning: Failed to unmarshal locales config, using defaults: %v\n", err)
	}

	// Tạo và khởi tạo manager
	manager := NewManager(localesConfig)

	// Đăng ký manager vào container
	app.Instance("locales", manager)
	app.Alias("i18n", "locales")

	// Thêm vào danh sách providers
	p.providers = append(p.providers, "locales")
	p.providers = append(p.providers, "i18n")
}

// Boot khởi động Locales provider.
//
// Phương thức này khởi động Locales provider sau khi tất cả các service provider đã được đăng ký.
// Trong trường hợp này, không cần thực hiện thêm tác vụ nào trong Boot vì các cấu hình
// đã được xử lý trong Register.
//
// Params:
//   - app: di.Application interface của ứng dụng
func (p *serviceProvider) Boot(app di.Application) {
	// Không cần thực hiện thêm tác vụ nào trong Boot
	// vì cấu hình đã được xử lý trong Register
	if app == nil {
		panic("Boot: di.Application cannot be nil")
	}
}

// Providers trả về danh sách các service mà provider này đăng ký.
//
// Locales provider đăng ký các dịch vụ "locales" và "i18n".
//
// Returns:
//   - []string: Mảng chứa tên của các service được đăng ký
func (p *serviceProvider) Providers() []string {
	return p.providers
}

// Requires trả về danh sách các dependency mà Locales provider phụ thuộc.
//
// Locales provider phụ thuộc vào config provider để đọc cấu hình.
//
// Trả về:
//   - []string: danh sách các service provider khác mà provider này yêu cầu
func (p *serviceProvider) Requires() []string {
	return []string{
		"config",
	}
}
