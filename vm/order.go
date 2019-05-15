package vm

type OrderRes struct {
	Openid       string `json:"openid"`
	OrderId      string `json:"order_id"`
	IsPay        int    `json:"is_pay"`
	TotalFee     string `json:"total_fee"`
	SiteId       string `json:"siteId"`
	Seller       string `json:"seller"`
	Time         string `json:"time"`
	Uid          string `json:"uid"`
	Taskid       string `json:"taskId"`
	ActivityName string `json:"activity_name"`
	ActivityId   string `json:"activity_id"`
	TaskState    string `json:"task_state"`
	FileName     string `json:"file_name"`
}

type OrderListRes struct {
	Count     int         `json:"count"`
	TotalFees float64     `json:"total_fees"`
	List      []*OrderRes `json:"list"`
}
