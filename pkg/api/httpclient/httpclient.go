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
	"bytes"
	"net/http"
	"io/ioutil"
	"io"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type httpClient struct {
	authConfig *auth.AuthConfig
	useFastHttp bool
}

func New(useFastHttp bool) *httpClient {
	return &httpClient{
		useFastHttp: useFastHttp,
	}
}

func (hc *httpClient) Auth(authConfig *auth.AuthConfig) *httpClient {
	hc.authConfig = authConfig
	return hc
}

func (hc *httpClient) Request(api api.API, req api.Request, result interface{}) error {
	if hc.useFastHttp {
		return hc.requestFastHttp(api, req, result)
	} else {
		return hc.requestDefaultHttp(api, req, result)
	}
}


func (hc *httpClient) requestDefaultHttp(api api.API, req api.Request, result interface{}) error {
	u, err := api.BaseURL()
	if err != nil {
		return errors.Wrapf(err, "set base URI")
	}
	payload := req.Payload()

	var body io.Reader
	if len(payload) > 0 {
		body = bytes.NewReader(payload)
	}
	rawReq, err := http.NewRequest(req.Method(), u.String(), body)
	if err != nil {
		return errors.Wrapf(err, "create POST request from url: %s", u.String())
	}
	if hc.authConfig != nil {
		header, err := auth.GenerateAuthHeaders(hc.authConfig, time.Now(), api, req)
		if err != nil {
			return errors.Wrap(err, "generate auth header")
		}
		rawReq.Header = *header
	}
	if len(payload) > 0 {
		rawReq.Header.Set("Content-Type", "application/json")
	}

	c := &http.Client{}
	resp, err := c.Do(rawReq)
	if err != nil {
		return errors.Wrapf(err, "send HTTP request with url: %s", u.String())
	}
	defer resp.Body.Close()

	// TODO: Don't use ioutil.ReadAll()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "read data fetched from url: %s", u.String())
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		return errors.Wrapf(err, "unmarshal data: %s", string(data))
	}
	return nil
}


func (hc *httpClient) requestFastHttp(api api.API, req api.Request, result interface{}) error {
	u, err := api.BaseURL()
	if err != nil {
		return errors.Wrapf(err, "set base URI")
	}
	payload := req.Payload()

	rawReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(rawReq)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	rawReq.Header.SetMethod(req.Method())
	rawReq.SetRequestURI(u.String())

	var body io.Reader
	var payloadLen = len(payload)
	if payloadLen > 0 {
		body = bytes.NewReader(payload)
		rawReq.SetBodyStream(body, payloadLen)
		rawReq.Header.Set("Content-Type", "application/json")
	}

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
