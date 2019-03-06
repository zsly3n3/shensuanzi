package handle

import (
	"bytes"
	"shensuanzi/commondata"
	"shensuanzi/datastruct"
	"shensuanzi/log"
	"shensuanzi/thirdparty/tls-sig-api-golang"
	"shensuanzi/tools"
)

func (app *AppHandler) AccountGenForIM(user_identifier string, appid int) datastruct.CodeType {
	privateKey, _ := tools.GetFileContentAsStringLines("datastruct/important/keys/private_key")
	publicKey, _ := tools.GetFileContentAsStringLines("datastruct/important/keys/public_key")
	var privateKey_str string
	var publicKey_str string
	var buf1 bytes.Buffer

	for _, v := range privateKey {
		buf1.WriteString(v)
		buf1.WriteString("\n")
	}
	privateKey_str = buf1.String()

	var buf2 bytes.Buffer
	for _, v := range publicKey {
		buf2.WriteString(v)
		buf2.WriteString("\n")
	}
	publicKey_str = buf2.String()

	expire := 10000
	userSig, err := TLSSigAPI.GenerateUsersigWithExpire(privateKey_str, appid, user_identifier, int64(expire))
	if err != nil {
		log.Error("GenerateUsersigWithExpire:%v", err.Error())
		return datastruct.GetDataFailed
	}
	err = TLSSigAPI.VerifyUsersig(publicKey_str, userSig, appid, user_identifier)
	if err != nil {
		log.Error("VerifyUsersig:%v", err.Error())
		return datastruct.GetDataFailed
	}
	log.Debug("userSig:%v", userSig)
	return datastruct.NULLError
}

func (app *AppHandler) GetServerInfoFromMemory() (string, bool) {
	serverInfo := commondata.GetServerInfo()
	serverInfo.RWMutex.RLock()
	defer serverInfo.RWMutex.RUnlock()
	return serverInfo.Version, serverInfo.IsMaintain
}

func (app *AppHandler) GetDirectDownloadApp() string {
	return app.dbHandler.GetDirectDownloadApp()
}

func (app *AppHandler) IsExistPhone(phone string, isFT bool) (interface{}, datastruct.CodeType) {
	return app.dbHandler.IsExistPhone(phone, isFT)
}

func (app *AppHandler) IsExistNickName(nickname string) (interface{}, datastruct.CodeType) {
	return app.dbHandler.IsExistNickName(nickname)
}
