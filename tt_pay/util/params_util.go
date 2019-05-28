package util

import (
	"errors"
	"regexp"
)

// 检查各种URL统一使用urlRegExp
// 所有在某些接口可选，在其他接口必选的参数均采用使用可选正则匹配（空字符串可通过）

const (
	// Regular expressions
	// 所有URL类参数都用RegExp_Url
	RegExp_AppId = "^[0-9a-zA-Z-_]{1,32}$"
	RegExp_MerchantId = "^[0-9a-zA-Z-_]{1,32}$"
	RegExp_Uid	= "^[0-9a-zA-Z-_]{1,32}$"
	RegExp_SignType = "^[0-9a-zA-Z]{1,10}$"
	RegExp_Version = "^[0-9].[0-9]$"
	RegExp_TimeStamp = "^[0-9]{1,19}$"
	RegExp_OutRefundNo = "^$|^[a-zA-Z0-9-_]{1,32}$"
	RegExp_RefundNo = "^$|^[a-zA-Z0-9-_]{1,64}$"
	RegExp_Url = `^(.*):(.*)$|(https|http):\/\/[-A-Za-z0-9+&@#\/%?=~_|!:,.;]+[-A-Za-z0-9+&@#\/%=~_|]`
	RegExp_RiskInfo = `^{(".*":.*,)*".*":.*}$`
	RegExp_OutOrderNo = `^$|^[0-9a-zA-Z-_]{1,32}$`
	RegExp_TradeNo = `^$|^[a-zA-Z0-9-_]{1,64}$`
	RegExp_TotalAmount = `^[1-9][0-9]*$`
)

var (
	appIdRegexp       *regexp.Regexp
	merchantIdRegexp  *regexp.Regexp
	uidRegExp         *regexp.Regexp
	signTypeRegExp    *regexp.Regexp
	versionRegExp     *regexp.Regexp
	timeStampRegExp   *regexp.Regexp
	outRefundNoRegExp *regexp.Regexp
	refundNoRegExp *regexp.Regexp
	urlRegExp *regexp.Regexp
	riskInfoRegExp *regexp.Regexp
	OutOrderNoRegExp *regexp.Regexp
	TradeNoRegExp *regexp.Regexp
	TotalAmountRegExp *regexp.Regexp
)

func init() {
	appIdRegexp = regexp.MustCompile(RegExp_AppId)
	merchantIdRegexp = regexp.MustCompile(RegExp_MerchantId)
	uidRegExp = regexp.MustCompile(RegExp_Uid)
	signTypeRegExp = regexp.MustCompile(RegExp_SignType)
	versionRegExp = regexp.MustCompile(RegExp_Version)
	timeStampRegExp = regexp.MustCompile(RegExp_TimeStamp)
	outRefundNoRegExp = regexp.MustCompile(RegExp_OutRefundNo)
	refundNoRegExp = regexp.MustCompile(RegExp_RefundNo)
	urlRegExp = regexp.MustCompile(RegExp_Url)
	riskInfoRegExp = regexp.MustCompile(RegExp_RiskInfo)
	OutOrderNoRegExp = regexp.MustCompile(RegExp_OutOrderNo)
	TradeNoRegExp = regexp.MustCompile(RegExp_TradeNo)
	TotalAmountRegExp = regexp.MustCompile(RegExp_TotalAmount)
}

// CheckAppId 检查AppId是否有效
func CheckAppId(appId string) error {
	isMatch := appIdRegexp.MatchString(appId)
	if !isMatch {
		return errors.New("invalid param: AppId")
	}
	return nil
}

func CheckMerchantId(merchantId string) error {
	isMatch := merchantIdRegexp.MatchString(merchantId)
	if !isMatch {
		return errors.New("invalid param: MerchantId")
	}
	return nil
}

func CheckAppSecret(appSecret string) error {
	if len(appSecret) == 0 {
		return errors.New("invalid param appSecret")
	}
	return nil
}

func CheckUid(uid string) error {
	isMatch := uidRegExp.MatchString(uid)
	if !isMatch {
		return errors.New("invalid param: Uid")
	}
	return nil
}

func CheckSignType(signType string) error {
	isMatch := signTypeRegExp.MatchString(signType)
	if !isMatch {
		return errors.New("invalid param: SignType")
	}
	return nil
}

func CheckVersion(version string) error {
	isMatch := versionRegExp.MatchString(version)
	if !isMatch {
		return errors.New("invalid param: Version")
	}
	return nil
}

func CheckTimeStamp(timeStamp string) error {
	isMatch := timeStampRegExp.MatchString(timeStamp)
	if !isMatch {
		return errors.New("invalid param: TimeStamp")
	}
	return nil
}

func CheckOutRefundNo(outRefundNo string ) error {
	isMatch := outRefundNoRegExp.MatchString(outRefundNo)
	if !isMatch {
		return errors.New("invalid param: OutRefundNo")
	}
	return nil
}

func CheckRefundNo(refundNo string) error {
	isMatch := refundNoRegExp.MatchString(refundNo)
	if !isMatch {
		return errors.New("invalid param: RefundNo")
	}
	return nil
}

func CheckTPDomain(tpDomain string ) error {
	isMatch := urlRegExp.MatchString(tpDomain)
	if !isMatch {
		return errors.New("invalid param: TPDomain")
	}
	return nil
}

func CheckNotifyUrl(url string) error {
	isMatch := urlRegExp.MatchString(url)
	if !isMatch {
		return errors.New("invalid param: NotifyUrl")
	}
	return nil
}

func CheckRefundAmount(num int) error {
	if num <= 0 {
		return errors.New("invalid param: RefundAmount")
	}
	return nil
}

// TODO: complete this part
func CheckRiskInfo(riskInfo string) error {
	// isMatch := riskInfoRegExp.MatchString(riskInfo)
	// if !isMatch {
	// 	return errors.New("invalid param: RiskInfo")
	// }
	return nil
}

func CheckOutOrderNo(outOrderNo string ) error {
	isMatch := OutOrderNoRegExp.MatchString(outOrderNo)
	if !isMatch {
		return errors.New("invalid param: OutOrderNo")
	}
	return nil
}

// CheckUidType 检查uid_type的合法性
func CheckUidType(uidType string) error {
	if len(uidType) == 0 {
		return errors.New("invalid param uidType")
	}
	return nil
}

func CheckTradeNo(tradeNo string ) error {
	isMatch := TradeNoRegExp.MatchString(tradeNo)
	if !isMatch {
		return errors.New("invalid param: TradeNo")
	}
	return nil
}

// CheckTotalAmount 检查totalAmount是否合法
func CheckTotalAmountFor1_0(totalAmount interface{}) error {
	if _, ok := totalAmount.(int); !ok {
		return errors.New("invalid param: TotalAmount: TotalAmount should be an int in version1.0")
	}
	if totalAmount.(int) <= 0 {
		return errors.New("invalid param: TotalAmount")
	}
	return nil
}

func CheckTotalAmountFor2_0(totalAmount interface{}) error {
	if _, ok := totalAmount.(string); !ok {
		return errors.New("invalid param: TotalAmount: TotalAmount should be a string in version2.0")
	}
	isMatch := TotalAmountRegExp.MatchString(totalAmount.(string))
	if !isMatch {
		return errors.New("invalid param: TotalAmount: TotalAmount should not has leading zero")
	}
	return nil
}

// CheckCurrency 检查币种是否合法
func CheckCurrency(currency string) error {
	if len(currency) == 0 {
		return errors.New("invalid param: Currency")
	}
	return nil
}

// CheckSubject 检查subject是否合法
func CheckSubject(subject string) error {
	if len(subject) == 0 {
		return errors.New("invalid param: Subject")
	}
	return nil
}

// CheckBody 检查body合法性
func CheckBody(body string) error {
	if len(body) == 0 {
		return errors.New("invalid param: Body")
	}
	return nil
}

// CheckProductCode 检查productCode是否合法
func CheckProductCode(productCode string) error {
	if len(productCode) == 0 {
		return errors.New("invalid param: ProductCode")
	}
	return nil
}

// CheckPaymentType 检查paymentType是否合法
func CheckPaymentType(paymentType string) error {
	if len(paymentType) == 0 {
		return errors.New("invalid param： PaymentType")
	}
	return nil
}

// CheckTradeTime 检查tradeTime是否合法
func CheckTradeTime(tradeTime string) error {
	if len(tradeTime) == 0 {
		return errors.New("invalid param: TradeTime")
	}
	return nil
}

// CheckValidTime 检查validTime是否合法
func CheckValidTime(validTime string) error {
	if len(validTime) == 0 {
		return errors.New("invalid param: validTime")
	}
	return nil
}

// CheckReturnUrl  检查returnUrl是否合法
func CheckReturnUrl(returnUrl string) error {
	if len(returnUrl) == 0 {
		return errors.New("invalid param: ReturnUrl")
	}
	return nil
}

// CheckServiceFee 检查serviceFee是否合法
func CheckServiceFee(serviceFee string) error {
	if len(serviceFee) == 0 {
		return errors.New("invalid param: ServiceFee")
	}
	return nil
}

// CheckSettlementProductCode 检查settlementProductCode是否合法
func CheckSettlementProductCode(settlementProductCode string) error {
	if len(settlementProductCode) == 0 {
		return errors.New("invalid param SettlementProductCode")
	}
	return nil
}

// CheckSettlementExt 检查settlementExt是否合法
func CheckSettlementExt(settlementExt string) error {
	if len(settlementExt) == 0 {
		return errors.New("invalid param SettlementExt")
	}
	return nil
}

// CheckSellerMerchantId 检查sellerMerchantId是否合法
func CheckSellerMerchantId(sellerMerchantId string) error {
	if len(sellerMerchantId) == 0 {
		return errors.New("invalid param SellerMerchantId")
	}
	return nil
}

// CheckRoyaltyParameters 检查toyaltyParameters是否合法
func CheckRoyaltyParameters(toyaltyParameters string) error {
	if len(toyaltyParameters) == 0 {
		return errors.New("invalid param RoyaltyParameters")
	}
	return nil
}

// CheckTransCode 检查transCode是否合法
func CheckTransCode(transCode string) error {
	if len(transCode) == 0 {
		return errors.New("invalid param TransCode")
	}
	return nil
}

func CheckCashDeskTradeType(tradeType string) error {
	if len(tradeType) == 0 {
		return errors.New("invalid param: CashDeskTradeType")
	}
	return nil
}