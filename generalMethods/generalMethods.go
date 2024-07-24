package generalMethods

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/LedbetterBlog/public-msg/allStruct"
	"github.com/LedbetterBlog/public-msg/database"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
)

// ValidateOrder 验证通用订单数据(和是否重复商户订单)
func ValidateOrder(MongoDBPoolManager *database.MongoDBPoolManager, order interface{}, requiredFields map[string]bool, mid string, collectionName string) (*allStruct.ValidationResult, error) {
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
		one, _ := MongoDBPoolManager.FindOne(collectionName, bson.M{"merchant_order_id": merchantOrderID, "mch_number": mid})
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
