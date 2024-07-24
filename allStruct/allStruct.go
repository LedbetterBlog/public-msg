package allStruct

// CollectOrderStruct 结构体
type CollectOrderStruct struct {
	MerchantOrderID string `json:"merchant_order_id"`
	Amount          int    `json:"amount"`
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	CustomerEmail   string `json:"customer_email"`
}

// PaymentOrderStruct 结构体
type PaymentOrderStruct struct {
	MerchantOrderID string `json:"merchant_order_id"`
	BenePhone       string `json:"bene_phone"`
	Amount          int    `json:"amount"`
	BeneName        string `json:"bene_name"`
	BeneEmail       string `json:"bene_email"`
	BeneIFSC        string `json:"bene_ifsc"`
	BeneAddress     string `json:"bene_address"`
	BeneBankAcct    string `json:"bene_bank_acct"`
}

// PayTmePaymentCallbackData 结构体
type PayTmePaymentCallbackData struct {
	Status        string  `json:"status"`
	TransactionId string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Date          string  `json:"date"`
	Utr           string  `json:"rrn"`
}

// PayTmeCollectCallbackData 结构体
type PayTmeCollectCallbackData struct {
	Status        int     `json:"status"`
	TransactionId string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Utr           string  `json:"rrn"`
}

// ValidationResult 结构体
type ValidationResult struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// MongoDbLocalStatusStruct 存入商户提交的订单数据，主要用来识别是否重复订单号
type MongoDbLocalStatusStruct struct {
	OrderID         string `json:"order_id" bson:"_id,omitempty"`
	MerchantNumber  string `json:"merchant_number" bson:"mch_number"`
	MerchantOrderID string `json:"merchant_order_id " bson:"merchant_order_id"`
	CreateTime      int64  `json:"create_time" bson:"create_time"`
	Amount          int    `json:"amount " bson:"amount"`
	Platform        string `json:"platform" bson:"platform"`
	Status          int    `json:"status" bson:"status"`
	CallbackStatus  int    `json:"callback_status" bson:"callback_status"`
}

// RedisCollectOrderDataStruct redisOrderDataStruct结构体用于解析 JSON 数据
type RedisCollectOrderDataStruct struct {
	OrderID         string `json:"order_id"`
	MerchantNumber  string `json:"merchant_number"`
	CreateTime      int64  `json:"create_time"`
	MerchantOrderID string `json:"merchant_order_id "`
	Amount          int    `json:"amount "`
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone" `
	CustomerEmail   string `json:"customer_email"`
	Platform        string `json:"platform"`
	Status          int    `json:"status"`
	PlatformOrderId string `json:"platform_oder_id"`
	RespMsg         string `json:"resp_msg"`
	UTR             string `json:"utr"`
	CallbackStatus  int    `json:"callback_status"`
}

// PayTmeCollectRespData 定义PayTme代收整体响应结构体
type PayTmeCollectRespData struct {
	PlatformOrderId string `json:"platform_order_id"`
	UPI             string `json:"upi_url"`
	RespMsg         string `json:"resp_msg"`
	Code            int    `json:"code"`
}

// RedisPaymentOrderDataStruct redisOrderDataStruct结构体用于解析 JSON 数据
type RedisPaymentOrderDataStruct struct {
	Amount          int    `json:"amount"`
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
	UTR             string `json:"utr"`
	CallbackStatus  int    `json:"callback_status"`
}

// PayTmePaymentRespData 定义PayTme代付整体响应结构体
type PayTmePaymentRespData struct {
	PlatformOrderId string `json:"platform_order_id"`
	RespMsg         string `json:"resp_msg"`
	Code            int    `json:"code"`
}

// PayTmeOrderStatus 定义PayTme代收状态查询返回结构
type PayTmeOrderStatus struct {
	PlatformOrderId string `json:"platform_order_id"`
	RespMsg         string `json:"resp_msg"`
	Code            int    `json:"code"`
	Utr             string `json:"utr"`
	Amount          int    `json:"amount"`
	Status          string `json:"status"`
}

// PayTmePaymentData 请求PayTme代付的请求结构体
type PayTmePaymentData struct {
	ID              string  `json:"id" bson:"_id,omitempty"`
	MchNumber       string  `json:"mch_number" bson:"mch_number"`
	MchOrderID      string  `json:"merchant_order_id" bson:"merchant_order_id"`
	Platform        string  `json:"platform" bson:"platform"`
	PlatformOrderId string  `json:"platform_order_id" bson:"platform_order_id"`
	Amount          float64 `json:"amount" bson:"amount"`
	Utr             string  `json:"utr" bson:"utr"`
	CreateTime      int64   `json:"create_time" bson:"create_time"`
	UpdateTime      int64   `json:"update_time" bson:"update_time"`
	Status          int     `json:"status" bson:"status"`
	CallbackStatus  int     `json:"callback_status" bson:"callback_status"`
}

// PayTmePayData 代收存入mongo的请求结构体
type PayTmePayData struct {
	ID              string `json:"id" bson:"_id,omitempty"`
	MchNumber       string `json:"mch_number" bson:"mch_number"`
	MchOrderID      string `json:"merchant_order_id" bson:"merchant_order_id"`
	Platform        string `json:"platform" bson:"platform"`
	PlatformOrderId string `json:"platform_order_id" bson:"platform_order_id"`
	Amount          int    `json:"amount" bson:"amount"`
	Utr             string `json:"utr" bson:"utr"`
	RespMsg         string `json:"resp_msg" bson:"resp_msg"`
	CreateTime      int64  `json:"create_time" bson:"create_time"`
	UpdateTime      int64  `json:"update_time" bson:"update_time"`
	Status          int    `json:"status" bson:"status"`
	CallbackStatus  int    `json:"callback_status" bson:"callback_status"`
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

// PayTmePayInStatusResponse 用于解析PayTme代收订单状态响应的结构体
type PayTmePayInStatusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Status                string `json:"status"`
		UserContactNumber     string `json:"userContactNumber"`
		MerchantTransactionId string `json:"merchantTransactionId"`
		Amount                int    `json:"amount"`
		MerchantId            string `json:"merchantId"`
		CreatedAt             string `json:"createdAt"`
		Utr                   string `json:"rrn"`
	} `json:"data"`
}

// PayTmePayOutStatusResponse 用于解析PayTme代付订单状态响应的结构体
type PayTmePayOutStatusResponse struct {
	Status string `json:"status"`
	Data   struct {
		Status     string `json:"status"`
		Amount     int    `json:"amount"`
		MerchantId string `json:"merchantId"`
		Utr        string `json:"rrn"`
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
