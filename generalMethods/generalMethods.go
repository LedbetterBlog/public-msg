package generalMethods

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/LedbetterBlog/public-msg/allStruct"
	"github.com/LedbetterBlog/public-msg/database"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// ValidateOrder 验证通用订单数据(和是否重复商户订单)
func ValidateOrder(ctx context.Context, MongoDBPoolManager *database.MongoDBPoolManager, order interface{}, requiredFields map[string]bool, collectionName string) (*allStruct.ValidationResult, error) {
	var errors []string
	v := reflect.ValueOf(order)

	for fieldName, isRequired := range requiredFields {
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			errors = append(errors, fmt.Sprintf("%s field is not valid", fieldName))
			continue
		}

		if isRequired && field.String() == "" {
			errors = append(errors, fmt.Sprintf("%s cannot be empty", fieldName))
		}

		// Check phone field length if necessary
		if fieldName == "CustomerPhone" || fieldName == "BenePhone" {
			if field.Len() != 10 {
				errors = append(errors, fmt.Sprintf("%s length must be 10 digits", fieldName))
			}
		}
	}

	// Check amount
	amountField := v.FieldByName("Amount")
	if amountField.IsValid() {
		switch amountField.Kind() {
		case reflect.Float32, reflect.Float64:
			if amountField.Float() <= 0 {
				errors = append(errors, "Amount must be greater than 0")
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if amountField.Int() <= 0 {
				errors = append(errors, "Amount must be greater than 0")
			}
		default:
			errors = append(errors, "Amount field has an unsupported type")
		}
	}

	// Check if merchant_order_id already exists
	merchantOrderIDField := v.FieldByName("MerchantOrderID")
	if merchantOrderIDField.IsValid() {
		merchantOrderID := merchantOrderIDField.String()
		one, err := MongoDBPoolManager.FindOne(ctx, collectionName, bson.M{"merchant_order_id": merchantOrderID})
		if err != nil {
			return nil, fmt.Errorf(": %w", err)
		}
		if one != nil {
			errors = append(errors, "merchant_order_id already exists")
		}
	}

	if len(errors) > 0 {
		return &allStruct.ValidationResult{Valid: false, Error: strings.Join(errors, ", ")}, nil
	}
	return &allStruct.ValidationResult{Valid: true}, nil
}

// GetMidMsg 校验商户信息和签名
func GetMidMsg(ctx context.Context, redisPoolManager *database.RedisPoolManager, mid, signature string, MerchantOrderID string, Amount int) (map[string]interface{}, error) {
	// 获取商户信息(通道配置什么都在这里)
	value, err := redisPoolManager.GetValue(ctx, fmt.Sprintf("star-pay:core:mch:%s", mid))
	if err != nil {
		return nil, err
	}

	var mchConfig map[string]interface{}
	err = json.Unmarshal([]byte(value), &mchConfig)
	if err != nil {
		return nil, err
	}

	if mchConfig["loginName"] == mid {
		result := fmt.Sprintf("mid=%s&key=%s&amount=%d&merchant_order_id=%s", mid, mchConfig["secretKey"], Amount, MerchantOrderID)
		hash := md5.Sum([]byte(result))
		md5Hash := strings.ToUpper(hex.EncodeToString(hash[:]))

		if strings.ToUpper(signature) == md5Hash {
			return map[string]interface{}{}, nil
		} else {
			return map[string]interface{}{"code": 402, "msg": "Signature error, please check mid, key and merchant_order_id"}, nil
		}
	}
	return map[string]interface{}{"code": 403, "msg": "商户名或密钥错误"}, nil
}

// CallbackOrder 回调
func CallbackOrder(ctx context.Context, callbackUrl string, data map[string]interface{}) {
	// 将数据编码为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error encoding data:", err)
		return
	}

	// 打印请求的地址和数据
	log.Println("Request mch callback url:", callbackUrl)
	log.Println("Request mch callback data:", string(jsonData))

	// 创建一个新的 POST 请求
	req, err := http.NewRequestWithContext(ctx, "POST", callbackUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 使用 http.Client 发送请求
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置超时时间为10秒
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing callback url body:", err)
		}
	}(resp.Body)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return
	}

	// 打印响应体内容（可选）
	log.Println("Response body:", string(body))
}
