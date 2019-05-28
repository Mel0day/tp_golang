package main

import (
	"code.byted.org/caijing_pay/tp_server_sdk/tp_golang/tt_pay"
	"code.byted.org/caijing_pay/tp_server_sdk/tp_golang/tt_pay/util"
	"context"
	"fmt"
)

func main() {
	util.SetDebugMode(true)

	p := "event_code=trade_status_sync&extension=%7B%22channel_code%22%3A%22LL_fangxinjie_1977%22%2C%22channel_merchant_order_no%22%3A%2211905210002203466466%22%2C%22channel_order_no%22%3A%222019052172521793%22%7D&merchant_id=fangxinjie&notify_id=SP2019052116131126729886806466&out_order_no=6693390424557490701&pay_channel=ANY_AGREEMENT_PAY&pay_time=1558426402&pay_type=ANY_AGREEMENT_PAY&real_amount=802133&sign=sWPeQGWvrooewI01IQsK4aWBe4OkIhzNKLRB7cq3sXoKX%2FTEduFqAPiuIpSpZDnyNQxREhYufXi8VyPxSgvc0poZ8m%2FixDON3XfmjSHF43TDbFTD%2BbJK0oLNcCcQ%2FhwjM%2BW8dt%2B9es%2BTM3iUv5SrNClkXdsDaxqVxtbOpS1kDS0%3D&sign_type=RSA&status=SUCCESS&total_amount=802133&trade_msg=success&trade_no=SP2019052116131126729886806466&trade_status=TRADE_SUCCESS"

	req := new(tt_pay.TradeNotifyRequest)
	req.SetParam(p)

	ctx := context.Background()
	resp, err := tt_pay.TradeNotify(ctx, req)

	if err != nil {
		fmt.Printf("Request failed! \nRequest:\n  [%#v]\nError:\n  [%s]\n", req, err)
		return
	}
	fmt.Printf("Request succeeded! \nRequest:\n [%#v]\nResponse:\n  [%#v]\n", req, resp)
}