package vm

type NotifyRes struct {
	Attach     string `xml:"attach"`
	MchId      string `xml:"mch_id"`
	OutTradeNo string `xml:"out_trade_no"`
	Sign       string `xml:"sign"`
	TotalFee   string `xml:"total_fee"`
	NonceStr   string `xml:"nonce_str"`
	SignType   string `xml:"sign_type"`
	TimeEnd    string `xml:"time_end"`

	AppId         string `xml:"app_id"`
	DeviceInfo    string `xml:"device_info"`
	ResultCode    string `xml:"result_code"`
	ErrCode       string `xml:"err_code"`
	ErrCodeDes    string `xml:"err_code_des"`
	Openid        string `xml:"openid"`
	TradeType     string `xml:"trade_type"`
	BankType      string `xml:"bank_type"`
	FeeType       string `xml:"fee_type"`
	CashFeeType   string `xml:"cash_fee_type"`
	TransactionId string `xml:"transaction_id"`
}

type ReturnNotifyRes struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

type WxPayReq struct {
	Openid     string   `json:"openid"`
	TotalFee   string   `json:"total_fee"`
	SiteId     string   `json:"siteId"`
	Taskid     string   `json:"taskId"`
	PrintId    string   `json:"print_id"`
	Uid        string   `json:"uid"`
	ActivityId string   `json:"activity_id"`
	Templates  []string `json:"templates"`
	Print      []string `json:"print"`
}

type WechatPayResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	Appid      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	PrepayId   string `xml:"prepay_id"`
	TradeType  string `xml:"trade_type"`
}

type WxPayRes struct {
	AppId     string `json:"appid"`
	TimeStamp string `json:"time_stamp"`
	NonceStr  string `json:"nonce_str"`
	Package   string `json:"package"`
	SignType  string `json:"sign_type"`
	Sign      string `json:"sign"`
}
