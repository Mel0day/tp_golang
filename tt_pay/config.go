package tt_pay

type Config struct {
	AppId      string
	AppSecret  string
	MerchantId string
	TPDomain   string // 请求支付域名 加http或者https前缀，比如：https://tp-pay.snssdk.com
	TPClientTimeoutMs int
}
