package datastruct

type FTRegisterBody struct {
	Phone         string `json:"phone"`
	Pwd           string `json:"pwd"`
	Identity      string `json:"identity"`
	ActualName    string `json:"actualname"`
	IdFrontCover  string `json:"front"`
	IdBehindCover string `json:"behind"`
	NickName      string `json:"nickname"`
	Mark          string `json:"mark"` //逗号分隔
	Desc          string `json:"desc"` //个人介绍
}
