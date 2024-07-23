package third

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/LedbetterBlog/public-msg/allStruct"
	"github.com/LedbetterBlog/public-msg/config"
	"github.com/LedbetterBlog/public-msg/database"
	"io"
	"log"
	"net/http"
	"time"
)

// PayTme 处理支付的结构体(这个是定义方法，无法放到allStruct)
type PayTme struct {
	SecretKey string
}

// PayIn 发起收款请求
func (p *PayTme) PayIn(ctx context.Context, nowOrder allStruct.RedisCollectOrderDataStruct) (allStruct.PayTmeCollectRespData, error) {
	url := "https://apis.paytme.com/v1/merchant/payin/direct-payin"

	payload := allStruct.PayTmePayInRequest{
		Amount:                float64(nowOrder.Amount / 100),
		Name:                  nowOrder.CustomerName,
		Email:                 nowOrder.CustomerEmail,
		UserContactNumber:     nowOrder.CustomerPhone,
		MerchantTransactionID: nowOrder.MerchantOrderID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return allStruct.PayTmeCollectRespData{RespMsg: err.Error(), Code: 400}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return allStruct.PayTmeCollectRespData{RespMsg: err.Error(), Code: 400}, err
	}

	req.Header.Set("x-api-key", p.SecretKey)
	req.Header.Set("Content-Type", "application/json")

	// 创建一个 http.Client，并设置超时时间
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置超时时间为10秒
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return allStruct.PayTmeCollectRespData{RespMsg: err.Error(), Code: 400}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return allStruct.PayTmeCollectRespData{RespMsg: err.Error(), Code: 400}, err
	}

	var payInResponse allStruct.PayTmePayInRespData
	err = json.Unmarshal(body, &payInResponse)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return allStruct.PayTmeCollectRespData{RespMsg: err.Error(), Code: 400}, err
	}

	if payInResponse.Code == 200 {
		return allStruct.PayTmeCollectRespData{
			PlatformOrderId: payInResponse.Data.PlatformOrderId,
			UPI:             payInResponse.Data.UpiUrl,
			RespMsg:         "create payment success",
			Code:            payInResponse.Code,
		}, nil
	} else {
		return allStruct.PayTmeCollectRespData{
			PlatformOrderId: "",
			UPI:             "",
			RespMsg:         payInResponse.Message,
			Code:            payInResponse.Code,
		}, nil
	}
}

// PayOut 发起支付请求
func (p *PayTme) PayOut(ctx context.Context, nowOrder allStruct.RedisPaymentOrderDataStruct) (allStruct.PayTmePaymentRespData, error) {
	url := "https://apis.paytme.com/v1/payout"

	payload := allStruct.PayTmePayoutRequest{
		Amount:                float64(nowOrder.Amount / 100),
		Name:                  nowOrder.BeneName,
		Email:                 nowOrder.BeneEmail,
		Phone:                 nowOrder.BenePhone,
		AccountNumber:         nowOrder.BeneBankAcct,
		BankIfsc:              nowOrder.BeneIFSC,
		AccountHolderName:     nowOrder.BeneName,
		BankName:              nowOrder.BeneName,
		UPI:                   "",
		Purpose:               "",
		MerchantTransactionID: nowOrder.MerchantOrderID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return allStruct.PayTmePaymentRespData{RespMsg: err.Error(), Code: 400}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return allStruct.PayTmePaymentRespData{RespMsg: err.Error(), Code: 400}, err
	}

	req.Header.Set("x-api-key", p.SecretKey)
	req.Header.Set("Content-Type", "application/json")

	// 创建一个 http.Client，并设置超时时间
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置超时时间为10秒
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return allStruct.PayTmePaymentRespData{RespMsg: err.Error(), Code: 400}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return allStruct.PayTmePaymentRespData{RespMsg: err.Error(), Code: 400}, err
	}

	var payoutResponse allStruct.PayTmePayOutResponse
	err = json.Unmarshal(body, &payoutResponse)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return allStruct.PayTmePaymentRespData{RespMsg: err.Error(), Code: 400}, err
	}

	if payoutResponse.Code == 200 {
		return allStruct.PayTmePaymentRespData{
			PlatformOrderId: payoutResponse.Data.ID,
			RespMsg:         "create payment success",
			Code:            payoutResponse.Code,
		}, nil
	} else {
		return allStruct.PayTmePaymentRespData{
			PlatformOrderId: "",
			RespMsg:         payoutResponse.Message,
			Code:            payoutResponse.Code,
		}, nil
	}
}

// PayTmePayIn payTme的代收逻辑
func PayTmePayIn(ctx context.Context, cfg *config.Config, redisPoolManager *database.RedisPoolManager, collectOrderData allStruct.RedisCollectOrderDataStruct) allStruct.CreateCollectOrderResp {
	var createOrderRsp allStruct.CreateCollectOrderResp
	payTME := PayTme{SecretKey: "3a791d70e5e82c436d0dc495516e2229"}
	result, err := payTME.PayIn(ctx, collectOrderData)
	if err != nil {
		log.Printf("Error processing payout: %v", err)
		return createOrderRsp
	}
	// 创建订单返回信息给客户
	createOrderRsp.PlatformOrderId = result.PlatformOrderId
	createOrderRsp.MerchantOrderId = collectOrderData.MerchantOrderID
	createOrderRsp.Code = result.Code
	createOrderRsp.UpiLink = result.UPI
	createOrderRsp.PaymentLink = "https://www.ez-pays.in/cashier/cashier-page?order_id=" + collectOrderData.OrderID
	createOrderRsp.Message = result.RespMsg
	// 如果有三方订单号，redis更新三方平台订单号和三方订单号
	starMember := "star-pay:collect:" + collectOrderData.OrderID
	// 更新了redis里的代收订单状态，和三方订单号
	collectOrderData.Status = cfg.OrderStatus.CommitStatus
	collectOrderData.PlatformOrderId = result.PlatformOrderId
	collectOrderData.RespMsg = result.RespMsg

	orderJSON, err := json.Marshal(collectOrderData)
	if err != nil {
		log.Printf("PaymentOrderData orderJSON is err: %v", err)
		return createOrderRsp
	}
	// 插入对应的订单信息
	err = redisPoolManager.SetValue(ctx, starMember, string(orderJSON))
	if err != nil {
		log.Printf("redisPoolManager SetValue is err: %v", err)
		return createOrderRsp
	}
	// 再往star-pay:third_platform_id里面插入三方id与系统id对应(如果有三方订单id的情况下)
	if result.PlatformOrderId != "" {
		thirdPlatformId := "star-pay:third_platform_id:" + result.PlatformOrderId
		err = redisPoolManager.SetValue(ctx, thirdPlatformId, collectOrderData.OrderID)
	}
	if err != nil {
		log.Printf("redisPoolManager SetValue is err: %v", err)
		return createOrderRsp
	}
	log.Printf("Collect orderID: %v-----Result: %+v", collectOrderData.OrderID, result)
	return createOrderRsp
}

// PayTmePayOut payTme的代付逻辑
func PayTmePayOut(ctx context.Context, redisPoolManager *database.RedisPoolManager, PaymentOrderData allStruct.RedisPaymentOrderDataStruct) {
	payTME := PayTme{SecretKey: "3a791d70e5e82c436d0dc495516e2229"}
	result, err := payTME.PayOut(ctx, PaymentOrderData)
	if err != nil {
		log.Printf("Error processing payout: %v", err)
		return
	}
	// 如果有三方订单号，redis更新三方平台订单号和三方订单号
	starMember := "star-pay:payment:" + PaymentOrderData.OrderID
	// 更新了redis里的代付订单状态，和三方订单号
	PaymentOrderData.Status = 6
	PaymentOrderData.PlatformOrderId = result.PlatformOrderId
	PaymentOrderData.RespMsg = result.RespMsg
	orderJSON, err := json.Marshal(PaymentOrderData)
	if err != nil {
		log.Printf("PaymentOrderData orderJSON is err: %v", err)
		return
	}
	err = redisPoolManager.SetValue(ctx, starMember, string(orderJSON))
	if err != nil {
		log.Printf("redisPoolManager SetValue is err: %v", err)
		return
	}
	// 再往star-pay:third_platform_id里面插入三方id与系统id对应(如果有三方订单id的情况下)
	if result.PlatformOrderId != "" {
		thirdPlatformId := "star-pay:third_platform_id:" + result.PlatformOrderId
		err = redisPoolManager.SetValue(ctx, thirdPlatformId, PaymentOrderData.OrderID)
	}
	log.Printf("Payment send order id: %v-----Result: %+v", PaymentOrderData.OrderID, result)
}
