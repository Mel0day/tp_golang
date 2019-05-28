package tt_pay

import (
	"context"
	"github.com/ttcjpay/tp_golang/tt_pay/util"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

var httpClient = http.Client{
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 60 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:        20000,
		MaxIdleConnsPerHost: 2000,
		IdleConnTimeout:     65 * time.Second,
		TLSHandshakeTimeout: 3 * time.Second,
	},
}

type TPRequest interface {
	Encode() (string, error)
	GetUrl() string
	GetLogId() string
}

type TPResponse interface {
	SetData(data *simplejson.Json)
	Decode()
}

type TPClient struct {
	TimeoutMs int
}

func NewTPClient(timeoutMs int) *TPClient {
	if timeoutMs <= 0 {
		timeoutMs = 6000
	}

	ret := new(TPClient)
	ret.TimeoutMs = timeoutMs

	return ret
}

// 返回tt_pay.Error
func (c *TPClient) Execute(ctx context.Context, req TPRequest, resp TPResponse) error {
	body, err := req.Encode()
	if err != nil {
		return err
	}

	_, respBytes, err := HttpPost(req.GetUrl(), "application/x-www-form-urlencoded", body, req.GetLogId(), c.TimeoutMs)
	if err != nil {
		return err
	}

	util.Debug("resp[%s]", string(respBytes))

	respJson, err := simplejson.NewJson(respBytes)
	if err != nil {
		return err
	}

	if respJson.Get("response").Get("code").MustString("") != "10000" {
		ret := new(Error)
		ret.Code = respJson.Get("response").Get("code").MustString("")
		ret.Msg = respJson.Get("response").Get("msg").MustString("")
		ret.SubCode = respJson.Get("response").Get("sub_code").MustString("")
		ret.SubMsg = respJson.Get("response").Get("sub_msg").MustString("")
		ret.Detail = "log_id:" + req.GetLogId()
		return ret
	}

	resp.SetData(respJson)
	resp.Decode()

	return nil
}

// 配置HttpClient
func SetHttpClient(c http.Client) {
	httpClient = c
	log.Printf("Set HttpClient to: %v", c)
}

func HttpPost(url, contentType, body string, logId string, timeoutMs int) (int, []byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		util.Debug("HttpTypePost NewRequest url[%s] body[%s] err[%s]\n", url, body, err)
		return 0, nil, err
	}
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Duration(timeoutMs)*time.Millisecond)
	req = req.WithContext(ctx)
	req.Header.Set("Content-type", contentType)
	req.Header.Set("X-Tt-Logid", logId)

	resp, err := httpClient.Do(req)
	if err != nil {
		util.Debug("HttpPost client.Do err: %v, url: %s\n", err, url)
		return 0, nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Debug("HttpPost ioutil.ReadAll err: %v, url: %s\n", err, url)
		return resp.StatusCode, nil, err
	}

	util.Debug("HttpTypePost url[%s] contentType[%s] body[%s] code[%d] resp body[%s]\n",
		url, contentType, body, resp.StatusCode, string(respBody))
	return resp.StatusCode, respBody, nil
}
