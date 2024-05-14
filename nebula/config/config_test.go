package config

// import (
// 	"os"
// 	"testing"
// 	"github.com/spf13/viper"
// 	"github.com/stretchr/testify/assert"
// )

// func varCleanup() {
// 	// This function can be called to reset global and env variables after each test
// 	os.Unsetenv("ENV")
// 	os.Unsetenv("CONFIG_DB_URL")
// 	os.Unsetenv("CONFIG_DB_KEY")
// 	os.Unsetenv("CONFIG_DB_PROVIDER")
// 	OpenAPIPath = "/nebula-data/openAPIspec.yaml"
// 	ConfigDBProvider = ""
// }

// func TestCheckMandatoryConfig(t *testing.T) {
// 	type args struct {
// 		paramName  string
// 		paramValue string
// 	}
// 	tests := []struct {
// 		name            string
// 		args            args
// 		expErrorMessage string
// 	}{
// 		{
// 			name: "Params OK",
// 			args: args{
// 				paramName:  "targetURL",
// 				paramValue: "target.url",
// 			},
// 			expErrorMessage: "",
// 		},
// 		{
// 			name: "Missing param",
// 			args: args{
// 				paramName:  "missingParam",
// 				paramValue: "missing.param",
// 			},
// 			expErrorMessage: "Mandatory parameter missing: targetURL",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			Conf = viper.New()
// 			NebulaConfigLoader := cLoader{
// 				provider: Conf,
// 			}
// 			Conf.Set(tt.args.paramName, tt.args.paramValue)
// 			err := NebulaConfigLoader.checkMandatoryConfig(Conf)
// 			if tt.expErrorMessage == "" {
// 				assert.Equal(t, nil, err)
// 			} else {
// 				assert.EqualErrorf(t, err, tt.expErrorMessage, "Unexpected error message")
// 			}
// 		})
// 	}
// }

// func TestLoadConfigLocal(t *testing.T) {
// 	tests := []struct {
// 		name                  string
// 		expSetConfigFileCalls int
// 		expReadInConfigCalls  int
// 	}{
// 		{
// 			name:                  "Dev env, all OK.",
// 			expSetConfigFileCalls: 1,
// 			expReadInConfigCalls:  1,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			Conf = viper.New()
// 			NebulaConfigLoader := cLoader{
// 				provider: Conf,
// 			}
// 			mockProvider := &mockConfigProvider{}
// 			mockProvider.On("SetConfigFile").Return(nil)
// 			mockProvider.On("ReadInConfig").Return(nil)

// 			NebulaConfigLoader.loadConfigLocal(mockProvider)

// 			mockProvider.AssertNumberOfCalls(t, "SetConfigFile", tt.expSetConfigFileCalls)
// 			mockProvider.AssertNumberOfCalls(t, "ReadInConfig", tt.expReadInConfigCalls)
// 		})
// 	}
// }

// func TestLoadOpenAPISpec(t *testing.T) {
// 	tests := []struct {
// 		name             string
// 		configDbProvider string
// 		expETCDCalls     int
// 		expConsulCalls   int
// 		expWriteFileCalls int
// 	}{
// 		{
// 			name:             "ETCD test",
// 			configDbProvider: "etcd3",
// 			expETCDCalls:     1,
// 			expConsulCalls:   0,
// 			expWriteFileCalls: 1,
// 		},
// 		{
// 			name:             "Consul test",
// 			configDbProvider: "consul",
// 			expETCDCalls:     0,
// 			expConsulCalls:   1,
// 			expWriteFileCalls: 1,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ConfigDBProvider = tt.configDbProvider
// 			mockedOs := new(mockOs)
// 			osWriteFile = mockedOs.WriteFile
// 			Conf = viper.New()
// 			NebulaConfigLoader := cLoader{
// 				provider: Conf,
// 			}

// 			mockOpenAPISpecLoader := &mockOpenAPISpecLoader{}
// 			mockOpenAPISpecLoader.On("loadOpenAPISpecETCD").Return(nil)
// 			mockOpenAPISpecLoader.On("loadOpenAPISpecConsul").Return(nil)

// 			mockedOs.On("WriteFile").Return(nil)

// 			NebulaConfigLoader.loadOpenAPISpec(mockOpenAPISpecLoader)

// 			mockOpenAPISpecLoader.AssertNumberOfCalls(t, "loadOpenAPISpecETCD", tt.expETCDCalls)
// 			mockOpenAPISpecLoader.AssertNumberOfCalls(t, "loadOpenAPISpecConsul", tt.expConsulCalls)
// 			mockedOs.AssertNumberOfCalls(t, "WriteFile", tt.expWriteFileCalls)

// 			varCleanup()
// 		})
// 	}
// }

// func TestLoadEnvVars(t *testing.T) {

// 	//logger.InitLogger()

// 	type args struct {
// 		env              string
// 		configDBUrl      string
// 		configDBKey      string
// 		configDBProvider string
// 	}
// 	tests := []struct {
// 		name            string
// 		args            args
// 		expErrorMessage string
// 		expDBProvider   string
// 		expOpenAPIPath  string
// 	}{
// 		{
// 			name: "Dev env",
// 			args: args{
// 				env:              "dev",
// 				configDBUrl:      "",
// 				configDBKey:      "",
// 				configDBProvider: "",
// 			},
// 			expErrorMessage: "",
// 			expDBProvider:   "",
// 			expOpenAPIPath:  "openAPIspec.yaml",
// 		},
// 		{
// 			name: "All params OK",
// 			args: args{
// 				env:              "staging",
// 				configDBUrl:      "db.url",
// 				configDBKey:      "db.key",
// 				configDBProvider: "etcd3",
// 			},
// 			expErrorMessage: "",
// 			expDBProvider:   "etcd3",
// 			expOpenAPIPath:  "/nebula-data/openAPIspec.yaml",
// 		},
// 		{
// 			name: "Invalid DB provider",
// 			args: args{
// 				env:              "staging",
// 				configDBUrl:      "db.url",
// 				configDBKey:      "db.key",
// 				configDBProvider: "invalid_db_provider",
// 			},
// 			expErrorMessage: "Config DB provider is not supported",
// 			expDBProvider:   "",
// 			expOpenAPIPath:  "/nebula-data/openAPIspec.yaml",
// 		},
// 		{
// 			name: "Missing DB provider defaults to etcd3",
// 			args: args{
// 				env:              "staging",
// 				configDBUrl:      "db.url",
// 				configDBKey:      "db.key",
// 				configDBProvider: "",
// 			},
// 			expErrorMessage: "",
// 			expDBProvider:   "etcd3",
// 			expOpenAPIPath:  "/nebula-data/openAPIspec.yaml",
// 		},
// 		{
// 			name: "Missing DB Url",
// 			args: args{
// 				env:              "staging",
// 				configDBUrl:      "",
// 				configDBKey:      "db.key",
// 				configDBProvider: "etcd3",
// 			},
// 			expErrorMessage: "Config DB not set",
// 			expDBProvider:   "",
// 			expOpenAPIPath:  "/nebula-data/openAPIspec.yaml",
// 		},
// 		{
// 			name: "Missing DB Key",
// 			args: args{
// 				env:              "staging",
// 				configDBUrl:      "db.url",
// 				configDBKey:      "",
// 				configDBProvider: "etcd3",
// 			},
// 			expErrorMessage: "Config DB microservice key not set",
// 			expDBProvider:   "",
// 			expOpenAPIPath:  "/nebula-data/openAPIspec.yaml",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			os.Setenv("ENV", tt.args.env)
// 			os.Setenv("CONFIG_DB_URL", tt.args.configDBUrl)
// 			os.Setenv("CONFIG_DB_KEY", tt.args.configDBKey)
// 			os.Setenv("CONFIG_DB_PROVIDER", tt.args.configDBProvider)
// 			if tt.expErrorMessage == "" {
// 				LoadEnvVars()
// 				assert.Equal(t, tt.args.configDBUrl, ConfigDBHost)
// 				assert.Equal(t, tt.args.configDBKey, ServiceConfigKey)
// 				assert.Equal(t, tt.expDBProvider, ConfigDBProvider)
// 				assert.Equal(t, tt.expOpenAPIPath, OpenAPIPath)
// 			} else {
// 				assert.PanicsWithValuef(t, tt.expErrorMessage, LoadEnvVars, "Unexpected panic message")
// 			}
// 			varCleanup()
// 		})
// 	}
// }

// func TestLoad(t *testing.T) {
// 	type args struct {
// 		env string
// 	}
// 	tests := []struct {
// 		name                 string
// 		args                 args
// 		expErrorMessage      string
// 		expLocalConfigCalls  int
// 		expNormalConfigCalls int
// 	}{
// 		{
// 			name: "Dev env, all OK.",
// 			args: args{
// 				env: "dev",
// 			},
// 			expErrorMessage:      "",
// 			expLocalConfigCalls:  1,
// 			expNormalConfigCalls: 0,
// 		},
// 		{
// 			name: "Staging env, all OK",
// 			args: args{
// 				env: "staging",
// 			},
// 			expErrorMessage:      "",
// 			expLocalConfigCalls:  0,
// 			expNormalConfigCalls: 1,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			os.Setenv("ENV", tt.args.env)
// 			mockProvider := &mockConfigProvider{}
// 			mockProvider.On("SetConfigType").Return(nil)
// 			mockProvider.On("SetConfigFile").Return(nil)
// 			mockProvider.On("AddRemoteProvider").Return(nil)
// 			mockProvider.On("ReadRemoteConfig").Return(nil)

// 			mockLoader := &mockConfigLoader{}
// 			mockLoader.On("checkMandatoryConfig").Return(nil)
// 			mockLoader.On("loadConfigLocal").Return(nil)
// 			mockLoader.On("loadOpenAPISpec").Return(nil)

// 			mockOpenAPISpecLoader := &mockOpenAPISpecLoader{}
// 			mockOpenAPISpecLoader.On("loadOpenAPISpecETCD").Return(nil)
// 			mockOpenAPISpecLoader.On("loadOpenAPISpecConsul").Return(nil)

// 			if tt.expErrorMessage == "" {
// 				load(mockProvider, mockLoader, mockOpenAPISpecLoader)
// 				mockProvider.AssertNumberOfCalls(t, "SetConfigType", tt.expNormalConfigCalls)
// 				mockProvider.AssertNumberOfCalls(t, "AddRemoteProvider", tt.expNormalConfigCalls)
// 				mockProvider.AssertNumberOfCalls(t, "ReadRemoteConfig", tt.expNormalConfigCalls)

// 				mockLoader.AssertNumberOfCalls(t, "checkMandatoryConfig", tt.expNormalConfigCalls)
// 				mockLoader.AssertNumberOfCalls(t, "loadConfigLocal", tt.expLocalConfigCalls)
// 				mockLoader.AssertNumberOfCalls(t, "loadOpenAPISpec", tt.expNormalConfigCalls)
// 			} else {
// 				assert.PanicsWithValuef(t, tt.expErrorMessage, LoadEnvVars, "Unexpected panic message")
// 			}
// 			varCleanup()
// 		})
// 	}
// }
