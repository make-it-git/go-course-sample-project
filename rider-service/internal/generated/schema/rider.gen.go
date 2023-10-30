// Package rider provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package rider

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// CreateOrder defines model for CreateOrder.
type CreateOrder struct {
	DropoffLocation Location `json:"dropoff_location"`
	IdempotencyKey  string   `json:"idempotency_key"`
	PickupLocation  Location `json:"pickup_location"`
}

// Error defines model for Error.
type Error struct {
	Message *string `json:"message,omitempty"`
}

// Location defines model for Location.
type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// Order defines model for Order.
type Order struct {
	// CompletedAt The date and time the ride order was completed
	CompletedAt *openapi_types.Date `json:"completed_at,omitempty"`

	// CreatedAt The date and time the ride order was created
	CreatedAt       openapi_types.Date `json:"created_at"`
	DropoffLocation Location           `json:"dropoff_location"`

	// Id The ID of the ride order
	Id             string   `json:"id"`
	PickupLocation Location `json:"pickup_location"`

	// TotalPrice The total price of the ride order
	TotalPrice int `json:"total_price"`
}

// XUserID defines model for X-User-ID.
type XUserID = int

// GetOrdersParams defines parameters for GetOrders.
type GetOrdersParams struct {
	XUserID XUserID `json:"X-User-ID"`
}

// PostOrdersParams defines parameters for PostOrders.
type PostOrdersParams struct {
	XUserID XUserID `json:"X-User-ID"`
}

// PostOrdersJSONRequestBody defines body for PostOrders for application/json ContentType.
type PostOrdersJSONRequestBody = CreateOrder

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get a list of ride orders
	// (GET /orders)
	GetOrders(w http.ResponseWriter, r *http.Request, params GetOrdersParams)
	// Create a ride order
	// (POST /orders)
	PostOrders(w http.ResponseWriter, r *http.Request, params PostOrdersParams)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Get a list of ride orders
// (GET /orders)
func (_ Unimplemented) GetOrders(w http.ResponseWriter, r *http.Request, params GetOrdersParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Create a ride order
// (POST /orders)
func (_ Unimplemented) PostOrders(w http.ResponseWriter, r *http.Request, params PostOrdersParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetOrders operation middleware
func (siw *ServerInterfaceWrapper) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetOrdersParams

	headers := r.Header

	// ------------- Required header parameter "X-User-ID" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-User-ID")]; found {
		var XUserID XUserID
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "X-User-ID", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-User-ID", runtime.ParamLocationHeader, valueList[0], &XUserID)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X-User-ID", Err: err})
			return
		}

		params.XUserID = XUserID

	} else {
		err := fmt.Errorf("Header parameter X-User-ID is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "X-User-ID", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetOrders(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostOrders operation middleware
func (siw *ServerInterfaceWrapper) PostOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostOrdersParams

	headers := r.Header

	// ------------- Required header parameter "X-User-ID" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-User-ID")]; found {
		var XUserID XUserID
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "X-User-ID", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-User-ID", runtime.ParamLocationHeader, valueList[0], &XUserID)
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X-User-ID", Err: err})
			return
		}

		params.XUserID = XUserID

	} else {
		err := fmt.Errorf("Header parameter X-User-ID is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "X-User-ID", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostOrders(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/orders", wrapper.GetOrders)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/orders", wrapper.PostOrders)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7xV3U7bTBB9ldV836VJQmlvfEdLVUVCAqFWqoRQtNjjZMH7w+yYNkJ+92p3QxwnJuVP",
	"vYvj2TNnzjmzfoDCamcNGvaQP4CTJDUyUnz6efDDIx1MT8KDMpDDAmWJBBkYqRHyjYoMCO8aRVhCztRg",
	"Br5YoJbhqFZG6UZDfpgBL104qAzjHAnatn2sjC2/EErGMwpdAh+yDokVxpclWWeralbbQrKyJvz3P2EF",
	"Ofw37gYZr/DGp491bQaqRO0soymWs1tchqMrKp5JmXmocaq4bdwr4NvN6S93eu0iZ7uzXK21sdc3WHAg",
	"9JXIDuig0Xs5x4EZ2gGQ041x+ji1ZMVNuQlkGn0dbMmgtmb+1Nutedc4m6eG5nnC16BtjYzlTHL0GX1B",
	"yiXS8H2BopSMQppSsNIoeIGCVInCBjzxS3qxhoAMKks6AEE4Bdmuz0UM2Vu6JYDn9HpbZof5TU+ErbZ4",
	"DfV+dZ4zYMuynjlSBQ5ziAUiFuwj0+359oo8byt6ZvVp7QYsNFGmsruMj404Pp+KylIiqMxcsPytfMBU",
	"XAeUi0A/JjTUQgb3SD4dPxxNRpMgi3VopFOQw9FoMjoKM0hexBCPI3D8OccYrBDxOMW0hBy+IZ+liqx3",
	"zV4Om9KVjLtLtr0KKnpnjU+b82EySQtkGE1sKp2rVRJvfOOT6d1NrBi1/1sO0pZ2V4kkksukbl/VU+U5",
	"mN8Z76PRvtFa0jINLaSoh+oycNYP6HRu/XsJddeg58+2XL5Io33SbH6e2jaFumfH4bu1WjfZVv2iu4tW",
	"uyF8UxTofdXU9TIo+/GFsdjHI32D9vNYSS0w1Wbw6V8QmBpGMrIWHuke6bF7L4DJMCE3L6dg3J8AAAD/",
	"/6OHTQT5CAAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}