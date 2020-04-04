package bifrost

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Brifrost interface {
	Post(*gin.Context, []byte) (*http.Response, error)
	Put(*gin.Context, []byte) (*http.Response, error)
	Get(*gin.Context) (*http.Response, error)
	Delete(*gin.Context) (*http.Response, error)
}

type Request struct {
	url     string
	ttl     int64
	headers Headers
}

type Headers map[string]string

type Response struct {
	Body       []byte
	StatusCode int
}

func NewRequest(ttl int64, url string, head Headers) Request {
	return Request{
		url:     url,
		ttl:     ttl,
		headers: head,
	}
}

func (req Request) Post(ctx *gin.Context, msg []byte) (*http.Response, error) {

	client := http.Client{
		Timeout: time.Duration(300 * time.Microsecond),
	}

	request, err := http.NewRequest(http.MethodPost, req.url, bytes.NewBuffer(msg))
	if err != nil {
		return nil, err
	}

	for key, val := range req.headers {
		request.Header.Add(key, val)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (req Request) Get(ctx *gin.Context) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Duration(300 * time.Microsecond),
	}

	request, err := http.NewRequest(http.MethodGet, req.url, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range req.headers {
		request.Header.Add(key, val)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (req Request) Put(ctx *gin.Context, msg []byte) (*http.Response, error) {

	client := http.Client{
		Timeout: time.Duration(300 * time.Microsecond),
	}

	request, err := http.NewRequest(http.MethodPut, req.url, bytes.NewBuffer(msg))
	if err != nil {
		return nil, err
	}

	for key, val := range req.headers {
		request.Header.Add(key, val)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (req Request) Delete(ctx *gin.Context) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Duration(300 * time.Microsecond),
	}

	request, err := http.NewRequest(http.MethodDelete, req.url, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range req.headers {
		request.Header.Add(key, val)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
