package allStruct

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
