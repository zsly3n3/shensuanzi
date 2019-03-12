package datastruct

type FTRegisterBody struct {
	Phone    string   `json:"phone"`
	Pwd      string   `json:"pwd"`
	NickName string   `json:"nickname"`
	Avatar   string   `json:"avatar"`
	Mark     string   `json:"mark"` //逗号分隔
	Desc     string   `json:"desc"` //个人介绍
	Platform Platform `json:"platform"`
}

type FTRegisterWithIDBody struct {
	FTRegisterBody
	FtIdentity
}

type FtLoginBody struct {
	Phone string `json:"phone"`
	Pwd   string `json:"pwd"`
}

type UpdateFtInfoBody struct {
	Avatar   string `json:"avatar"`
	NickName string `json:"nickname"`
}

type UpdateFtMarkBody struct {
	Mark string `json:"mark"`
}

type UpdateFtIntroductionBody struct {
	Imgs []string `json:"imgs"`
	Desc string   `json:"desc"`
}

type UpdateFtAutoReplyBody struct {
	AutoReply  string   `json:"autoreply"`
	QuickReply []string `json:"quickreply"`
}

type FtIdentity struct {
	Identity      string `json:"identity"`
	ActualName    string `json:"actualname"`
	IdFrontCover  string `json:"front"`
	IdBehindCover string `json:"behind"`
}

type RemoveWithIdBody struct {
	Id int `json:"id"`
}
