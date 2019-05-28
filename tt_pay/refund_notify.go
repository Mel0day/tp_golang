package tt_pay

import (
	"code.byted.org/caijing_pay/tp_server_sdk/tp_golang/tt_pay/consts"
	"code.byted.org/caijing_pay/tp_server_sdk/tp_golang/tt_pay/util"
	"context"
	"errors"
	"net/url"
)

type RefundNotifyRequest struct {
	Param       string
}

func RefundNotify(ctx context.Context, req *RefundNotifyRequest) (*RefundNotifyResponse, error) {
	params, err := url.ParseQuery(req.Param)

	if err != nil {
		util.Debug("Parse params failed: err[%s]", err)
		return nil, err
	}

	resp := new(RefundNotifyResponse)
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

// 将回调的param参数赋值给该实例成员变量
func (req *RefundNotifyRequest) SetParam(s string) {
	req.Param = s
}

type RefundNotifyResponse struct {
	Param	map[string]string

	NotifyId string
	SignType string
	Sign     string
	AppId    string //非必须

	EventCode string
	OutRefundNo string //商户的退款单号
	RefundNo string //头条支付的退款单号
	RefundAmount string //退款金额
	RefundTime string
	MerchantId string
	RefundStatus string
}

func (resp *RefundNotifyResponse) Decode() {
	resp.NotifyId = resp.Get("notify_id")
	resp.SignType = resp.Get("sign_type")
	resp.Sign = resp.Get("sign")
	resp.AppId = resp.Get("app_id")
	resp.EventCode = resp.Get("event_code") // 这是啥？
	resp.OutRefundNo = resp.Get("out_refund_no")
	resp.RefundNo = resp.Get("refund_no")
	resp.RefundAmount= resp.Get("refund_amount")
	resp.RefundTime = resp.Get("refund_time")
	resp.MerchantId = resp.Get("merchant_id")
	resp.RefundStatus = resp.Get("refund_status")
}

// 提取Param内的值
func (resp *RefundNotifyResponse) Get(key string) string{
	return resp.Param[key]
}

