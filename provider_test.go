package locales

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.fork.vn/config/mocks"
	di_mocks "go.fork.vn/di/mocks"
)

func TestNewServiceProvider(t *testing.T) {
	t.Run("should create new service provider", func(t *testing.T) {
		provider := NewServiceProvider()

		assert.NotNil(t, provider)
		assert.Implements(t, (*ServiceProvider)(nil), provider)
	})
}

func TestServiceProvider_Register(t *testing.T) {
	t.Run("should register locales service successfully", func(t *testing.T) {
		// Create mocks
		mockApp := di_mocks.NewMockApplication(t)
		mockContainer := di_mocks.NewMockContainer(t)
		mockConfig := mocks.NewMockManager(t)

		// Setup expectations
		mockApp.EXPECT().Container().Return(mockContainer)
		mockApp.EXPECT().MustMake("config").Return(mockConfig)

		// Mock config unmarshal
		mockConfig.EXPECT().UnmarshalKey("locales", mock.AnythingOfType("*locales.Config")).Return(nil).Run(func(key string, target interface{}) {
			// Simulate successful config unmarshaling
			if config, ok := target.(*Config); ok {
				config.DefaultLocale = "en"
				config.SupportedLocales = []string{"en", "vi"}
			}
		})

		// Mock app instance registration
		mockApp.EXPECT().Instance("locales", mock.Anything).Return()
		mockApp.EXPECT().Alias("i18n", "locales").Return()

		// Test
		provider := NewServiceProvider()

		// This should not panic
		provider.Register(mockApp)

		// Verify expectations
		mockApp.AssertExpectations(t)
		mockConfig.AssertExpectations(t)
	})

	t.Run("should handle nil container", func(t *testing.T) {
		mockApp := di_mocks.NewMockApplication(t)
		mockApp.EXPECT().Container().Return(nil)

		provider := NewServiceProvider()

		// This should not panic and should return early
		provider.Register(mockApp)

		mockApp.AssertExpectations(t)
	})

	t.Run("should panic when config service not found", func(t *testing.T) {
		mockApp := di_mocks.NewMockApplication(t)
		mockContainer := di_mocks.NewMockContainer(t)

		mockApp.EXPECT().Container().Return(mockContainer)
		mockApp.EXPECT().MustMake("config").Return("not-a-config-manager")

		provider := NewServiceProvider()

		assert.Panics(t, func() {
			provider.Register(mockApp)
		})
	})

	t.Run("should handle config unmarshal error gracefully", func(t *testing.T) {
		mockApp := di_mocks.NewMockApplication(t)
		mockContainer := di_mocks.NewMockContainer(t)
		mockConfig := mocks.NewMockManager(t)

		mockApp.EXPECT().Container().Return(mockContainer)
		mockApp.EXPECT().MustMake("config").Return(mockConfig)

		// Mock config unmarshal with error
		mockConfig.EXPECT().UnmarshalKey("locales", mock.AnythingOfType("*locales.Config")).Return(assert.AnError)

		// Mock app instance registration (should still proceed)
		mockApp.EXPECT().Instance("locales", mock.Anything).Return()
		mockApp.EXPECT().Alias("i18n", "locales").Return()

		provider := NewServiceProvider()

		// This should not panic even with unmarshal error
		provider.Register(mockApp)

		mockApp.AssertExpectations(t)
		mockConfig.AssertExpectations(t)
	})
}

func TestServiceProvider_Boot(t *testing.T) {
	t.Run("should boot successfully with valid app", func(t *testing.T) {
		mockApp := di_mocks.NewMockApplication(t)

		provider := NewServiceProvider()

		// This should not panic
		provider.Boot(mockApp)
	})

	t.Run("should panic when app is nil", func(t *testing.T) {
		provider := NewServiceProvider()

		assert.Panics(t, func() {
			provider.Boot(nil)
		})
	})
}

func TestServiceProvider_Providers(t *testing.T) {
	t.Run("should return empty providers initially", func(t *testing.T) {
		provider := NewServiceProvider()

		providers := provider.Providers()
		assert.Empty(t, providers)
	})

	t.Run("should return providers after registration", func(t *testing.T) {
		// Create mocks
		mockApp := di_mocks.NewMockApplication(t)
		mockContainer := di_mocks.NewMockContainer(t)
		mockConfig := mocks.NewMockManager(t)

		// Setup expectations
		mockApp.EXPECT().Container().Return(mockContainer)
		mockApp.EXPECT().MustMake("config").Return(mockConfig)
		mockConfig.EXPECT().UnmarshalKey("locales", mock.AnythingOfType("*locales.Config")).Return(nil)
		mockApp.EXPECT().Instance("locales", mock.Anything).Return()
		mockApp.EXPECT().Alias("i18n", "locales").Return()

		provider := NewServiceProvider()
		provider.Register(mockApp)

		providers := provider.Providers()
		assert.Len(t, providers, 2)
		assert.Contains(t, providers, "locales")
		assert.Contains(t, providers, "i18n")
	})
}

func TestServiceProvider_Requires(t *testing.T) {
	t.Run("should return config as required dependency", func(t *testing.T) {
		provider := NewServiceProvider()

		requires := provider.Requires()
		assert.Len(t, requires, 1)
		assert.Contains(t, requires, "config")
	})
}

func TestServiceProvider_Interface(t *testing.T) {
	t.Run("should implement ServiceProvider interface", func(t *testing.T) {
		provider := NewServiceProvider()

		assert.Implements(t, (*ServiceProvider)(nil), provider)
	})

	t.Run("should implement di.ServiceProvider interface", func(t *testing.T) {
		provider := NewServiceProvider()

		// This should compile without issues since ServiceProvider extends di.ServiceProvider
		var _ ServiceProvider = provider
	})
}

func TestServiceProvider_Integration(t *testing.T) {
	t.Run("should work with real config data", func(t *testing.T) {
		// Create mocks
		mockApp := di_mocks.NewMockApplication(t)
		mockContainer := di_mocks.NewMockContainer(t)
		mockConfig := mocks.NewMockManager(t)

		// Setup expectations with realistic config data
		mockApp.EXPECT().Container().Return(mockContainer)
		mockApp.EXPECT().MustMake("config").Return(mockConfig)

		mockConfig.EXPECT().UnmarshalKey("locales", mock.AnythingOfType("*locales.Config")).Return(nil).Run(func(key string, target interface{}) {
			// Simulate realistic config
			if config, ok := target.(*Config); ok {
				config.DefaultLocale = "vi"
				config.LocalesDirectory = "resources/lang"
				config.SupportedLocales = []string{"en", "vi", "fr"}
				config.Bundle = "messages"
			}
		})

		mockApp.EXPECT().Instance("locales", mock.Anything).Return()
		mockApp.EXPECT().Alias("i18n", "locales").Return()

		provider := NewServiceProvider()

		// Register
		provider.Register(mockApp)

		// Boot
		provider.Boot(mockApp)

		// Check providers
		providers := provider.Providers()
		assert.Contains(t, providers, "locales")
		assert.Contains(t, providers, "i18n")

		// Check requirements
		requires := provider.Requires()
		assert.Contains(t, requires, "config")

		mockApp.AssertExpectations(t)
		mockConfig.AssertExpectations(t)
	})
}

func TestServiceProvider_ErrorScenarios(t *testing.T) {
	t.Run("should handle config service returning wrong type", func(t *testing.T) {
		mockApp := di_mocks.NewMockApplication(t)
		mockContainer := di_mocks.NewMockContainer(t)

		mockApp.EXPECT().Container().Return(mockContainer)
		mockApp.EXPECT().MustMake("config").Return("wrong-type")

		provider := NewServiceProvider()

		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Should panic when config service returns wrong type")
	})

	t.Run("should handle nil config manager", func(t *testing.T) {
		mockApp := di_mocks.NewMockApplication(t)
		mockContainer := di_mocks.NewMockContainer(t)

		mockApp.EXPECT().Container().Return(mockContainer)
		mockApp.EXPECT().MustMake("config").Return(nil)

		provider := NewServiceProvider()

		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Should panic when config service returns nil")
	})
}
