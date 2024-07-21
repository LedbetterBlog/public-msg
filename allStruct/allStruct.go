package allStruct

// Order 结构体
type Order struct {
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

// PaymentData PaymentData结构体用于解析 JSON 数据
type PaymentData struct {
	Amount          int64  `json:"amount"`
	BeneAddress     string `json:"bene_address"`
	BeneBankAcct    string `json:"bene_bank_acct"`
	BeneEmail       string `json:"bene_email"`
	BeneIFSC        string `json:"bene_ifsc"`
	BeneName        string `json:"bene_name"`
	BenePhone       string `json:"bene_phone"`
	CreateTime      int64  `json:"create_time"`
	MerchantOrderID string `json:"merchant_order_id"`
	OrderID         string `json:"order_id"`
	Platform        string `json:"platform"`
	Status          int    `json:"status"`
	PlatformOrderId string `json:"platform_oder_id"`
	RespMsg         string `json:"resp_msg"`
}
type PaymentRespData struct {
	PlatformOrderId string `json:"platform_order_id"`
	RspMsg          string `json:"rsp_msg"`
	Code            int    `json:"code"`
}

type PayTmePaymentData struct {
	ID              string  `json:"id" bson:"_id,omitempty"`
	PlatformOrderId string  `json:"platform_order_id" bson:"platform_order_id"`
	Amount          float64 `json:"amount" bson:"amount"`
	Utr             string  `json:"utr" bson:"utr"`
	Date            string  `json:"date" bson:"date"`
	Status          int     `json:"status" bson:"status"`
}

// PayTmePayoutRequest 用于创建支付请求的结构体
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

// PayTmePayoutResponse 用于解析支付响应的结构体
type PayTmePayoutResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID string `json:"_id"`
	} `json:"data"`
}
