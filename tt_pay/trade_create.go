package tt_pay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"code.byted.org/caijing_pay/tp_server_sdk/tp_golang/tt_pay/consts"
	"code.byted.org/caijing_pay/tp_server_sdk/tp_golang/tt_pay/util"
	"github.com/bitly/go-simplejson"
)

func TradeCreate(ctx context.Context, req *TradeCreateRequest) (*TradeCreateResponse, error) {
	// 检查参数
	if err := req.checkParams(); err != nil {
		util.Debug("Parameter verification failed: err[%s]", err)
		return nil, err
	}
	// 新建Response
	resp := NewTradeCreateResponse(req)
	// 1.0需要与财经后端通信
	if req.Version == "1.0" {
		// 新建Client
		tpCli := NewTPClient(req.TPClientTimeoutMs)

		// 发请求
		err := tpCli.Execute(ctx, req, resp)
		if err != nil {
			util.Debug("Request Execution failed: err[%s]", err)
			return nil, err
		}
	}
	return resp, nil
}

// 预下单Request
type TradeCreateRequest struct {
	Config
	Version   string

	SignType  string

	OutTradeNo		string
	WithdrawTradeNo string
	Uid         string
	UidType     string
	OutOrderNo  string
	TradeNo 	string
	ProductCode string
	PaymentType string
	ChannelInfo 	string
	timestamp string //非商户生成，自动填充参数
	TotalAmount interface{}
	NotifyUrl  string
	// 拉起收银台所需额外参数
	CashdeskTradeType string
	CashdeskLimitPay  string
	CashdeskExts      string


	Subject     string
	Body        string
	TradeTime  string // 时间戳，合法检查
	ValidTime  string
	ProductId	string
	Currency    string

	OutSubscribeId	string
	ThisTermStart	string
	ThisTermEnd		string
	SubscribeStart	string
	subscribeEnd	string
	SubscribePeriodNo string
	SubscribeCycle	string
	CycleUnit		string
	SubscribeType	string
	RealAmount		string

	// 区分内部版本的字段（比如是否需要收银台预下单）
	version			string
	AlipayUrl		string
	WxUrl			string
	WxType			string

	ReturnUrl  string // 前端回调地址
	ServiceFee string // 平台手续费
	RiskInfo   string // 给出必要的事例，必须字段ip、device_id等

	SettlementProductCode string
	SettlementExt         string
	SellerMerchantId      string
	RoyaltyParameters     string
	TransCode             string

	path       string
	bizContent *simplejson.Json
}

// New函数内赋默认值，目前含默认值（或仅支持一个值的）参数包括：
// Version = "1.0"
// SignType = "MD5"
// path = "gateway"
// Config.TPDomain = "https://tp-pay.snssdk.com"（正式）或 "https://tp-pay-test.snssdk.com" （测试）
// 另外，注意初始化bizContent，以免出现nil指针错误
func NewTradeCreateRequest(config Config) *TradeCreateRequest {
	ret := new(TradeCreateRequest)
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
func (req *TradeCreateRequest) Encode() (string, error) {
	// 盖时间戳
	req.timestamp = fmt.Sprintf("%d", time.Now().Unix())
	//加签
	req.bizContent.Set("out_order_no", req.OutOrderNo)
	req.bizContent.Set("uid", req.Uid)
	req.bizContent.Set("uid_type", req.UidType)
	req.bizContent.Set("merchant_id", req.MerchantId)
	req.bizContent.Set("total_amount", req.TotalAmount)
	req.bizContent.Set("currency", req.Currency)
	req.bizContent.Set("subject", req.Subject)
	req.bizContent.Set("body", req.Body)
	req.bizContent.Set("product_code", req.ProductCode)
	req.bizContent.Set("payment_type", req.PaymentType)
	req.bizContent.Set("trade_time", req.TradeTime)
	req.bizContent.Set("valid_time", req.ValidTime)
	req.bizContent.Set("notify_url", req.NotifyUrl)
	req.bizContent.Set("return_url", req.ReturnUrl)
	req.bizContent.Set("service_fee", req.ServiceFee)
	req.bizContent.Set("risk_info", req.RiskInfo)
	req.bizContent.Set("settlement_prouct_code", req.SettlementProductCode)
	req.bizContent.Set("settlement_ext", req.SettlementExt)
	req.bizContent.Set("seller_merchant_id", req.SellerMerchantId)
	req.bizContent.Set("royalty_parameters", req.RoyaltyParameters)
	req.bizContent.Set("tarns_code", req.TransCode)

	// Json encode
	bizContentBytes, _ := req.bizContent.Encode()

	signParams := make(map[string]interface{})
	signParams["app_id"] = req.Config.AppId
	signParams["method"] = consts.MethodTradeCreate
	signParams["format"] = "JSON"
	signParams["charset"] = "utf-8"
	signParams["sign_type"] = req.SignType
	signParams["timestamp"] = req.timestamp
	signParams["version"] = req.Version
	signParams["biz_content"] = string(bizContentBytes)

	sign := util.BuildMd5WithSalt(signParams, req.Config.AppSecret)
	// URL Encode
	values := url.Values{}
	values.Set("app_id", req.Config.AppId)
	values.Set("method", consts.MethodTradeCreate)
	values.Set("format", "JSON")
	values.Set("charset", "utf-8")
	values.Set("sign_type", req.SignType)
	values.Set("sign", sign)
	values.Set("timestamp", req.timestamp)
	values.Set("version", req.Version)
	values.Set("biz_content", string(bizContentBytes))
	return values.Encode(), nil
}

// GetLogId 生成该次请求logid
func (req *TradeCreateRequest) GetLogId() string {
	return fmt.Sprintf("%s_%s_%s_%s", req.AppId, req.MerchantId, req.OutOrderNo, req.timestamp)
}

// GetUrl 获取请求url地址
func (req *TradeCreateRequest) GetUrl() string {
	return req.Config.TPDomain + "/" + req.path
}

type TradeCreateResponse struct {
	Data    *simplejson.Json
	TradeNo string              `json:"trade_no"`
	req     *TradeCreateRequest // 包含拉起收银台所需参数
}

func NewTradeCreateResponse(req *TradeCreateRequest) *TradeCreateResponse {
	ret := new(TradeCreateResponse)
	ret.Data = simplejson.New()
	ret.req = req
	return ret
}

// 返回拉起sdk收银台已经签名好的参数对
func (resp *TradeCreateResponse) GetCashdeskSdkParams() (map[string]string, error) {
	cashDeskParams := make(map[string]interface{})

	cashDeskParams["app_id"] = resp.req.AppId
	cashDeskParams["sign_type"] = resp.req.SignType
	cashDeskParams["out_trade_no"] = resp.req.OutTradeNo
	cashDeskParams["withdraw_trade_no"] = resp.req.WithdrawTradeNo
	cashDeskParams["merchant_id"] = resp.req.MerchantId
	cashDeskParams["uid"] = resp.req.Uid
	cashDeskParams["out_order_no"] = resp.req.OutOrderNo
	cashDeskParams["trade_no"] = resp.TradeNo
	cashDeskParams["product_code"] = resp.req.ProductCode
	cashDeskParams["payment_type"] = resp.req.PaymentType
	cashDeskParams["channeel_info"] = resp.req.ChannelInfo
	cashDeskParams["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
	if _, ok := resp.req.TotalAmount.(int); ok {
		cashDeskParams["total_amount"] = strconv.Itoa(resp.req.TotalAmount.(int))
	} else {
		cashDeskParams["total_amount"] = resp.req.TotalAmount.(string)
	}
	cashDeskParams["notify_url"] = resp.req.NotifyUrl
	cashDeskParams["trade_type"] = resp.req.CashdeskTradeType
	cashDeskParams["limit_pay"] = resp.req.CashdeskLimitPay
	cashDeskParams["trans_code"] = resp.req.TransCode
	cashDeskParams["exts"] = resp.req.CashdeskExts

	// 2.0新增参数
	cashDeskParams["subject"] = resp.req.Subject
	cashDeskParams["body"] = resp.req.Body
	cashDeskParams["tradde_time"] = resp.req.TradeTime
	cashDeskParams["valid_time"] = resp.req.ValidTime
	cashDeskParams["product_id"] = resp.req.ProductId
	cashDeskParams["currency"] = resp.req.Currency

	cashDeskParams["out_subscribe_id"] = resp.req.OutSubscribeId
	cashDeskParams["this_term_start"] = resp.req.ThisTermStart
	cashDeskParams["this_term_end"] = resp.req.ThisTermEnd
	cashDeskParams["subscribe_start"] = resp.req.SubscribeStart
	cashDeskParams["subscribe_end"] = resp.req.subscribeEnd
	cashDeskParams["subscribe_period_no"] = resp.req.SubscribePeriodNo
	cashDeskParams["subscribe_cycle"] = resp.req.SubscribeCycle
	cashDeskParams["cycle_unit"] = resp.req.CycleUnit
	cashDeskParams["subscribe_type"] = resp.req.SubscribeType
	cashDeskParams["real_amount"] = resp.req.RealAmount

	cashDeskParams["version"] = resp.req.version

	cashDeskParams["alipay_url"] = resp.req.AlipayUrl
	cashDeskParams["wx_url"] = resp.req.WxUrl
	cashDeskParams["wx_type"] = resp.req.WxType

	cashDeskParams["sign"] = util.BuildMd5WithSalt(cashDeskParams, resp.req.AppSecret)

	// convert to map[string]string
	returnParams := make(map[string]string)
	for key, val := range cashDeskParams {
		returnParams[key] = val.(string)
	}
	return returnParams, nil
}

// TODO: Double check this part
// 返回拉起H5收银台url
func (resp *TradeCreateResponse) GetCashdeskWapUrl() (string, error) {
	cashDeskParams, err := resp.GetCashdeskSdkParams()
	if err != nil {
		util.Debug("GetCashdeskWapUrl failed: err[%s] in GetCashdeskSdkParams()", err)
		return "", err
	}
	// url encode cashDeskParams
	paramsForDecode := make(map[string][]string)
	for key, val := range cashDeskParams {
		paramsForDecode[key] = []string{val}
	}
	query := url.Values(paramsForDecode).Encode()
	return resp.req.TPDomain + "/cashdesk?" + query, nil
}

func (resp *TradeCreateResponse) Decode() {
	respBytes, _ := resp.Data.Get("response").Encode()
	json.Unmarshal(respBytes, resp)
}

func (resp *TradeCreateResponse) SetData(data *simplejson.Json) {
	resp.Data = data
}

func (req *TradeCreateRequest) checkParams() error {
	if req.Version == "1.0" {
		return req.checkParamsFor1_0()
	}
	if req.Version == "2.0" {
		return req.checkParamsFor2_0()
	}
	return errors.New("invalid param: Version")
}


func (req *TradeCreateRequest) checkParamsFor2_0() error {
	// 以下六个含默认值参数
	// 为空则补全默认值，否则信任商户的设定
	// TODO: 考虑下赋默认值还是商户传
	if req.SignType == "" {
		req.SignType = "MD5"
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

	if req.version == "" {
		req.version = "2.0"
	}

	if req.ProductCode == "" {
		req.ProductCode = "pay"
	}

	if err := util.CheckCashDeskTradeType(req.CashdeskTradeType); err != nil {
		return err
	}

	// 其他参数，正则检验
	if err := util.CheckAppId(req.AppId); err != nil {
		return err
	}

	if err := util.CheckMerchantId(req.MerchantId); err != nil {
		return err
	}

	// OutOrderNo在此接口内必选，需要额外确认不为空
	if req.OutOrderNo == "" {
		return errors.New("OutOrderNo shouldn't be blank when create a new trade")
	}

	if err := util.CheckOutOrderNo(req.OutOrderNo); err != nil {
		return err
	}

	if err := util.CheckUid(req.Uid); err != nil {
		return err
	}

	// 可选参数，确认不空后正则化查验
	if len(req.UidType) > 0 {
		if err := util.CheckUidType(req.UidType); err != nil {
			return err
		}
	}

	if err := util.CheckTotalAmountFor2_0(req.TotalAmount); err != nil {
		return err
	}

	if err := util.CheckCurrency(req.Currency); err != nil {
		return err
	}

	if err := util.CheckSubject(req.Subject); err != nil {
		return err
	}

	if err := util.CheckBody(req.Body); err != nil {
		return err
	}

	// 可选
	if len(req.ProductCode) > 0 {
		if err := util.CheckProductCode(req.ProductCode); err != nil {
			return err
		}
	}

	// 可选
	if len(req.PaymentType) > 0 {
		if err := util.CheckPaymentType(req.PaymentType); err != nil {
			return err
		}
	}

	if err := util.CheckTradeTime(req.TradeTime); err != nil {
		return err
	}

	if err := util.CheckValidTime(req.ValidTime); err != nil {
		return err
	}

	if err := util.CheckNotifyUrl(req.NotifyUrl); err != nil {
		return err
	}

	// 可选
	if len(req.ReturnUrl) > 0 {
		if err := util.CheckReturnUrl(req.ReturnUrl); err != nil {
			return nil
		}
	}

	// 可选
	if len(req.ServiceFee) > 0 {
		if err := util.CheckServiceFee(req.ServiceFee); err != nil {
			return nil
		}
	}

	if err := util.CheckRiskInfo(req.RiskInfo); err != nil {
		return err
	}

	if len(req.SettlementProductCode) > 0 {
		if err := util.CheckSettlementProductCode(req.SettlementProductCode); err != nil {
			return nil
		}
	}

	if len(req.SettlementExt) > 0 {
		if err := util.CheckSettlementExt(req.SettlementExt); err != nil {
			return nil
		}
	}

	if len(req.SellerMerchantId) > 0 {
		if err := util.CheckSellerMerchantId(req.SellerMerchantId); err != nil {
			return nil
		}
	}

	if len(req.RoyaltyParameters) > 0 {
		if err := util.CheckRoyaltyParameters(req.RoyaltyParameters); err != nil {
			return nil
		}
	}

	if len(req.TransCode) > 0 {
		if err := util.CheckTransCode(req.TransCode); err != nil {
			return nil
		}
	}

	if len(req.CashdeskExts) > 0 {
		// TODO 参数合法性校验
		req.bizContent.Set("exts", req.CashdeskExts)
	}

	if len(req.CashdeskLimitPay) > 0 {
		// TODO 参数合法性校验
		req.bizContent.Set("limit_pay", req.CashdeskLimitPay)
	}

	if len(req.CashdeskTradeType) > 0 {
		// TODO 参数合法性校验
		req.bizContent.Set("trade_type", req.CashdeskTradeType)
	}

	return nil
}

func (req *TradeCreateRequest) checkParamsFor1_0() error {
	// 以下六个含默认值参数
	// 为空则补全默认值，否则信任商户的设定
	// TODO: 考虑下赋默认值还是商户传
	if req.SignType == "" {
		req.SignType = "MD5"
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

	if err := util.CheckCashDeskTradeType(req.CashdeskTradeType); err != nil {
		return err
	}

	// 其他参数，正则检验
	if err := util.CheckAppId(req.AppId); err != nil {
		return err
	}

	if err := util.CheckMerchantId(req.MerchantId); err != nil {
		return err
	}

	// OutOrderNo在此接口内必选，需要额外确认不为空
	if req.OutOrderNo == "" {
		return errors.New("OutOrderNo shouldn't be blank when create a new trade")
	}

	if err := util.CheckOutOrderNo(req.OutOrderNo); err != nil {
		return err
	}

	if err := util.CheckUid(req.Uid); err != nil {
		return err
	}

	// 可选参数，确认不空后正则化查验
	if len(req.UidType) > 0 {
		if err := util.CheckUidType(req.UidType); err != nil {
			return err
		}
	}

	if err := util.CheckTotalAmountFor1_0(req.TotalAmount); err != nil {
		return err
	}

	if err := util.CheckCurrency(req.Currency); err != nil {
		return err
	}

	if err := util.CheckSubject(req.Subject); err != nil {
		return err
	}

	if err := util.CheckBody(req.Body); err != nil {
		return err
	}

	// 可选
	if len(req.ProductCode) > 0 {
		if err := util.CheckProductCode(req.ProductCode); err != nil {
			return err
		}
	}

	// 可选
	if len(req.PaymentType) > 0 {
		if err := util.CheckPaymentType(req.PaymentType); err != nil {
			return err
		}
	}

	if err := util.CheckTradeTime(req.TradeTime); err != nil {
		return err
	}

	if err := util.CheckValidTime(req.ValidTime); err != nil {
		return err
	}

	if err := util.CheckNotifyUrl(req.NotifyUrl); err != nil {
		return err
	}

	// 可选
	if len(req.ReturnUrl) > 0 {
		if err := util.CheckReturnUrl(req.ReturnUrl); err != nil {
			return nil
		}
	}

	// 可选
	if len(req.ServiceFee) > 0 {
		if err := util.CheckServiceFee(req.ServiceFee); err != nil {
			return nil
		}
	}

	if err := util.CheckRiskInfo(req.RiskInfo); err != nil {
		return err
	}

	if len(req.SettlementProductCode) > 0 {
		if err := util.CheckSettlementProductCode(req.SettlementProductCode); err != nil {
			return nil
		}
	}

	if len(req.SettlementExt) > 0 {
		if err := util.CheckSettlementExt(req.SettlementExt); err != nil {
			return nil
		}
	}

	if len(req.SellerMerchantId) > 0 {
		if err := util.CheckSellerMerchantId(req.SellerMerchantId); err != nil {
			return nil
		}
	}

	if len(req.RoyaltyParameters) > 0 {
		if err := util.CheckRoyaltyParameters(req.RoyaltyParameters); err != nil {
			return nil
		}
	}

	if len(req.TransCode) > 0 {
		if err := util.CheckTransCode(req.TransCode); err != nil {
			return nil
		}
	}

	if len(req.CashdeskExts) > 0 {
		// TODO 参数合法性校验
		req.bizContent.Set("exts", req.CashdeskExts)
	}

	if len(req.CashdeskLimitPay) > 0 {
		// TODO 参数合法性校验
		req.bizContent.Set("limit_pay", req.CashdeskLimitPay)
	}

	if len(req.CashdeskTradeType) > 0 {
		// TODO 参数合法性校验
		req.bizContent.Set("trade_type", req.CashdeskTradeType)
	}

	return nil
}


