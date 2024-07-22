package allStruct

// CollectOrderStruct 结构体
type CollectOrderStruct struct {
	MerchantOrderID string `json:"merchant_order_id"`
	Amount          int64  `json:"amount"`
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	CustomerEmail   string `json:"customer_email"`
}

// PaymentOrderStruct 结构体
type PaymentOrderStruct struct {
	MerchantOrderID string `json:"merchant_order_id"`
	BenePhone       string `json:"bene_phone"`
	Amount          int64  `json:"amount"`
	BeneName        string `json:"bene_name"`
	BeneEmail       string `json:"bene_email"`
	BeneIFSC        string `json:"bene_ifsc"`
	BeneAddress     string `json:"bene_address"`
	BeneBankAcct    string `json:"bene_bank_acct"`
}

// PayTmeCallbackData 结构体
type PayTmeCallbackData struct {
	Status        string  `json:"status"`
	TransactionId string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Date          string  `json:"date"`
	Utr           string  `json:"rrn"`
}

// ValidationResult 结构体
type ValidationResult struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// RedisCollectOrderDataStruct redisOrderDataStruct结构体用于解析 JSON 数据
type RedisCollectOrderDataStruct struct {
	OrderID         string `json:"order_id"`
	MerchantNumber  string `json:"merchant_number"`
	CreateTime      int64  `json:"create_time"`
	MerchantOrderID string `json:"merchant_order_id"`
	Amount          int64  `json:"amount"`
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	CustomerEmail   string `json:"customer_email"`
	Platform        string `json:"platform"`
	Status          int    `json:"status"`
	PlatformOrderId string `json:"platform_oder_id"`
	RespMsg         string `json:"resp_msg"`
}

// PayTmeCollectRespData 定义PayTme代收整体响应结构体
type PayTmeCollectRespData struct {
	PlatformOrderId string `json:"platform_order_id"`
	UPI             string `json:"upi_url"`
	RspMsg          string `json:"rsp_msg"`
	Code            int    `json:"code"`
}

// RedisPaymentOrderDataStruct redisOrderDataStruct结构体用于解析 JSON 数据
type RedisPaymentOrderDataStruct struct {
	Amount          int64  `json:"amount"`
	BeneAddress     string `json:"bene_address"`
	BeneBankAcct    string `json:"bene_bank_acct"`
	BeneEmail       string `json:"bene_email"`
	BeneIFSC        string `json:"bene_ifsc"`
	BeneName        string `json:"bene_name"`
	BenePhone       string `json:"bene_phone"`
	CreateTime      int64  `json:"create_time"`
	MerchantNumber  string `json:"merchant_number"`
	MerchantOrderID string `json:"merchant_order_id"`
	OrderID         string `json:"order_id"`
	Platform        string `json:"platform"`
	Status          int    `json:"status"`
	PlatformOrderId string `json:"platform_oder_id"`
	RespMsg         string `json:"resp_msg"`
}

// PayTmePaymentRespData 定义PayTme代付整体响应结构体
type PayTmePaymentRespData struct {
	PlatformOrderId string `json:"platform_order_id"`
	RspMsg          string `json:"rsp_msg"`
	Code            int    `json:"code"`
}

// PayTmePaymentData 请求PayTme代付的请求结构体
type PayTmePaymentData struct {
	ID              string  `json:"id" bson:"_id,omitempty"`
	MchNumber       string  `json:"mch_number" bson:"mch_number"`
	MchOrderID      string  `json:"mch_order_id" bson:"mch_order_id"`
	Platform        string  `json:"platform" bson:"platform"`
	PlatformOrderId string  `json:"platform_order_id" bson:"platform_order_id"`
	Amount          float64 `json:"amount" bson:"amount"`
	Utr             string  `json:"utr" bson:"utr"`
	CreateTime      int64   `json:"create_time" bson:"create_time"`
	UpdateTime      int64   `json:"update_time" bson:"update_time"`
	Status          int     `json:"status" bson:"status"`
}

// PayTmePayInRequest 用于创建PayTme支付请求的结构体
type PayTmePayInRequest struct {
	Amount                float64 `json:"amount"`
	Name                  string  `json:"name"`
	Email                 string  `json:"email"`
	UserContactNumber     string  `json:"userContactNumber"`
	MerchantTransactionID string  `json:"merchantTransactionId"`
}

// PayTmePayoutRequest 用于创建PayTme支付请求的结构体
type PayTmePayoutRequest struct {
	Amount                float64 `json:"amount"`
	Name                  string  `json:"name"`
	Email                 string  `json:"email"`
	Phone                 string  `json:"phone"`
	AccountNumber         string  `json:"accountNumber"`
	BankIfsc              string  `json:"bankIfsc"`
	AccountHolderName     string  `json:"accountHolderName"`
	BankName              string  `json:"bankName"`
	UPI                   string  `json:"upi"`
	Purpose               string  `json:"purpose"`
	MerchantTransactionID string  `json:"merchantTransactionId"`
}

// PayTmePayOutResponse 用于解析PayTme支付响应的结构体
type PayTmePayOutResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID string `json:"_id"`
	} `json:"data"`
}

// PayTmePayInRespData 定义PayTme代收响应结构体
type PayTmePayInRespData struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    struct {
		UpiUrl          string `json:"upiurl"`
		PlatformOrderId string `json:"transaction_id"`
	} `json:"data"`
}

// CreateCollectOrderResp 定义代收创建订单响应结构体
type CreateCollectOrderResp struct {
	Code            int    `json:"code"`
	Message         string `json:"message"`
	MerchantOrderId string `json:"merchant_order_id"`
	PlatformOrderId string `json:"platform_order_id"`
	PaymentLink     string `json:"payment_link"`
	UpiLink         string `json:"upi_link"`
}
