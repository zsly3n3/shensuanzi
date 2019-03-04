package datastruct

/*用户信息冷数据*/
type ColdUserInfo struct {
	Id           int64    `xorm:"not null pk autoincr bigint COMMENT('自增编号')"`
	Phone        string   `xorm:"VARCHAR(20) not null COMMENT('手机号')"`
	Pwd          string   `xorm:"VARCHAR(30) not null COMMENT('密码')"`
	NickName     string   `xorm:"VARCHAR(100) not null COMMENT('昵称')"`
	Avatar       string   `xorm:"VARCHAR(500) not null COMMENT('头像')"`
	Sex          Sex      `xorm:"TINYINT(1) not null COMMENT('性别')"`
	ActualName   string   `xorm:"VARCHAR(100) not null COMMENT('真实姓名')"`
	DateBirth    string   `xorm:"VARCHAR(16) not null COMMENT('出生日期')"`
	BirthPlace   string   `xorm:"VARCHAR(100) not null COMMENT('出生地')"`
	Calendar     Calendar `xorm:"TINYINT(1) not null COMMENT('0是公历,1是农历')"`
	Registration Platform `xorm:"TINYINT(1) not null COMMENT('注册平台')"`
	CreatedAt    int64    `xorm:"bigint not null COMMENT('注册时间')"`
}

/*用户信息热数据*/
type HotUserInfo struct {
	Id            int64        `xorm:"bigint not null pk COMMENT('用户Id')"`
	Token         string       `xorm:"CHAR(40) not null COMMENT('令牌值')"`
	LoginTime     int64        `xorm:"bigint not null COMMENT('登录时间')"`
	GoldCount     int64        `xorm:"bigint not null COMMENT('当前金币数')"`
	DepositTotal  float64      `xorm:"decimal(16,2) not null COMMENT('总充值额')"`
	Balance       float64      `xorm:"decimal(16,2) not null COMMENT('当前佣金额')"`
	BalanceTotal  float64      `xorm:"decimal(16,2) not null COMMENT('总佣金额')"`
	PayGoldTotal  float64      `xorm:"decimal(16,2) not null COMMENT('金币消费总额')"`
	PurchaseTotal float64      `xorm:"decimal(16,2) not null COMMENT('直接购买消费总额')"`
	AccountState  AccountState `xorm:"TINYINT(1) not null COMMENT('账户状态,0是正常,1是黑名单,2是冻结')"`
	IMPrivateKey  string       `xorm:"CHAR(255) not null COMMENT('IM私钥')"`
}

/*微信用户信息*/
type WXUserInfo struct {
	Id               int64  `xorm:"bigint not null pk COMMENT('用户Id')"`
	UUID             string `xorm:"CHAR(50) not null COMMENT('微信UUID')"`
	PayOpenidForGzh  string `xorm:"CHAR(50) not null COMMENT('微信公众号OpenId')"`
	PayOpenidForKfpt string `xorm:"CHAR(50) not null COMMENT('微信开放平台OpenId')"`
	PayeeOpenid      string `xorm:"CHAR(50) not null COMMENT('提现OpenId')"`
	GzhAppid         string `xorm:"CHAR(22) not null COMMENT('微信公众号Appid')"`
	KfptAppid        string `xorm:"CHAR(22) not null COMMENT('微信开放平台Appid')"`
}

/*广告信息表*/
type AdInfo struct {
	Id       int      `xorm:"not null pk autoincr INT(4) COMMENT('自增编号')"`
	ImgUrl   string   `xorm:"VARCHAR(255) not null COMMENT('广告图片')"`
	IsJumpTo bool     `xorm:"TINYINT(1) not null COMMENT('是否跳转')"`
	JumpTo   string   `xorm:"VARCHAR(255) not null COMMENT('跳转地址')"`
	Platform Platform `xorm:"TINYINT(1) not null COMMENT('展示平台')"`
	IsHidden bool     `xorm:"TINYINT(1) not null COMMENT('是否隐藏')"`
}
