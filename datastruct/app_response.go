package datastruct

type RespFtInfo struct {
	Id           int          `json:"id"`
	AccountState AccountState `json:"accountstate"`
	NickName     string       `json:"nickname"`
	Avatar       string       `json:"avatar"`
	Mark         string       `json:"mark"`       //逗号分隔
	EnableFree   bool         `json:"enablefree"` //ui状态
	AuthIcon     *AuthIcon    `json:"authicon"`
}

type AuthIcon struct {
	IdCard bool `json:"idcard"` //身份证图标是否点亮
	IsCP   bool `json:"cp"`     //消费保障图标是否点亮
	IsHR   bool `json:"hr"`     //金牌推荐是否点亮
}

type RespFtLogin struct {
	Token        string      `json:"token"`
	IMPrivateKey string      `json:"imkey"`
	FtInfo       *RespFtInfo `json:"info"`
}
