// Copyright (C) 2017 Kazumasa Kohtaka <kkohtaka@gmail.com> All right reserved
// This file is available under the MIT license.

package httpclient

import (
	"time"

	"github.com/json-iterator/go"
	"github.com/noda-sin/go-bitflyer/pkg/api"
	"github.com/noda-sin/go-bitflyer/pkg/api/auth"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type httpClient struct {
	authConfig *auth.AuthConfig
}

func New() *httpClient {
	return &httpClient{}
}

func (hc *httpClient) Auth(authConfig *auth.AuthConfig) *httpClient {
	hc.authConfig = authConfig
	return hc
}

func (hc *httpClient) Request(api api.API, req api.Request, result interface{}) error {
	u, err := api.BaseURL()
	if err != nil {
		return errors.Wrapf(err, "set base URI")
	}
	payload := req.Payload()

	rawReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(rawReq)

	rawReq.Header.SetMethod(req.Method())
	rawReq.SetRequestURI(u.String())

	if len(payload) > 0 {
		rawReq.SetBody(payload)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err != nil {
		return errors.Wrapf(err, "create POST request from url: %s", u.String())
	}
	if hc.authConfig != nil {
		header, err := auth.GenerateAuthHeaders(hc.authConfig, time.Now(), api, req)
		if err != nil {
			return errors.Wrap(err, "generate auth header")
		}
		for key, _ := range *header {
			value := header.Get(key)
			rawReq.Header.Set(key, value)
		}
	}
	if len(payload) > 0 {
		rawReq.Header.Set("Content-Type", "application/json")
	}

	err = fasthttp.Do(rawReq, resp)

	if err != nil {
		return errors.Wrapf(err, "send HTTP request with url: %s", u.String())
	}

	data := resp.Body()

	err = json.Unmarshal(data, result)
	if err != nil {
		return errors.Wrapf(err, "unmarshal data: %s", string(data))
	}
	return nil
}
