package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LedbetterBlog/public-msg/allStruct"
	"github.com/LedbetterBlog/public-msg/config"
	"github.com/LedbetterBlog/public-msg/database"
	"github.com/LedbetterBlog/public-msg/third"
	"log"
	"time"
)

// SendThirdPayment 发送到对应的三方平台(目前只有paytme)--代付
func SendThirdPayment(ctx context.Context, redisPoolManager *database.RedisPoolManager, orderData allStruct.RedisPaymentOrderDataStruct) {
	// 发送到对应的三方平台(目前只有paytme)
	if orderData.Platform == "PAYTME" {
		third.PayTmePayOut(ctx, redisPoolManager, orderData)
	}
	switch orderData.Platform {
	case "PAYTME":
		third.PayTmePayOut(ctx, redisPoolManager, orderData)
		break
	}

}

// GetPaymentOrderID 生成代付订单 ID 并保存数据
func GetPaymentOrderID(ctx context.Context, cfg *config.Config, redisPoolManager *database.RedisPoolManager, order allStruct.PaymentOrderStruct, mid string) (string, error) {
	LocalStatus := cfg.OrderStatus.LocalStatus
	var paymentData allStruct.RedisPaymentOrderDataStruct
	orderID, err := redisPoolManager.GetUniqueID(ctx)
	if err != nil {
		return "", err
	}
	paymentData.OrderID = orderID
	paymentData.MerchantNumber = mid
	paymentData.CreateTime = time.Now().Unix()
	paymentData.Platform = "PAYTME"
	paymentData.Status = LocalStatus
	paymentData.Amount = order.Amount
	paymentData.BeneName = order.BeneName
	paymentData.BeneEmail = order.BeneEmail
	paymentData.BenePhone = order.BenePhone
	paymentData.BeneIFSC = order.BeneIFSC
	paymentData.BeneBankAcct = order.BeneBankAcct
	paymentData.BeneAddress = order.BeneAddress
	paymentData.MerchantOrderID = order.MerchantOrderID

	orderJSON, err := json.Marshal(paymentData)
	if err != nil {
		return "", err
	}
	// 代付订单数据
	err = redisPoolManager.SetValue(ctx, fmt.Sprintf("star-pay:payment:%s", orderID), string(orderJSON))
	if err != nil {
		return "", err
	}
	// 当前代付订单的键
	err = redisPoolManager.SetSAddValue(ctx, "star-pay:payment_id", orderID)
	if err != nil {
		return "", err
	}

	// 打印接收到的数据
	log.Printf("Create order data: %+v", string(orderJSON))
	return orderID, nil
}
