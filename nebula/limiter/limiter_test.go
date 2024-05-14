package limiter

import (
	"github.com/getkin/kin-openapi/routers"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/config"
	"net/http"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	tests := []struct {
		name        string
		distributed bool
		expErr      error
	}{
		/*
			Unable to test errors for now as Fatalv exits the program before any can be returned
			When it can be tested add additional tests to check for errors
		*/
		{
			name:        "Rate Limits distributed",
			distributed: true,
			expErr:      nil,
		},
		{
			name:        "Rate Limits not distributed",
			distributed: false,
			expErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Conf = viper.New()
			config.Conf.Set("rateLimits", map[string]interface{}{
				"distributed": tt.distributed,
			})

			store, err := Store(123456789)

			// Check that a store object is returned containing some kind of config, be it Redis or another in-memory DB.
			assert.NotEmpty(t, store)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func createStore() *limiter.Store {
	// Create a basic store object
	store, _ := memorystore.New(&memorystore.Config{
		Tokens:   123456789,
		Interval: time.Second,
	})
	return &store
}

func TestMiddlewareLimiter(t *testing.T) {
	type args struct {
		limiterkey  string
		limiterType string
		s           *limiter.Store
		validator   routers.Router
	}
	tests := []struct {
		name      string
		args      args
		expResult *Middleware
		expError  error
	}{
		/*
			Unable to test errors for now as Fatalv exits the program before any can be returned
			When it can be tested add additional tests to check for errors
		*/
		{
			name: "Middleware Limiter created",
			args: args{
				limiterkey:  "userid",
				limiterType: "path",
				s:           createStore(),
				validator:   getTestValidator(),
			},
			expResult: &Middleware{},
			expError:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw, err := MiddlewareLimiter(tt.args.limiterkey, tt.args.limiterType, tt.args.s, tt.args.validator)
			assert.IsType(t, mw, tt.expResult)
			assert.Equal(t, tt.expError, err)
		})
	}
}

func TestSpikeArrestLimiter(t *testing.T) {
	type args struct {
		s              *limiter.Store
		spikeArrestVal routers.Router
	}
	tests := []struct {
		name      string
		args      args
		expResult *Middleware
		expError  error
	}{
		/*
			Unable to test errors for now as Fatalv exits the program before any can be returned
			When it can be tested add additional tests to check for errors
		*/
		{
			name: "SpikeArrest Limiter created",
			args: args{
				s:              createStore(),
				spikeArrestVal: getTestValidator(),
			},
			expResult: &Middleware{},
			expError:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw, err := SpikeArrestLimiter(tt.args.s, tt.args.spikeArrestVal)
			assert.IsType(t, mw, tt.expResult)
			assert.Equal(t, tt.expError, err)
		})
	}
}

func TestCreateKeyFunc(t *testing.T) {
	type args struct {
		limiterkey  string
		limitertype string
		validator   routers.Router
	}
	tests := []struct {
		name             string
		args             args
		limiterValHeader []string
		url              string
		expLimiterValue  string
		expError         error
	}{
		{
			name: "limiterValue created from header",
			args: args{
				limiterkey:  "Limit",
				limitertype: "header",
				validator:   getTestValidator(),
			},
			limiterValHeader: []string{"50"},
			url:              "/go-test-service/test/123456789",
			expLimiterValue:  "50",
			expError:         nil,
		},
		{
			name: "limiterValue created from path",
			args: args{
				limiterkey:  "limit",
				limitertype: "path",
				validator:   getTestValidator(),
			},
			limiterValHeader: []string{},
			url:              "/go-test-service/test/123456789/50",
			expLimiterValue:  "50",
			expError:         nil,
		},
		{
			name: "Invalid path",
			args: args{
				limiterkey:  "limit",
				limitertype: "path",
				validator:   getTestValidator(),
			},
			limiterValHeader: []string{},
			url:              "/go-test-service/fake/path/test/123456789/50",
			expLimiterValue:  "",
			expError:         &routers.RouteError{Reason: "no matching operation was found"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.url, nil)

			if tt.args.limitertype == "header" {
				req.Header[tt.args.limiterkey] = tt.limiterValHeader
			}

			keyfunc, _ := CreateKeyFunc(tt.args.limiterkey, tt.args.limitertype, tt.args.validator)
			result, err := keyfunc(req)

			assert.Equal(t, tt.expLimiterValue, result)
			assert.Equal(t, tt.expError, err)
		})
	}
}

func TestCreateKeyFuncSpikeArrest(t *testing.T) {
	tests := []struct {
		name                 string
		spikeArrestValidator routers.Router
		route                string
		expSpikeArrestKey    string
		expError             error
	}{
		{
			name:                 "KeyFunc Spike Arrest created OK",
			spikeArrestValidator: getTestValidator(),
			route:                "http://example.com/go-test-service/test/1234567890",
			expSpikeArrestKey:    "go-test-service",
			expError:             nil,
		},
		{
			name:                 "Invalid Route",
			spikeArrestValidator: getTestValidator(),
			route:                "http://example.com/invalid-route/test/1234567890",
			expSpikeArrestKey:    "",
			expError:             &routers.RouteError{Reason: "no matching operation was found"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.route, nil)
			keyfunc := CreateKeyFuncSpikeArrest(tt.spikeArrestValidator)
			result, err := keyfunc(req)

			assert.Equal(t, tt.expSpikeArrestKey, result)
			assert.Equal(t, tt.expError, err)

		})
	}
}
