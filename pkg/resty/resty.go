package resty

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zuttrax/rdlf-components/pkg/logs"
)

type Headers map[string]string
type Parameters map[string]string

const componentName = "resty-component"

type Resty interface {
	Post(context.Context, []byte) (*http.Response, error)
	Put(context.Context, []byte) (*http.Response, error)
	Get(context.Context) (*http.Response, error)
	Delete(context.Context) (*http.Response, error)
	AddParameterValueByKey(string, string)
}

type Response struct {
	Body       []byte
	StatusCode int
}

type Request struct {
	BaseURL string
	Trace   bool
	Client  http.Client
	Headers Headers
	Log     logs.Logger
}

func NewRequest(ttl int64, trace bool, baseURL string, head Headers) Request {
	return Request{
		BaseURL: baseURL,
		Client: http.Client{
			Timeout: time.Duration(ttl * time.Hour.Milliseconds()),
		},
		Trace:   trace,
		Headers: head,
		Log:     logs.InitializeLog(componentName, "debug"),
	}
}

type Endpoint struct {
	request    Request
	path       string
	parameters Parameters
}

func NewEndpoint(request Request, path string) Endpoint {
	return Endpoint{
		request: request,
		path:    path,
	}
}

func (e Endpoint) Post(ctx context.Context, msg []byte) (*http.Response, error) {
	return e.doRequest(ctx, msg, http.MethodPost)
}

func (e Endpoint) Put(ctx context.Context, msg []byte) (*http.Response, error) {
	return e.doRequest(ctx, msg, http.MethodPut)
}

func (e Endpoint) Get(ctx context.Context) (*http.Response, error) {
	return e.doRequest(ctx, nil, http.MethodGet)
}

func (e Endpoint) Delete(ctx context.Context) (*http.Response, error) {
	return e.doRequest(ctx, nil, http.MethodDelete)
}

func (e *Endpoint) AddParameterValueByKey(key string, value string) {
	parameters := make(map[string]string)
	for k, v := range e.parameters {
		parameters[k] = v
	}

	parameters[key] = value

	e.parameters = parameters
}

func (e Endpoint) doRequest(ctx context.Context, msg []byte, method string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", e.request.BaseURL, e.path)

	e.request.logContext(ctx)

	httpRequest, err := http.NewRequest(method, url, bytes.NewBuffer(msg))
	if err != nil {
		return nil, err
	}

	addQueryParamsToRequest(httpRequest, e.parameters)

	addHeadersToRequest(httpRequest, e.request.Headers)

	return e.request.Client.Do(httpRequest)
}

func addQueryParamsToRequest(httpRequest *http.Request, params Parameters) {
	queryParams := httpRequest.URL.Query()
	for key, value := range params {
		queryParams.Add(key, value)
	}
	httpRequest.URL.RawQuery = queryParams.Encode()
}

func addHeadersToRequest(httpRequest *http.Request, headers Headers) {
	for key, value := range headers {
		httpRequest.Header.Add(key, value)
	}
}

func (r Request) logContext(ctx context.Context) {
	if r.Trace {
		endpoint := ctx.Value("Endpoint").(string)
		token := ctx.Value("X-Token").(string)
		clientIP := ctx.Value("Client-IP").(string)

		r.Log.Info(fmt.Sprintf("client-ip: %s x-token: %s endpoint: %s", clientIP, token, endpoint))
	}
}
