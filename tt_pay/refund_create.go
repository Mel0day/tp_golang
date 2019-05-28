package tt_pay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/ttcjpay/tp_golang/tt_pay/consts"
	"github.com/ttcjpay/tp_golang/tt_pay/util"
	"net/url"
	"time"
)

// 退款申请接口
func RefundCreate(ctx context.Context, req *RefundCreateRequest) (*RefundCreateResponse, error) {
	tpCli := NewTPClient(req.TPClientTimeoutMs)
	if err := req.checkParams(); err != nil {
		util.Debug("Parameter verification failed: err[%s]", err)
		return nil, err
	}
	resp := NewRefundCreateResponse()
	err := tpCli.Execute(ctx, req, resp)
	if err != nil {
		util.Debug("Request Execution failed: err[%s]", err)
		return nil, err
	}
	return resp, nil
}

// TODO: 补全参数
// 退款申请Request
type RefundCreateRequest struct {
	Config
	Uid        string
	OutOrderNo string
	TradeNo    string
	SignType   string
	Version    string
	timestamp  string

	OutRefundNo  string //商户退款单号
	RefundAmount int // 注意合法性检查
	NotifyUrl    string
	RiskInfo     string // 给出必要的事例，必须字段ip、device_id等

	path 		string
	bizContent *simplejson.Json
}

// New函数内赋默认值，目前含默认值（或仅支持一个值的）参数包括：
// Version = "1.0"
// SignType = "MD5"
// path = "gateway"
// Config.TPDomain = "https://tp-pay.snssdk.com"（正式）或 "https://tp-pay-test.snssdk.com" （测试）
// 另外，注意初始化bizContent，以免出现nil指针错误
func NewRefundCreateRequest(config Config) *RefundCreateRequest {
	ret := new(RefundCreateRequest)
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

// 编码函数：
// 	盖时间戳、加签、序列化
func (req *RefundCreateRequest) Encode() (string, error) {
	// 盖时间戳
	req.timestamp = fmt.Sprintf("%d", time.Now().Unix())
	// 加签
	req.bizContent.Set("out_order_no", req.OutOrderNo)
	req.bizContent.Set("trade_no", req.TradeNo)
	req.bizContent.Set("merchant_id", req.Config.MerchantId)
	req.bizContent.Set("uid", req.Uid)
	req.bizContent.Set("out_refund_no", req.OutRefundNo)
	req.bizContent.Set("refund_amount", req.RefundAmount)
	req.bizContent.Set("notify_url", req.NotifyUrl)
	req.bizContent.Set("risk_info", req.RiskInfo)

	bizContentBytes, _ := req.bizContent.Encode()

	signParams := make(map[string]interface{})
	signParams["app_id"] = req.Config.AppId
	signParams["method"] = consts.MethodRefundCreate
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
	values.Set("method", consts.MethodRefundCreate)
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
func (req *RefundCreateRequest) GetLogId() string {
	id := ""
	if len(req.TradeNo) != 0 {
		id = req.TradeNo
	}
	if len(req.OutOrderNo) != 0 {
		id = req.OutOrderNo
	}
	return fmt.Sprintf("%s_%s_%s_%s", req.Config.AppId, req.Config.MerchantId, id, req.timestamp)
}

// 获取请求url地址
func (req *RefundCreateRequest) GetUrl() string {
	return req.Config.TPDomain + "/" + req.path
}

// 比如提供该接口，方便业务方设置可选参数，比如product_code、payment_type等
func (req *RefundCreateRequest) SetBizContentKV(key string, val interface{}) {
	req.bizContent.Set(key, val)
}

type RefundCreateResponse struct {
	Data         *simplejson.Json
	OutOrderNo   string `json:"out_order_no"`
	OutRefundNo  string `json:"out_refund_no"`
	RefundNo     string `json:"refund_no"`
	RefundAmount string `json:"refund_amount"` // 注意合法性检查
}

func NewRefundCreateResponse() *RefundCreateResponse{
	ret := new(RefundCreateResponse)
	ret.Data = simplejson.New()
	return ret
}

// 将响应json数据反序列化为对应接口
func (resp *RefundCreateResponse) Decode() {
	respBytes, _ := resp.Data.Get("response").Encode()
	json.Unmarshal(respBytes, resp)
}

// 设置原始响应
func (resp *RefundCreateResponse) SetData(data *simplejson.Json) {
	resp.Data = data
}


// 目前只查验大写字母开头的参数(用户必传参数)
func (req *RefundCreateRequest) checkParams() error {
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

	// OutRefundNo在此接口内必选，需要额外确认不为空
	if req.OutRefundNo == "" {
		return errors.New("OutRefundNo shouldn't be blank in this interface")
	}

	if err := util.CheckOutRefundNo(req.OutRefundNo); err != nil {
		return err
	}

	if err := util.CheckRefundAmount(req.RefundAmount); err != nil {
		return err
	}

	if err := util.CheckNotifyUrl(req.NotifyUrl); err != nil {
		return err
	}

	if err := util.CheckRiskInfo(req.RiskInfo); err != nil {
		return err
	}

	return nil
}
