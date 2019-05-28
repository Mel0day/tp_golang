package main

import (
	"context"
	"fmt"
	"github.com/ttcjpay/tp_golang/tt_pay"
	"github.com/ttcjpay/tp_golang/tt_pay/util"
	"time"
)

func main() {
	util.SetDebugMode(true)

	conf := tt_pay.Config{
		AppId:"__________________", // 支付方分配给业务方的ID，用于获取 签名/验签 的密钥信息
		AppSecret:"______________", // 支付方密钥
		MerchantId:"_____________", // 支付方分配给业务方的商户编号
		TPDomain:"https://tp-pay-test.snssdk.com",
		TPClientTimeoutMs:6000,
	}
	req := tt_pay.NewTradeCreateRequest(conf)
	req.Version = "2.0"
	req.OutOrderNo = fmt.Sprintf("%d", time.Now().Unix())
	req.Uid = "123"
	// 注意：调用1.0版本时TotalAmount为int类型，2.0版本时为string类型
	req.TotalAmount = "1"
	req.Currency = "CNY"
	req.Subject = "测试订单"
	req.Body = "测试订单内容"
	req.TradeTime = fmt.Sprintf("%d", time.Now().Unix())
	req.ValidTime = "3600"
	req.NotifyUrl = "https://google.com"
	req.RiskInfo = `{"ip":"127.0.0.1", "device_id":"122333"}`
	// 支付方式（必填）：可选值：SDK|H5。
	// SDK：业务方App必须是头条主端App，或者具备头条主端支付SDK及ToutiaoJSBridge的能力
	// H5：业务方App不具备SDK支付的能力如果非法，默认为H5支付方式
	req.CashdeskTradeType = "SDK"

	ctx := context.Background()
	resp, err := tt_pay.TradeCreate(ctx, req)

	if err != nil {
		fmt.Printf("Request failed! \nRequest:\n  [%#v]\nError:\n  [%s]\n", req, err)
		return
	}
	cashDeskParams, _ := resp.GetCashdeskSdkParams()
	fmt.Printf("Request succeeded! \nRequest:\n [%#v]\nResponse:\n  [%#v]\nGet Cashdesk Parameters:\n  [%#v]", req, resp, cashDeskParams)
}
