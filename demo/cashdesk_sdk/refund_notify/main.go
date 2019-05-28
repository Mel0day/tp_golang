package main

import (
	"context"
	"fmt"
	"github.com/ttcjpay/tp_golang/tt_pay"
	"github.com/ttcjpay/tp_golang/tt_pay/util"
)

func main() {
	util.SetDebugMode(true)

	p := "app_id=800083521500&event_code=refund_status_sync&extension=%7B%22ChannelExt%22%3A%22%7B%5C%22channel_code%5C%22%3A%5C%22%5C%22%2C%5C%22channel_merchant_id%5C%22%3A%5C%22%5C%22%2C%5C%22channel_merchant_order_no%5C%22%3A%5C%2221905190003612575711%5C%22%2C%5C%22channel_order_no%5C%22%3A%5C%22%5C%22%7D%22%7D&merchant_id=temai&notify_id=201905192142581229312524&origin_trade_no=SP2019042413070723909545087715&out_refund_no=3341661567434507918&refund_amount=6900&refund_no=SR2019051921404724920765587715&refund_status=REFUND_SUCCESS&refund_time=1558273381&sign=xiEn4l349%2B4ndkun4eCUb8cVXYeuRPVm6NOrqjVG4AwS69xWL6jY5uXdg0xALWUls4b2u8KXeQ%2BOw3gyEn37I1wgDDXecWNTmQe90FOlUkRM0oSTZGpi0RrEKIdy18yrisbhmoSoY66ztWbpsSj3jhOZ%2FxhNvu9VMnJg9hR3%2BZw%3D&sign_type=RSA&status=SUCCESS"

	req := new(tt_pay.RefundNotifyRequest)
	req.SetParam(p)

	ctx := context.Background()
	resp, err := tt_pay.RefundNotify(ctx, req)

	if err != nil {
		fmt.Printf("Request failed! \nRequest:\n  [%#v]\nError:\n  [%s]\n", req, err)
		return
	}
	fmt.Printf("Request succeeded! \nRequest:\n [%#v]\nResponse:\n  [%#v]\n", req, resp)
}
