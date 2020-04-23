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

const componentName = "resty-component"

type Resty interface {
	Post(context.Context, ...RestyOptions) (*http.Response, error)
	Put(context.Context, ...RestyOptions) (*http.Response, error)
	Get(context.Context, ...RestyOptions) (*http.Response, error)
	Delete(context.Context, ...RestyOptions) (*http.Response, error)
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

func (r Request) Post(ctx context.Context, opt ...RestyOptions) (*http.Response, error) {
	return r.doRequest(ctx, http.MethodPost, opt...)
}

func (r Request) Put(ctx context.Context, opt ...RestyOptions) (*http.Response, error) {
	return r.doRequest(ctx, http.MethodPut, opt...)
}

func (r Request) Get(ctx context.Context, opt ...RestyOptions) (*http.Response, error) {
	return r.doRequest(ctx, http.MethodGet, opt...)
}

func (r Request) Delete(ctx context.Context, opt ...RestyOptions) (*http.Response, error) {
	return r.doRequest(ctx, http.MethodDelete, opt...)
}

func (r Request) doRequest(ctx context.Context, method string, opt ...RestyOptions) (*http.Response, error) {
	r.logContext(ctx)

	var restyOpt options
	for i := range opt {
		opt[i](&restyOpt)
	}

	url := fmt.Sprintf("%s%s", r.BaseURL, restyOpt.path)

	httpRequest, err := http.NewRequest(method, url, bytes.NewBuffer(restyOpt.body))
	if err != nil {
		return nil, err
	}

	addQueryParamsToRequest(httpRequest, restyOpt.params)

	addHeadersToRequest(httpRequest, r.Headers)

	return r.Client.Do(httpRequest)
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
