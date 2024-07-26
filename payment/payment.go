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
func SendThirdPayment(ctx context.Context, cfg *config.Config, redisPoolManager *database.RedisPoolManager, MongoDBPoolManager *database.MongoDBPoolManager, orderData allStruct.RedisPaymentOrderDataStruct) {
	// 发送到对应的三方平台(目前只有paytme)
	switch orderData.Platform {
	case "PAYTME":
		third.PayTmePayOut(ctx, cfg, redisPoolManager, MongoDBPoolManager, orderData)
	}

}

// GetPaymentOrderID 生成代付订单 ID 并保存数据
func GetPaymentOrderID(ctx context.Context, cfg *config.Config, redisPoolManager *database.RedisPoolManager, MongoDBPoolManager *database.MongoDBPoolManager, order allStruct.PaymentOrderStruct, mid string) (string, error) {
	LocalStatus := cfg.OrderStatus.LocalStatus
	WaitCallBackStatus := cfg.OrderStatus.WaitCallBackStatus
	CreateTime := time.Now().Unix()
	var paymentData allStruct.RedisPaymentOrderDataStruct
	orderID, err := redisPoolManager.GetUniqueID(ctx)
	if err != nil {
		return "", err
	}
	paymentData.OrderID = orderID
	paymentData.MerchantNumber = mid
	paymentData.CreateTime = CreateTime
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
	// 插入mongo(用来查是否重复订单)
	var mongoDbLocalStatusStruct allStruct.MongoDbLocalStatusStruct
	mongoDbLocalStatusStruct.OrderID = orderID
	mongoDbLocalStatusStruct.MerchantOrderID = order.MerchantOrderID
	mongoDbLocalStatusStruct.Status = LocalStatus
	mongoDbLocalStatusStruct.Amount = order.Amount
	mongoDbLocalStatusStruct.Platform = "PAYTME"
	mongoDbLocalStatusStruct.CreateTime = CreateTime
	mongoDbLocalStatusStruct.MerchantNumber = mid
	mongoDbLocalStatusStruct.CallbackStatus = WaitCallBackStatus
	mongoDbLocalStatusStruct.OrderType = 1
	_, err = MongoDBPoolManager.InsertData("payment_order_test", mongoDbLocalStatusStruct)
	if err != nil {
		log.Printf("MongoDBPoolManager insert paymentData is err: %v", err)
		return "", err
	}
	// 当前代付订单的键
	err = redisPoolManager.SetSAddValue(ctx, "star-pay:payment_id", orderID)
	if err != nil {
		return "", err
	}

	// 打印接收到的数据
	log.Printf("Create payment order data: %+v", string(orderJSON))

	return orderID, nil
}
