package vm

type PrintRes struct {
	PrintId    string `json:"print_id"`
	//ActivityId string `json:"activity_id"`
	Address    string `json:"address"`
	Lon        string `json:"lon"`
	Lat        string `json:"lat"`
	Mark       string `json:"mark"`
	//CreateTime string `json:"create_time"`
}

//type PrintListRes struct {
//	Code int         `json:"code"`
//	List []*PrintRes `json:"list"`
//}
