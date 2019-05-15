package vm

type OpenidRes struct {
	Id         int    `json:"id"`
	ActivityId string `json:"activity_id"`
	WxOpenid   string `json:"wx_openid"`
	Taskid     string `json:"taskId"`
	IsPay      int    `json:"is_pay"`
	State      string `json:"state"`
	FileName   string `json:"file_name"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type TaskRes struct {
	//Task  interface{} `bson:"task" json:"task"`
	ActivityId        string `json:"activity_id"`
	TaskId            string `json:"taskId"`
	IsPay             int    `json:"is_pay"`
	State             string `json:"state"`
	CreatedAt         string `json:"created_at"`
	FileName1         string `json:"file_name"`
	Mode              string `json:"mode"`
	OriginalTotal     int    `json:"original_total"`
	Total             int    `json:"total"`
	TemplateName      string `json:"template_name"`
	Enablegreenscreen bool   `json:"enablegreenscreen"`
}
