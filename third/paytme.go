package third

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/LedbetterBlog/public-msg/allStruct"
	"github.com/LedbetterBlog/public-msg/config"
	"github.com/LedbetterBlog/public-msg/database"
	"go.mongodb.org/mongo-driver/bson"
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

// PayInStatus 代收状态请求
func (p *PayTme) PayInStatus(ctx context.Context, nowOrder allStruct.RedisCollectOrderDataStruct) (allStruct.PayTmeOrderStatus, error) {
	baseURL := "https://apis.paytme.com/v1/merchant/payin/"
	transactionId := nowOrder.PlatformOrderId // 例如，你的 transactionId
	// 构造完整的 URL
	fullURL := fmt.Sprintf("%s%s", baseURL, transactionId)
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return allStruct.PayTmeOrderStatus{RespMsg: err.Error(), Code: 400}, err
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
		return allStruct.PayTmeOrderStatus{RespMsg: err.Error(), Code: 400}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return allStruct.PayTmeOrderStatus{RespMsg: err.Error(), Code: 400}, err
	}

	var payTmePayInStatusResponse allStruct.PayTmePayInStatusResponse
	err = json.Unmarshal(body, &payTmePayInStatusResponse)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return allStruct.PayTmeOrderStatus{RespMsg: err.Error(), Code: 400}, err
	}

	if payTmePayInStatusResponse.Code == 200 {
		switch payTmePayInStatusResponse.Data.Status {
		case "success":
			return allStruct.PayTmeOrderStatus{
				PlatformOrderId: payTmePayInStatusResponse.Data.MerchantId,
				RespMsg:         "order status is success",
				Code:            payTmePayInStatusResponse.Code,
				Utr:             payTmePayInStatusResponse.Data.Utr,
				Amount:          payTmePayInStatusResponse.Data.Amount,
				Status:          payTmePayInStatusResponse.Data.Status,
			}, nil
		case "failed":
			return allStruct.PayTmeOrderStatus{
				PlatformOrderId: payTmePayInStatusResponse.Data.MerchantId,
				RespMsg:         "order status is failed",
				Code:            payTmePayInStatusResponse.Code,
				Utr:             payTmePayInStatusResponse.Data.Utr,
				Amount:          payTmePayInStatusResponse.Data.Amount,
				Status:          payTmePayInStatusResponse.Data.Status,
			}, nil
		default:
			return allStruct.PayTmeOrderStatus{
				PlatformOrderId: payTmePayInStatusResponse.Data.MerchantId,
				RespMsg:         "order status is not found",
				Code:            payTmePayInStatusResponse.Code,
				Utr:             payTmePayInStatusResponse.Data.Utr,
				Amount:          payTmePayInStatusResponse.Data.Amount,
				Status:          payTmePayInStatusResponse.Data.Status,
			}, nil
		}
	}
	return allStruct.PayTmeOrderStatus{Code: 400}, err
}

// PayOutStatus 代付状态请求
func (p *PayTme) PayOutStatus(ctx context.Context, nowOrder allStruct.RedisPaymentOrderDataStruct) (allStruct.PayTmeOrderStatus, error) {
	baseURL := "https://apis.paytme.com/v1/payout/status/"
	transactionId := nowOrder.PlatformOrderId // 例如，你的 transactionId
	// 构造完整的 URL
	fullURL := fmt.Sprintf("%s%s", baseURL, transactionId)
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return allStruct.PayTmeOrderStatus{RespMsg: err.Error(), Code: 400}, err
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
		return allStruct.PayTmeOrderStatus{RespMsg: err.Error(), Code: 400}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return allStruct.PayTmeOrderStatus{RespMsg: err.Error(), Code: 400}, err
	}

	var payTmePayInStatusResponse allStruct.PayTmePayOutStatusResponse
	err = json.Unmarshal(body, &payTmePayInStatusResponse)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return allStruct.PayTmeOrderStatus{RespMsg: err.Error(), Code: 400}, err
	}

	if payTmePayInStatusResponse.Status == "200" {
		switch payTmePayInStatusResponse.Data.Status {
		case "success":
			return allStruct.PayTmeOrderStatus{
				PlatformOrderId: payTmePayInStatusResponse.Data.MerchantId,
				RespMsg:         "payment order status is success",
				Utr:             payTmePayInStatusResponse.Data.Utr,
				Amount:          payTmePayInStatusResponse.Data.Amount,
				Status:          payTmePayInStatusResponse.Data.Status,
			}, nil
		case "failed":
			return allStruct.PayTmeOrderStatus{
				PlatformOrderId: payTmePayInStatusResponse.Data.MerchantId,
				RespMsg:         "payment order status is fail",
				Utr:             payTmePayInStatusResponse.Data.Utr,
				Amount:          payTmePayInStatusResponse.Data.Amount,
				Status:          payTmePayInStatusResponse.Data.Status,
			}, nil
		default:
			return allStruct.PayTmeOrderStatus{
				PlatformOrderId: payTmePayInStatusResponse.Data.MerchantId,
				RespMsg:         "payment order status is not found",
				Utr:             payTmePayInStatusResponse.Data.Utr,
				Amount:          payTmePayInStatusResponse.Data.Amount,
				Status:          payTmePayInStatusResponse.Data.Status,
			}, nil
		}
	}
	return allStruct.PayTmeOrderStatus{Code: 400}, err
}

// PayTmePayIn payTme的代收逻辑
func PayTmePayIn(ctx context.Context, cfg *config.Config, redisPoolManager *database.RedisPoolManager, MongoDBPoolManager *database.MongoDBPoolManager, collectOrderData allStruct.RedisCollectOrderDataStruct) allStruct.CreateCollectOrderResp {
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
	// 数据序列化更新mongodb的payment_order_test表
	filter := bson.M{"_id": collectOrderData.OrderID} // 使用 email 作为过滤条件
	update := bson.M{"$set": bson.M{"platform_order_id": result.PlatformOrderId, "resp_msg": result.RespMsg,
		"update_time": time.Now().Unix(), "status": cfg.OrderStatus.CommitStatus}}
	_, err = MongoDBPoolManager.UpdateData("payment_order_test", filter, update)
	if err != nil {
		log.Printf("MongoDBPoolManager update collect status is err: %v", err)
	}
	log.Printf("Send collect orderID: %v-----Result: %+v", collectOrderData.OrderID, result)
	return createOrderRsp
}

// PayTmePayOut payTme的代付逻辑
func PayTmePayOut(ctx context.Context, cfg *config.Config, redisPoolManager *database.RedisPoolManager, MongoDBPoolManager *database.MongoDBPoolManager, PaymentOrderData allStruct.RedisPaymentOrderDataStruct) {
	payTME := PayTme{SecretKey: "3a791d70e5e82c436d0dc495516e2229"}
	result, err := payTME.PayOut(ctx, PaymentOrderData)
	if err != nil {
		log.Printf("Error processing payout: %v", err)
		return
	}
	// 如果有三方订单号，redis更新三方平台订单号和三方订单号
	starMember := "star-pay:payment:" + PaymentOrderData.OrderID
	// 更新了redis里的代付订单状态，和三方订单号
	PaymentOrderData.Status = cfg.OrderStatus.CommitStatus
	PaymentOrderData.PlatformOrderId = result.PlatformOrderId
	PaymentOrderData.RespMsg = result.RespMsg
	orderJSON, err := json.Marshal(PaymentOrderData)
	if err != nil {
		log.Printf("PaymentOrderData orderJSON is err: %v", err)
		return
	}
	err = redisPoolManager.SetValue(ctx, starMember, string(orderJSON))
	if err != nil {
		log.Printf("redisPoolManager PayTmePayOut SetValue is err: %v", err)
		return
	}
	// 再往star-pay:third_platform_id里面插入三方id与系统id对应(如果有三方订单id的情况下)
	if result.PlatformOrderId != "" {
		thirdPlatformId := "star-pay:third_platform_id:" + result.PlatformOrderId
		err = redisPoolManager.SetValue(ctx, thirdPlatformId, PaymentOrderData.OrderID)
	}

	// 数据序列化更新mongodb的payment_order_test表
	filter := bson.M{"_id": PaymentOrderData.OrderID} // 使用 email 作为过滤条件
	update := bson.M{"$set": bson.M{"platform_order_id": result.PlatformOrderId, "status": PaymentOrderData.Status,
		"resp_msg": result.RespMsg, "update_time": time.Now().Unix()}}
	_, err = MongoDBPoolManager.UpdateData("payment_order_test", filter, update)
	if err != nil {
		log.Printf("MongoDBPoolManager update PayTmePayOut orderJSON is err: %v", err)
		return
	}
	log.Printf("Payment send order id: %v-----Result: %+v", PaymentOrderData.OrderID, result)
}

// PayTmePayInStatus payTme的代收状态逻辑
func PayTmePayInStatus(ctx context.Context, collectOrderData allStruct.RedisCollectOrderDataStruct) allStruct.PayTmeOrderStatus {
	payTME := PayTme{SecretKey: "3a791d70e5e82c436d0dc495516e2229"}
	result, err := payTME.PayInStatus(ctx, collectOrderData)
	if err != nil {
		log.Printf("Error processing payout: %v", err)
		return result
	}
	log.Printf("Send collect orderID: %v-----Result: %+v", collectOrderData.OrderID, result)
	return result
}

// PayTmePayOutStatus payTme的代付状态逻辑
func PayTmePayOutStatus(ctx context.Context, collectOrderData allStruct.RedisPaymentOrderDataStruct) allStruct.PayTmeOrderStatus {
	payTME := PayTme{SecretKey: "3a791d70e5e82c436d0dc495516e2229"}
	result, err := payTME.PayOutStatus(ctx, collectOrderData)
	if err != nil {
		log.Printf("Error processing payout: %v", err)
		return result
	}
	log.Printf("Send collect orderID: %v-----Result: %+v", collectOrderData.OrderID, result)
	return result
}
