package collect

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

// SendThirdCollect 发送到对应的三方平台(目前只有paytme)--代收
func SendThirdCollect(ctx context.Context, cfg *config.Config, redisPoolManager *database.RedisPoolManager, collectOrderData allStruct.RedisCollectOrderDataStruct) allStruct.CreateCollectOrderResp {
	var createOrderRsp allStruct.CreateCollectOrderResp
	switch collectOrderData.Platform {
	case "PAYTME":
		createOrderRsp = third.PayTmePayIn(ctx, cfg, redisPoolManager, collectOrderData)
		return createOrderRsp
	case "SBI":
		createOrderRsp = third.PayTmePayIn(ctx, cfg, redisPoolManager, collectOrderData)
		return createOrderRsp
	}
	return createOrderRsp
}

// GetCollectOrderID 生成代收订单 ID 并保存数据
func GetCollectOrderID(ctx context.Context, cfg *config.Config, redisPoolManager *database.RedisPoolManager, order allStruct.CollectOrderStruct, mid string) (allStruct.RedisCollectOrderDataStruct, error) {
	LocalStatus := cfg.OrderStatus.LocalStatus
	var collectData allStruct.RedisCollectOrderDataStruct
	orderID, err := redisPoolManager.GetUniqueID(ctx)
	if err != nil {
		return collectData, err
	}
	// 还缺个校验通道，然后分配（这里暂时默认PAYTME）
	collectData.Platform = "PAYTME"
	collectData.Amount = order.Amount
	collectData.OrderID = orderID
	collectData.MerchantNumber = mid
	collectData.CreateTime = time.Now().Unix()
	collectData.Status = LocalStatus
	collectData.CustomerName = order.CustomerName
	collectData.CustomerEmail = order.CustomerEmail
	collectData.CustomerPhone = order.CustomerPhone
	collectData.MerchantOrderID = order.MerchantOrderID

	orderJSON, err := json.Marshal(collectData)
	if err != nil {
		return collectData, err
	}
	// 代收订单数据
	err = redisPoolManager.SetValue(ctx, fmt.Sprintf("star-pay:collect:%s", orderID), string(orderJSON))
	if err != nil {
		return collectData, err
	}
	// 设置当前代收订单查询的键
	err = redisPoolManager.SetSAddValue(ctx, "star-pay:collect_id", orderID)
	if err != nil {
		return collectData, err
	}

	// 打印接收到的数据
	log.Printf("Create collect order data: %+v", string(orderJSON))
	return collectData, nil
}
