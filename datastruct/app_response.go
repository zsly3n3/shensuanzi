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
	IdCard IdCardState `json:"idcard"` //身份证图标是否点亮
	IsCP   bool        `json:"cp"`     //消费保障图标是否点亮
	IsHR   bool        `json:"hr"`     //金牌推荐是否点亮
}

type RespFtLogin struct {
	Token        string      `json:"token"`
	IMPrivateKey string      `json:"imkey"`
	FtInfo       *RespFtInfo `json:"info"`
}

type RespAppraise struct {
	Name        string  `json:"name"`
	Avatar      string  `json:"avatar"`
	Score       float64 `json:"score"`
	Mark        string  `json:"mark"`
	Desc        string  `json:"desc"`
	ProductName string  `json:"productname"`
	IsAnonym    bool    `json:"isanonym"`
	CreatedAt   int64   `json:"createdat"`
}
