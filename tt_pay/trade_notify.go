package tt_pay

import (
	"code.byted.org/caijing_pay/tp_server_sdk/tp_golang/tt_pay/consts"
	"code.byted.org/caijing_pay/tp_server_sdk/tp_golang/tt_pay/util"
	"context"
	"errors"
	"net/url"
)

type TradeNotifyRequest struct {
	Param string
}

func TradeNotify(ctx context.Context, req *TradeNotifyRequest) (*TradeNotifyResponse, error) {
	params, err := url.ParseQuery(req.Param)

	if err != nil {
		util.Debug("Parse params failed: err[%s]", err)
		return nil, err
	}

	resp := new(TradeNotifyResponse)
	resp.Param = make(map[string]string)
	signMap := make(map[string]interface{})

	// TODO：这里收到的Params里有没有可能一个key有多个value？
	for key, val := range params {
		resp.Param[key] = val[0]
		signMap[key] = interface{}(val[0])
	}

	resp.Decode()

	sign := resp.Get("sign")

	if valid := util.VerifyMd5WithRsa(signMap, sign, consts.Tt_pay_public_key); !valid {
		util.Debug("VerifyMd5WithRsa failed!")
		return nil, errors.New("sign verification failed: invalid sign")
	}

	return resp, nil
}

// SetParam 将回调的param参数赋值给该实例成员变量
func (req *TradeNotifyRequest) SetParam(s string) {
	req.Param = s
}

type TradeNotifyResponse struct {
	Param map[string]string

	NotifyId    string
	SignType    string
	Sign        string
	AppId       string //非必须
	EventCode   string //非必须
	MerchantId  string
	OutOrderNo  string
	TradeNo     string
	TotalAmount string
	PayChannel  string
	PayTime     string
	PayType     string
	TradeStatus string
	TradeMsg    string
	Extension   string `json:"extension"`
}

func (resp *TradeNotifyResponse) Decode() {
	resp.NotifyId = resp.Get("notify_id")
	resp.SignType = resp.Get("sign_type")
	resp.Sign = resp.Get("sign")
	resp.AppId = resp.Get("app_id")
	resp.EventCode = resp.Get("event_code")
	resp.OutOrderNo = resp.Get("out_order_no")
	resp.TradeNo = resp.Get("trade_no")
	resp.TotalAmount = resp.Get("total_amount")
	resp.PayChannel = resp.Get("pay_channel")
	resp.MerchantId= resp.Get("merchant_id")
	resp.PayTime = resp.Get("pay_time")
	resp.PayType = resp.Get("pay_type")
	resp.TradeStatus = resp.Get("trade_status")
	resp.TradeMsg = resp.Get("trade_msg")
	resp.Extension = resp.Get("extension")
}

// 设置原始响应
func (resp *TradeNotifyResponse) Get(key string) string {
	return resp.Param[key]
}
