package vm

type UserRes struct {
	Id           int    `json:"id"`
	UserId       string `json:"user_id"`
	Phone        string `json:"phone"`
	UserName     string `json:"user_name"`
	Gender       int    `json:"gender"`
	Birthday     string `json:"birthday"`
	Address      string `json:"address"`
	Wxopenid     string `json:"wxopenid"`
	Head         string `json:"head"`
	Introduction string `json:"introduction"`
	CreateTime   string `json:"create_time"`
	UpdateTime   string `json:"update_time"`
}
type CreateUserReq struct {
	NickName string `json:"user_name"`
	//Openid  string `json:"openid"`
	Passwd string `json:"passwd"`
}
type User struct {
	Pwd        string `json:"pwd"`
	UserId     string `json:"user_id"`
	CreateTime string `json:"create_time"`
}
