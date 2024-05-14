package limiter

import (
	"errors"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/config"
)

func TestLimiterType(t *testing.T) {
	tests := []struct {
		name           string
		limiterName    string
		rateLimits     map[string]interface{}
		expLimiterType string
		expErr         error
	}{
		{
			name:        "Limiter name found and has assigned type",
			limiterName: "userid",
			rateLimits: map[string]interface{}{
				"userid": map[string]interface{}{
					"type": "path",
				},
			},
			expLimiterType: "path",
			expErr:         nil,
		},
		{
			name:        "Limiter name found no assigned type",
			limiterName: "userid",
			rateLimits: map[string]interface{}{
				"userid": map[string]interface{}{},
			},
			expLimiterType: "",
			expErr:         errors.New("limiter does not have an assigned type"),
		},
		{
			name:           "Limiter name not found",
			limiterName:    "userid",
			rateLimits:     map[string]interface{}{},
			expLimiterType: "",
			expErr:         errors.New("limiter name not found in configuration"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Conf = viper.New()
			config.Conf.Set("rateLimits", map[string]interface{}{"limiters": tt.rateLimits})
			limiter, err := LimiterType(tt.limiterName)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expLimiterType, limiter)
		})
	}
}

func TestLimiterTPS(t *testing.T) {
	tests := []struct {
		name          string
		limiterName   string
		rateLimits    map[string]interface{}
		expLimiterTPS uint64
		expErr        error
	}{
		{
			name:        "Limiter name found and TPS Defined",
			limiterName: "userid",
			rateLimits: map[string]interface{}{
				"userid": map[string]interface{}{
					"tps": 50,
				},
			},
			expLimiterTPS: 50,
			expErr:        nil,
		},
		{
			name:        "Limiter name found no TPS defined",
			limiterName: "userid",
			rateLimits: map[string]interface{}{
				"userid": map[string]interface{}{},
			},
			expLimiterTPS: 0,
			expErr:        errors.New("limiter TPS not defined"),
		},
		{
			name:          "Limiter name not found",
			limiterName:   "userid",
			rateLimits:    map[string]interface{}{},
			expLimiterTPS: 0,
			expErr:        errors.New("limiter name not found in configuration"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Conf = viper.New()
			config.Conf.Set("rateLimits", map[string]interface{}{"limiters": tt.rateLimits})
			tps, err := LimiterTPS(tt.limiterName)
			assert.Equal(t, tt.expErr, err)
			assert.Equal(t, tt.expLimiterTPS, tps)
		})
	}
}

func TestSpikeArrestTPS(t *testing.T) {
	tests := []struct {
		name          string
		rateLimits    map[string]interface{}
		expLimiterTPS uint64
		expErr        error
	}{
		{
			name: "Get Spike Arrest Limits OK.",
			rateLimits: map[string]interface{}{
				"spikeArrest": 1000,
			},
			expLimiterTPS: 1000,
			expErr:        nil,
		},
		{
			name: "Get Spike Arrest Limits Invalid Type.",
			rateLimits: map[string]interface{}{
				"spikeArrest": "1000",
			},
			expLimiterTPS: 0,
			expErr:        errors.New("invalid type for spikearrest value in config file"),
		},
		{
			name:          "No Spike arrest specified in config.",
			rateLimits:    map[string]interface{}{},
			expLimiterTPS: 1000,
			expErr:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Conf = viper.New()
			config.Conf.Set("rateLimits", tt.rateLimits)
			tps, err := SpikeArrestTPS()
			assert.Equal(t, tt.expLimiterTPS, tps)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func TestCreateSpikeArrest(t *testing.T) {
	tests := []struct {
		name          string
		rateLimitsKey string
		rateLimitsVal map[string]interface{}
		expErr        error
	}{
		{
			name:          "Test create Spike Arrest OK",
			rateLimitsKey: "rateLimits",
			rateLimitsVal: map[string]interface{}{
				"distributed": false,
				"spikeArrest": 1000,
			},
			expErr: nil,
		},
		{
			name:          "rateLimits key exists but val is null",
			rateLimitsKey: "rateLimits",
			rateLimitsVal: nil,
			expErr:        nil,
		},
		{
			name:          "RateLimits key does not exist",
			rateLimitsKey: "",
			rateLimitsVal: nil,
			expErr:        nil,
		},
		{
			name:          "SpikeArrest Value doesn't exist",
			rateLimitsKey: "rateLimits",
			rateLimitsVal: map[string]interface{}{
				"distributed": false,
				// Spike arrest defaults to 1000 if one isn't specified
			},
			expErr: nil,
		},
		{
			name:          "Invalid spike arrest type",
			rateLimitsKey: "rateLimits",
			rateLimitsVal: map[string]interface{}{
				"distributed": false,
				"spikeArrest": "1000",
			},
			expErr: errors.New("invalid type for spikearrest value in config file"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Conf = viper.New()
			if tt.rateLimitsKey != "" { // only create a rateLimits key in config if one if provided in a test
				config.Conf.Set(tt.rateLimitsKey, tt.rateLimitsVal)
			}
			r := mux.NewRouter()
			doc, _ := openapi3.NewLoader().LoadFromFile("../testdata/openAPIspec_test.yaml")
			specValidatorRouter, _ := gorillamux.NewRouter(doc)
			err := CreateSpikeArrest(specValidatorRouter, r)
			assert.Equal(t, tt.expErr, err)

		})
	}
}

func TestCreateAndLoadLimiters(t *testing.T) {
	tests := []struct {
		name          string
		rateLimitsKey string
		rateLimitsVal map[string]interface{}
		expErr        error
	}{
		{
			name:          "Limiters key exists and MW Limiter is created OK",
			rateLimitsKey: "rateLimits",
			rateLimitsVal: map[string]interface{}{
				"distributed": false,
				"limiters": map[string]interface{}{
					"userid": map[string]interface{}{
						"type": "path",
						"tps":  50,
					},
					"limit": map[string]interface{}{
						"type": "header",
						"tps":  20,
					},
				},
			},
			expErr: nil,
		},
		{
			name:          "Limiters key does not exist",
			rateLimitsKey: "rateLimits",
			rateLimitsVal: map[string]interface{}{
				"distributed": false,
				// No limiters key to see here!
			},
			expErr: nil,
		},
		{
			name:          "rateLimits key exists but val is null",
			rateLimitsKey: "rateLimits",
			rateLimitsVal: nil,
			expErr:        nil,
		},
		{
			name:          "rateLimits key does not exist",
			rateLimitsKey: "",
			rateLimitsVal: nil,
			expErr:        nil,
		},
		{
			name:          "Limiters key exists missing type",
			rateLimitsKey: "rateLimits",
			rateLimitsVal: map[string]interface{}{
				"distributed": false,
				"limiters": map[string]interface{}{
					"userid": map[string]interface{}{
						"tps": 50,
					},
					"limit": map[string]interface{}{
						"tps": 20,
					},
				},
			},
			expErr: errors.New("limiter does not have an assigned type"),
		},
		{
			name:          "Limiters key exists incorrect TPS type",
			rateLimitsKey: "rateLimits",
			rateLimitsVal: map[string]interface{}{
				"distributed": false,
				"limiters": map[string]interface{}{
					"userid": map[string]interface{}{
						"type": "path",
						"tps":  "invalid_tps",
					},
					"limit": map[string]interface{}{
						"type": "header",
						"tps":  "invalid_tps",
					},
				},
			},
			expErr: errors.New("limiter TPS not of valid type"),
		},
		{
			name:          "Limiters key exists missing TPS",
			rateLimitsKey: "rateLimits",
			rateLimitsVal: map[string]interface{}{
				"distributed": false,
				"limiters": map[string]interface{}{
					"userid": map[string]interface{}{
						"type": "path",
					},
					"limit": map[string]interface{}{
						"type": "header",
					},
				},
			},
			expErr: errors.New("limiter TPS not defined"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Conf = viper.New()
			if tt.rateLimitsKey != "" { // only create a rateLimits key in config if one if provided in a test
				config.Conf.Set(tt.rateLimitsKey, tt.rateLimitsVal)
			}
			r := mux.NewRouter()
			doc, _ := openapi3.NewLoader().LoadFromFile("../testdata/openAPIspec_test.yaml")
			specValidatorRouter, _ := gorillamux.NewRouter(doc)

			err := CreateAndLoadLimiters(specValidatorRouter, r)
			assert.Equal(t, tt.expErr, err)
		})
	}
}
