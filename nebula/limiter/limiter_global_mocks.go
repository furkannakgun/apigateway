package limiter

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/stretchr/testify/mock"
)

// Mock for limiter.Store interface so that it can be used for testing methods against different return values provided by the below methods

type mockStore interface {
	Take(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error)
	Get(ctx context.Context, key string) (tokens, remaining uint64, err error)
	Set(ctx context.Context, key string, tokens uint64, interval time.Duration) error
	Burst(ctx context.Context, key string, tokens uint64) error
	Close(ctx context.Context) error
	On(methodName string, arguments ...interface{}) *mock.Call
	AssertNumberOfCalls(t mock.TestingT, methodName string, expectedCalls int) bool
}

// Mock for all methods that return no errors
type mockStoreOk struct {
	mock.Mock
}

func (m *mockStoreOk) Take(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error) {
	return 123456789, 123456788, 987654321, true, nil
}

func (m *mockStoreOk) Get(ctx context.Context, key string) (tokens, remaining uint64, err error) {
	return 123456789, 123456788, nil
}

func (m *mockStoreOk) Set(ctx context.Context, key string, tokens uint64, interval time.Duration) error {
	return nil
}

func (m *mockStoreOk) Burst(ctx context.Context, key string, tokens uint64) error {
	return nil
}

func (m *mockStoreOk) Close(ctx context.Context) error {
	return nil
}

// Mock for when the Take() method returns an error
type mockStoreTakeErr struct {
	// reuse all the above methods from mockStoreOk unless overriden in this struct
	mockStoreOk
	mock.Mock
}

func (m *mockStoreTakeErr) Take(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error) {
	return 123456789, 123456788, 987654321, true, errors.New("error taking token from the provided key")
}

// Mock for when the Take() method has no tokens left
type mockStoreNoTokensLeft struct {
	// reuse all the above methods from mockStoreOk unless overriden in this struct
	mockStoreOk
	mock.Mock
}

func (m *mockStoreNoTokensLeft) Take(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error) {
	return 123456789, 123456788, 987654321, false, nil
}

// Retreive validator from test OpenAPISpec file
func getTestValidator() routers.Router {
	doc, _ := openapi3.NewLoader().LoadFromFile("../testdata/openAPIspec_test.yaml")
	validator, _ := gorillamux.NewRouter(doc)
	return validator
}

// Mocks to return a basic KeyFunc object for testing purposes

// Returns a KeyFunc OK
func mockKeyFuncOk() KeyFunc {
	return func(r *http.Request) (string, error) {
		return "127.0.0.1", nil
	}
}

// Returns an error retrieving IP
func mockKeyFuncErr() KeyFunc {
	return func(r *http.Request) (string, error) {
		return "", errors.New("error generating keyfunc")
	}
}

// Returns an empty string for an IP.
func mockKeyFuncEmptyStr() KeyFunc {
	return func(r *http.Request) (string, error) {
		return "", nil
	}
}
