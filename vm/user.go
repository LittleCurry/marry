package vm

type UserRes struct {
	UserName   string `json:"user_name"`
	Openid  string `json:"openid"`


}
type CreateUserReq struct{
	NickName   string `json:"user_name"`
	//Openid  string `json:"openid"`
	Passwd  string `json:"passwd"`
}
type User struct {
	Pwd string `json:"pwd"`
	UserId string `json:"user_id"`
	CreateTime string `json:"create_time"`
}