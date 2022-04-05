package rest

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type AuthedType int

const (
	None AuthedType = iota
	JWT
	BASIC
)

type ApiClient struct {
	ctx    context.Context
	client *http.Client
}

type ClientInfo struct {
	BaseUrl  string
	UserName string
	PassWord string
	JwtAuth  string
	AuthType AuthedType // 1: Jwt, 2: BasicAuth, 0: default
}

type ApiInfo struct {
	name    string
	path    string
	method  string
	reqBody interface{}
	resData interface{}
}

func NewApi(name, path, method string, reqBody, resData interface{}) *ApiInfo {
	return &ApiInfo{
		name:    name,
		path:    path,
		method:  method,
		reqBody: reqBody,
		resData: resData,
	}
}

func NewClient(ctx context.Context, tlsConfig *tls.Config, connectTimeout time.Duration, requestTimeout time.Duration) *ApiClient {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   connectTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ResponseHeaderTimeout: requestTimeout,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		TLSClientConfig:       tlsConfig,
	}

	return &ApiClient{
		ctx: ctx,
		client: &http.Client{
			Timeout:   requestTimeout,
			Transport: transport,
		}}
}

func (clInfo *ClientInfo) makeEndpoint(path string) (endpoint *url.URL, err error) {
	btu := []byte(clInfo.BaseUrl)
	if btu[len(btu)-1] == '/' {
		newUrl := btu[:len(btu)-2]
		clInfo.BaseUrl = string(newUrl)
	}
	urlPath := clInfo.BaseUrl + "/" + path

	if endpoint, err = url.Parse(urlPath); err != nil {
		return nil, errors.New(fmt.Sprintf("path service is invalid: %s", path))
	}
	return
}

func (ac *ApiClient) Call(clientInfo *ClientInfo, api *ApiInfo) error {
	var requestId string
	var err error
	var payload io.Reader
	var endpoint *url.URL
	var resp *http.Response
	if payload, err = toJsonReader(api.reqBody); err != nil {
		return errors.New(fmt.Sprintf("request body is invalid: %s", err))
	}

	if endpoint, err = clientInfo.makeEndpoint(api.path); err != nil {
		return err
	}

	req, _ := http.NewRequest(api.method, endpoint.String(), payload)
	req.Header.Add("X-Request-ID", requestId)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Service-Name", api.name)
	switch clientInfo.AuthType {
	case 1:
		req.Header.Add("Authorization", clientInfo.JwtAuth)
		break
	case 2:
		req.SetBasicAuth(clientInfo.UserName, clientInfo.PassWord)
		break
	default:
		break
	}

	if resp, err = ac.client.Do(req); err != nil {
		return errors.New(fmt.Sprintf("call api error: %s, %s, %s", requestId, endpoint, err))
	}

	if err = json.NewDecoder(resp.Body).Decode(api.resData); err != nil {
		return errors.New(fmt.Sprintf("response data is invallid: %s, %s", requestId, err))
	}
	resp.Body.Close()

	return nil
}

func toJsonReader(v interface{}) (io.Reader, error) {
	data, err := json.Marshal(v)
	return bytes.NewBuffer(data), err
}
