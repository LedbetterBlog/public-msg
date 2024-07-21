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

// Payout 发起支付请求
func (p *PayTme) Payout(ctx context.Context, nowOrder allStruct.PaymentData) (allStruct.PaymentRespData, error) {
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
		return allStruct.PaymentRespData{RspMsg: err.Error(), Code: 400}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return allStruct.PaymentRespData{RspMsg: err.Error(), Code: 400}, err
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
		return allStruct.PaymentRespData{RspMsg: err.Error(), Code: 400}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return allStruct.PaymentRespData{RspMsg: err.Error(), Code: 400}, err
	}

	var payoutResponse allStruct.PayTmePayoutResponse
	err = json.Unmarshal(body, &payoutResponse)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return allStruct.PaymentRespData{RspMsg: err.Error(), Code: 400}, err
	}

	if payoutResponse.Code == 200 {
		return allStruct.PaymentRespData{
			PlatformOrderId: payoutResponse.Data.ID,
			RspMsg:          "create payment success",
			Code:            payoutResponse.Code,
		}, nil
	} else {
		return allStruct.PaymentRespData{
			PlatformOrderId: "",
			RspMsg:          payoutResponse.Message,
			Code:            payoutResponse.Code,
		}, nil
	}
}
