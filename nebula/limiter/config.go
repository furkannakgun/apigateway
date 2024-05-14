package limiter

import (
	"errors"

	"github.com/getkin/kin-openapi/routers"
	"github.com/gorilla/mux"

	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/config"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

func LimiterType(name string) (string, error) {
	limiter := config.Conf.Get("rateLimits").(map[string]interface{})["limiters"].(map[string]interface{})[name]
	if limiter == nil {
		err := errors.New("limiter name not found in configuration")
		return "", err
	}
	limiterType := limiter.(map[string]interface{})["type"]
	if limiterType == nil {
		err := errors.New("limiter does not have an assigned type")
		return "", err
	}

	return limiterType.(string), nil
}

func LimiterTPS(name string) (uint64, error) {
	limiter := config.Conf.Get("rateLimits").(map[string]interface{})["limiters"].(map[string]interface{})[name]
	if limiter == nil {
		err := errors.New("limiter name not found in configuration")
		return 0, err
	}
	tps := limiter.(map[string]interface{})["tps"]
	if tps == nil {
		err := errors.New("limiter TPS not defined")
		return 0, err
	}
	if limiterTPS, ok := tps.(int); ok {
		return uint64(limiterTPS), nil
	} else {
		return 0, errors.New("limiter TPS not of valid type")
	}
}

func SpikeArrestTPS() (uint64, error) {
	rateLimitsConfig, ok := config.Conf.Get("rateLimits").(map[string]interface{})
	if !ok || rateLimitsConfig == nil { // if rateLimits field is not found or null, default to 1000 TPS
		return 1000, nil
	}

	spikeArrest, ok := rateLimitsConfig["spikearrest"]
	if !ok || spikeArrest == nil { // spikearrest field not found or null, default to 1000 TPS
		return 1000, nil
	}

	spikeArrestVal, ok := spikeArrest.(int)
	if !ok { // if a spike arrest is found but there are errors parsing it, throw err instead of 1000 TPS
		return 0, errors.New("invalid type for spikearrest value in config file")
	}

	return uint64(spikeArrestVal), nil
}

func CreateSpikeArrest(validator routers.Router, r *mux.Router) error {
	logger := structlog.GetLogger("nebula")
	// Set rate limits for the entire API (Spike Arrest)
	spikeArrestMW, err := NewSpikeArrestLimiter(validator)
		if err != nil {
			logger.Errorv(err)
			return err
		}
	r.Use(spikeArrestMW.Handle)
	return nil
}

func CreateAndLoadLimiters(validator routers.Router, r *mux.Router) error {
	logger := structlog.GetLogger("nebula")
	var limiters interface{}
	rateLimitsConfig, ok := config.Conf.Get("rateLimits").(map[string]interface{})
	if ok && rateLimitsConfig != nil {
		limiters = config.Conf.Get("rateLimits").(map[string]interface{})["limiters"]
		// If limiters key does not exist
		if limiters == nil {
			return nil
		}
	} else { // If rateLimits key does not exist
		return nil
	}
	// Set rate limits for individual routes/requests
	for limiterName := range limiters.(map[string]interface{}) {
		middleware, err := NewMiddlewareLimiter(limiterName, validator)
		if err != nil {
			logger.Errorv(err)
			return err
		}
		r.Use(middleware.Handle)
	}
	return nil
}
