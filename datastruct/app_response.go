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
	CreatedAt   int64   `json:"time"`
}

type RespRegisterMsg struct {
	QRCode    string `json:"qrcode"`
	GZH_Name  string `json:"gzhname"`
	CreatedAt int64  `json:"time"`
	Type      int    `json:"type"`
}

type RespRefundFTMsg struct {
	TmpOrderInfoFT
	RefundType UserOrderRefundType `json:"refundtype"`
}

type RespRightsFinishedFTMsg struct {
	TmpOrderInfoFT
	IsAgree bool `json:"isagree"`
}

type TmpOrderInfoFT struct {
	OrderId     int64  `json:"orderid"`
	NickName    string `json:"nickname"`
	ProductName string `json:"productname"`
	CreatedAt   int64  `json:"time"`
	Type        int    `json:"type"`
}

type RespRefundUserMsg struct {
	TmpOrderInfoFT
	RefundResultType UserOrderRefundResultType `json:"refundresulttype"`
}

type RespDndList struct {
	Id        int    `json:"id"`
	NickName  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	CreatedAt int64  `json:"time"`
}

type TmpProductInfo struct {
	Id          int     `json:"id"`
	ProductName string  `json:"name"`
	ProductDesc string  `json:"desc"`
	Price       float64 `json:"price"`
}

type RespProductInfo struct {
	Total        int               `json:"total"`
	OnSaleCount  int               `json:"onsalecount"`
	OffSaleCount int               `json:"offsalecount"`
	OnSale       []*TmpProductInfo `json:"onsale"`
	OffSale      []*TmpProductInfo `json:"offsale"`
}

type RespOrderForFt struct {
	DataType    int    `json:"datatype"`
	OrderId     int64  `json:"orderid"`
	ProductName string `json:"productname"`
	ProductDesc string `json:"productdesc"`
	NickName    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	CreatedAt   int64  `json:"time"`
}

type RespPurchaseOrderForFt struct {
	Order       *RespOrderForFt `json:"order"`
	IsAppraised bool            `json:"isappraised"`
}

type RespRefundingOrderForFt struct {
	Order *RespOrderForFt `json:"order"`
}

type RespRefundFinishedOrderForFt struct {
	Order      *RespOrderForFt     `json:"order"`
	RefundType UserOrderRefundType `json:"refundtype"`
}

type RespRightingOrderForFt struct {
	Order *RespOrderForFt `json:"order"`
}

type RespRightFinishedOrderForFt struct {
	Order   *RespOrderForFt `json:"order"`
	IsAgree bool            `json:"isagree"`
}

type RespFtFinance struct {
	NotCheckAmount  float64 `json:"notcheckamount"`
	CheckedAmount   float64 `json:"checkedamount"`
	TotalAmount     float64 `json:"totalamount"`     //总订单金额
	CheckedBalance  float64 `json:"checkedbalance"`  //已结算的总收益
	NotCheckBalance float64 `json:"notcheckbalance"` //未结算的总收益
	TotalBalance    float64 `json:"totalbalance"`    //总收益
	Balance         float64 `json:"balance"`         //可提现的金额
}
