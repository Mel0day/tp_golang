# 财经支付
简要介绍

## 目录
- [安装](#安装)
- [依赖版本](#依赖版本)
- [快速接入](#快速接入)
- [API示例](#API示例)

## 安装
安装财经支付SDK前，请确保安装GO并正确设置环境变量以及WorkSpace。
下载SDK：
```go
go get -u github.com/ttcjpay/tp_golang/tt_pay
```
导入SDK：
```go
import "github.com/ttcjpay/tp_golang/tt_pay"
```
要设置Debug模式，还需导入：
```go
import "github.com/ttcjpay/tp_golang/tt_pay/util"
```
## 依赖版本

## 快速接入
本部分以预下单接口为例，演示如何快速接入财经支付SDK。其他接口的接入方法与预下单接口类似，主要差别在于商户传入参数不同。具体请见：各接口demo https://github.com/ttcjpay/tp_golang/tree/master/demo,财经支付接入文档 http://lf6-ttcdn-tos.pstatp.com/obj/caijing-tp-cashdesk-bucket1/cashdesk_pay.pdf. 

#### 步骤一：设置
调用SetDebugMode可设置是否进入Debug模式，传入true打开debug模式，传输false关闭debug模式
```go
util.SetDebugMode(true)
```
#### 步骤二：创建预下单请求
使用SDK提供的New函数创建请求时，会自动设置部分参数为默认值，这些参数包括：
 - Version = "1.0"
 - SignType = "MD5"
 - CashdeskTradeType = "H5" (只在预下单时需要，默认为H5)
 - Config.TPDomain = "https://tp-pay.snssdk.com"
 
要自定义参数，请仔细查阅接入手册https://tp-pay.snssdk.com/vip-develop/down.html
```go
conf := tt_pay.Config{
    AppId:"__________________", // 支付方分配给业务方的ID，用于获取 签名/验签 的密钥信息
    AppSecret:"______________", // 支付方密钥
    MerchantId:"_____________", // 支付方分配给业务方的商户编号
    TPDomain:"https://tp-pay-test.snssdk.com",
    TPClientTimeoutMs:6000,
}
req := tt_pay.NewTradeCreateRequest(conf)
req.OutOrderNo = fmt.Sprintf("%d", time.Now().Unix())
req.Uid = "123"
req.TotalAmount = 1
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
```
#### 步骤三：创建context
```go
ctx := context.Background()
```
#### 步骤四：调用接口，处理业务逻辑
```go
// 调用预下单接口，获取下单凭据
resp, err := tt_pay.TradeCreate(ctx, req)
// 你的业务逻辑
DoSomething()
```

## API示例
完整的各API使用Demo可以在 https://github.com/ttcjpay/tp_golang/tree/master/demo 找到。
完整的接入文档可在 http://lf6-ttcdn-tos.pstatp.com/obj/caijing-tp-cashdesk-bucket1/cashdesk_pay.pdf 找到。
注意：为使接口正常工作，请务必按照开发文档要求正确传入所有参数，漏传、多传以及不正确的参数格式将使接口失效！

// TODO: 贴每个demo的链接
### 预下单接口 
```go
// 预下单接口
resp, err := tt_pay.TradeCerate(ctx, req)
```
### 订单查询接口
```go
// 订单查询接口
resp, err := tt_pay.TradeQuery(ctx, req)
```
### 退款申请接口

```go
// 退款申请接口
resp, err := tt_pay.RefundCreate(ctx, req)
```

### 退款查询接口
```go
// 退款查询接口
resp, err := tt_pay.RefundQuery(ctx, req)
```

### 支付回调接口
```go
// 支付回调接口
resp, err := tt_pay.TradeNotify(ctx, req)
```

### 退款回调接口
```go
// 退款回调接口
resp, err := tt_pay.RefundNotify(ctx, req)
```

