package server

// import (
// 	"context"

// 	"github.com/fsnotify/fsnotify"
// 	"github.com/gorilla/mux"
// 	capi "github.com/hashicorp/consul/api"
// 	"github.com/stretchr/testify/mock"
// 	clientv3 "go.etcd.io/etcd/client/v3"
// )

// // Mock for config.Load

// type configMock struct {
// 	mock.Mock
// }

// func (m *configMock) Load() {
// 	m.Called()
// }

// // Mock for router.Configure

// type routerMock struct {
// 	mock.Mock
// }

// func (m *routerMock) Configure(r *mux.Router) {
// 	m.Called()
// }

// type mockLocalWatcher struct {
// 	mock.Mock
// }

// func (m *mockLocalWatcher) OnConfigChange(fn func(fsnotify.Event)) {
// 	m.Called()
// }

// func (m *mockLocalWatcher) WatchConfig() {
// 	m.Called()
// }

// type mockEtcdWatcher struct {
// 	mock.Mock
// }

// func (m *mockEtcdWatcher) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
// 	m.Called()
// 	return nil
// }

// func (m *mockEtcdWatcher) Close() error {
// 	m.Called()
// 	return nil
// }

// type mockConsulWatcher struct {
// 	mock.Mock
// }

// func (m *mockConsulWatcher) WatchTree(ctx context.Context, path string) (<-chan capi.KVPairs, error) {
// 	m.Called()
// 	return nil, nil
// }

// // Mocks for all of the watchers

// type mockConfigWatcher struct {
// 	mock.Mock
// }

// func (m *mockConfigWatcher) watchConfigETCD(client etcdWatcher) {
// 	m.Called()
// }

// func (m *mockConfigWatcher) watchConfigConsul(client consulWatcher) {
// 	m.Called()
// }

// func (m *mockConfigWatcher) watchConfigLocal(provider localWatcher) {
// 	m.Called()
// }
