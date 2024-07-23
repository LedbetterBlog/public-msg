package third

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/LedbetterBlog/public-msg/allStruct"
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
