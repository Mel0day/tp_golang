package main


import (
	"context"
	"fmt"
	"github.com/ttcjpay/tp_golang/tt_pay"
	"github.com/ttcjpay/tp_golang/tt_pay/util"
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
	req := tt_pay.NewRefundQueryRequest(conf)
	req.Uid = "123"
	req.OutRefundNo = "test12345"

	ctx := context.Background()

	resp, err := tt_pay.RefundQuery(ctx, req)

	if err != nil {
		fmt.Printf("Request failed! \nRequest:\n  [%#v]\nError:\n  [%s]\n", req, err)
		return
	}
	fmt.Printf("Request succeeded! \nRequest:\n [%#v]\nResponse:\n  [%#v]\n", req, resp)
}
