package limiter

import (
	"net/http"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/routers"
	"github.com/gomodule/redigo/redis"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
	"github.com/sethvargo/go-redisstore"

	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/nebula/config"
	"github.vodafone.com/VFGroup-NetworkArchitecture-NAAP/structlog"
)

// Memory store: creates a store in memory to save limiter rates
func Store(tokens uint64) (*limiter.Store, error) {
	logger := structlog.GetLogger("nebula")
	var store limiter.Store
	var err error
	rateLimits, ok := config.Conf.Get("rateLimits").(map[string]interface{})
	distributed, found := rateLimits["distributed"].(bool)
	storeDistributed := false
	if ok && rateLimits != nil && found{
		storeDistributed = distributed
	}
	if storeDistributed {
		host := config.Conf.Get("rateLimits").(map[string]interface{})["host"]
		passwd := config.Conf.Get("rateLimits").(map[string]interface{})["password"]
		store, err = redisstore.New(&redisstore.Config{
			Tokens:   tokens,
			Interval: time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", host.(string),
					redis.DialPassword(passwd.(string)))
			},
		})
	} else {
		store, err = memorystore.New(&memorystore.Config{
			Tokens:   tokens,
			Interval: time.Second,
		})
	}
	if err != nil {
		logger.Fatalv(err)
		return nil, err
	}

	return &store, nil
}

// Returns an HTTP middleware based on the the limiter value and the store
// TODO: limiterValue refers to a header now. It has to be broaden
func MiddlewareLimiter(limiterkey string, limiterType string, s *limiter.Store, validator routers.Router) (*Middleware, error) {
	logger := structlog.GetLogger("nebula")
	keyFunc, err := CreateKeyFunc(limiterkey, limiterType, validator)
	if err != nil {
		logger.Fatalv(err)
		return nil, err
	}

	middleware, err := NewMiddleware(*s, keyFunc)
	if err != nil {
		logger.Fatalv(err)
		return nil, err
	}

	return middleware, nil
}

func SpikeArrestLimiter(s *limiter.Store, validator routers.Router) (*Middleware, error) {
	logger := structlog.GetLogger("nebula")
	keyfunc := CreateKeyFuncSpikeArrest(validator)
	middleware, err := NewMiddleware(*s, keyfunc)
	if err != nil {
		logger.Fatalv(err)
		return nil, err
	}
	return middleware, nil
}

func CreateKeyFuncSpikeArrest(validator routers.Router) KeyFunc {
	keyFunc := KeyFunc(func(r *http.Request) (string, error) {
		route, _, err := validator.FindRoute(r)
		if err != nil {
			return "", err
		}
		// Use API name from the path as a key
		path_segments := strings.Split(route.Path, "/")
		spike_arrest_key := path_segments[1] // API name is assumed to be the first segment in a path
		return spike_arrest_key, nil
	})
	return keyFunc
}

// Obtain a key from the request to limit the request rate
// TODO: Broaden function to the use cases we are using in Apigee

func CreateKeyFunc(limiterKey string, limiterType string, validator routers.Router) (KeyFunc, error) {
	keyFunc := KeyFunc(func(r *http.Request) (string, error) {
		var limiterValue string
		if limiterType == "header" {
			limiterValue = r.Header.Get(limiterKey)
		} else if limiterType == "path" {
			_, pathParams, err := validator.FindRoute(r)
			if err != nil {
				return "", err
			}

			limiterValue = pathParams[limiterKey]
		}

		return limiterValue, nil
	})

	return keyFunc, nil
}

// Returns a Middleware struct to pass to mux router
// Args:
//     key: header used as key for the limiter (each value will have tps rate limit)
//     tps: Rate to limit the request from base on key

func NewMiddlewareLimiter(limiterKey string, validator routers.Router) (*Middleware, error) {
	logger := structlog.GetLogger("nebula")
	limiterType, err := LimiterType(limiterKey)
	if err != nil {
		logger.Errorv(err)
		return nil, err
	}

	tps, err := LimiterTPS(limiterKey)
	if err != nil {
		logger.Errorv(err)
		return nil, err
	}

	store, err := Store(tps)
	if err != nil {
		logger.Errorv(err)
		return nil, err
	}

	middleware, err := MiddlewareLimiter(limiterKey, limiterType, store, validator)
	if err != nil {
		logger.Errorv(err)
		return nil, err
	}

	return middleware, nil
}

func NewSpikeArrestLimiter(validator routers.Router) (*Middleware, error) {
	logger := structlog.GetLogger("nebula")

	spikeArrestTPS, err := SpikeArrestTPS()
	if err != nil {
		logger.Errorv(err)
		return nil, err
	}

	store, err := Store(spikeArrestTPS)
	if err != nil {
		logger.Errorv(err)
		return nil, err
	}

	middleware, err := SpikeArrestLimiter(store, validator)
	if err != nil {
		logger.Errorv(err)
		return nil, err
	}
	return middleware, nil
}
