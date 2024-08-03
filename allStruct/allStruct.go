package allStruct

// PayInOrderStruct 结构体
type PayInOrderStruct struct {
	MerchantOrderID string `json:"merchant_order_id"`
	Amount          int    `json:"amount"`
	CustomerName    string `json:"customer_name"`
	CustomerPhone   string `json:"customer_phone"`
	CustomerEmail   string `json:"customer_email"`
}

// PayOutOrderStruct 结构体
type PayOutOrderStruct struct {
	MerchantOrderID string `json:"merchant_order_id"`
	BenePhone       string `json:"bene_phone"`
	Amount          int    `json:"amount"`
	BeneName        string `json:"bene_name"`
	BeneEmail       string `json:"bene_email"`
	BeneIFSC        string `json:"bene_ifsc"`
	BeneAddress     string `json:"bene_address"`
	BeneBankAcct    string `json:"bene_bank_acct"`
}

// PayTmePayOutCallbackData 结构体
type PayTmePayOutCallbackData struct {
	Status        string  `json:"status"`
	TransactionId string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
	Date          string  `json:"date"`
	Utr           string  `json:"rrn"`
}

// PayInCallbackData 结构体
type PayInCallbackData struct {
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
	OrderID         string  `json:"order_id" bson:"_id,omitempty"`
	MerchantNumber  string  `json:"merchant_number" bson:"mch_number"`
	MerchantOrderID string  `json:"merchant_order_id " bson:"merchant_order_id"`
	CreateTime      int64   `json:"create_time" bson:"create_time"`
	Amount          int     `json:"amount " bson:"amount"`
	UserName        string  `json:"user_name" bson:"user_name"`
	UserPhone       string  `json:"user_phone" bson:"user_phone"`
	UserEmail       string  `json:"user_email" bson:"user_email"`
	Platform        string  `json:"platform" bson:"platform"`
	BankIFSC        string  `json:"bank_ifsc" bson:"bank_ifsc"`
	UserBankAcct    string  `json:"user_bank_acct" bson:"user_bank_acct"`
	UserAddress     string  `json:"user_address" bson:"user_address"`
	Chnl            int     `json:"chnl" bson:"chnl"`
	ChnlFeeRatio    float64 `json:"chnl_fee_ratio" bson:"chnl_fee_ratio"`
	MchFeeRatio     float64 `json:"mch_fee_ratio" bson:"mch_fee_ratio"`
	Status          int     `json:"status" bson:"status"`
	CallbackStatus  int     `json:"callback_status" bson:"callback_status"`
	OrderType       int     `json:"order_type" bson:"order_type"`
	MchSettleAmount float64 `json:"mch_settle_amount" bson:"mch_settle_amount"`
}

// RedisPayInOrderDataStruct redisOrderDataStruct结构体用于解析 JSON 数据
type RedisPayInOrderDataStruct struct {
	OrderID         string  `json:"order_id"`
	MerchantNumber  string  `json:"merchant_number"`
	CreateTime      int64   `json:"create_time"`
	MerchantOrderID string  `json:"merchant_order_id "`
	Amount          int     `json:"amount "`
	UserName        string  `json:"user_name"`
	UserPhone       string  `json:"user_phone" `
	UserEmail       string  `json:"user_email"`
	Platform        string  `json:"platform"`
	Status          int     `json:"status"`
	PlatformOrderId string  `json:"platform_oder_id"`
	RespMsg         string  `json:"resp_msg"`
	UTR             string  `json:"utr"`
	CallbackStatus  int     `json:"callback_status"`
	Chnl            int     `json:"chnl"`
	ChnlFeeRatio    float64 `json:"chnl_fee_ratio"`
	MchFeeRatio     float64 `json:"mch_fee_ratio"`
}

//{
//    "_id": ""订单号,
//    "mch_number": ""商户名字,
//    "merchant_order_id": ""商户订单Id,
//    "create_time": NumberLong("")订单创建时间,
//    "amount": NumberInt("")金额,
//    "platform": ""平台,
//    "status": NumberInt("")订单状态（0-支付成功，1-支付失败，2-支付超时（有回调一样会变成成功），5-本地创建订单成功，6-已经提交到三方）,
//    "callback_status": NumberInt("")回调商户状态（0-未回调，1-回调）,
//    "order_type": NumberInt("")订单类型（0-代收，1-代付）,
//    "platform_order_id": ""三方平台订单id,
//    "resp_msg": ""三方回复信息,
//    "chnl":""通道,
//    "chnl_fee_ratio":""通道费率,
//    "user_name":""用户名称,
//    "user_phone":""用户手机号,
//    "user_email":""用户邮箱,
//    "ifsc":""银行代号,
//    "bank_acc":""银行账号,
//    "payment_link":""系统支付链接,
//    "upi_link":""UPI支付链接,
//    "mch_fee_ratio":""商户费率,
//    "update_time": NumberLong("")更新记录时间
//}

// PayTmePayInRespData 定义PayTme代收整体响应结构体
type PayTmePayInRespData struct {
	PlatformOrderId string `json:"platform_order_id"`
	UPI             string `json:"upi_url"`
	RespMsg         string `json:"resp_msg"`
	Code            int    `json:"code"`
}

// RedisPayOutOrderDataStruct redisOrderDataStruct结构体用于解析 JSON 数据
type RedisPayOutOrderDataStruct struct {
	Amount          int     `json:"amount"`
	UserAddress     string  `json:"user_address"`
	UserBankAcct    string  `json:"user_bank_acct"`
	UserEmail       string  `json:"user_email"`
	BankIFSC        string  `json:"bank_ifsc"`
	UserName        string  `json:"user_name"`
	UserPhone       string  `json:"user_phone"`
	CreateTime      int64   `json:"create_time"`
	MerchantNumber  string  `json:"merchant_number"`
	MerchantOrderID string  `json:"merchant_order_id"`
	OrderID         string  `json:"order_id"`
	Platform        string  `json:"platform"`
	Status          int     `json:"status"`
	PlatformOrderId string  `json:"platform_oder_id"`
	RespMsg         string  `json:"resp_msg"`
	UTR             string  `json:"utr"`
	CallbackStatus  int     `json:"callback_status"`
	Chnl            int     `json:"chnl"`
	ChnlFeeRatio    float64 `json:"chnl_fee_ratio"`
	MchFeeRatio     float64 `json:"mch_fee_ratio"`
}

// PayTmePayOutRespData 定义PayTme代付整体响应结构体
type PayTmePayOutRespData struct {
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

// ProcessOrderStatusData 状态返回结构体
type ProcessOrderStatusData struct {
	Utr    string `json:"utr"`
	Status int    `json:"status"`
}

// PayTmePayInRequest 用于创建PayTme支付请求的结构体
type PayTmePayInRequest struct {
	Amount                float64 `json:"amount"`
	Name                  string  `json:"name"`
	Email                 string  `json:"email"`
	UserContactNumber     string  `json:"userContactNumber"`
	MerchantTransactionID string  `json:"merchantTransactionId"`
}

// PayTmePayOutRequest 用于创建PayTme支付请求的结构体
type PayTmePayOutRequest struct {
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
		ID string `json:"id"`
	} `json:"data"`
}

// PayTmePayInResponse 定义PayTme代收响应结构体
type PayTmePayInResponse struct {
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

// CreatePayInOrderResp 定义代收创建订单响应结构体
type CreatePayInOrderResp struct {
	Code            int    `json:"code"`
	Message         string `json:"message"`
	MerchantOrderId string `json:"merchant_order_id"`
	PlatformOrderId string `json:"platform_order_id"`
	PayOutLink      string `json:"payout_link"`
	UpiLink         string `json:"upi_link"`
}
