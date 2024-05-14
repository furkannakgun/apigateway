package server

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/spf13/viper"
// 	"github.com/stretchr/testify/assert"
// 	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/config"
// 	clientv3 "go.etcd.io/etcd/client/v3"
// )

// func TestGetServer(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		expResult *http.Server
// 	}{
// 		{
// 			name:      "Test server generation. All OK.",
// 			expResult: &http.Server{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			server := getServer()
// 			assert.IsType(t, tt.expResult, server)
// 		})
// 	}
// }

// func TestReloadRouter(t *testing.T) {
// 	tests := []struct {
// 		name              string
// 		expLoadCalls      int
// 		expConfigureCalls int
// 	}{
// 		{
// 			name:              "Reload Router",
// 			expLoadCalls:      1,
// 			expConfigureCalls: 1,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			configMock := configMock{}
// 			routerMock := routerMock{}

// 			configMock.On("Load").Return(nil)
// 			routerMock.On("Configure").Return(nil)

// 			routerConfigurator = routerMock.Configure
// 			configLoader = configMock.Load
// 			reloadRouter()

// 			configMock.AssertNumberOfCalls(t, "Load", tt.expLoadCalls)
// 			routerMock.AssertNumberOfCalls(t, "Configure", tt.expConfigureCalls)
// 		})
// 	}
// }
// func TestWatchConfig(t *testing.T) {
// 	tests := []struct {
// 		name                      string
// 		env                       string
// 		configDBProvider          string
// 		expwatchConfigLocalCalls  int
// 		expWatchConfigETCDCalls   int
// 		expWatchConfigConsulCalls int
// 	}{
// 		{
// 			name:                      "With Local configuration",
// 			env:                       "dev",
// 			configDBProvider:          "",
// 			expwatchConfigLocalCalls:  1,
// 			expWatchConfigETCDCalls:   0,
// 			expWatchConfigConsulCalls: 0,
// 		},
// 		{
// 			name:                      "With ETCD configuration",
// 			env:                       "staging",
// 			configDBProvider:          "etcd3",
// 			expwatchConfigLocalCalls:  0,
// 			expWatchConfigETCDCalls:   1,
// 			expWatchConfigConsulCalls: 0,
// 		},
// 		{
// 			name:                      "With Consul configuration",
// 			env:                       "staging",
// 			configDBProvider:          "consul",
// 			expwatchConfigLocalCalls:  0,
// 			expWatchConfigETCDCalls:   0,
// 			expWatchConfigConsulCalls: 1,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Setenv("ENV", tt.env)
// 			config.ConfigDBProvider = tt.configDBProvider

// 			nebulaConfigWatcher := &mockConfigWatcher{}
// 			nebulaConfigWatcher.On("watchConfigETCD").Return(nil)
// 			nebulaConfigWatcher.On("watchConfigConsul").Return(nil)
// 			nebulaConfigWatcher.On("watchConfigLocal").Return(nil)

// 			watchConfig(nebulaConfigWatcher)

// 			nebulaConfigWatcher.AssertNumberOfCalls(t, "watchConfigLocal", tt.expwatchConfigLocalCalls)
// 			nebulaConfigWatcher.AssertNumberOfCalls(t, "watchConfigETCD", tt.expWatchConfigETCDCalls)
// 			nebulaConfigWatcher.AssertNumberOfCalls(t, "watchConfigConsul", tt.expWatchConfigConsulCalls)
// 		})
// 	}
// }

// func TestWatchConfigConsul(t *testing.T) {
// 	tests := []struct {
// 		name                    string
// 		expConfigWatchTreeCalls int
// 	}{
// 		{
// 			name:                    "Test Watch Consul Config",
// 			expConfigWatchTreeCalls: 1,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			nebulaConfigWatcher := &cWatcher{}
// 			mockWatcher := &mockConsulWatcher{}
// 			mockWatcher.On("WatchTree").Return(nil)

// 			nebulaConfigWatcher.watchConfigConsul(mockWatcher)
// 			mockWatcher.AssertNumberOfCalls(t, "WatchTree", tt.expConfigWatchTreeCalls)

// 		})
// 	}
// }

// func TestChannelIsClosedETCD(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		wrespCancelled bool
// 		expResponse    bool
// 	}{
// 		{
// 			name:           "ETCD Channel closed",
// 			wrespCancelled: true,
// 			expResponse:    true,
// 		},
// 		{
// 			name:           "Event Reception",
// 			wrespCancelled: false,
// 			expResponse:    false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockWatchChannel := clientv3.WatchResponse{
// 				Canceled: tt.wrespCancelled,
// 			}
// 			result := channelIsClosedETCD(mockWatchChannel)

// 			assert.Equal(t, tt.expResponse, result)
// 		})
// 	}
// }
// func TestWatchConfigETCD(t *testing.T) {
// 	tests := []struct {
// 		name                string
// 		expConfigWatchCalls int
// 		expConfigCloseCalls int
// 	}{
// 		{
// 			name:                "Test Watch ETCD Config",
// 			expConfigWatchCalls: 1,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			nebulaConfigWatcher := &cWatcher{}
// 			mockWatcher := &mockEtcdWatcher{}

// 			mockWatcher.On("Watch").Return(nil)

// 			nebulaConfigWatcher.watchConfigETCD(mockWatcher)

// 			mockWatcher.AssertNumberOfCalls(t, "Watch", tt.expConfigWatchCalls)
// 		})
// 	}
// }

// func TestWatchConfigLocal(t *testing.T) {
// 	tests := []struct {
// 		name                 string
// 		expConfigChangeCalls int
// 		expConfigWatchCalls  int
// 	}{
// 		{
// 			name:                 "Test Watch Local Config function",
// 			expConfigChangeCalls: 1,
// 			expConfigWatchCalls:  1,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			config.Conf = viper.New()

// 			nebulaConfigWatcher := &cWatcher{}

// 			mockWatcher := &mockLocalWatcher{}
// 			mockWatcher.On("OnConfigChange").Return(nil)
// 			mockWatcher.On("WatchConfig").Return(nil)

// 			nebulaConfigWatcher.watchConfigLocal(mockWatcher)

// 			mockWatcher.AssertNumberOfCalls(t, "OnConfigChange", tt.expConfigChangeCalls)
// 			mockWatcher.AssertNumberOfCalls(t, "WatchConfig", tt.expConfigWatchCalls)
// 		})
// 	}
// }
// func TestCorsMiddleware(t *testing.T) {
// 	tests := []struct {
// 		name string
// 	}{
// 		{
// 			name: "Create Cors middleware",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			request := httptest.NewRequest("GET", "/", nil)
// 			recorder := httptest.NewRecorder()

// 			handlerFunction := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

// 			corsMiddleware(handlerFunction).ServeHTTP(
// 				recorder,
// 				request,
// 			)
// 			assert.Equal(t, recorder.Header().Get("Access-Control-Allow-Origin"), "*")
// 			assert.Equal(t, recorder.Header().Get("Access-Control-Allow-Methods"), "GET, PUT, POST, DELETE, OPTIONS")
// 			assert.Equal(t, recorder.Header().Get("Access-Control-Allow-Headers"), "*")
// 			assert.Equal(t, recorder.Header().Get("Content-Type"), "application/json")
// 		})
// 	}
// }
