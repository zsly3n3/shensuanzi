package datastruct

type CodeType int //错误码
const (
	NULLError                   CodeType = iota //无错误
	ParamError                                  //参数错误,数据为空或者类型不对等
	LoginFailed                                 //登录失败,如无此账号或者密码错误等
	JsonParseFailedFromPostBody                 //来自post请求中的Body解析json失败
	GetDataFailed                               //获取数据失败
	UpdateDataFailed                            //修改数据失败
	VersionError                                //客户端与服务器版本不一致
	TokenError                                  //没有Token或者值为空,或者不存在此Token
	JsonParseFailedFromPutBody                  //来自put请求中的Body解析json失败
	WXCodeInvalid                               //无效的微信code
	PlatformInvalid                             //无效的平台参数
)

type Platform int //平台
const (
	APP Platform = iota
	H5
	PC
)

type PayPlatform int //付费平台
const (
	WXPay PayPlatform = iota
)

type Sex int

const (
	Female Sex = iota
	Male
	Secret //保密
)

type Calendar int

const (
	GongLi Calendar = iota //公历
	Nongli                 //农历
)

type AccountState int

const (
	Normal    AccountState = iota //正常
	BlackList                     //黑名单
	Freeze                        //冻结
)
