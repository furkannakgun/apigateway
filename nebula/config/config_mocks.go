package config

// import (
// 	"io/fs"

// 	"github.com/stretchr/testify/mock"
// )

// // Create a mock configProvider type and add all the required methods to it
// type mockConfigProvider struct {
// 	mock.Mock
// }

// func (m *mockConfigProvider) SetConfigType(string) {
// 	m.Called()
// }
// func (m *mockConfigProvider) SetConfigFile(string) {
// 	m.Called()
// }
// func (m *mockConfigProvider) ReadInConfig() error {
// 	m.Called()
// 	return nil
// }
// func (m *mockConfigProvider) AddRemoteProvider(string, string, string) error {
// 	m.Called()
// 	return nil
// }
// func (m *mockConfigProvider) ReadRemoteConfig() error {
// 	m.Called()
// 	return nil
// }
// func (m *mockConfigProvider) Get(string) interface{} {
// 	m.Called()
// 	return nil
// }

// // Create mock configLoader type and add all the required methods to it
// type mockConfigLoader struct {
// 	mock.Mock
// }

// func (m *mockConfigLoader) checkMandatoryConfig(configProvider) error {
// 	m.Called()
// 	return nil
// }
// func (m *mockConfigLoader) loadConfigLocal(configProvider) {
// 	m.Called()
// }
// func (m *mockConfigLoader) loadOpenAPISpec(openAPISpecLoader) {
// 	m.Called()
// }

// // Create mock openAPISpecLoader type and add all the required methods to it
// type mockOpenAPISpecLoader struct {
// 	mock.Mock
// }

// func (m *mockOpenAPISpecLoader) loadOpenAPISpecETCD() []byte {
// 	m.Called()
// 	return nil
// }
// func (m *mockOpenAPISpecLoader) loadOpenAPISpecConsul() []byte {
// 	m.Called()
// 	return nil
// }

// // Create mock for the 'os' package so we can detect calls to the os.WriteFile method
// type mockOs struct {
// 	mock.Mock
// }

// func (m *mockOs) WriteFile(name string, data []byte, perm fs.FileMode) error {
// 	m.Called()
// 	return nil
// }
