package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"order/client/dtos"
	"order/pkg/rest"
	"testing"
	"time"
)

func Test_Usms(t *testing.T) {
	const usmsPath = ""
	ctx := context.Background()
	tls := tls.Config{}
	connectTimeOut := 5 * time.Second
	requestTimeOut := 30 * time.Second
	client := rest.NewClient(ctx, &tls, connectTimeOut, requestTimeOut)
	clientInfo := rest.ClientInfo{
		BaseUrl: "",
	}
	var resData dtos.UserInfoData
	apiInfo := rest.NewApi("USMS", usmsPath, http.MethodPost, nil, &resData)
	err := client.Call(&clientInfo, apiInfo)
	if err != nil {
		t.Fatal(err)
	}
}
