package tt_pay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"code.byted.org/caijing_pay/tp_server_sdk/tp_golang/tt_pay/consts"
	"code.byted.org/caijing_pay/tp_server_sdk/tp_golang/tt_pay/util"
	simplejson "github.com/bitly/go-simplejson"
)

// 订单查询接口
func TradeQuery(ctx context.Context, req *TradeQueryRequest) (*TradeQueryResponse, error) {
	tpCli := NewTPClient(req.TPClientTimeoutMs)
	if err := req.checkParams(); err != nil {
		util.Debug("Parameter verification failed: err[%s]", err)
		return nil, err
	}
	resp := NewTradeQueryResponse()
	err := tpCli.Execute(ctx, req, resp)
	if err != nil {
		util.Debug("Request Execution failed: err[%s]", err)
		return nil, err
	}
	return resp, nil
}

// TODO: 补全参数
// 订单查询Request
type TradeQueryRequest struct {
	Config
	Uid        string
	OutOrderNo string
	TradeNo    string
	SignType   string
	Version    string
	timestamp  string
	path 		string
	bizContent *simplejson.Json
}

// New函数内赋默认值，目前含默认值（或仅支持一个值的）参数包括：
// Version = "1.0"
// SignType = "MD5"
// path = "gateway"
// Config.TPDomain = "https://tp-pay.snssdk.com"（正式）或 "https://tp-pay-test.snssdk.com" （测试）
// 另外，注意初始化bizContent，以免出现nil指针错误
func NewTradeQueryRequest(config Config) *TradeQueryRequest {
	ret := new(TradeQueryRequest)
	ret.Config = config
	ret.Version = "1.0"
	ret.SignType = "MD5"
	ret.path = consts.TPPath
	if len(ret.Config.TPDomain) == 0 {
		ret.Config.TPDomain = consts.TPDomain
	}
	ret.bizContent = simplejson.New()
	return ret
}

func (req *TradeQueryRequest) Encode() (string, error) {
	// 盖时间戳
	req.timestamp = fmt.Sprintf("%d", time.Now().Unix())
	// 加签
	req.bizContent.Set("merchant_id", req.MerchantId)
	req.bizContent.Set("uid", req.Uid)
	// TODO: 补全参数
	req.bizContent.Set("uid_type", "")
	req.bizContent.Set("out_order_no", req.OutOrderNo)
	req.bizContent.Set("trade_no", req.TradeNo)

	bizContentBytes, _ := req.bizContent.Encode()

	signParams := make(map[string]interface{})
	signParams["app_id"] = req.Config.AppId
	signParams["method"] = consts.MethodTradeQuery
	signParams["format"] = "JSON"
	signParams["charset"] = "utf-8"
	signParams["sign_type"] = req.SignType
	signParams["timestamp"] = req.timestamp
	signParams["version"] = req.Version
	signParams["biz_content"] = string(bizContentBytes)

	sign := util.BuildMd5WithSalt(signParams, req.Config.AppSecret)
	// 序列化
	values := url.Values{}
	values.Set("app_id", req.Config.AppId)
	values.Set("method", consts.MethodTradeQuery)
	values.Set("format", "JSON")
	values.Set("charset", "utf-8")
	values.Set("sign_type", req.SignType)
	values.Set("sign", sign)
	values.Set("timestamp", req.timestamp)
	values.Set("version", req.Version)
	values.Set("biz_content", string(bizContentBytes))

	return values.Encode(), nil
}

// 生成该次请求logid
// out_order_no 和 trade_no 哪个不空用哪个，都不空优先用out_order_no
func (req *TradeQueryRequest) GetLogId() string {
	id := ""
	if len(req.TradeNo) != 0 {
		id = req.TradeNo
	}
	if len(req.OutOrderNo) != 0 {
		id = req.OutOrderNo
	}
	return fmt.Sprintf("%s_%s_%s_%s", req.Config.AppId, req.Config.MerchantId, id, req.timestamp)
}

func (req *TradeQueryRequest) GetUrl() string {
	return req.Config.TPDomain + "/" + req.path
}

func (req *TradeQueryRequest) SetBizContentKV(key string, val interface{}) {
	req.bizContent.Set(key, val)
}

type TradeQueryResponse struct {
	Data        *simplejson.Json

	TradeNo     string `json:"trade_no"`
	OutOrderNo  string `json:"out_order_no"`
	MerchantId  string `json:"merchant_id"`
	Uid         string `json:"uid"`
	TradeStatus string `json:"trade_status"`
	TradeName   string `json:"trade_name"`
	TradeDesc   string `json:"trade_desc"`
	TotalAmount string `json:"total_amount"`
	Curreny     string `json:"currency"`
	PayChannel  string `json:"pay_channel"`
}

func NewTradeQueryResponse() *TradeQueryResponse {
	ret := new(TradeQueryResponse)
	ret.Data = simplejson.New()
	return ret
}

func (resp *TradeQueryResponse) Decode() {
	respBytes, _ := resp.Data.Get("response").Encode()
	json.Unmarshal(respBytes, resp)
}

func (resp *TradeQueryResponse) SetData(data *simplejson.Json) {
	resp.Data = data
}

func (req *TradeQueryRequest) checkParams() error {
	// 以下五个含默认值参数自动补全
	if req.SignType == "" {
		req.SignType = "MD5"
	}

	if req.Version == "" {
		req.Version = "1.0"
	}

	if req.path == "" {
		req.path = consts.TPPath
	}

	if req.TPDomain == "" {
		req.TPDomain = consts.TPDomain
	}

	if req.bizContent == nil {
		req.bizContent = simplejson.New()
	}

	if err := util.CheckAppId(req.AppId); err != nil {
		return err
	}

	if err := util.CheckMerchantId(req.MerchantId); err != nil {
		return err
	}

	if err := util.CheckUid(req.Uid); err != nil {
		return err
	}

	// 二选一参数判断
	if req.OutOrderNo == "" && req.TradeNo == "" {
		return errors.New("OutOrderNo and TradeNo can't both be blank")
	}

	if len(req.OutOrderNo) > 0 {
		if err := util.CheckOutOrderNo(req.OutOrderNo); err != nil {
			return err
		}
	}

	if len(req.TradeNo) > 0 {
		if err := util.CheckTradeNo(req.TradeNo); err != nil {
			return err
		}
	}

	return nil
}