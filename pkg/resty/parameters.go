package resty

import "net/http"

type Parameters map[string]string

func NewParameters() Parameters {
	return make(map[string]string)
}

func (p Parameters) Add(key string, value string) Parameters {
	parametersResult := make(map[string]string)
	for k, v := range p {
		parametersResult[k] = v
	}

	parametersResult[key] = value

	return parametersResult
}

func addQueryParamsToRequest(httpRequest *http.Request, inputParams Parameters) {
	params := httpRequest.URL.Query()

	for key, value := range inputParams {
		params.Add(key, value)
	}

	httpRequest.URL.RawQuery = params.Encode()
}
