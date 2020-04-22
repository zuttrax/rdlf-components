package resty_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zuttrax/rdlf-components/pkg/resty"
)

type LogStub struct{}

func (l LogStub) Info(msg string) {

}

func (l LogStub) Error(err error) {

}

func (l LogStub) Fatal(err error) {

}

//nolint
func addContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, "Client-IP", "192.168.0.1")
	ctx = context.WithValue(ctx, "X-Token", "token")
	ctx = context.WithValue(ctx, "Endpoint", "path/path")

	return ctx
}

func TestRequest_Post(t *testing.T) {
	type fields struct {
		headers resty.Headers
	}
	type args struct {
		ctx context.Context
		msg []byte
	}
	tt := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			"success case with headers",
			fields{
				headers: map[string]string{
					"token": "token",
				},
			},
			args{
				ctx: context.Background(),
				msg: []byte("test"),
			},
			&http.Response{
				Status:     http.StatusText(http.StatusNoContent),
				StatusCode: http.StatusNoContent,
			},
			false,
		},
		{
			"success case without headers",
			fields{
				headers: nil,
			},
			args{
				ctx: context.Background(),
				msg: []byte("test"),
			},
			&http.Response{
				Status:     http.StatusText(http.StatusNoContent),
				StatusCode: http.StatusNoContent,
			},
			false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()

			req := resty.Request{
				BaseURL: server.URL,
				Trace:   true,
				Client:  *server.Client(),
				Headers: tc.fields.headers,
				Log:     LogStub{},
			}

			endpoint := resty.NewEndpoint(req, "")

			tc.args.ctx = addContext(tc.args.ctx)

			got, err := endpoint.Post(tc.args.ctx, tc.args.msg)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			assert.Equal(t, tc.want.StatusCode, got.StatusCode)
		})
	}
}

func TestRequest_Put(t *testing.T) {
	type fields struct {
		headers resty.Headers
	}
	type args struct {
		ctx context.Context
		msg []byte
	}
	tt := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			"success case with headers",
			fields{
				headers: map[string]string{
					"token": "token",
				},
			},
			args{
				ctx: context.Background(),
				msg: []byte("test"),
			},
			&http.Response{
				Status:     http.StatusText(http.StatusNoContent),
				StatusCode: http.StatusOK,
				Body:       nil,
			},
			false,
		},
		{
			"success case without headers",
			fields{
				headers: nil,
			},
			args{
				ctx: context.Background(),
				msg: []byte("test"),
			},
			&http.Response{
				Status:     http.StatusText(http.StatusNoContent),
				StatusCode: http.StatusOK,
				Body:       nil,
			},
			false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			req := resty.Request{
				BaseURL: server.URL,
				Trace:   false,
				Client:  *server.Client(),
				Headers: tc.fields.headers,
			}

			endpoint := resty.NewEndpoint(req, "")

			got, err := endpoint.Put(tc.args.ctx, tc.args.msg)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			assert.Equal(t, tc.want.StatusCode, got.StatusCode)
		})
	}
}

func TestRequest_Get(t *testing.T) {
	type fields struct {
		headers    resty.Headers
		parameters resty.Parameters
	}
	type args struct {
		ctx context.Context
	}
	tt := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			"success case with headers and parameters",
			fields{
				headers: map[string]string{
					"token": "token",
				},
				parameters: map[string]string{
					"param1": "value",
					"param2": "value",
				},
			},
			args{
				ctx: context.Background(),
			},
			&http.Response{
				Status:     http.StatusText(http.StatusNoContent),
				StatusCode: http.StatusOK,
				Body:       nil,
			},
			false,
		},
		{
			"success case with headers but without parameters",
			fields{
				headers: map[string]string{
					"token": "token",
				},
				parameters: nil,
			},
			args{
				ctx: context.Background(),
			},
			&http.Response{
				Status:     http.StatusText(http.StatusNoContent),
				StatusCode: http.StatusOK,
				Body:       nil,
			},
			false,
		},
		{
			"success case without headers and parameters",
			fields{
				headers:    nil,
				parameters: nil,
			},
			args{
				ctx: context.Background(),
			},
			&http.Response{
				Status:     http.StatusText(http.StatusNoContent),
				StatusCode: http.StatusOK,
				Body:       nil,
			},
			false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			req := resty.Request{
				BaseURL: server.URL,
				Trace:   false,
				Client:  *server.Client(),
				Headers: tc.fields.headers,
			}

			endpoint := resty.NewEndpoint(req, "")

			queryParams := resty.NewParameters()
			for k, v := range tc.fields.parameters {
				queryParams = queryParams.Add(k, v)
			}

			got, err := endpoint.Get(tc.args.ctx, queryParams)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			assert.Equal(t, tc.want.StatusCode, got.StatusCode)
		})
	}
}

func TestRequest_Delete(t *testing.T) {
	type fields struct {
		headers resty.Headers
	}
	type args struct {
		ctx context.Context
	}
	tt := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			"success case to delete",
			fields{
				headers: map[string]string{
					"token": "token",
				},
			},
			args{
				ctx: context.Background(),
			},
			&http.Response{
				Status:     http.StatusText(http.StatusNoContent),
				StatusCode: http.StatusNoContent,
				Body:       nil,
			},
			false,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()

			req := resty.Request{
				BaseURL: server.URL,
				Trace:   false,
				Client:  *server.Client(),
				Headers: tc.fields.headers,
			}

			endpoint := resty.NewEndpoint(req, "")

			got, err := endpoint.Delete(tc.args.ctx)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			assert.Equal(t, tc.want.StatusCode, got.StatusCode)
		})
	}
}
