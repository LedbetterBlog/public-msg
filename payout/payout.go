package payout

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

// SendThirdPayout 发送到对应的三方平台(目前只有paytme)--代付
func SendThirdPayout(ctx context.Context, cfg *config.Config, redisPoolManager *database.RedisPoolManager, MongoDBPoolManager *database.MongoDBPoolManager, orderData allStruct.RedisPayOutOrderDataStruct) {
	// 发送到对应的三方平台(目前只有paytme)
	switch orderData.Platform {
	case "PAYTME":
		third.PayTmePayOut(ctx, cfg, redisPoolManager, MongoDBPoolManager, orderData)
	}

}

// GetPayoutOrderID 生成代付订单 ID 并保存数据
func GetPayoutOrderID(ctx context.Context, cfg *config.Config, redisPoolManager *database.RedisPoolManager, MongoDBPoolManager *database.MongoDBPoolManager, order allStruct.PayOutOrderStruct, mid string) (string, error) {
	LocalStatus := cfg.OrderStatus.LocalStatus
	WaitCallBackStatus := cfg.OrderStatus.WaitCallBackStatus
	CreateTime := time.Now().Unix()
	var payoutData allStruct.RedisPayOutOrderDataStruct
	orderID, err := redisPoolManager.GetUniqueID(ctx)
	if err != nil {
		return "", err
	}
	payoutData.OrderID = orderID
	payoutData.MerchantNumber = mid
	payoutData.CreateTime = CreateTime
	payoutData.Platform = "PAYTME"
	payoutData.Status = LocalStatus
	payoutData.Amount = order.Amount
	payoutData.UserName = order.BeneName
	payoutData.UserEmail = order.BeneEmail
	payoutData.UserPhone = order.BenePhone
	payoutData.BankIFSC = order.BeneIFSC
	payoutData.UserBankAcct = order.BeneBankAcct
	payoutData.UserAddress = order.BeneAddress
	payoutData.MerchantOrderID = order.MerchantOrderID
	orderJSON, err := json.Marshal(payoutData)
	if err != nil {
		return "", err
	}
	// 代付订单数据
	err = redisPoolManager.SetValue(ctx, fmt.Sprintf("star-pay:payout:%s", orderID), string(orderJSON))
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
	mongoDbLocalStatusStruct.UserName = order.BeneName
	mongoDbLocalStatusStruct.UserEmail = order.BeneEmail
	mongoDbLocalStatusStruct.UserPhone = order.BenePhone
	mongoDbLocalStatusStruct.BankIFSC = order.BeneIFSC
	mongoDbLocalStatusStruct.UserBankAcct = order.BeneBankAcct
	mongoDbLocalStatusStruct.UserAddress = order.BeneAddress
	mongoDbLocalStatusStruct.Chnl = 0
	mongoDbLocalStatusStruct.ChnlFeeRatio = 0.0177
	mongoDbLocalStatusStruct.MchFeeRatio = 0.022
	RateAmount := float64(order.Amount) * mongoDbLocalStatusStruct.MchFeeRatio
	mongoDbLocalStatusStruct.MchSettleAmount = float64(order.Amount) + RateAmount
	_, err = MongoDBPoolManager.InsertData(ctx, "payment_order_test", mongoDbLocalStatusStruct)
	if err != nil {
		log.Printf("MongoDBPoolManager insert paymentData is err: %v", err)
		return "", err
	}
	// 当前代付订单的键
	err = redisPoolManager.SetSAddValue(ctx, "star-pay:payout_id", orderID)
	if err != nil {
		return "", err
	}

	// 打印接收到的数据
	log.Printf("Create payout order data: %+v", string(orderJSON))

	return orderID, nil
}
